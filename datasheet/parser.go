package datasheet

import (
	"encoding/binary"
	"fmt"
	"github.com/new-world-tools/new-world-tools/reader"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
)

var headerSize = reflect.TypeOf(Header{}).NumField() * 4

func Parse(dataSheetFile *DataSheetFile) (*DataSheet, error) {
	meta, err := dataSheetFile.GetMeta()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(dataSheetFile.GetPath())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dataSheet, err := ParseBody(f, meta)
	if err != nil {
		return nil, err
	}

	return dataSheet, nil
}

func ParseMeta(r io.ReadSeeker) (*Meta, error) {
	var data []byte
	var u32 uint32
	var str string
	var err error

	meta := &Meta{}

	header, err := ParseHeader(r)
	if err != nil {
		return nil, err
	}
	meta.Header = header

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	meta.Field2 = data

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	meta.Field3 = data

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	meta.ColumnCount = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	meta.RowCount = u32

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	meta.Field6 = data

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	meta.Field7 = data

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	meta.Field8 = data

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	meta.Field9 = data

	columnsIndex := make([]*Column, meta.ColumnCount)
	for i := 0; i < int(meta.ColumnCount); i++ {
		column := &Column{}

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		column.Crc32 = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		column.Offset = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		column.ColumnType = u32

		columnsIndex[i] = column
	}
	meta.ColumnsIndex = columnsIndex

	rowsIndex := make([]*Row, meta.RowCount)
	for i := 0; i < int(meta.RowCount); i++ {
		row := &Row{
			Cells: make([]*Cell, meta.ColumnCount),
		}

		for j := 0; j < int(meta.ColumnCount); j++ {
			cell := &Cell{}

			u32, err = reader.ReadUint32(r, binary.LittleEndian)
			if err != nil {
				return nil, err
			}
			cell.Offset = u32

			data, err = reader.ReadBytes(r, 4)
			if err != nil {
				return nil, err
			}
			cell.Field2 = data

			row.Cells[j] = cell
		}

		rowsIndex[i] = row
	}
	meta.RowsIndex = rowsIndex

	str, err = reader.ReadNullTerminatedString(r)
	if err != nil {
		return nil, err
	}
	meta.WorksheetName = str

	uniqueId, err := reader.ReadNullTerminatedStringByOffset(r, int64(uint32(headerSize)+meta.Header.BodyOffset+meta.Header.UniqueIdOffset))
	if err != nil {
		return nil, err
	}
	meta.UniqueId = uniqueId

	typ, err := reader.ReadNullTerminatedStringByOffset(r, int64(uint32(headerSize)+meta.Header.BodyOffset+meta.Header.TypeOffset))
	if err != nil {
		return nil, err
	}
	meta.Type = typ

	return meta, nil
}

func ParseHeader(r io.Reader) (*Header, error) {
	var data []byte
	var u32 uint32
	var err error

	header := &Header{}

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	header.Signature = data

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	header.UniqueIdCrc32 = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	header.UniqueIdOffset = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	header.TypeCrc32 = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	header.TypeOffset = u32

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	header.Field6 = data

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	header.BodyLength = u32

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	header.Field8 = data

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	header.Field9 = data

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	header.Field10 = data

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	header.Field11 = data

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	header.Field12 = data

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	header.Field13 = data

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	header.Field14 = data

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	header.BodyOffset = u32

	return header, nil
}

func ParseBody(r io.ReadSeeker, meta *Meta) (*DataSheet, error) {
	dataSheet := &DataSheet{
		UniqueId: meta.UniqueId,
		Type:     meta.Type,
		Columns:  make([]ColumnData, meta.ColumnCount),
		Rows:     make([][]string, meta.RowCount),
	}

	for i, column := range meta.ColumnsIndex {
		name, err := reader.ReadNullTerminatedStringByOffset(r, int64(uint32(headerSize)+meta.Header.BodyOffset+column.Offset))
		if err != nil {
			return nil, err
		}
		dataSheet.Columns[i] = ColumnData{
			Name:       name,
			ColumnType: ColumnType(column.ColumnType),
		}
		if dataSheet.Columns[i].ColumnType > ColumnTypeBoolean {
			log.Printf("New ColumnType: %d", dataSheet.Columns[i].ColumnType)
		}
	}

	for i, row := range meta.RowsIndex {
		dataSheet.Rows[i] = make([]string, meta.ColumnCount)
		for j, cell := range row.Cells {
			value, err := reader.ReadNullTerminatedStringByOffset(r, int64(uint32(headerSize)+meta.Header.BodyOffset+cell.Offset))
			if err != nil {
				return nil, err
			}
			dataSheet.Rows[i][j] = value
		}
	}

	return dataSheet, nil
}

type DataSheet struct {
	Type     string
	UniqueId string
	Columns  []ColumnData
	Rows     [][]string
}

func (dataSheet *DataSheet) GetColumnIndexes() map[string]int {
	indexes := make(map[string]int, len(dataSheet.Columns))
	for i, column := range dataSheet.Columns {
		indexes[column.Name] = i
	}

	return indexes
}

func (dataSheet *DataSheet) GetCellValueByColumnName(row []string, columnName string) (string, error) {
	for i, column := range dataSheet.Columns {
		if strings.ToLower(column.Name) == strings.ToLower(columnName) {
			return row[i], nil
		}
	}

	return "", fmt.Errorf("column %q was not found", columnName)
}

type ColumnType int32

const (
	ColumnTypeString ColumnType = iota + 1
	ColumnTypeNumber
	ColumnTypeBoolean
)

type ColumnData struct {
	Name       string
	ColumnType ColumnType
}

type Meta struct {
	Header        *Header
	Field2        []byte `yaml:",flow"`
	Field3        []byte `yaml:"-"`
	ColumnCount   uint32
	RowCount      uint32
	Field6        []byte `yaml:"-"`
	Field7        []byte `yaml:"-"`
	Field8        []byte `yaml:"-"`
	Field9        []byte `yaml:"-"`
	ColumnsIndex  []*Column
	RowsIndex     []*Row
	WorksheetName string
	UniqueId      string
	Type          string
}

type Header struct {
	Signature      []byte `yaml:",flow"`
	UniqueIdCrc32  uint32
	UniqueIdOffset uint32
	TypeCrc32      uint32
	TypeOffset     uint32
	Field6         []byte `yaml:",flow"`
	BodyLength     uint32
	Field8         []byte `yaml:"-"`
	Field9         []byte `yaml:"-"`
	Field10        []byte `yaml:"-"`
	Field11        []byte `yaml:"-"`
	Field12        []byte `yaml:"-"`
	Field13        []byte `yaml:"-"`
	Field14        []byte `yaml:"-"`
	BodyOffset     uint32
}

type Column struct {
	Crc32      uint32
	Offset     uint32
	ColumnType uint32
}

type Row struct {
	Cells []*Cell
}

type Cell struct {
	Offset uint32
	Field2 []byte `yaml:",flow"`
}
