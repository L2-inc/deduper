package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type softID struct {
	name string
	size int64
}

var totalDupes int
var wastedSpace int64
var deletePrefix *string
var report *bool

func validateDirs() {
	for _, dir := range flag.Args() {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("invalid dir %v", err)
			os.Exit(2)
		}
	}

}

func confirmDupes(k softID, files []string) {
	hardID := make(map[string][]string)
	for _, path := range files {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		h := md5.New()

		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}
		f.Close()
		md5 := string(h.Sum(nil)[:])
		hardID[md5] = append(hardID[md5], path)
	}
	for _, paths := range hardID {
		toDelete := []string{}
		if 2 > len(paths) {
			continue
		}
		totalDupes--
		wastedSpace = wastedSpace - k.size
		for i, p := range paths {
			totalDupes++
			wastedSpace += k.size
			if *report {
				fmt.Printf(" duplicate %d: %s\n", i, p)
			}
			if *deletePrefix != "" && strings.HasPrefix(p, *deletePrefix) {
				toDelete = append(toDelete, p)
			}
		}
		if len(toDelete) == len(paths) {
			fmt.Println("delete prefix needs to be more restrictive.  all copies of a file are")
			for _, p := range toDelete {
				fmt.Printf("\t%s\n", p)
			}
			continue
		}
		for i, p := range toDelete {
			fmt.Printf(" deleting copy %d at %s\n", i, p)
		}
	}
}

func main() {
	deletePrefix = flag.String("delete-prefix", "", "delete dupes that start with this prefix")
	report = flag.Bool("report", false, "print out report only.  This is on unless 'delete-prefix' flag is specified")
	flag.Parse()
	if *deletePrefix != "" {
		*report = true
	}

	validateDirs()

	var allFiles int
	var allDirs int
	var totalSize int64
	uniqueFiles := make(map[softID][]string)
	for _, dir := range flag.Args() {
		filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if f == nil {
				fmt.Printf("Invalid path %s\n", path)
				return nil
			}
			if !f.IsDir() {
				s := softID{f.Name(), f.Size()}
				if s.size == 0 {
					return nil
				}
				totalSize += s.size
				allFiles++
				uniqueFiles[s] = append(uniqueFiles[s], path)
			}
			allDirs++
			return nil
		})
	}
	for k, v := range uniqueFiles {
		if k.size == 0 || len(v) < 2 {
			continue
		}
		confirmDupes(k, v)
	}
	fmt.Printf("\nTotal dupes %d.  Total bytes wasted %d\n", totalDupes, wastedSpace)
	fmt.Printf("\nTotal files %d.  Total bytes %d. Total dirs %d\n", allFiles,
		totalSize, allDirs)
}
