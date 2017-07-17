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
	for _, dir := range flag.Args() {
		filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() {
				s := sameFile{f.Name(), f.Size()}
				if s.size == 0 {
					return nil
				}
				sameFiles[s] = append(sameFiles[s], path)
			}
			return nil
		})
	}
	for k, v := range sameFiles {
		if k.size == 0 || len(v) < 2 {
			continue
		}
		fingerPrint := make(map[string]int)
		toDelete := []string{}
		fmt.Printf("%s %d ==> %d\n", k.name, k.size, len(v))
		for _, path := range v {
			fmt.Printf("\t%s\n", path)
			f, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			h := md5.New()

			if _, err := io.Copy(h, f); err != nil {
				log.Fatal(err)
			}
			md5 := string(h.Sum(nil)[:])
			fmt.Printf("%x\n", md5)
			fingerPrint[md5]++
      fmt.Printf("has prefix %s %v\n", path, strings.HasPrefix(path, *deletePrefix))
			if *deletePrefix != "" && strings.HasPrefix(path, *deletePrefix) {
				toDelete = append(toDelete, path)
			}
		}
		if 1 < len(fingerPrint) {
			for _, filePath := range toDelete {
				fmt.Printf("going to delete %v\n", filePath)
			}
		}
	}
}
