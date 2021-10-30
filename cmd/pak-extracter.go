package main

import (
	"flag"
	"fmt"
	"github.com/new-world-tools/extracter/pak"
	"github.com/new-world-tools/go-oodle"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	_, err := os.Stat("oo2core_9_win64.dll")
	if os.IsNotExist(err) {
		err := oodle.Download()
		if err != nil {
			log.Fatalf("no oo2core_9_win64.dll")
		}
	}

	assetsDirPtr := flag.String("assets", "C:\\Program Files (x86)\\Steam\\steamapps\\common\\New World\\assets", "directory path")
	outputDirPtr := flag.String("output", "./extract", "directory path")
	filterPtr := flag.String("filter", "", "comma separated file extensions")
	flag.Parse()

	assetsDir, err := filepath.Abs(filepath.Clean(*assetsDirPtr))
	if err != nil {
		log.Fatalf("filepath.Abs: %s", err)
	}

	_, err = os.Stat(assetsDir)
	if os.IsNotExist(err) {
		log.Fatalf("'%s' does not exist", assetsDir)
	}

	outputDir, err := filepath.Abs(filepath.Clean(*outputDirPtr))
	if err != nil {
		log.Fatalf("filepath.Abs: %s", err)
	}

	filters := map[string]bool{}
	filterParts := strings.Split(*filterPtr, ",")
	for _, ext := range filterParts {
		filters[fmt.Sprintf(".%s", strings.TrimPrefix(ext, "."))] = true
	}

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("MkdirAll: %s", err)
	}

	pakFiles, err := pak.FindAll(assetsDir)
	if err != nil {
		log.Fatalf("pak.FindAll: %s", err)
	}

	for _, pakFile := range pakFiles {
		log.Printf("Working: %s", pakFile.GetPath())

		files, err := pakFile.GetFiles()
		if err != nil {
			log.Fatalf("pakFile.GetFiles: %s", err)
		}

		for _, file := range files {
			ext := filepath.Ext(file.Name)
			if filters[ext] {
				continue
			}
			fpath := filepath.ToSlash(filepath.Clean(filepath.Join(outputDir, strings.ReplaceAll(filepath.Dir(pakFile.GetPath()), assetsDir, ""), file.Name)))
			err = os.MkdirAll(filepath.Dir(fpath), 0755)
			if err != nil {
				log.Fatalf("os.MkdirAll: %s", err)
			}

			dest, err := os.Create(fpath)
			if err != nil {
				log.Fatalf("os.Create: %s", err)
			}

			reader, err := file.Decompress()
			if err != nil {
				log.Fatalf("file.Decompress: %s", err)
			}

			_, err = io.Copy(dest, reader)
			if err != nil {
				log.Printf("io.Copy: %s", err)
			}

			dest.Close()
			reader.Close()
		}

		pakFile.Close()
	}
}
