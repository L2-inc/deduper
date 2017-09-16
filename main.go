package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
)

type aspect struct {
	name string
	size int64
}

type trait struct {
	size  int64
	paths []string
}

type cmdOpt struct {
	forReal        bool
	report         bool
	deletePrefix   string
	ignoreSuffixes string
	dirs           []string
}

func validateDirs(dirs []string) bool {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("invalid dir %v", err)
			return false
		}
	}
	return true
}

func hardID(paths []string) map[string][]string {
	md5sum := make(map[string][]string)
	for _, path := range paths {
		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			fmt.Printf("cannot open %v. %v\n", path, err)
			continue
		}
		h := md5.New()

		if _, err := io.Copy(h, f); err != nil {
			fmt.Printf("cannot compute md5 for %v. %v", path, err)
			continue
		}
		sum := fmt.Sprintf("%x", h.Sum(nil))
		md5sum[sum] = append(md5sum[sum], path)
	}
	return md5sum
}

func (t trait) confirmDupes(forReal bool) bool {
	if t.size == 0 || len(t.paths) < 2 {
		return false
	}
	md5sums := hardID(t.paths)
	uniqueSums := len(md5sums)
	if uniqueSums == 1 {
		return true
	} else if forReal {
		return false
	}
	fmt.Printf(" expect exactly 1 md5sum but found %d out of %d with size %d bytes\n", uniqueSums, len(t.paths), t.size)
	for s, p := range md5sums {
		fmt.Printf("\t %s\n", s)
		for _, path := range p {
			fmt.Printf("\t\t%s\n", path)
		}
	}
	return false
}

func (t trait) purge(reportOnly bool, prefix string, rm func(string) error) int {
	toDelete := []string{}
	for i, p := range t.paths {
		if reportOnly {
			fmt.Printf(" duplicate %d: %s\n", i, p)
		}
		if prefix != "" && strings.HasPrefix(p, prefix) {
			toDelete = append(toDelete, p)
		}
	}
	if len(toDelete) == len(t.paths) {
		fmt.Println("delete prefix needs to be more restrictive.  all copies of a file are")
		for _, p := range toDelete {
			fmt.Printf("\t%s\n", p)
		}
		return 0
	}
	for i, p := range toDelete {
		fmt.Printf(" deleting copy %d at '%s'\n", i, p)
		if reportOnly {
			continue
		}
		if err := rm(p); err != nil {
			panic(err)
		}
	}
	return len(toDelete)
}

func processArgs() cmdOpt {
	deletePrefix := flag.String("delete-prefix", "", "delete dupes that start with this prefix")
	report := flag.Bool("report", false, "print out report only.  This is on if 'delete-prefix' flag is omitted.  If on, nothing is deleted.")
	forReal := flag.Bool("for-real", false, "minimal output; dry-run without this")
	ignoreSuffixes := flag.String("ignore-suffixes", "", "ignore all files with these comma separated suffixes")
	flag.Parse()
	if *deletePrefix != "" && !*forReal {
		*report = true
	}
	work := flag.Args()
	return cmdOpt{*forReal, *report, *deletePrefix, *ignoreSuffixes, work}
}

func (c cmdOpt) compileData() (size int64, count int, simFiles map[aspect][]string) {
	simFiles = make(map[aspect][]string)
	suffixes := strings.Split(c.ignoreSuffixes, ",")
	for _, dir := range c.dirs {
		filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if f == nil {
				fmt.Printf("Invalid path %s\n", path)
				return nil
			}
			for _, suffix := range suffixes {
				if strings.HasSuffix(path, suffix) {
					return nil
				}
			}
			s := aspect{f.Name(), f.Size()}
			if f.IsDir() || f.Mode()&os.ModeSymlink != 0 || s.size == 0 {
				return nil
			}
			size += s.size
			count++
			simFiles[s] = append(simFiles[s], path)
			return nil
		})
	}
	return size, count, simFiles
}

func reportStats(all int, size int64, dupes int, saved int64) {
	fmt.Printf("\n%d dupe files deleted.  Total saved %s\n", dupes,
		humanize.Bytes(uint64(saved)))
	fmt.Printf("\nTotal files %d.  Total size %s\n", all, humanize.Bytes(uint64(size)))
}

func (c cmdOpt) doWork() (allFiles int, totalSize int64, totalDupes int, spaceSaved int64) {
	totalSize, allFiles, similarFiles := c.compileData()
	for a, paths := range similarFiles {
		t := trait{a.size, paths}
		if !t.confirmDupes(c.forReal) {
			continue
		}
		deleted := t.purge(c.report, c.deletePrefix, os.Remove)
		totalDupes += deleted
		spaceSaved += int64(deleted) * a.size
	}
	return allFiles, totalSize, totalDupes, spaceSaved
}

func main() {
	options := processArgs()
	if !validateDirs(options.dirs) {
		os.Exit(2)
	}
	if 0 == len(options.dirs) {
		fmt.Printf("Usage: ./deduper <options> <folder>\n\n")
		flag.Usage()
		os.Exit(0)
	}

	reportStats(options.doWork())
}
