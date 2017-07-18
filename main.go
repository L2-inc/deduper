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

type sameFile struct {
	name string
	size int64
}

func main() {
	deletePrefix := flag.String("delete-prefix", "", "delete dupes that start with this prefix")
	flag.Parse()
	for _, dir := range flag.Args() {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("invalid dir %v", err)
			os.Exit(2)
		}
	}
	sameFiles := make(map[sameFile][]string)
	all_files := 0
	all_dirs := 0
	var all_file_size int64
	for _, dir := range flag.Args() {
		filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if f == nil {
				fmt.Printf("Invalid path %s\n", path)
				return nil
			}
			if !f.IsDir() {
				s := sameFile{f.Name(), f.Size()}
				if s.size == 0 {
					return nil
				}
				all_file_size += s.size
				all_files++
				sameFiles[s] = append(sameFiles[s], path)
			}
			all_dirs++
			return nil
		})
	}
	total_dupes := 0
	var size_wasted int64
	for k, v := range sameFiles {
		if k.size == 0 || len(v) < 2 {
			continue
		}
		fingerPrint := make(map[string][]string)
		for _, path := range v {
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
			fingerPrint[md5] = append(fingerPrint[md5], path)
		}
		for _, paths := range fingerPrint {
			toDelete := []string{}
			if 2 > len(paths) {
				continue
			}
			total_dupes--
			size_wasted = size_wasted - k.size
			for i, p := range paths {
				total_dupes++
				size_wasted = size_wasted + k.size
				fmt.Printf(" duplicate %d: %s\n", i, p)
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
		}
	}
	fmt.Printf("\nTotal dupes %d.  Total bytes wasted %d\n", total_dupes, size_wasted)
	fmt.Printf("\nTotal files %d.  Total bytes %d. Total dirs %d\n", all_files,
		all_file_size, all_dirs)
}
