package datasheet

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

type DataSheetFile struct {
	path    string
	relPath string
	meta    *Meta
}

func (dataSheetFile DataSheetFile) GetPath() string {
	return dataSheetFile.path
}

func (dataSheetFile DataSheetFile) GetRelPath() string {
	return dataSheetFile.relPath
}

func (dataSheetFile *DataSheetFile) GetMeta() (*Meta, error) {
	if dataSheetFile.meta == nil {
		f, err := os.Open(dataSheetFile.GetPath())
		if err != nil {
			return nil, err
		}
		defer f.Close()

		meta, err := ParseMeta(f)
		if err != nil {
			return nil, err
		}

		dataSheetFile.meta = meta
	}

	return dataSheetFile.meta, nil
}

func (dataSheetFile *DataSheetFile) GetData() (*DataSheet, error) {
	ds, err := Parse(dataSheetFile)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func NewDataSheet(path string, relPath string) *DataSheetFile {
	return &DataSheetFile{
		path:    path,
		relPath: relPath,
	}
}

func FindAll(root string) ([]*DataSheetFile, error) {
	rePak := regexp.MustCompile(`.datasheet$`)

	files := []*DataSheetFile{}

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil || !rePak.MatchString(info.Name()) {
			return err
		}

		path, err = filepath.Abs(filepath.Clean(path))
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(filepath.FromSlash(root), filepath.FromSlash(path))
		if err != nil {
			return err
		}

		files = append(files, NewDataSheet(path, relPath))

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].GetPath() < files[j].GetPath()
	})

	return files, nil
}
