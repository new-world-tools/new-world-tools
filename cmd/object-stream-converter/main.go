package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/new-world-tools/new-world-tools/asset"
	"github.com/new-world-tools/new-world-tools/azcs"
	"github.com/new-world-tools/new-world-tools/profiler"
	"github.com/new-world-tools/new-world-tools/structure"
	workerpool "github.com/zelenin/go-worker-pool"
	"gopkg.in/yaml.v3"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

const (
	defaultThreads int64 = 3
	maxThreads     int64 = 50
)

const (
	formatYml  = "yml"
	formatJson = "json"
)

var formats = map[string]bool{
	//formatYml:  true,
	formatJson: true,
}

var (
	pool             *workerpool.Pool
	outputDir        string
	withIndents      bool
	indentsSize      int
	resolveHashValue bool
	debug            bool
	debugLog         string
	assetMap         map[string]*asset.AssetInfo
	pr               *profiler.Profiler
)

type DebugData struct {
	mu                sync.Mutex
	NotResolvedHashes map[string]bool
	NotResolvedTypes  map[string]bool
}

var debugData = &DebugData{
	NotResolvedHashes: map[string]bool{},
	NotResolvedTypes:  map[string]bool{},
}

func main() {
	//f, err := os.Create("cpu3.prof")
	//if err != nil {
	//	log.Fatal("could not create CPU profile: ", err)
	//}
	//defer f.Close()
	//
	//if err := pprof.StartCPUProfile(f); err != nil {
	//	log.Fatal("could not start CPU profile: ", err)
	//}
	//defer pprof.StopCPUProfile()

	pr = profiler.New()

	inputDirPtr := flag.String("input", ".\\extract", "directory path")
	outputDirPtr := flag.String("output", ".\\json", "directory path")
	threadsPtr := flag.Int64("threads", defaultThreads, fmt.Sprintf("1-%d", maxThreads))
	withIndentsPtr := flag.Bool("with-indents", false, "enable indents in json")
	indentsSizePtr := flag.Int("indents-size", 4, "indents size")
	resolveHashValuePtr := flag.Bool("resolve-hash-value", false, "")
	poolCapacityPtr := flag.Int64("pool-capacity", 1000, "pool capacity")
	debugPtr := flag.Bool("debug", false, "")
	debugLogPtr := flag.String("debug-log", ".\\debug.log", "debug log path")
	formatPtr := flag.String("format", "json", "yml or json")
	assetCatalogInputPtr := flag.String("asset-catalog", "", "assetcatalog.catalog path")
	flag.Parse()

	format := *formatPtr
	if formats[format] != true {
		log.Fatalf("Unsupported format: %s", format)
	}

	threads := *threadsPtr
	if threads < 1 || threads > maxThreads {
		threads = defaultThreads
	}
	log.Printf("The number of threads is set to %d", threads)

	withIndents = *withIndentsPtr
	indentsSize = *indentsSizePtr
	resolveHashValue = *resolveHashValuePtr
	poolCapacity := *poolCapacityPtr
	debug = *debugPtr
	debugLog = *debugLogPtr
	assetCatalogInput := *assetCatalogInputPtr

	inputDir, err := filepath.Abs(filepath.Clean(*inputDirPtr))
	if err != nil {
		log.Fatalf("filepath.Abs: %s", err)
	}

	_, err = os.Stat(inputDir)
	if os.IsNotExist(err) {
		log.Fatalf("'%s' does not exist", inputDir)
	}

	outputDir, err = filepath.Abs(filepath.Clean(*outputDirPtr))
	if err != nil {
		log.Fatalf("filepath.Abs: %s", err)
	}

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("MkdirAll: %s", err)
	}

	if assetCatalogInput != "" {
		assetCatalogInput, err = filepath.Abs(filepath.Clean(assetCatalogInput))
		if err != nil {
			log.Fatalf("filepath.Abs: %s", err)
		}

		_, err = os.Stat(assetCatalogInput)
		if os.IsNotExist(err) {
			log.Fatalf("'%s' does not exist", assetCatalogInput)
		}

		f, err := os.Open(assetCatalogInput)
		if err != nil {
			log.Fatalf("os.Open: %s", err)
		}

		log.Printf("Parsing the catalog...")
		cat, err := asset.ParseAssetCatalog(f)
		if err != nil {
			log.Fatalf("asset.ParseAssetCatalog: %s", err)
		}

		assetMap = make(map[string]*asset.AssetInfo, cat.AssetIdToInfoNumEntries)
		for _, ref := range cat.AssetIdToInfo {
			assetInfo, err := ref.Load(f, cat)
			if err != nil {
				log.Fatalf("ref.Load: %s", err)
			}
			assetMap[assetInfo.AssetId.String()] = assetInfo
		}

		f.Close()
	}

	pool = workerpool.NewPool(threads, poolCapacity)

	go func() {
		errorChan := pool.Errors()

		for {
			err, ok := <-errorChan
			if !ok {
				break
			}

			taskId := err.(workerpool.TaskError).Id
			err = errors.Unwrap(err)
			log.Fatalf("task #%d err: %s", taskId, err)
		}
	}()

	var id int64
	err = filepath.Walk(inputDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		path, err = filepath.Abs(filepath.Clean(path))
		if err != nil {
			return err
		}

		isAzcsFile, isCompressed, err := azcs.IsAzcsFile(path)
		if err != nil {
			return err
		}

		if !isAzcsFile {
			return nil
		}

		baseName := filepath.Base(path)
		if strings.Contains(baseName, ".dds") {
			return nil
		}

		id++

		relPath, err := filepath.Rel(inputDir, path)
		if err != nil {
			return err
		}

		output := filepath.Join(outputDir, relPath)
		if format == formatJson {
			output += ".json"
		}
		if format == formatYml {
			output += ".yml"
		}

		job := Job{
			Input:        path,
			Output:       output,
			RelPath:      relPath,
			IsCompressed: isCompressed,
			Format:       format,
		}

		addTask(id, job)

		return nil
	})
	if err != nil {
		log.Fatalf("filepath.Walk: %s", err)
	}

	pool.Close()
	pool.Wait()

	log.Printf("PeakMemory: %0.1fMb Duration: %s", float64(pr.GetPeakMemory())/1024/1024, pr.GetDuration().String())
	if debug {
		if len(debugData.NotResolvedHashes) > 0 {
			log.Printf("Not resolved hashes: %s", strings.Join(sortMap(debugData.NotResolvedHashes), ", "))
		}
		if len(debugData.NotResolvedTypes) > 0 {
			log.Printf("Not resolved types: %s", strings.Join(sortMap(debugData.NotResolvedTypes), ", "))
		}
	}
	if debugLog != "" {
		err = storeDebug(debugLog)
		if err != nil {
			log.Fatalf("storeDebug: %s", err)
		}
	}
}

func addTask(id int64, job Job) {
	pool.AddTask(workerpool.NewTask(id, func(id int64) error {
		log.Printf("Working: %s", job.RelPath)
		//defer log.Printf("Done: [#%06d] %s", id, job.Input)

		rc, err := azcs.GetReader(job.Input, job.IsCompressed)
		if err != nil {
			return err
		}
		defer rc.Close()

		stream, err := azcs.Parse(bufio.NewReader(rc))
		if err != nil {
			log.Fatalf("azcs.Parse: %s", err)
		}

		streamNode, err := azcs.ResolveStream(stream, resolveType, resolveHash, assetMap)
		if err != nil {
			log.Fatalf("azcs.ResolveStream: %s", err)
		}

		if resolveHashValue {
			hook(streamNode)
		}

		err = os.MkdirAll(filepath.Dir(job.Output), 0777)
		if err != nil {
			return err
		}

		var of *os.File

		if job.Format == formatYml {
			of, err = os.Create(job.Output)
			if err != nil {
				return err
			}

			enc := yaml.NewEncoder(of)

			enc.SetIndent(indentsSize)

			err = enc.Encode(streamNode)
			if err != nil {
				return err
			}
			enc.Close()
			of.Close()
		}

		if job.Format == formatJson {
			of, err = os.Create(job.Output)
			if err != nil {
				return err
			}

			enc := json.NewEncoder(of)

			if withIndents {
				enc.SetIndent("", strings.Repeat(" ", indentsSize))
			}

			err = enc.Encode(streamNode)
			if err != nil {
				return err
			}
			of.Close()
		}

		return nil
	}))
}

type Job struct {
	Input        string
	Output       string
	RelPath      string
	IsCompressed bool
	Format       string
}

func sortMap(data map[string]bool) []string {
	values := make([]string, len(data))
	i := 0
	for value, _ := range data {
		values[i] = value
		i++
	}

	sort.Strings(values)

	return values
}

func resolveHash(element *azcs.Element) string {
	hash := element.NameCrc

	if azcs.DefaultHashRegistry.Has(hash) {
		return azcs.DefaultHashRegistry.Get(hash)
	}

	formattedHash := fmt.Sprintf("0x%08x", hash)

	debugData.mu.Lock()
	debugData.NotResolvedHashes[formattedHash] = true
	debugData.mu.Unlock()

	return formattedHash
}

func resolveType(element *azcs.Element) string {
	typ := element.Type.String()
	if !element.SpecializedType.IsNil() {
		typ = element.SpecializedType.String()
	}

	if azcs.DefaultTypeRegistry.Has(typ) {
		return azcs.DefaultTypeRegistry.Get(typ)
	}

	debugData.mu.Lock()
	debugData.NotResolvedTypes[typ] = true
	debugData.mu.Unlock()

	return typ
}

func hook(data any) {
	node, ok := data.(*structure.OrderedMap[string, any])
	if ok {
		_, v, _ := node.GetByPosition(0)
		if v == "Crc32" && resolveHashValue {
			value, ok := node.Get("value")
			if ok {
				hash := value.(uint32)
				if hash != 0x0 && azcs.DefaultHashRegistry.Has(hash) {
					node.Add("__value", azcs.DefaultHashRegistry.Get(hash))
				}
			}
		} else {
			for node.Has() {
				_, value := node.Next()
				hook(value)
			}
		}

		return
	}

	list, ok := data.([]any)
	if ok {
		for _, value := range list {
			hook(value)
		}
		return
	}
}

func storeDebug(debugLog string) error {
	var oldDebugData = DebugData{
		NotResolvedHashes: map[string]bool{},
		NotResolvedTypes:  map[string]bool{},
	}

	f, err := os.Open(debugLog)
	if err == nil {
		json.NewDecoder(f).Decode(&oldDebugData)
		f.Close()
	}

	for k, _ := range oldDebugData.NotResolvedHashes {
		debugData.NotResolvedHashes[k] = true
	}
	for k, _ := range oldDebugData.NotResolvedTypes {
		debugData.NotResolvedTypes[k] = true
	}

	f, err = os.Create(debugLog)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(debugData)
	if err != nil {
		return err
	}

	return nil
}
