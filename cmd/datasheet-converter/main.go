package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/expr-lang/expr"
	"github.com/goccy/go-yaml"
	expr2 "github.com/new-world-tools/new-world-tools/cmd/datasheet-converter/expr"
	"github.com/new-world-tools/new-world-tools/datasheet"
	"github.com/new-world-tools/new-world-tools/localization"
	"github.com/new-world-tools/new-world-tools/profiler"
	"github.com/new-world-tools/new-world-tools/store"
	"github.com/new-world-tools/new-world-tools/structure"
	workerpool "github.com/zelenin/go-worker-pool"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	defaultThreads int64 = 3
	maxThreads     int64 = 10
)

var (
	pool             *workerpool.Pool
	localizationData *store.Store[string, string]
	resolveExpr      bool
	resolveContext   *ResolveContext
	inputDir         string
	outputDir        string
	format           string
	withIndents      bool
	keepStructure    bool
	pr               *profiler.Profiler
)

const (
	formatCsv  = "csv"
	formatJson = "json"
	formatYaml = "yaml"
)

var formats = map[string]bool{
	formatCsv:  true,
	formatJson: true,
	formatYaml: true,
}

type ResolveContext struct {
	CharacterLevel      int64
	GearScore           int64
	StatusEffectPotency map[string]float64
	ConsumablePotency   map[string]float64
	PerkScaling         map[string]*expr2.Scaling
}

func main() {
	var err error
	pr = profiler.New()

	inputDirPtr := flag.String("input", ".\\extract\\sharedassets\\springboardentitites\\datatables", "directory path")
	localizationDirPtr := flag.String("localization", "", "localization path")
	resolveExprPtr := flag.Bool("resolve-expr", false, "resolve expressions in strings")
	outputDirPtr := flag.String("output", ".\\datasheets", "directory path")
	formatPtr := flag.String("format", "csv", "csv, json, yaml")
	threadsPtr := flag.Int64("threads", defaultThreads, fmt.Sprintf("1-%d", maxThreads))
	withIndentsPtr := flag.Bool("with-indents", false, "enable indents in json")
	keepStructurePtr := flag.Bool("keep-structure", false, "keep original file structure")
	flag.Parse()

	format = *formatPtr
	localizationDir := *localizationDirPtr
	resolveExpr = *resolveExprPtr
	withIndents = *withIndentsPtr
	keepStructure = *keepStructurePtr

	if formats[format] != true {
		log.Fatalf("Unsupported format: %s", format)
	}

	threads := *threadsPtr
	if threads < 1 || threads > maxThreads {
		threads = defaultThreads
	}
	log.Printf("The number of threads is set to %d", threads)

	inputDir, err = filepath.Abs(filepath.Clean(*inputDirPtr))
	if err != nil {
		log.Fatalf("filepath.Abs: %s", err)
	}

	_, err = os.Stat(inputDir)
	if os.IsNotExist(err) {
		log.Fatalf("'%s' does not exist", inputDir)
	}

	if localizationDir != "" {
		localizationDir, err = filepath.Abs(filepath.Clean(localizationDir))
		if err != nil {
			log.Fatalf("filepath.Abs: %s", err)
		}

		_, err = os.Stat(localizationDir)
		if os.IsNotExist(err) {
			log.Fatalf("'%s' does not exist", localizationDir)
		}

		localizationData, err = localization.New(localizationDir)
		if err != nil {
			log.Fatalf("localization.New: %s", err)
		}
	}

	outputDir, err = filepath.Abs(filepath.Clean(*outputDirPtr))
	if err != nil {
		log.Fatalf("filepath.Abs: %s", err)
	}

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("MkdirAll: %s", err)
	}

	dsStore, err := datasheet.NewStore(inputDir)
	if err != nil {
		log.Fatalf("datasheet.NewStore: %s", err)
	}

	resolveContext = &ResolveContext{
		CharacterLevel:      65,
		GearScore:           700,
		StatusEffectPotency: make(map[string]float64),
		ConsumablePotency:   make(map[string]float64),
		PerkScaling:         make(map[string]*expr2.Scaling),
	}

	if resolveExpr {
		statusEffectPotency, consumablePotency, err := expr2.GetConsumablePotencies(dsStore)
		if err != nil {
			log.Fatalf("expr2.GetConsumablePotencies: %s", err)
		}
		resolveContext.StatusEffectPotency = statusEffectPotency
		resolveContext.ConsumablePotency = consumablePotency

		perkScaling, _, err := expr2.GetPerkMultipliers(dsStore)
		if err != nil {
			log.Fatalf("expr2.GetPerkMultipliers: %s", err)
		}
		resolveContext.PerkScaling = perkScaling
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

	keys := map[string]bool{}

	var id int64
	for _, file := range dsStore.GetAll() {
		f, err := os.Open(file.GetPath())
		if err != nil {
			log.Fatalf("os.Open: %s", err)
		}

		meta, err := datasheet.ParseMeta(f)
		if err != nil {
			log.Fatalf("datasheet.ParseMeta err: %s", err)
		}

		f.Close()

		key := fmt.Sprintf("%s.%s", meta.Type, meta.UniqueId)
		_, ok := keys[key]
		if ok {
			log.Printf("duplicate key: %q (%s)", key, file.GetPath())
			continue
		}
		keys[key] = true
		id++
		addTask(id, file)
	}

	pool.Close()
	pool.Wait()

	log.Printf("PeakMemory: %0.1fMb Duration: %s", float64(pr.GetPeakMemory())/1024/1024, pr.GetDuration().String())
}

func addTask(id int64, file *datasheet.DataSheetFile) {
	pool.AddTask(workerpool.NewTask(id, func(id int64) error {
		log.Printf("Working: %s", file.GetPath())
		ds, err := datasheet.Parse(file)
		if err != nil {
			return err
		}

		var outputPath string
		if keepStructure {
			relPath, err := filepath.Rel(inputDir, file.GetPath())
			if err != nil {
				return err
			}
			outputPath = strings.TrimSuffix(filepath.Join(outputDir, relPath), ".datasheet")
		} else {
			outputPath = filepath.Join(outputDir, ds.Type, ds.UniqueId)
		}

		if format == formatCsv {
			outputPath = outputPath + ".csv"
			err = storeToCsv(ds, outputPath)
			if err != nil {
				return err
			}
		}

		if format == formatJson {
			outputPath = outputPath + ".json"
			err = storeToJson(ds, outputPath)
			if err != nil {
				return err
			}
		}

		if format == formatYaml {
			outputPath = outputPath + ".yml"
			err = storeToYaml(ds, outputPath)
			if err != nil {
				return err
			}
		}

		return nil
	}))
}

func storeToCsv(ds *datasheet.DataSheet, path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.UseCRLF = false
	//csvWriter.Comma = ';'

	record := make([]string, len(ds.Columns))
	for i, column := range ds.Columns {
		record[i] = fmt.Sprintf("%s", column.Name)
	}

	err = csvWriter.Write(record)
	if err != nil {
		return err
	}

	for i, row := range ds.Rows {
		record := make([]string, len(row))
		for j, cell := range row {
			if ds.Columns[j].ColumnType == datasheet.ColumnTypeString {
				cell = resolveValue(cell, row, ds)
			}
			record[j] = toString(normalizeCellValue(ds.Columns[j], cell))
		}
		err = csvWriter.Write(record)
		if err != nil {
			return err
		}

		if i%100 == 0 {
			csvWriter.Flush()
		}
	}

	csvWriter.Flush()

	err = csvWriter.Error()
	if err != nil {
		return err
	}

	return nil
}

func storeToJson(ds *datasheet.DataSheet, path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	result := make([]*structure.OrderedMap[string, any], len(ds.Rows))
	for i, row := range ds.Rows {
		record := structure.NewOrderedMap[string, any]()
		for j, cell := range row {
			if ds.Columns[j].ColumnType == datasheet.ColumnTypeString {
				cell = resolveValue(cell, row, ds)
			}
			record.Add(fmt.Sprintf("%s", ds.Columns[j].Name), normalizeCellValue(ds.Columns[j], cell))
		}

		result[i] = record
	}

	encoder := json.NewEncoder(file)
	if withIndents {
		encoder.SetIndent("", "    ")
	}

	err = encoder.Encode(result)
	if err != nil {
		return err
	}

	return nil
}

func storeToYaml(ds *datasheet.DataSheet, path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	result := make([]*structure.OrderedMap[string, any], len(ds.Rows))
	for i, row := range ds.Rows {
		record := structure.NewOrderedMap[string, any]()
		for j, cell := range row {
			if ds.Columns[j].ColumnType == datasheet.ColumnTypeString {
				cell = resolveValue(cell, row, ds)
			}
			record.Add(fmt.Sprintf("%s", ds.Columns[j].Name), normalizeCellValue(ds.Columns[j], cell))
		}

		result[i] = record
	}

	encoder := yaml.NewEncoder(file, yaml.Indent(2))

	err = encoder.Encode(result)
	if err != nil {
		return err
	}

	return nil
}

func normalizeCellValue(column datasheet.ColumnData, str string) any {
	if column.ColumnType == datasheet.ColumnTypeString {
		return str
	}

	if column.ColumnType == datasheet.ColumnTypeNumber {
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil
		}

		//ugly rounding fix
		return floatFix(val)
	}

	if column.ColumnType == datasheet.ColumnTypeBoolean {
		val, err := strconv.ParseBool(str)
		if err != nil {
			return nil
		}

		return val
	}

	return str
}

func toString(val interface{}) string {
	str, ok := val.(string)
	if ok {
		return str
	}

	f, ok := val.(float64)
	if ok {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}

	b, ok := val.(bool)
	if ok {
		return strconv.FormatBool(b)
	}

	if val == nil {
		return ""
	}

	log.Fatalf("not supported value")

	return ""
}

var exprRe = regexp.MustCompile(`{\[[^-.0-9(]*([^]]+)\]}`)

func resolveValue(key string, row []string, ds *datasheet.DataSheet) string {
	if localizationData == nil || !strings.HasPrefix(key, "@") {
		return key
	}

	if localizationData.Has(key) {
		val := localizationData.Get(key)
		if resolveExpr {
			matches := exprRe.FindAllStringSubmatch(val, -1)
			if len(matches) > 0 {
				for _, match := range matches {
					exprStr := match[1]
					env := map[string]any{}
					if ds.Type == "StatusEffectData" {
						id, err := ds.GetCellValueByColumnName(row, "StatusID")
						if err != nil {
							break
						}
						env["perkMultiplier"] = float64(1)
						potency, ok := resolveContext.StatusEffectPotency[strings.ToLower(id)]
						if ok {
							env["ConsumablePotency"] = potency * float64(resolveContext.CharacterLevel)
						}
						exprStr = normalizeExpr(exprStr)
					}
					if ds.Type == "MasterItemDefinitions" {
						id, err := ds.GetCellValueByColumnName(row, "ItemID")
						if err != nil {
							break
						}
						potency, ok := resolveContext.ConsumablePotency[strings.ToLower(id)]
						if ok {
							env["ConsumablePotency"] = potency * float64(resolveContext.CharacterLevel)
						}
						exprStr = normalizeExpr(exprStr)
					}
					if ds.Type == "PerkData" {
						id, err := ds.GetCellValueByColumnName(row, "PerkID")
						if err != nil {
							break
						}
						scaling, ok := resolveContext.PerkScaling[strings.ToLower(id)]
						if ok {
							env["perkMultiplier"] = scaling.GetScaling(resolveContext.GearScore)
						} else {
							env["perkMultiplier"] = float64(1)
						}
						exprStr = normalizeExpr(exprStr)
					}

					program, err := expr.Compile(exprStr, expr.Env(env), expr.AsFloat64())
					if err != nil {
						log.Printf("expr.Compile: %s", err)
						continue
					}
					output, err := expr.Run(program, env)
					if err != nil {
						log.Fatalf("expr.Run: %s", err)
					}
					f64, ok := output.(float64)
					if !ok {
						log.Fatalf("not float64: %v", output)
					}
					f64 = floatFix(f64)
					val = strings.Replace(val, match[0], strconv.FormatFloat(f64, 'f', -1, 64), 1)
				}
			}
		}
		return val
	}

	return key
}

func floatFix(f64 float64) float64 {
	pow := math.Pow(10, 6)
	return math.Round(f64*pow) / pow
}

var re1 = regexp.MustCompile(`{([\w]+)}`)
var re2 = regexp.MustCompile(`[\d]([\s]+{[\w]+})`)

func normalizeExpr(exprStr string) string {
	var matches [][]string

	matches = re2.FindAllStringSubmatch(exprStr, 1)
	if len(matches) > 0 {
		exprStr = strings.Replace(exprStr, matches[0][1], " *"+matches[0][1], 1)
	}

	matches = re1.FindAllStringSubmatch(exprStr, 1)
	if len(matches) > 0 {
		exprStr = strings.Replace(exprStr, matches[0][0], matches[0][1], 1)
	}

	return exprStr
}
