package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/new-world-tools/new-world-tools/azcs"
	"github.com/new-world-tools/new-world-tools/profiler"
	"github.com/new-world-tools/new-world-tools/reader"
	azcs2 "github.com/new-world-tools/new-world-tools/reader/azcs"
	"github.com/new-world-tools/new-world-tools/structure"
	workerpool "github.com/zelenin/go-worker-pool"
	"io"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	defaultThreads int64 = 3
	maxThreads     int64 = 10
)

var (
	pool             *workerpool.Pool
	outputDir        string
	withIndents      bool
	resolveHashValue bool
	debug            bool
	pr               *profiler.Profiler
)

type DebugData struct {
	mu                sync.Mutex
	notResolvedHashes map[string]bool
	notResolvedTypes  map[string]bool
}

var debugData = &DebugData{
	notResolvedHashes: map[string]bool{},
	notResolvedTypes:  map[string]bool{},
}

const typeField = "__type"
const valueField = "__value"

func main() {
	pr = profiler.New()

	inputDirPtr := flag.String("input", ".\\extract", "directory path")
	outputDirPtr := flag.String("output", ".\\json", "directory path")
	threadsPtr := flag.Int64("threads", defaultThreads, fmt.Sprintf("1-%d", maxThreads))
	withIndentsPtr := flag.Bool("with-indents", false, "enable indents in json")
	resolveHashValuePtr := flag.Bool("resolve-hash-value", false, "")
	poolCapacityPtr := flag.Int64("pool-capacity", 1000, "pool capacity")
	debugPtr := flag.Bool("debug", false, "")
	flag.Parse()

	threads := *threadsPtr
	if threads < 1 || threads > maxThreads {
		threads = defaultThreads
	}
	log.Printf("The number of threads is set to %d", threads)

	withIndents = *withIndentsPtr
	resolveHashValue = *resolveHashValuePtr
	poolCapacity := *poolCapacityPtr
	debug = *debugPtr

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

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		size := info.Size()

		if size < 5 {
			return nil
		}

		dataSize := 5

		br := bufio.NewReader(f)
		data, err := br.Peek(dataSize)
		if err != nil && err != io.EOF {
			return err
		}

		var isCompressedFile bool
		var isAzcs bool

		if isCompressed(data) {
			isCompressedFile = true
			isAzcs = true
		}

		if isUncompressed(data) {
			isCompressedFile = false
			isAzcs = true
		}

		if !isAzcs {
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

		job := Job{
			Input:        path,
			Output:       filepath.Join(outputDir, relPath) + ".json",
			RelPath:      relPath,
			IsCompressed: isCompressedFile,
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
		if len(debugData.notResolvedHashes) > 0 {
			log.Printf("Not resolved hashes: %s", strings.Join(sortMap(debugData.notResolvedHashes), ", "))
		}
		if len(debugData.notResolvedTypes) > 0 {
			log.Printf("Not resolved types: %s", strings.Join(sortMap(debugData.notResolvedTypes), ", "))
		}
	}
}

func addTask(id int64, job Job) {
	pool.AddTask(workerpool.NewTask(id, func(id int64) error {
		log.Printf("Working: %s", job.RelPath)
		//defer log.Printf("Done: [#%06d] %s", id, job.Input)

		f, err := os.Open(job.Input)
		if err != nil {
			return err
		}
		defer f.Close()

		var r io.Reader

		r = f

		if job.IsCompressed {
			r, err = azcs2.NewReader(r)
			if err != nil {
				return err
			}
		}

		stream, err := azcs.Parse(r)
		if err != nil {
			return err
		}

		if len(stream.Elements) > 1 {
			log.Fatalf("too much elements")
		}

		data := resolveNode(stream.Elements[0])

		err = os.MkdirAll(filepath.Dir(job.Output), 0777)
		if err != nil {
			return err
		}

		jf, err := os.Create(job.Output)
		if err != nil {
			return err
		}

		enc := json.NewEncoder(jf)

		if withIndents {
			enc.SetIndent("", "    ")
		}

		err = enc.Encode(data)
		if err != nil {
			return err
		}

		return nil
	}))
}

var uncompressedSignatures = [][]byte{
	{0x00, 0x00, 0x00, 0x00, 0x03},
	{0x00, 0x00, 0x00, 0x00, 0x02},
	{0x00, 0x00, 0x00, 0x00, 0x01},
}
var azcsSig = []byte{0x41, 0x5a, 0x43, 0x53}

func isUncompressed(data []byte) bool {
	for _, uncompressedSignature := range uncompressedSignatures {
		if len(data) >= len(uncompressedSignature) && bytes.Equal(uncompressedSignature, data[:len(uncompressedSignature)]) {
			return true
		}
	}

	return false
}

func isCompressed(data []byte) bool {
	if len(data) < len(azcsSig) {
		return false
	}

	return bytes.Equal(azcsSig, data[:len(azcsSig)])
}

type Job struct {
	Input        string
	Output       string
	RelPath      string
	IsCompressed bool
}

type float32s interface {
	float32 | JsonFloat32
}

type float64s interface {
	float64 | JsonFloat64
}

type floats interface {
	float32s | float64s
}

func createFloatArray[T floats](data []byte) ([]T, error) {
	var z T
	dataTypeSize := int(reflect.Indirect(reflect.ValueOf(z)).Type().Size())

	fs := make([]T, len(data)/dataTypeSize)
	for i := 0; i < len(fs); i++ {
		var f T
		buf := bytes.NewReader(data[i*4 : (i+1)*4])
		err := binary.Read(buf, binary.BigEndian, &f)
		if err != nil {
			return nil, err
		}
		fs[i] = f
	}

	return fs, nil
}

func maxFloat[T floats](data []T) T {
	var max T
	for i, f := range data {
		if i == 0 {
			max = f
		} else {
			if f > max {
				max = f
			}
		}
	}

	return max
}

//func colLen(floats []JsonFloat32) JsonFloat32 {
//	lengthSq := floats[0]*floats[0] + floats[1]*floats[1] + floats[2]*floats[2]
//	return JsonFloat32(math.Sqrt(float64(lengthSq)))
//}

func resolveNode(element *azcs.Element) any {
	node := structure.NewOrderedMap[string, any]()

	node.Add(typeField, resolveType(element))

	switch element.ResolveType().String() {
	case
		// Transform
		"5d9958e9-9f1e-4985-b532-fffde75fedfd":
		switch element.Version {
		case 0:
			if len(element.Data) != 48 {
				log.Fatalf("wrong size: %d", len(element.Data))
			}
			col0, err := createFloatArray[JsonFloat32](element.Data[0:12])
			if err != nil {
				log.Fatalf("createFloatArrray: %s", err)
			}
			col1, err := createFloatArray[JsonFloat32](element.Data[12:24])
			if err != nil {
				log.Fatalf("createFloatArrray: %s", err)
			}
			col2, err := createFloatArray[JsonFloat32](element.Data[24:36])
			if err != nil {
				log.Fatalf("createFloatArrray: %s", err)
			}
			col3, err := createFloatArray[JsonFloat32](element.Data[36:48])
			if err != nil {
				log.Fatalf("createFloatArrray: %s", err)
			}

			tnode := structure.NewOrderedMap[string, any]()

			tnode.Add("rotation/scale", []JsonFloat32{col0[0], col0[1], col0[2], col1[0], col1[1], col1[2], col2[0], col2[1], col2[2]})
			tnode.Add("translation", col3)

			node.Add(valueField, tnode)

		case 1:
			if len(element.Data) != 40 {
				log.Fatalf("wrong size: %d", len(element.Data))
			}
			rotation, err := createFloatArray[JsonFloat32](element.Data[0:16])
			if err != nil {
				log.Fatalf("createFloatArrray: %s", err)
			}
			vectorScale, err := createFloatArray[JsonFloat32](element.Data[16:28])
			if err != nil {
				log.Fatalf("createFloatArrray: %s", err)
			}
			translation, err := createFloatArray[JsonFloat32](element.Data[28:40])
			if err != nil {
				log.Fatalf("createFloatArrray: %s", err)
			}
			scale := maxFloat[JsonFloat32](vectorScale)

			tnode := structure.NewOrderedMap[string, any]()

			tnode.Add("rotation", rotation)
			tnode.Add("scale", scale)
			tnode.Add("translation", translation)

			node.Add(valueField, tnode)

		default:
			log.Fatalf("unsupported version: %d", element.Version)
		}

		return node

	case
		// Color
		"7894072a-9050-4f0f-901b-34b1a0d29417":
		f32s, err := createFloatArray[JsonFloat32](element.Data)
		if err != nil {
			log.Fatalf("createFloatArrray: %s", err)
		}

		node.Add(valueField, f32s)

		return node

	// Asset
	case "77a19d40-8731-4d3c-9041-1b43047366a4":
		buf := bytes.NewBuffer(element.Data)

		data, err := reader.ReadBytes(buf, 16)
		if err != nil {
			log.Fatalf("reader.ReadBytes: %s", err)
		}
		id, _ := uuid.FromBytes(data)
		node.Add("guid", id.String())

		data, err = reader.ReadBytes(buf, 16)
		if err != nil {
			log.Fatalf("reader.ReadBytes: %s", err)
		}
		id, _ = uuid.FromBytes(data)
		node.Add("subId", id.String())

		data, err = reader.ReadBytes(buf, 16)
		if err != nil {
			log.Fatalf("reader.ReadBytes: %s", err)
		}
		id, _ = uuid.FromBytes(data)
		node.Add("type", id.String())

		u64, err := reader.ReadUint64(buf, binary.BigEndian)
		if err != nil {
			log.Fatalf("reader.ReadUint64: %s", err)
		}
		if u64 > 0 {
			data, err = reader.ReadBytes(buf, int(u64))
			if err != nil {
				log.Fatalf("reader.ReadBytes: %s", err)
			}
			node.Add("hint", string(data))
		}

		return node
	}

	_, v, _ := node.GetByPosition(0)

	switch v {
	case "bool":
		var b bool
		l := len(element.Data)
		if l != 1 {
			log.Fatalf("unsupported bool size: %d", l)
		}

		switch element.Data[0] {
		case 0x00:
			b = false

		case 0x01:
			b = true

		default:
			b = true
		}

		return b

	case "AZStd::string":
		return string(element.Data)

	case "AZ::Uuid":
		id, _ := uuid.FromBytes(element.Data)
		return id.String()

	case
		"unsigned char",
		"unsigned int",
		"unsigned short",
		"AZ::u64":
		l := len(element.Data)

		switch l {
		case 1:
			return element.Data[0]

		case 2:
			return binary.BigEndian.Uint16(element.Data)

		case 4:
			return binary.BigEndian.Uint32(element.Data)

		case 8:
			return binary.BigEndian.Uint64(element.Data)

		default:
			log.Fatalf("unsupported data size: %d", l)
		}

	case
		"int",
		"AZ::s64":
		l := len(element.Data)

		switch l {
		case 1:
			var i8 int8
			buf := bytes.NewReader(element.Data)
			err := binary.Read(buf, binary.BigEndian, &i8)
			if err != nil {
				log.Fatalf("binary.Read: %s", err)
			}
			return i8

		case 2:
			var i16 int16
			buf := bytes.NewReader(element.Data)
			err := binary.Read(buf, binary.BigEndian, &i16)
			if err != nil {
				log.Fatalf("binary.Read: %s", err)
			}
			return i16

		case 4:
			var i32 int32
			buf := bytes.NewReader(element.Data)
			err := binary.Read(buf, binary.BigEndian, &i32)
			if err != nil {
				log.Fatalf("binary.Read: %s", err)
			}
			return i32

		case 8:
			var i64 int64
			buf := bytes.NewReader(element.Data)
			err := binary.Read(buf, binary.BigEndian, &i64)
			if err != nil {
				log.Fatalf("binary.Read: %s", err)
			}
			return i64

		default:
			log.Fatalf("unsupported data size: %d", l)
		}

	case
		"float",
		"double":
		l := len(element.Data)

		switch l {
		case 4:
			var f32 JsonFloat32
			buf := bytes.NewReader(element.Data)
			err := binary.Read(buf, binary.BigEndian, &f32)
			if err != nil {
				log.Fatalf("binary.Read: %s", err)
			}
			return f32

		case 8:
			var f64 JsonFloat64
			buf := bytes.NewReader(element.Data)
			err := binary.Read(buf, binary.BigEndian, &f64)
			if err != nil {
				log.Fatalf("binary.Read: %s", err)
			}
			return f64

		default:
			log.Fatalf("unsupported data size: %d", l)
		}

	case
		"AZStd::vector",
		"AZStd::unordered_set":
		nodes := make([]any, len(element.Elements))

		for i, element := range element.Elements {
			nodes[i] = resolveNode(element)
		}

		return nodes

	case
		"Quaternion",
		"Vector2",
		"Vector3":
		l := len(element.Data)

		f32s := make([]JsonFloat32, l/4)
		for i := 0; i < len(f32s); i++ {
			var f32 JsonFloat32
			buf := bytes.NewReader(element.Data[i*4 : (i+1)*4])
			err := binary.Read(buf, binary.BigEndian, &f32)
			if err != nil {
				log.Fatalf("binary.Read: %s", err)
			}
			f32s[i] = f32
		}

		return f32s

	case
		"Amazon::Pervasives::UID",
		"Amazon::Hub::ActorRef",
		"BitSet":
		return element.Data

	case
		"AZStd::array",
		"AZStd::fixed_vector",
		"AZStd::list",
		"AZStd::map",
		"AZStd::unordered_map":
		values := make([]any, len(element.Elements))
		for i, element := range element.Elements {
			key := resolveHash(element)
			value := resolveNode(element)

			if key != "element" {
				log.Fatalf("wrong key: %s", key)
			}

			values[i] = value
		}
		return values

	case
		"AZStd::intrusive_ptr",
		"AZStd::shared_ptr",
		"AZStd::unique_ptr":
		if len(element.Elements) == 0 {
			return nil
		}
		if len(element.Elements) != 1 {
			log.Fatalf("wrong elements count: %d", len(element.Elements))
		}
		for _, element := range element.Elements {
			key := resolveHash(element)
			value := resolveNode(element)

			if key != "element" {
				log.Fatalf("wrong key: %s", key)
			}

			return value
		}

	default:
		if len(element.Data) > 0 {
			_, v, _ := node.GetByPosition(0)

			vs := v.(string)

			checkId, err := uuid.FromString(vs)
			if err != nil || checkId.IsNil() {
				if len(element.Data) > 24 {
					log.Fatalf("unsupported data type: %s, type: %s, stype: %s", vs, element.Type.String(), element.SpecializedType.String())
				}
				log.Fatalf("unsupported data type: %s, type: %s, stype: %s data: %x", vs, element.Type.String(), element.SpecializedType.String(), element.Data)
			} else {
				node.Add(valueField, element.Data)
			}
		} else {
			for _, element := range element.Elements {
				key := resolveHash(element)
				value := resolveNode(element)

				node.Add(key, value)
			}
			if v == "Crc32" && resolveHashValue {
				value, ok := node.Get("value")
				if ok {
					hash := value.(uint32)
					if azcs.DefaultHashRegistry.Has(hash) {
						node.Add(valueField, azcs.DefaultHashRegistry.Get(hash))
					}
				}
			}
		}
	}

	return node
}

func resolveHash(element *azcs.Element) string {
	hash := element.NameCrc

	if azcs.DefaultHashRegistry.Has(hash) {
		return azcs.DefaultHashRegistry.Get(hash)
	}

	formattedHash := fmt.Sprintf("0x%08x", hash)

	debugData.mu.Lock()
	debugData.notResolvedHashes[formattedHash] = true
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
	debugData.notResolvedTypes[typ] = true
	debugData.mu.Unlock()

	return typ
}

type JsonFloat64 float64

func (v JsonFloat64) MarshalJSON() ([]byte, error) {
	f64 := float64(v)
	var s string
	switch {
	case math.IsInf(f64, 1):
		s = "+Inf"
	case math.IsInf(f64, -1):
		s = "-Inf"
	case math.IsNaN(f64):
		s = "NaN"
	default:
		s = strconv.FormatFloat(f64, 'f', -1, 64)
		return []byte(s), nil
	}
	return []byte(`"` + s + `"`), nil
}

func (v *JsonFloat64) UnmarshalJSON(b []byte) error {
	switch {
	case bytes.Equal(b, []byte(`"+Inf"`)):
		*v = JsonFloat64(math.Inf(1))
	case bytes.Equal(b, []byte(`"-Inf"`)):
		*v = JsonFloat64(math.Inf(-1))
	case bytes.Equal(b, []byte(`"NaN"`)):
		*v = JsonFloat64(math.NaN())
	default:
		n, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			return err
		}
		*v = JsonFloat64(n)
	}
	return nil
}

type JsonFloat32 float32

func (v JsonFloat32) MarshalJSON() ([]byte, error) {
	f64 := float64(v)
	var s string
	switch {
	case math.IsInf(f64, 1):
		s = "+Inf"
	case math.IsInf(f64, -1):
		s = "-Inf"
	case math.IsNaN(f64):
		s = "NaN"
	default:
		s = strconv.FormatFloat(f64, 'f', -1, 32)
		return []byte(s), nil
	}
	return []byte(`"` + s + `"`), nil
}

func (v *JsonFloat32) UnmarshalJSON(b []byte) error {
	switch {
	case bytes.Equal(b, []byte(`"+Inf"`)):
		*v = JsonFloat32(math.Inf(1))
	case bytes.Equal(b, []byte(`"-Inf"`)):
		*v = JsonFloat32(math.Inf(-1))
	case bytes.Equal(b, []byte(`"NaN"`)):
		*v = JsonFloat32(math.NaN())
	default:
		n, err := strconv.ParseFloat(string(b), 32)
		if err != nil {
			return err
		}
		*v = JsonFloat32(n)
	}
	return nil
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
