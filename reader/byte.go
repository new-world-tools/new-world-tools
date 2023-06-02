package reader

import (
	"bytes"
	"encoding/binary"
	"io"
)

func ReadBytes(r io.Reader, count int) ([]byte, error) {
	buf := make([]byte, count)

	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func SkipBytes(r io.Reader, count int64) error {
	_, err := io.CopyN(io.Discard, r, count)
	if err != nil {
		return err
	}

	return nil
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

func ReadNullTerminatedString(r io.Reader) (string, error) {
	buf := bytes.NewBuffer(nil)

	for {
		b := make([]byte, 1)
		_, err := r.Read(b)
		if err != nil {
			return "", err
		}

		if b[0] == 0x00 {
			break
		}

		_, err = buf.Write(b)
		if err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

func ReadNullTerminatedStringByOffset(r io.ReadSeeker, offset int64) (string, error) {
	pos, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}

	defer r.Seek(pos, io.SeekStart)

	_, err = r.Seek(offset, io.SeekStart)
	if err != nil {
		return "", err
	}

	str, err := ReadNullTerminatedString(r)
	if err != nil {
		return "", err
	}

	return str, nil
}
