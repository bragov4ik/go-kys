package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	args := flag.Args()

	var files []string

	for _, v := range args {
		info, err := os.Stat(v)

		if os.IsNotExist(err) {
			log.Fatal("File does not exist.")
		}

		if info.IsDir() {
			files, err = WalkMatch(v, "*.go")
		} else {
			files = append(files, v)
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(files)
}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
