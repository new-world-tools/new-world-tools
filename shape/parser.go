package shape

import (
	"bytes"
	"encoding/binary"
	"github.com/new-world-tools/new-world-tools/reader"
	"io"
)

type VertexContainer struct {
	Version       uint32 // 0
	VerticesCount uint32
	Vertices      [][3]float32
	MetaDataCount uint32
	MetaData      []MetaDataElement
	Field6        uint32 //  0
	Field7        uint32 // 0
	Flags         []byte
	Field9        uint32 // 0
}

type MetaDataElement struct {
	Key   string
	Value string
}

func Parse(r io.Reader) (*VertexContainer, error) {
	var u32 uint32
	var f32 float32
	var data []byte
	var err error

	cont := &VertexContainer{}

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cont.Version = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cont.VerticesCount = u32

	cont.Vertices = make([][3]float32, cont.VerticesCount)

	for i := 0; i < int(cont.VerticesCount); i++ {
		cont.Vertices[i] = [3]float32{}

		for j := 0; j < 3; j++ {
			data, err := reader.ReadBytes(r, 4)
			if err != nil {
				return nil, err
			}

			buf := bytes.NewReader(data)
			err = binary.Read(buf, binary.LittleEndian, &f32)
			if err != nil {
				return nil, err
			}

			cont.Vertices[i][j] = f32
		}
	}

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cont.MetaDataCount = u32

	for i := 0; i < int(cont.MetaDataCount); i++ {
		element := MetaDataElement{}

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}

		data, err = reader.ReadBytes(r, int(u32))
		if err != nil {
			return nil, err
		}
		element.Key = string(data)

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}

		data, err = reader.ReadBytes(r, int(u32))
		if err != nil {
			return nil, err
		}
		element.Value = string(data)
	}

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cont.Field6 = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cont.Field7 = u32

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	cont.Flags = data

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cont.Field9 = u32

	return cont, nil
}
