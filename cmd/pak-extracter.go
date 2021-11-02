package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/new-world-tools/extracter/pak"
	"github.com/new-world-tools/go-oodle"
	workerpool "github.com/zelenin/go-worker-pool"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultThreads = 5
	maxThreads     = 10
)

var (
	pool      *workerpool.Pool
	filters   map[string]bool
	assetsDir string
	outputDir string
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
	threadsPtr := flag.Int("threads", defaultThreads, fmt.Sprintf("1-%d", maxThreads))
	flag.Parse()

	assetsDir, err = filepath.Abs(filepath.Clean(*assetsDirPtr))
	if err != nil {
		log.Fatalf("filepath.Abs: %s", err)
	}

	_, err = os.Stat(assetsDir)
	if os.IsNotExist(err) {
		log.Fatalf("'%s' does not exist", assetsDir)
	}

	outputDir, err = filepath.Abs(filepath.Clean(*outputDirPtr))
	if err != nil {
		log.Fatalf("filepath.Abs: %s", err)
	}

	filters = map[string]bool{}
	filterParts := strings.Split(*filterPtr, ",")
	for _, ext := range filterParts {
		filters[fmt.Sprintf(".%s", strings.TrimPrefix(ext, "."))] = true
	}

	threads := *threadsPtr
	if threads < 1 || threads > maxThreads {
		threads = defaultThreads
	}
	log.Printf("The number of threads is set to %d", threads)

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("MkdirAll: %s", err)
	}

	pakFiles, err := pak.FindAll(assetsDir)
	if err != nil {
		log.Fatalf("pak.FindAll: %s", err)
	}

	pool = workerpool.NewPool(5, 1000)

	go func() {
		errorChan := pool.Errors()

		for {
			err, ok := <-errorChan
			if !ok {
				break
			}

			taskId := err.(workerpool.TaskError).Id
			err = errors.Unwrap(err)
			log.Printf("task #%d err: %s", taskId, err)
		}
	}()

	var id int64
	for _, pakFile := range pakFiles {
		id++
		log.Printf("Working: %s", pakFile.GetPath())

		files, err := pakFile.GetFiles()
		if err != nil {
			log.Fatalf("pakFile.GetFiles: %s", err)
		}

		for _, file := range files {
			addTask(id, pakFile, file)
		}

		pakFile.Close()
	}

	pool.Wait()
}

func addTask(id int64, pakFile *pak.Pak, file *pak.File) {
	pool.AddTask(workerpool.NewTask(id, func(id int64) error {
		ext := filepath.Ext(file.Name)
		if filters[ext] {
			return nil
		}
		fpath := filepath.ToSlash(filepath.Clean(filepath.Join(outputDir, strings.ReplaceAll(filepath.Dir(pakFile.GetPath()), assetsDir, ""), file.Name)))
		err := os.MkdirAll(filepath.Dir(fpath), 0755)
		if err != nil {
			return err
		}

		dest, err := os.Create(fpath)
		if err != nil {
			return err
		}
		defer dest.Close()

		reader, err := file.Decompress()
		if err != nil {
			return err
		}
		defer reader.Close()

		_, err = io.Copy(dest, reader)
		if err != nil {
			return err
		}

		return nil
	}))
}
