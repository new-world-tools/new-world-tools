package pak

import (
	"archive/zip"
	"io/fs"
	"path/filepath"
	"regexp"
)

type Pak struct {
	path      string
	zipReader *zip.ReadCloser
}

func (pak Pak) GetPath() string {
	return pak.path
}

func (pak *Pak) GetFiles() ([]*File, error) {
	return Parse(pak)
}

func (pak *Pak) Close() error {
	if pak.zipReader != nil {
		return pak.zipReader.Close()
	}

	return nil
}

func NewPak(path string) *Pak {
	return &Pak{
		path: path,
	}
}

func FindAll(root string) ([]*Pak, error) {
	rePak := regexp.MustCompile(`.pak$`)

	files := []*Pak{}

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil || !rePak.MatchString(info.Name()) {
			return err
		}

		path, err = filepath.Abs(filepath.Clean(path))
		if err != nil {
			return err
		}

		files = append(files, NewPak(path))

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
