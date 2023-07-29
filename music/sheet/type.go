package sheet

import "fmt"

type MusicSheet struct {
	Pattern   uint8 // Novice, Skilled, Expert
	NoteCount uint32
	Notes     []Note
}

var keyMapping = map[uint8]string{
	0x01: "w",
	0x02: "a",
	0x03: "s",
	0x04: "d",
	0x05: "sp",
	0x06: "lmb",
	0x07: "rmb",
}

var trackMapping = map[uint8]uint8{
	0xfe: 1,
	0xff: 2,
	0x00: 3,
	0x01: 4,
	0x02: 5,
}

type Note struct {
	Key   string
	Track uint8
	Time  float32
}

func (note Note) String() string {
	return fmt.Sprintf("key: %q, track: %d, time: %.3f", note.Key, note.Track, note.Time)
}
