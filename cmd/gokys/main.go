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

		die(err)

		if info.IsDir() {
			files, err = WalkMatch(v, "*.go")
			die(err)
		} else {
			files = append(files, v)
		}
	}

	xmlConfig, err := os.Open("../../config.xml")
	die(err)
	// defer the closing of our xmlFile so that we can parse it later on
	defer func(xmlConfig *os.File) {
		err := xmlConfig.Close()
		die(err)
	}(xmlConfig)

	byteValue, err := ioutil.ReadAll(xmlConfig)
	die(err)

	var config kys.Config

	// get config from config.xml file
	err = xml.Unmarshal(byteValue, &config)
	die(err)

	fset := token.NewFileSet()
	scores := kys.Info{}

	for _, file := range files {
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		die(err)
		kys.GetInfo(node, &scores, &config)
	}

	pp.Println(scores)

}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		die(err)
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
	die(err)
	return matches, nil
}

// error handling
func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
