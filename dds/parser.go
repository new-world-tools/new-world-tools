package dds

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	bin "github.com/zelenin/go-binary"
)

var signature = []byte{'D', 'D', 'S', ' '}

type Meta struct {
	Signature   []byte
	Header      *Header
	HeaderDxt10 *HeaderDxt10
}

type Header struct {
	Size              uint32
	Flags             uint32
	Height            uint32
	Width             uint32
	PitchOrLinearSize uint32
	Depth             uint32
	MipMapCount       uint32
	Reserved1         [11]uint32
	PixelFormat       *PixelFormat
	Caps              uint32
	Caps2             uint32
	Caps3             uint32
	Caps4             uint32
	Reserved2         uint32
}

type HeaderDxt10 struct {
	DxgiFormat        uint32
	ResourceDimension uint32
	MiscFlag          uint32
	ArraySize         uint32
	MiscFlags2        uint32
}

type PixelFormat struct {
	Size        uint32
	Flags       uint32
	FourCc      []byte
	RgbBitCount uint32
	RBitMask    uint32
	GBitMask    uint32
	BBitMask    uint32
	ABitMask    uint32
}

func ParseMeta(r io.Reader) (*Meta, error) {
	var (
		data []byte
		u32  uint32
		err  error
	)

	meta := &Meta{}

	buf := bin.NewReader(r, binary.LittleEndian)

	header := &Header{}

	data, err = buf.ReadBytes(4)
	if err != nil {
		return nil, fmt.Errorf("failed to read signature: %w", err)
	}
	if !bytes.Equal(signature, data) {
		return nil, fmt.Errorf("invalid DDS signature: expected %s, got %s", signature, data)
	}
	meta.Signature = data

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	header.Size = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	header.Flags = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	header.Height = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	header.Width = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	header.PitchOrLinearSize = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	header.MipMapCount = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	header.Depth = u32

	reserved1 := [11]uint32{}
	for i := range reserved1 {
		u32, err = buf.ReadUint32()
		if err != nil {
			return nil, fmt.Errorf("failed to read uint32: %w", err)
		}
		reserved1[i] = u32
	}
	header.Reserved1 = reserved1

	pixelFormat := &PixelFormat{}

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	pixelFormat.Size = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	pixelFormat.Flags = u32

	data, err = buf.ReadBytes(4)
	if err != nil {
		return nil, fmt.Errorf("failed to read signature: %w", err)
	}
	pixelFormat.FourCc = data

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	pixelFormat.RgbBitCount = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	pixelFormat.RBitMask = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	pixelFormat.GBitMask = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	pixelFormat.BBitMask = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	pixelFormat.ABitMask = u32

	header.PixelFormat = pixelFormat

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	header.Caps = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	header.Caps2 = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	header.Caps3 = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	header.Caps4 = u32

	u32, err = buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read uint32: %w", err)
	}
	header.Reserved2 = u32

	meta.Header = header

	if bytes.Equal(meta.Header.PixelFormat.FourCc, []byte{'D', 'X', '1', '0'}) {
		headerDxt10 := &HeaderDxt10{}

		u32, err = buf.ReadUint32()
		if err != nil {
			return nil, fmt.Errorf("failed to read uint32: %w", err)
		}
		headerDxt10.DxgiFormat = u32

		u32, err = buf.ReadUint32()
		if err != nil {
			return nil, fmt.Errorf("failed to read uint32: %w", err)
		}
		headerDxt10.ResourceDimension = u32

		u32, err = buf.ReadUint32()
		if err != nil {
			return nil, fmt.Errorf("failed to read uint32: %w", err)
		}
		headerDxt10.MiscFlag = u32

		u32, err = buf.ReadUint32()
		if err != nil {
			return nil, fmt.Errorf("failed to read uint32: %w", err)
		}
		headerDxt10.ArraySize = u32

		u32, err = buf.ReadUint32()
		if err != nil {
			return nil, fmt.Errorf("failed to read uint32: %w", err)
		}
		headerDxt10.MiscFlags2 = u32

		meta.HeaderDxt10 = headerDxt10
	}

	return meta, nil
}
