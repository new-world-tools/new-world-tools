package datasheet

import (
	"io/fs"
	"path/filepath"
	"regexp"
)

type DataSheetFile struct {
	path string
}

func (dataSheet DataSheetFile) GetPath() string {
	return dataSheet.path
}

func NewDataSheet(path string) *DataSheetFile {
	return &DataSheetFile{
		path: path,
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

		files = append(files, NewDataSheet(path))

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
