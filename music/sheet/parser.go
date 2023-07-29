package sheet

import (
	"encoding/binary"
	"fmt"
	"github.com/new-world-tools/new-world-tools/reader"
	"io"
	"math"
)

func Parse(r io.Reader) (*MusicSheet, error) {
	var u8 uint8
	var u32 uint32
	var err error

	sheet := &MusicSheet{
		Pattern:   0,
		NoteCount: 0,
		Notes:     nil,
	}

	u8, err = reader.ReadUint8(r)
	if err != nil {
		return nil, err
	}
	sheet.Pattern = u8

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	sheet.NoteCount = u32

	sheet.Notes = make([]Note, sheet.NoteCount)

	for i := uint32(0); i < sheet.NoteCount; i++ {
		note := Note{
			Key:   "",
			Track: 0,
			Time:  0,
		}

		u8, err = reader.ReadUint8(r)
		if err != nil {
			return nil, err
		}
		key, ok := keyMapping[u8]
		if !ok {
			return nil, fmt.Errorf("unknown key: 0x%02x", u8)
		}
		note.Key = key

		u8, err = reader.ReadUint8(r)
		if err != nil {
			return nil, err
		}
		track, ok := trackMapping[u8]
		if !ok {
			return nil, fmt.Errorf("unknown string: 0x%02x", u8)
		}
		note.Track = track

		data, err := reader.ReadBytes(r, 4)
		if err != nil {
			return nil, err
		}
		note.Time = math.Float32frombits(binary.LittleEndian.Uint32(data))

		sheet.Notes[i] = note
	}

	return sheet, nil
}
