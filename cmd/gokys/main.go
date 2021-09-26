package main

import (
	"flag"
	"fmt"
	_ "github.com/bragov4ik/go-kys/pkg/kys"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var fileNames string
	var directory string

	flag.StringVar(&fileNames, "f", "", "Files to count metrics for")
	flag.StringVar(&directory, "d", "", "Name of the directory for recursive scanning")
	flag.Parse()

	var files []string
	var err error

	if directory == "" {
		files = strings.Split(fileNames, " ")
	} else {
		files, err = WalkMatch(directory, "*.go")
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
