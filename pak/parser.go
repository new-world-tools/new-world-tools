package pak

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/flate"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"github.com/zelenin/go-oodle-lz"
	"io"
	"time"
)

var ErrUnsupportedMethod = errors.New("unsupported method")

func Parse(pak *Pak) ([]*File, error) {
	zipReader, err := zip.OpenReader(pak.GetPath())
	if err != nil {
		return nil, err
	}

	files := make([]*File, len(zipReader.File))
	for i, archivedFile := range zipReader.File {
		files[i] = &File{
			zipFile: archivedFile,
			Name:    archivedFile.Name,
		}
	}

	return files, nil
}

type File struct {
	zipFile *zip.File
	Name    string
}

func (file *File) Decompress() (io.ReadCloser, error) {
	var r io.Reader
	var rc io.ReadCloser
	var err error

	switch file.zipFile.Method {
	case 0x00:
		rc, err = file.zipFile.Open()
		if err != nil {
			return nil, err
		}

	case 0x08:
		r, err = file.zipFile.OpenRaw()
		if err != nil {
			return nil, err
		}

		bufReader := bufio.NewReaderSize(r, 4096)

		sigData, err := bufReader.Peek(2)
		if err == nil {
			if isZlib(sigData) {
				rc, err = zlib.NewReader(bufReader)
				if err != nil {
					return nil, err
				}
			} else {
				rc = flate.NewReader(bufReader)
			}
		} else {
			rc = flate.NewReader(bufReader)
		}

	case 0x0f:
		r, err = file.zipFile.OpenRaw()
		if err != nil {
			return nil, err
		}

		// works on linux, bugged on windows
		//rc, err = oodle.NewReader(r, int64(file.zipFile.UncompressedSize64))
		//if err != nil {
		//	return nil, err
		//}

		data, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		data, err = oodle.Decompress(data, int(file.zipFile.UncompressedSize64))
		if err != nil {
			return nil, err
		}

		rc = io.NopCloser(bytes.NewBuffer(data))

	default:
		return nil, ErrUnsupportedMethod
	}

	return rc, nil
}

func (file *File) GetModifiedTime() time.Time {
	return file.zipFile.Modified
}

type zlibHeader struct {
	cmf struct {
		cm    uint8 // 8
		cinfo uint8 // <=7
	}
	flg struct {
		fcheck uint8
		fdict  uint8 // 0-1
		flevel uint8 // 0-3
	}
}

func isZlib(sigData []byte) bool {
	cmfByte := sigData[0]
	flgByte := sigData[1]
	zh := &zlibHeader{
		cmf: struct {
			cm    uint8
			cinfo uint8
		}{
			cm:    (cmfByte >> 0) & 0b1111,
			cinfo: (cmfByte >> 4) & 0b1111,
		},
		flg: struct {
			fcheck uint8
			fdict  uint8
			flevel uint8
		}{
			fcheck: (flgByte >> 0) & 0b11111,
			fdict:  (flgByte >> 5) & 0b1,
			flevel: (flgByte >> 6) & 0b11,
		},
	}
	return zh.cmf.cm == 0x08 && zh.cmf.cinfo <= 0x07 && zh.flg.fdict <= 0x01 && zh.flg.flevel <= 0x03 && binary.BigEndian.Uint16(sigData)%31 == 0
}
