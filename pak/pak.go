package pak

import (
	"archive/zip"
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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

var rePak = regexp.MustCompile(`.pak$`)
var reSortPak = regexp.MustCompile(`^((.*)[^\d]+)([\d].*).pak$`)

func FindAll(root string) ([]*Pak, error) {
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

	// @todo natural sorting
	sort.Slice(files, func(i, j int) bool {
		aPath := files[i].GetPath()
		bPath := files[j].GetPath()

		aMatches := reSortPak.FindAllStringSubmatch(aPath, -1)
		bMatches := reSortPak.FindAllStringSubmatch(bPath, -1)

		if aMatches == nil || bMatches == nil {
			return aPath < bPath
		}

		if aMatches[0][1] != bMatches[0][1] {
			return aPath < bPath
		}

		intA, err := strconv.Atoi(aMatches[0][3])
		if err != nil {
			return aPath < bPath
		}

		intB, err := strconv.Atoi(bMatches[0][3])
		if err != nil {
			return aPath < bPath
		}

		return intA < intB
	})

	return files, nil
}
