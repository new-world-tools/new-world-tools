package pak

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/new-world-tools/go-oodle"
	"io"
)

var ErrUnsupportedMethod = errors.New("Unsupported method.")

func Parse(pak *Pak) ([]*File, error) {
	zipReader, err := zip.OpenReader(pak.GetPath())
	if err != nil {
		return nil, err
	}

	files := []*File{}
	for _, archivedFile := range zipReader.File {
		files = append(files, &File{
			zipFile: archivedFile,
			Name:    archivedFile.Name,
		})
	}

	return files, nil
}

type File struct {
	zipFile *zip.File
	Name    string
}

func (file *File) Decompress() (io.ReadCloser, error) {
	if file.zipFile.Method == 0x00 || file.zipFile.Method == 0x08 {
		reader, err := file.zipFile.Open()
		if err != nil {
			return nil, err
		}

		return reader, nil
	}

	if file.zipFile.Method == 0x0f {
		reader, err := file.zipFile.OpenRaw()
		if err != nil {
			return nil, err
		}

		data, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		data, err = oodle.Decompress(data, int64(file.zipFile.UncompressedSize64))
		if err != nil {
			return nil, err
		}

		return io.NopCloser(bytes.NewBuffer(data)), nil
	}

	return nil, ErrUnsupportedMethod
}
