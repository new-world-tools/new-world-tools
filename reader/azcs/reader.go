package azcs

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"github.com/new-world-tools/new-world-tools/reader"
	"io"
)

const signature = "AZCS"

type header struct {
	signature        string
	compressorId     uint32
	uncompressedSize uint64
}

type zLibHeader struct {
	numSeekPoints uint32
}

type zLibSeekPoint struct {
	compressedOffset   uint64
	uncompressedOffset uint64
}

func NewReader(r io.Reader) (io.Reader, error) {
	buf := bufio.NewReaderSize(r, 1024*1024)
	r = buf

	headerData := &header{}

	data, err := reader.ReadBytes(r, len([]byte(signature)))
	if err != nil {
		return nil, err
	}

	if string(data) != signature {
		return nil, fmt.Errorf("wrong signature - %q. Must be %q", string(data), signature)
	}
	headerData.signature = signature

	compressorId, err := reader.ReadUint32(r, binary.BigEndian)
	if err != nil {
		return nil, err
	}
	headerData.compressorId = compressorId

	uncompressedSize, err := reader.ReadUint64(r, binary.BigEndian)
	if err != nil {
		return nil, err
	}
	headerData.uncompressedSize = uncompressedSize

	switch headerData.compressorId {
	case 0x73887d3a:
		return handleZlib(r)

	case 0x72fd505e:
		return nil, fmt.Errorf("zstd is not implemented")
	}

	return nil, fmt.Errorf("unsupported commpressorId: 0x%08x", headerData.compressorId)
}

func handleZlib(r io.Reader) (io.Reader, error) {
	numSeekPoints, err := reader.ReadUint32(r, binary.BigEndian)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	numSeekPointsSize := numSeekPoints * 16

	buf := bytes.NewBuffer(data[:len(data)-int(numSeekPointsSize)])

	zr, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}

	return zr, nil
}
