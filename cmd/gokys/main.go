package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"go/parser"
	"go/token"

	"github.com/bragov4ik/go-kys/pkg/wmfp"
)

var cfgpath = flag.String("c", "config.xml", "XML config")

func readCfg() wmfp.Config {
	file, err := os.Open(*cfgpath)
	die(err)
	defer func(cfg *os.File) { die(cfg.Close()) }(file)
	bytes, err := ioutil.ReadAll(file)
	die(err)

	var cfg wmfp.Config
	die(xml.Unmarshal(bytes, &cfg))
	return cfg
}

func getFiles(args []string) []string {
	var files []string

	for _, v := range args {
		info, err := os.Stat(v)
		die(err)

		if info.IsDir() {
			files, err = walkMatch(v, "*.go")
			die(err)
		} else {
			files = append(files, v)
		}
	}
	return files
}

func main() {
	flag.Parse()
	cfg := readCfg()
	files := getFiles(flag.Args())

	measurers := make(chan *wmfp.MeasurerWMFP, len(files))
	for _, file := range files {
		go measureFile(file, cfg, measurers)
	}

	fmt.Println(combineMeasures(measurers, len(files)))
}

func measureFile(file string, config wmfp.Config, results chan *wmfp.MeasurerWMFP) {
	measurer := wmfp.NewMeasurerWMFP(&config)
	node, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.ParseComments)
	die(err)
	measurer.ParseFile(node)
	results <- &measurer
}

func combineMeasures(results chan *wmfp.MeasurerWMFP, mNum int) float64 {
	var nextMeasurer *wmfp.MeasurerWMFP
	var totalScore float64
	for i := 0; i < mNum; i++ {
		nextMeasurer = <-results
		totalScore += nextMeasurer.Finish()
	}
	return totalScore
}

func walkMatch(root, pattern string) ([]string, error) {
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

// error handling
func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
