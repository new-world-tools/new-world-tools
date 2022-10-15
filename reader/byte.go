package reader

import (
	"encoding/binary"
	"io"
)

func ReadBytes(r io.Reader, size int) ([]byte, error) {
	buf := make([]byte, size)

	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func ReadUint8(r io.Reader) (uint8, error) {
	b, err := ReadBytes(r, 1)
	if err != nil {
		return 0, err
	}

	return b[0], nil
}

func ReadUint16(r io.Reader, byteOrder binary.ByteOrder) (uint16, error) {
	b, err := ReadBytes(r, 2)
	if err != nil {
		return 0, err
	}

	return byteOrder.Uint16(b), nil
}

func ReadUint32(r io.Reader, byteOrder binary.ByteOrder) (uint32, error) {
	b, err := ReadBytes(r, 4)
	if err != nil {
		return 0, err
	}

	return byteOrder.Uint32(b), nil
}

func ReadUint64(r io.Reader, byteOrder binary.ByteOrder) (uint64, error) {
	b, err := ReadBytes(r, 8)
	if err != nil {
		return 0, err
	}

	return byteOrder.Uint64(b), nil
}
