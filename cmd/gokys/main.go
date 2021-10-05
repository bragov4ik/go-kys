package main

import (
	"encoding/xml"
	"flag"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"go/parser"
	"go/token"

	"github.com/bragov4ik/go-kys/pkg/kys"
	"github.com/k0kubun/pp"
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

		if err != nil {
			log.Fatal(err)
		}

		if info.IsDir() {
			files, err = WalkMatch(v, "*.go")
			if err != nil {
				log.Fatal(err)
			}
		} else {
			files = append(files, v)
		}
	}

	xmlConfig, err := os.Open("../../config.xml")
	if err != nil {
		log.Fatal(err)
	}
	// defer the closing of our xmlFile so that we can parse it later on
	defer func(xmlConfig *os.File) {
		err := xmlConfig.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(xmlConfig)

	byteValue, err := ioutil.ReadAll(xmlConfig)
	if err != nil {
		log.Fatal(err)
	}

	var config kys.Config

	// get config from config.xml file
	err = xml.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	scores := kys.Info{}

	for _, file := range files {
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}
		kys.GetInfo(node, &scores, &config)
	}

	pp.Println(scores)

}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
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
