package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"flag"
	"fmt"
	"github.com/new-world-tools/new-world-tools/hash"
	"github.com/new-world-tools/new-world-tools/pak"
	"github.com/new-world-tools/new-world-tools/profiler"
	"github.com/new-world-tools/new-world-tools/reader"
	"github.com/new-world-tools/new-world-tools/reader/azcs"
	workerpool "github.com/zelenin/go-worker-pool"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var rePak = regexp.MustCompile(`.pak$`)

const (
	defaultThreads int64 = 3
	maxThreads     int64 = 10
)

var (
	pool           *workerpool.Pool
	inputPath      string
	outputDir      string
	hashSumFile    string
	decompressAzcs bool
	fixLuac        bool
	hashRegistry   *hash.Registry
	pr             *profiler.Profiler
)

var (
	excludeRe       *regexp.Regexp
	includeRe       *regexp.Regexp
	includePriority bool
)

func main() {
	pr = profiler.New()

	var err error

	inputPathPtr := flag.String("input", "", "directory or .pak path")
	outputDirPtr := flag.String("output", "./extract", "directory path")
	threadsPtr := flag.Int64("threads", defaultThreads, fmt.Sprintf("1-%d", maxThreads))
	hashSumFilePtr := flag.String("hash", "", "hash sum path")
	decompressAzcsPtr := flag.Bool("decompress-azcs", false, "decompress AZCS (Amazon Object Stream)")
	fixLuacPtr := flag.Bool("fix-luac", false, "fix .luac header for unluac")
	excludePtr := flag.String("exclude", "", "regexp")
	includePtr := flag.String("include", "", "regexp")
	includePriorityPtr := flag.Bool("include-priority", false, "include flag priority")
	flag.Parse()

	if *excludePtr != "" {
		re, err := regexp.Compile(*excludePtr)
		if err != nil {
			log.Fatalf("regexp.Compile: %s", err)
		}
		excludeRe = re
	}
	if *includePtr != "" {
		re, err := regexp.Compile(*includePtr)
		if err != nil {
			log.Fatalf("regexp.Compile: %s", err)
		}
		includeRe = re
	}
	includePriority = *includePriorityPtr

	inputPath = *inputPathPtr

	inputPath, err = filepath.Abs(filepath.Clean(inputPath))
	if err != nil {
		log.Fatalf("filepath.Abs: %s", err)
	}

	fi, err := os.Stat(inputPath)
	if os.IsNotExist(err) {
		log.Fatalf("'%s' does not exist", inputPath)
	}

	isDir := fi.IsDir()

	outputDir, err = filepath.Abs(filepath.Clean(*outputDirPtr))
	if err != nil {
		log.Fatalf("filepath.Abs: %s", err)
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
	fixLuac = *fixLuacPtr

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("MkdirAll: %s", err)
	}

	var pakFiles []*pak.Pak
	if isDir {
		pakFiles, err = pak.FindAll(inputPath)
		if err != nil {
			log.Fatalf("pak.FindAll: %s", err)
		}
	} else {
		if !rePak.MatchString(inputPath) {
			log.Fatalf("Not valid .pak: %s", inputPath)
		}
		pakFiles = []*pak.Pak{pak.NewPak(inputPath)}
	}

	pool = workerpool.NewPool(threads, 1000)

	go func() {
		errorChan := pool.Errors()

		for {
			err, ok := <-errorChan
			if !ok {
				break
			}

			log.Printf("err: %s", err)
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
			if match(file.Name) {
				addTask(id, pakFile, file)
			}
		}

		pakFile.Close()
	}

	pool.Stop()
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

var luacSig = []byte{0x04, 0x00, 0x1b, 0x4c, 0x75, 0x61}

type TaskError struct {
	Pak  string
	Path string
	Err  error
}

func (err *TaskError) Error() string {
	return fmt.Sprintf("[%s:/%s] %s", err.Pak, err.Path, err.Err)
}

func newTaskError(pak string, path string, err error) *TaskError {
	return &TaskError{
		Pak:  pak,
		Path: path,
		Err:  err,
	}
}

func addTask(id int64, pakFile *pak.Pak, file *pak.File) {
	pool.AddTask(func(ctx context.Context) error {
		var err error

		basePath := inputPath
		if rePak.MatchString(inputPath) {
			basePath = filepath.Dir(inputPath)
		}

		fpath := filepath.ToSlash(filepath.Clean(filepath.Join(outputDir, strings.ReplaceAll(filepath.Dir(pakFile.GetPath()), basePath, ""), file.Name)))
		err = os.MkdirAll(filepath.Dir(fpath), 0755)
		if err != nil {
			return newTaskError(pakFile.GetPath(), file.Name, err)
		}

		dest, err := os.Create(fpath)
		if err != nil {
			return newTaskError(pakFile.GetPath(), file.Name, err)
		}

		decompressReader, err := file.Decompress()
		if err != nil {
			return newTaskError(pakFile.GetPath(), file.Name, err)
		}
		defer decompressReader.Close()

		var r io.Reader

		bufReader := bufio.NewReaderSize(decompressReader, 1024*1024)

		sigData, err := bufReader.Peek(8)
		if err != nil && err != io.EOF {
			return newTaskError(pakFile.GetPath(), file.Name, err)
		}

		r = bufReader

		if decompressAzcs {
			if bytes.Equal(azcs.Signature, sigData[:len(azcs.Signature)]) {
				r, err = azcs.NewReader(r)
				if err != nil {
					return newTaskError(pakFile.GetPath(), file.Name, err)
				}
			}
		}

		if fixLuac {
			if bytes.Equal(luacSig, sigData[:len(luacSig)]) {
				err = reader.SkipBytes(r, 2)
				if err != nil {
					return newTaskError(pakFile.GetPath(), file.Name, err)
				}
			}
		}

		if hashSumFile == "" {
			_, err = io.Copy(dest, r)
			if err != nil {
				return newTaskError(pakFile.GetPath(), file.Name, err)
			}
		} else {
			hasher := sha1.New()
			reader := io.TeeReader(r, hasher)

			_, err = io.Copy(dest, reader)
			if err != nil {
				return newTaskError(pakFile.GetPath(), file.Name, err)
			}

			hashRegistry.Add(file.Name, hasher.Sum(nil))
		}

		dest.Close()

		err = os.Chtimes(fpath, time.Now(), file.GetModifiedTime())
		if err != nil {
			return err
		}

		return nil
	})
}

func match(fileName string) bool {
	if excludeRe == nil && includeRe != nil {
		return includeRe.MatchString(fileName)
	}

	if excludeRe != nil && includeRe == nil {
		return !excludeRe.MatchString(fileName)
	}

	if excludeRe == nil || includeRe == nil {
		return true
	}

	if excludeRe != nil && includeRe != nil {
		if includePriority {
			if includeRe.MatchString(fileName) {
				return true
			} else {
				return !excludeRe.MatchString(fileName)
			}
		} else {
			if excludeRe.MatchString(fileName) {
				return false
			} else {
				return includeRe.MatchString(fileName)
			}
		}
	}

	if includePriority {
		return includeRe.MatchString(fileName) || !excludeRe.MatchString(fileName)
	} else {
		return !excludeRe.MatchString(fileName) && includeRe.MatchString(fileName)
	}
}
