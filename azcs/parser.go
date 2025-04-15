package azcs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/new-world-tools/new-world-tools/reader"
	azcsReader "github.com/new-world-tools/new-world-tools/reader/azcs"
	"io"
	"os"
)

const streamTag uint8 = 0x00

const (
	ST_BINARYFLAG_MASK             = 0xF8
	ST_BINARY_VALUE_SIZE_MASK      = 0x07
	ST_BINARYFLAG_ELEMENT_HEADER   = 1 << 3
	ST_BINARYFLAG_HAS_VALUE        = 1 << 4
	ST_BINARYFLAG_EXTRA_SIZE_FIELD = 1 << 5
	ST_BINARYFLAG_HAS_NAME         = 1 << 6
	ST_BINARYFLAG_HAS_VERSION      = 1 << 7
	ST_BINARYFLAG_ELEMENT_END      = 0
)

var (
	EOE = fmt.Errorf("end of element")
)

type Stream struct {
	Version  uint32
	Elements []*Element
}

type Element struct {
	Type            uuid.UUID
	SpecializedType uuid.UUID
	NameCrc         uint32
	Version         uint8
	DataSize        uint32
	Data            []byte
	Elements        []*Element
}

func (element Element) ResolveType() uuid.UUID {
	if !element.SpecializedType.IsNil() {
		return element.SpecializedType
	}

	return element.Type
}

func Parse(r io.Reader) (*Stream, error) {
	u8, err := reader.ReadUint8(r)
	if err != nil {
		return nil, err
	}
	if u8 != streamTag {
		return nil, errors.New("not valid stream")
	}

	u32, err := reader.ReadUint32(r, binary.BigEndian)
	if err != nil {
		return nil, err
	}

	stream := &Stream{
		Version:  u32,
		Elements: []*Element{},
	}

	for {
		element, err := readElement(r, stream)
		if err == EOE {
			return stream, nil
		}
		if err != nil {
			return nil, err
		}
		stream.Elements = append(stream.Elements, element)
	}

	return stream, nil
}

func readElement(r io.Reader, stream *Stream) (*Element, error) {
	var u8 uint8
	var u16 uint16
	var u32 uint32
	var data []byte
	var id uuid.UUID
	var err error

	element := &Element{
		Elements: []*Element{},
	}

	u8, err = reader.ReadUint8(r)
	if err != nil {
		return nil, err
	}
	flags := u8

	if flags == ST_BINARYFLAG_ELEMENT_END {
		return nil, EOE
	}

	if flags&ST_BINARYFLAG_HAS_NAME > 0 {
		u32, err = reader.ReadUint32(r, binary.BigEndian)
		if err != nil {
			return nil, err
		}
		element.NameCrc = u32
	}

	if flags&ST_BINARYFLAG_HAS_VERSION > 0 {
		u8, err = reader.ReadUint8(r)
		if err != nil {
			return nil, err
		}
		element.Version = u8
	}

	data, err = reader.ReadBytes(r, 16)
	if err != nil {
		return nil, err
	}
	id, err = uuid.FromBytes(data)
	if err != nil {
		return nil, err
	}
	element.Type = id

	if stream.Version == 2 {
		data, err = reader.ReadBytes(r, 16)
		if err != nil {
			return nil, err
		}
		id, err = uuid.FromBytes(data)
		if err != nil {
			return nil, err
		}
		element.SpecializedType = id
	}

	if flags&ST_BINARYFLAG_HAS_VALUE > 0 {
		valueBytes := flags & ST_BINARY_VALUE_SIZE_MASK
		if flags&ST_BINARYFLAG_EXTRA_SIZE_FIELD > 0 {
			switch valueBytes {
			case 1:
				u8, err = reader.ReadUint8(r)
				if err != nil {
					return nil, err
				}
				element.DataSize = uint32(u8)

			case 2:
				u16, err = reader.ReadUint16(r, binary.BigEndian)
				if err != nil {
					return nil, err
				}
				element.DataSize = uint32(u16)

			case 4:
				u32, err = reader.ReadUint32(r, binary.BigEndian)
				if err != nil {
					return nil, err
				}
				element.DataSize = u32

			default:
				return nil, fmt.Errorf("unsupported valueBytes: %d", valueBytes)
			}
		} else {
			element.DataSize = uint32(valueBytes)
		}
	}

	if element.DataSize > 0 {
		data, err = reader.ReadBytes(r, int(element.DataSize))
		if err != nil {
			return nil, err
		}
		element.Data = data
	}

	for {
		childElement, err := readElement(r, stream)
		if err == EOE {
			return element, nil
		}
		if err != nil {
			return nil, err
		}
		element.Elements = append(element.Elements, childElement)
	}

	return nil, errors.New("unexpected end")
}

func IsAzcsFile(path string) (bool, bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, false, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false, false, err
	}
	if fi.IsDir() {
		return false, false, nil
	}

	sampleSize := 5

	size := int(fi.Size())
	if size < sampleSize {
		return false, false, nil
	}

	sampleData, err := reader.ReadBytes(f, sampleSize)
	if err != nil {
		return false, false, err
	}

	if isCompressed(sampleData) {
		return true, true, nil
	}

	if isUncompressed(sampleData) {
		return true, false, nil
	}

	return false, false, nil
}

var uncompressedSignatures = [][]byte{
	{0x00, 0x00, 0x00, 0x00, 0x03},
	{0x00, 0x00, 0x00, 0x00, 0x02},
	{0x00, 0x00, 0x00, 0x00, 0x01},
}

func isUncompressed(data []byte) bool {
	for _, uncompressedSignature := range uncompressedSignatures {
		if len(data) >= len(uncompressedSignature) && bytes.Equal(uncompressedSignature, data[:len(uncompressedSignature)]) {
			return true
		}
	}

	return false
}

func isCompressed(data []byte) bool {
	if len(data) < len(azcsReader.Signature) {
		return false
	}

	return bytes.Equal(azcsReader.Signature, data[:len(azcsReader.Signature)])
}

func GetReader(path string, isCompressed bool) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser

	rc = f

	if isCompressed {
		r, err := azcsReader.NewReader(rc)
		if err != nil {
			return nil, err
		}
		rc = io.NopCloser(r)
	}

	return rc, nil
}
