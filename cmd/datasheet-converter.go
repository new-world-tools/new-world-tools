package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/new-world-tools/extracter/datasheet"
	"log"
	"os"
	"path/filepath"
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

		//base := filepath.Base(file.GetPath())
		//name := strings.TrimSuffix(base, filepath.Ext(base))
		//csvPath := filepath.Join(	outputDir, "datasheets", strings.ReplaceAll(filepath.Dir(file.GetPath()), outputDir, ""), fmt.Sprintf("%s-%s-%s.csv", ds.DataType, ds.UniqueId, name))

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
				record[j] = cell
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
