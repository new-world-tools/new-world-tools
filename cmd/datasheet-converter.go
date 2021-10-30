package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/new-world-tools/extracter/datasheet"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	inputDirPtr := flag.String("input", ".\\assets", "directory path")
	outputDirPtr := flag.String("output", ".\\assets\\datasheets", "directory path")
	flag.Parse()

	inputDir, err := filepath.Abs(filepath.Clean(*inputDirPtr))
	if err != nil {
		log.Fatalf("filepath.Abs: %s", err)
	}

	_, err = os.Stat(inputDir)
	if os.IsNotExist(err) {
		log.Fatalf("'%s' does not exist", inputDir)
	}

	outputDir, err := filepath.Abs(filepath.Clean(*outputDirPtr))
	if err != nil {
		log.Fatalf("filepath.Abs: %s", err)
	}

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("MkdirAll: %s", err)
	}

	files, err := datasheet.FindAll(inputDir)
	if err != nil {
		log.Fatalf("datasheet.FindAll: %s", err)
	}

	for _, file := range files {
		log.Printf("Working: %s", file.GetPath())
		ds, _ := datasheet.Parse(file)
		if err != nil {
			log.Fatalf("datasheet.Parse: %s", err)
		}

		csvPath := filepath.Join(outputDir, "datasheets", ds.DataType, fmt.Sprintf("%s.csv", ds.UniqueId))
		err = os.MkdirAll(filepath.Dir(csvPath), 0755)
		if err != nil {
			log.Fatalf("os.MkdirAll: %s", err)
		}

		file, err := os.Create(csvPath)
		if err != nil {
			log.Fatalf("os.Create: %s", err)
		}

		csvWriter := csv.NewWriter(file)
		csvWriter.UseCRLF = false
		//csvWriter.Comma = ';'

		record := make([]string, len(ds.Columns))
		for i, column := range ds.Columns {
			record[i] = fmt.Sprintf("%s", column.Name)
		}

		err = csvWriter.Write(record)
		if err != nil {
			log.Fatalf("csvWriter.Write: %s", err)
		}

		for i, row := range ds.Rows {
			record := make([]string, len(row))
			for j, cell := range row {
				record[j] = normalizeCellValue(ds.Columns[j], cell)
			}
			err = csvWriter.Write(record)
			if err != nil {
				log.Fatalf("csvWriter.Write: %s", err)
			}

			if i%100 == 0 {
				csvWriter.Flush()
			}
		}

		csvWriter.Flush()

		err = csvWriter.Error()
		if err != nil {
			log.Fatalf(" csvWriter.Error: %s", err)
		}

		file.Close()
	}
}

func normalizeCellValue(column datasheet.Column, str string) string {
	if column.ColumnType == datasheet.ColumnTypeString {
		return str
	}

	if column.ColumnType == datasheet.ColumnTypeNumber {
		if strings.Index(str, ".") > 0 || strings.Index(str, "E") > 0 {
			f, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return str
			}

			return strconv.FormatFloat(f, 'f', -1, 64)
		}
	}

	if column.ColumnType == datasheet.ColumnTypeBoolean {
		if str == "1" {
			return "true"
		}
		if str == "2" {
			return "false"
		}
	}

	return str
}
