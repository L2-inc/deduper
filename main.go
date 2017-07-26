package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type aspect struct {
	name string
	size int64
}

type trait struct {
	size  int64
	paths []string
}

var totalDupes int
var spaceSaved int64
var deletePrefix *string
var report, quiet *bool
var allFiles int
var totalSize int64
var similarFiles map[aspect][]string

func validateDirs(dirs []string) {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("invalid dir %v", err)
			os.Exit(2)
		}
	}
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

func (t trait) confirmDupes(quiet bool) bool {
	if t.size == 0 || len(t.paths) < 2 {
		return false
	}
	md5sums := hardID(t.paths)
	uniqueSums := len(md5sums)
	if uniqueSums != 1 {
		if quiet {
			return false
		}
		fmt.Printf(" expect exactly 1 md5sum but found %d with size %d\n", uniqueSums, t.size)
		for s, p := range md5sums {
			fmt.Printf("\t %s\n", s)
			for _, path := range p {
				fmt.Printf("\t\t%s\n", path)
			}
		}
		return false
	}
	return true
}

func (t trait) deleteDupes(verbose bool, prefix string) int {
	toDelete := []string{}
	for i, p := range t.paths {
		if verbose {
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
		fmt.Printf(" deleting copy %d at %s\n", i, p)
	}
	return len(toDelete)
}

func walker(path string, f os.FileInfo, err error) error {
	if f == nil {
		fmt.Printf("Invalid path %s\n", path)
		return nil
	}
	if f.IsDir() {
		return nil
	}
	if f.Mode()&os.ModeSymlink != 0 {
		return nil
	}
	s := aspect{f.Name(), f.Size()}
	if s.size == 0 {
		return nil
	}
	totalSize += s.size
	allFiles++
	similarFiles[s] = append(similarFiles[s], path)
	return nil
}

func processArgs() {
	deletePrefix = flag.String("delete-prefix", "", "delete dupes that start with this prefix")
	report = flag.Bool("report", false, "print out report only.  This is on unless 'delete-prefix' flag is specified")
	quiet = flag.Bool("quiet", false, "minimal output")
	flag.Parse()
	if *deletePrefix != "" && !*quiet {
		*report = true
	}
}

func compileData() {
	similarFiles = make(map[aspect][]string)
	for _, dir := range flag.Args() {
		filepath.Walk(dir, walker)
	}
}

func reportStats() {
	fmt.Printf("\nTotal dupes %d.  Total bytes wasted %d\n", totalDupes, spaceSaved)
	fmt.Printf("\nTotal files %d.  Total bytes %d\n", allFiles, totalSize)
}

func main() {
	processArgs()
	validateDirs(flag.Args())

	compileData()
	for a, paths := range similarFiles {
		t := trait{a.size, paths}
		if !t.confirmDupes(*quiet) {
			continue
		}
		deleted := t.deleteDupes(*report, *deletePrefix)
		totalDupes += deleted
		spaceSaved += int64(deleted) * a.size
	}
	reportStats()
}
