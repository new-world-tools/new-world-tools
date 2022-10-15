package main

import (
	"bufio"
	"crypto/sha1"
	"errors"
	"flag"
	"fmt"
	"github.com/new-world-tools/go-oodle"
	"github.com/new-world-tools/new-world-tools/hash"
	"github.com/new-world-tools/new-world-tools/pak"
	"github.com/new-world-tools/new-world-tools/profiler"
	"github.com/new-world-tools/new-world-tools/reader/azcs"
	workerpool "github.com/zelenin/go-worker-pool"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	defaultThreads int64 = 3
	maxThreads     int64 = 10
	azcsSignature        = "AZCS"
)

var (
	pool           *workerpool.Pool
	filters        map[string]bool
	assetsDir      string
	outputDir      string
	hashSumFile    string
	decompressAzcs bool
	hashRegistry   *hash.Registry
	pr             *profiler.Profiler
)

func main() {
	pr = profiler.New()

	var err error

	if !oodle.IsDllExist() {
		err := oodle.Download()
		if err != nil {
			log.Fatalf("no oo2core_9_win64.dll")
		}
	}

	assetsDirPtr := flag.String("assets", "C:\\Program Files (x86)\\Steam\\steamapps\\common\\New World\\assets", "directory path")
	outputDirPtr := flag.String("output", "./extract", "directory path")
	filterPtr := flag.String("filter", "", "comma separated file extensions")
	threadsPtr := flag.Int64("threads", defaultThreads, fmt.Sprintf("1-%d", maxThreads))
	hashSumFilePtr := flag.String("hash", "", "hash sum path")
	decompressAzcsPtr := flag.Bool("decompress-azcs", false, "hash sum path")
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

	hashSumFile = *hashSumFilePtr
	if hashSumFile != "" {
		hashSumFile, err = filepath.Abs(filepath.Clean(hashSumFile))
		if err != nil {
			log.Fatalf("filepath.Abs: %s", err)
		}
		hashRegistry = hash.NewRegistry()
	}

	decompressAzcs = *decompressAzcsPtr

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("MkdirAll: %s", err)
	}

	pakFiles, err := pak.FindAll(assetsDir)
	if err != nil {
		log.Fatalf("pak.FindAll: %s", err)
	}

	pool = workerpool.NewPool(threads, 1000)

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

	if hashSumFile != "" {
		log.Printf("Writing %s", hashSumFile)

		hashes := hashRegistry.Hashes()
		sort.Slice(hashes, func(i, j int) bool {
			return hashes[i].FileName < hashes[j].FileName
		})

		hashSumsFile, err := os.Create(hashSumFile)
		if err != nil {
			log.Fatalf("os.Create: %s", err)
		}
		defer hashSumsFile.Close()

		for _, fileHash := range hashes {
			_, err = hashSumsFile.WriteString(fmt.Sprintf("%x *%s\n", fileHash.Hash, fileHash.FileName))
			if err != nil {
				log.Fatalf("hashSumsFile.WriteString: %s", err)
			}
		}
	}

	log.Printf("PeakMemory: %0.1fMb Duration: %s", float64(pr.GetPeakMemory())/1024/1024, pr.GetDuration().String())
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

		decompressReader, err := file.Decompress()
		if err != nil {
			return err
		}
		defer decompressReader.Close()

		var r io.Reader

		if decompressAzcs {
			bufReader := bufio.NewReaderSize(decompressReader, 16)

			data, err := bufReader.Peek(4)
			if err != nil {
				return err
			}

			if string(data) == azcsSignature {
				r, err = azcs.NewReader(bufReader)
				if err != nil {
					return err
				}
			} else {
				r = bufReader
			}
		} else {
			r = decompressReader
		}

		if hashSumFile == "" {
			_, err = io.Copy(dest, r)
			if err != nil {
				return err
			}
		} else {
			hasher := sha1.New()
			reader := io.TeeReader(r, hasher)

			_, err = io.Copy(dest, reader)
			if err != nil {
				return err
			}

			hashRegistry.Add(file.Name, hasher.Sum(nil))
		}

		return nil
	}))
}
