package datasheet

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
)

func Parse(dataSheetFile *DataSheetFile) (*DataSheet, error) {
	file, err := os.Open(dataSheetFile.GetPath())
	if err != nil {
		return nil, err
	}

	defer file.Close()

	meta, err := parseMeta(file)
	if err != nil {
		return nil, err
	}

	dataSheet, err := parseBody(file, meta)
	if err != nil {
		return nil, err
	}

	return dataSheet, nil
}

func parseMeta(r *os.File) (*meta, error) {
	var headerSize int32

	signature, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	field2, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	uniqueIdOffset, err := readInt32(r)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	field4, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	dataTypeOffset, err := readInt32(r)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	field6, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	bodyLength, err := readInt32(r)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	field8, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	field9, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	field10, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	field11, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	field12, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	field13, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	field14, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}
	headerSize += 4

	bodyOffset, err := readInt32(r)
	if err != nil {
		return nil, err
	}
	headerSize += 4
	bodyOffset += headerSize

	crc32, err := readUint32(r)
	if err != nil {
		return nil, err
	}

	field17, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}

	columnCount, err := readInt32(r)
	if err != nil {
		return nil, err
	}

	rowCount, err := readInt32(r)
	if err != nil {
		return nil, err
	}

	field20, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}

	field21, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}

	field22, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}

	field23, err := readBytes(r, 4)
	if err != nil {
		return nil, err
	}

	columns := []*column{}
	for i := int32(0); i < columnCount; i++ {
		field1, err := readBytes(r, 4)
		if err != nil {
			return nil, err
		}

		offset, err := readInt32(r)
		if err != nil {
			return nil, err
		}

		columnType, err := readInt32(r)
		if err != nil {
			return nil, err
		}

		columns = append(columns, &column{
			field1:     field1,
			offset:     offset,
			columnType: columnType,
		})
	}

	rows := []*row{}
	for i := int32(0); i < rowCount; i++ {
		row := &row{
			cells: []*cell{},
		}
		for j := int32(0); j < columnCount; j++ {
			offset, err := readInt32(r)
			if err != nil {
				return nil, err
			}

			field2, err := readBytes(r, 4)
			if err != nil {
				return nil, err
			}

			row.cells = append(row.cells, &cell{
				offset: offset,
				field2: field2,
			})
		}

		rows = append(rows, row)
	}

	return &meta{
		signature:      signature,
		field2:         field2,
		uniqueIdOffset: uniqueIdOffset,
		field4:         field4,
		dataTypeOffset: dataTypeOffset,
		field6:         field6,
		bodyLength:     bodyLength,
		field8:         field8,
		field9:         field9,
		field10:        field10,
		field11:        field11,
		field12:        field12,
		field13:        field13,
		field14:        field14,
		bodyOffset:     bodyOffset,

		crc32:       crc32,
		field17:     field17,
		columnCount: columnCount,
		rowCount:    rowCount,
		field20:     field20,
		field21:     field21,
		field22:     field22,
		field23:     field23,

		columns: columns,
		rows:    rows,
	}, nil
}

func parseBody(r *os.File, meta *meta) (*DataSheet, error) {
	// OUTPUT
	_, err := readNullTerminatedString(r)
	if err != nil {
		return nil, err
	}

	uniqueId, err := readNullTerminatedStringByOffset(r, meta.bodyOffset+meta.uniqueIdOffset)
	if err != nil {
		return nil, err
	}

	dataType, err := readNullTerminatedStringByOffset(r, meta.bodyOffset+meta.dataTypeOffset)
	if err != nil {
		return nil, err
	}

	dataSheet := &DataSheet{
		UniqueId: uniqueId,
		DataType: dataType,
		Columns:  make([]Column, meta.columnCount),
		Rows:     make([][]string, meta.rowCount),
	}

	for i, column := range meta.columns {
		name, err := readNullTerminatedStringByOffset(r, meta.bodyOffset+column.offset)
		if err != nil {
			return nil, err
		}
		dataSheet.Columns[i] = Column{
			Name:       name,
			ColumnType: ColumnType(column.columnType),
		}
		if dataSheet.Columns[i].ColumnType > ColumnTypeBoolean {
			log.Printf("New ColumnType: %d", dataSheet.Columns[i].ColumnType)
		}
	}

	for i, row := range meta.rows {
		dataSheet.Rows[i] = make([]string, meta.columnCount)
		for j, cell := range row.cells {
			value, err := readNullTerminatedStringByOffset(r, meta.bodyOffset+cell.offset)
			if err != nil {
				return nil, err
			}
			dataSheet.Rows[i][j] = value
		}
	}

	return dataSheet, nil
}

func readBytes(r io.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)

	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func readInt32(r io.Reader) (int32, error) {
	b, err := readBytes(r, 4)
	if err != nil {
		return 0, err
	}

	return int32(binary.LittleEndian.Uint32(b)), nil
}

func readUint32(r io.Reader) (uint32, error) {
	b, err := readBytes(r, 4)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(b), nil
}

func readNullTerminatedString(r io.Reader) (string, error) {
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

func readNullTerminatedStringByOffset(r *os.File, offset int32) (string, error) {
	pos, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}

	defer r.Seek(pos, io.SeekStart)

	_, err = r.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return "", err
	}

	str, err := readNullTerminatedString(r)
	if err != nil {
		return "", err
	}

	return str, nil
}

func skipBytes(r *bytes.Reader, n int) error {
	_, err := readBytes(r, n)
	return err
}

type DataSheet struct {
	UniqueId string
	DataType string
	Columns  []Column
	Rows     [][]string
}

type ColumnType int32

const (
	ColumnTypeString ColumnType = iota + 1
	ColumnTypeNumber
	ColumnTypeBoolean
)

type Column struct {
	Name       string
	ColumnType ColumnType
}

type meta struct {
	signature      []byte
	field2         []byte
	uniqueIdOffset int32
	field4         []byte
	dataTypeOffset int32
	field6         []byte
	bodyLength     int32
	field8         []byte
	field9         []byte
	field10        []byte
	field11        []byte
	field12        []byte
	field13        []byte
	field14        []byte
	bodyOffset     int32
	crc32          uint32
	field17        []byte
	columnCount    int32
	rowCount       int32
	field20        []byte
	field21        []byte
	field22        []byte
	field23        []byte
	columns        []*column
	rows           []*row
}

type column struct {
	field1     []byte
	offset     int32
	columnType int32
}

type row struct {
	cells []*cell
}

type cell struct {
	offset int32
	field2 []byte
}
