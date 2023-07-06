package azcs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/new-world-tools/new-world-tools/reader"
	"github.com/new-world-tools/new-world-tools/structure"
	"math"
	"reflect"
	"strconv"
)

const (
	TypeField  = "__type"
	ValueField = "__value"
)

type JsonFloat64 float64

func (v JsonFloat64) MarshalJSON() ([]byte, error) {
	f64 := float64(v)
	var s string
	switch {
	case math.IsInf(f64, 1):
		s = "+Inf"
	case math.IsInf(f64, -1):
		s = "-Inf"
	case math.IsNaN(f64):
		s = "NaN"
	default:
		s = strconv.FormatFloat(f64, 'f', -1, 64)
		return []byte(s), nil
	}
	return []byte(`"` + s + `"`), nil
}

func (v *JsonFloat64) UnmarshalJSON(b []byte) error {
	switch {
	case bytes.Equal(b, []byte(`"+Inf"`)):
		*v = JsonFloat64(math.Inf(1))
	case bytes.Equal(b, []byte(`"-Inf"`)):
		*v = JsonFloat64(math.Inf(-1))
	case bytes.Equal(b, []byte(`"NaN"`)):
		*v = JsonFloat64(math.NaN())
	default:
		n, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			return err
		}
		*v = JsonFloat64(n)
	}
	return nil
}

type JsonFloat32 float32

func (v JsonFloat32) MarshalJSON() ([]byte, error) {
	f64 := float64(v)
	var s string
	switch {
	case math.IsInf(f64, 1):
		s = "+Inf"
	case math.IsInf(f64, -1):
		s = "-Inf"
	case math.IsNaN(f64):
		s = "NaN"
	default:
		s = strconv.FormatFloat(f64, 'f', -1, 32)
		return []byte(s), nil
	}
	return []byte(`"` + s + `"`), nil
}

func (v *JsonFloat32) UnmarshalJSON(b []byte) error {
	switch {
	case bytes.Equal(b, []byte(`"+Inf"`)):
		*v = JsonFloat32(math.Inf(1))
	case bytes.Equal(b, []byte(`"-Inf"`)):
		*v = JsonFloat32(math.Inf(-1))
	case bytes.Equal(b, []byte(`"NaN"`)):
		*v = JsonFloat32(math.NaN())
	default:
		n, err := strconv.ParseFloat(string(b), 32)
		if err != nil {
			return err
		}
		*v = JsonFloat32(n)
	}
	return nil
}

type float32s interface {
	float32 | JsonFloat32
}

type float64s interface {
	float64 | JsonFloat64
}

type floats interface {
	float32s | float64s
}

func createFloatArray[T floats](data []byte) ([]T, error) {
	var z T
	dataTypeSize := int(reflect.Indirect(reflect.ValueOf(z)).Type().Size())

	fs := make([]T, len(data)/dataTypeSize)
	for i := 0; i < len(fs); i++ {
		var f T
		buf := bytes.NewReader(data[i*4 : (i+1)*4])
		err := binary.Read(buf, binary.BigEndian, &f)
		if err != nil {
			return nil, err
		}
		fs[i] = f
	}

	return fs, nil
}

func maxFloat[T floats](data []T) T {
	var max T
	for i, f := range data {
		if i == 0 {
			max = f
		} else {
			if f > max {
				max = f
			}
		}
	}

	return max
}

func ResolveStream(stream *Stream, typeResolver TypeResolver, hashResolver HashResolver) (any, error) {
	if len(stream.Elements) > 1 {
		return nil, fmt.Errorf("too much elements")
	}

	node, err := resolveNode(stream.Elements[0], typeResolver, hashResolver)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func resolveNode(element *Element, typeResolver TypeResolver, hashResolver HashResolver) (any, error) {
	node := structure.NewOrderedMap[string, any]()

	node.Add(TypeField, typeResolver(element))

	switch element.ResolveType().String() {
	case
		// Transform
		"5d9958e9-9f1e-4985-b532-fffde75fedfd":
		switch element.Version {
		case 0:
			if len(element.Data) != 48 {
				return nil, fmt.Errorf("wrong size: %d", len(element.Data))
			}
			col0, err := createFloatArray[JsonFloat32](element.Data[0:12])
			if err != nil {
				return nil, fmt.Errorf("createFloatArrray: %s", err)
			}
			col1, err := createFloatArray[JsonFloat32](element.Data[12:24])
			if err != nil {
				return nil, fmt.Errorf("createFloatArrray: %s", err)
			}
			col2, err := createFloatArray[JsonFloat32](element.Data[24:36])
			if err != nil {
				return nil, fmt.Errorf("createFloatArrray: %s", err)
			}
			col3, err := createFloatArray[JsonFloat32](element.Data[36:48])
			if err != nil {
				return nil, fmt.Errorf("createFloatArrray: %s", err)
			}

			tnode := structure.NewOrderedMap[string, any]()

			tnode.Add("rotation/scale", []JsonFloat32{col0[0], col0[1], col0[2], col1[0], col1[1], col1[2], col2[0], col2[1], col2[2]})
			tnode.Add("translation", col3)

			node.Add(ValueField, tnode)

		case 1:
			if len(element.Data) != 40 {
				return nil, fmt.Errorf("wrong size: %d", len(element.Data))
			}
			rotation, err := createFloatArray[JsonFloat32](element.Data[0:16])
			if err != nil {
				return nil, fmt.Errorf("createFloatArrray: %s", err)
			}
			vectorScale, err := createFloatArray[JsonFloat32](element.Data[16:28])
			if err != nil {
				return nil, fmt.Errorf("createFloatArrray: %s", err)
			}
			translation, err := createFloatArray[JsonFloat32](element.Data[28:40])
			if err != nil {
				return nil, fmt.Errorf("createFloatArrray: %s", err)
			}
			scale := maxFloat[JsonFloat32](vectorScale)

			tnode := structure.NewOrderedMap[string, any]()

			tnode.Add("rotation", rotation)
			tnode.Add("scale", scale)
			tnode.Add("translation", translation)

			node.Add(ValueField, tnode)

		default:
			return nil, fmt.Errorf("unsupported version: %d", element.Version)
		}

		return node, nil

	case
		// Color
		"7894072a-9050-4f0f-901b-34b1a0d29417":
		f32s, err := createFloatArray[JsonFloat32](element.Data)
		if err != nil {
			return nil, fmt.Errorf("createFloatArrray: %s", err)
		}

		node.Add(ValueField, f32s)

		return node, nil

	// Asset
	case "77a19d40-8731-4d3c-9041-1b43047366a4":
		buf := bytes.NewBuffer(element.Data)

		data, err := reader.ReadBytes(buf, 16)
		if err != nil {
			return nil, fmt.Errorf("reader.ReadBytes: %s", err)
		}
		id, _ := uuid.FromBytes(data)
		node.Add("guid", id.String())

		data, err = reader.ReadBytes(buf, 16)
		if err != nil {
			return nil, fmt.Errorf("reader.ReadBytes: %s", err)
		}
		id, _ = uuid.FromBytes(data)
		node.Add("subId", id.String())

		data, err = reader.ReadBytes(buf, 16)
		if err != nil {
			return nil, fmt.Errorf("reader.ReadBytes: %s", err)
		}
		id, _ = uuid.FromBytes(data)
		node.Add("type", id.String())

		u64, err := reader.ReadUint64(buf, binary.BigEndian)
		if err != nil {
			return nil, fmt.Errorf("reader.ReadUint64: %s", err)
		}
		if u64 > 0 {
			data, err = reader.ReadBytes(buf, int(u64))
			if err != nil {
				return nil, fmt.Errorf("reader.ReadBytes: %s", err)
			}
			node.Add("hint", string(data))
		}

		return node, nil
	}

	_, v, _ := node.GetByPosition(0)

	switch v {
	case "bool":
		var b bool
		l := len(element.Data)
		if l != 1 {
			return nil, fmt.Errorf("unsupported bool size: %d", l)
		}

		switch element.Data[0] {
		case 0x00:
			b = false

		case 0x01:
			b = true

		default:
			b = true
		}

		return b, nil

	case "AZStd::string":
		return string(element.Data), nil

	case "AZ::Uuid":
		id, _ := uuid.FromBytes(element.Data)
		return id.String(), nil

	case
		"unsigned char",
		"unsigned int",
		"unsigned short",
		"AZ::u64":
		l := len(element.Data)

		switch l {
		case 1:
			return element.Data[0], nil

		case 2:
			return binary.BigEndian.Uint16(element.Data), nil

		case 4:
			return binary.BigEndian.Uint32(element.Data), nil

		case 8:
			return binary.BigEndian.Uint64(element.Data), nil

		default:
			return nil, fmt.Errorf("unsupported data size: %d", l)
		}

	case
		"int",
		"AZ::s64":
		l := len(element.Data)

		switch l {
		case 1:
			var i8 int8
			buf := bytes.NewReader(element.Data)
			err := binary.Read(buf, binary.BigEndian, &i8)
			if err != nil {
				return nil, fmt.Errorf("binary.Read: %s", err)
			}
			return i8, nil

		case 2:
			var i16 int16
			buf := bytes.NewReader(element.Data)
			err := binary.Read(buf, binary.BigEndian, &i16)
			if err != nil {
				return nil, fmt.Errorf("binary.Read: %s", err)
			}
			return i16, nil

		case 4:
			var i32 int32
			buf := bytes.NewReader(element.Data)
			err := binary.Read(buf, binary.BigEndian, &i32)
			if err != nil {
				return nil, fmt.Errorf("binary.Read: %s", err)
			}
			return i32, nil

		case 8:
			var i64 int64
			buf := bytes.NewReader(element.Data)
			err := binary.Read(buf, binary.BigEndian, &i64)
			if err != nil {
				return nil, fmt.Errorf("binary.Read: %s", err)
			}
			return i64, nil

		default:
			return nil, fmt.Errorf("unsupported data size: %d", l)
		}

	case
		"float",
		"double":
		l := len(element.Data)

		switch l {
		case 4:
			var f32 JsonFloat32
			buf := bytes.NewReader(element.Data)
			err := binary.Read(buf, binary.BigEndian, &f32)
			if err != nil {
				return nil, fmt.Errorf("binary.Read: %s", err)
			}
			return f32, nil

		case 8:
			var f64 JsonFloat64
			buf := bytes.NewReader(element.Data)
			err := binary.Read(buf, binary.BigEndian, &f64)
			if err != nil {
				return nil, fmt.Errorf("binary.Read: %s", err)
			}
			return f64, nil

		default:
			return nil, fmt.Errorf("unsupported data size: %d", l)
		}

	case
		"AZStd::vector",
		"AZStd::unordered_set":
		nodes := make([]any, len(element.Elements))

		for i, element := range element.Elements {
			value, err := resolveNode(element, typeResolver, hashResolver)
			if err != nil {
				return nil, err
			}
			nodes[i] = value
		}

		return nodes, nil

	case
		"Quaternion",
		"Vector2",
		"Vector3":
		l := len(element.Data)

		f32s := make([]JsonFloat32, l/4)
		for i := 0; i < len(f32s); i++ {
			var f32 JsonFloat32
			buf := bytes.NewReader(element.Data[i*4 : (i+1)*4])
			err := binary.Read(buf, binary.BigEndian, &f32)
			if err != nil {
				return nil, fmt.Errorf("binary.Read: %s", err)
			}
			f32s[i] = f32
		}

		return f32s, nil

	case
		"Amazon::Pervasives::UID",
		"Amazon::Hub::ActorRef",
		"BitSet":
		return element.Data, nil

	case
		"AZStd::array",
		"AZStd::fixed_vector",
		"AZStd::list",
		"AZStd::map",
		"AZStd::unordered_map",
		"AZStd::unordered_flat_map",
		"5f9f78d5-bdf7-5531-961d-8a91dfa2e126",
		"c2fd8c07-90d3-5d02-bab7-b1fac968c43f",
		"fdbd40b4-8a70-5b23-bce7-a717ba039a86",
		"9c719dd5-f8d3-59d3-b55b-627422922a43":
		values := make([]any, len(element.Elements))
		for i, element := range element.Elements {
			key := hashResolver(element)
			value, err := resolveNode(element, typeResolver, hashResolver)
			if err != nil {
				return nil, err
			}

			if key != "element" {
				return nil, fmt.Errorf("wrong key: %s", key)
			}

			values[i] = value
		}
		return values, nil

	case
		"AZStd::intrusive_ptr",
		"AZStd::shared_ptr",
		"AZStd::unique_ptr":
		if len(element.Elements) == 0 {
			return nil, nil
		}
		if len(element.Elements) != 1 {
			return nil, fmt.Errorf("wrong elements count: %d", len(element.Elements))
		}
		for _, element := range element.Elements {
			key := hashResolver(element)
			value, err := resolveNode(element, typeResolver, hashResolver)
			if err != nil {
				return nil, err
			}

			if key != "element" {
				return nil, fmt.Errorf("wrong key: %s", key)
			}

			return value, nil
		}

	case
		"EntityId":
		if len(element.Elements) != 1 {
			return nil, fmt.Errorf("wrong elements count: %d", len(element.Elements))
		}
		for _, element := range element.Elements {
			key := hashResolver(element)
			value, err := resolveNode(element, typeResolver, hashResolver)
			if err != nil {
				return nil, err
			}

			if key != "id" {
				return nil, fmt.Errorf("wrong key: %s", key)
			}

			return value, nil
		}

	default:
		if len(element.Data) > 0 {
			_, v, _ := node.GetByPosition(0)

			vs := v.(string)

			checkId, err := uuid.FromString(vs)
			if err != nil || checkId.IsNil() {
				if len(element.Data) > 24 {
					return nil, fmt.Errorf("unsupported data type: %s, type: %s, stype: %s", vs, element.Type.String(), element.SpecializedType.String())
				}
				return nil, fmt.Errorf("unsupported data type: %s, type: %s, stype: %s data: %x", vs, element.Type.String(), element.SpecializedType.String(), element.Data)
			} else {
				node.Add(ValueField, element.Data)
			}
		} else {
			for _, element := range element.Elements {
				key := hashResolver(element)
				value, err := resolveNode(element, typeResolver, hashResolver)
				if err != nil {
					return nil, err
				}

				node.Add(key, value)
			}
		}
	}

	return node, nil
}

type TypeResolver func(element *Element) string
type HashResolver func(element *Element) string

var (
	DefaultTypeResolver TypeResolver = func(element *Element) string {
		typ := element.Type.String()
		if !element.SpecializedType.IsNil() {
			typ = element.SpecializedType.String()
		}

		if DefaultTypeRegistry.Has(typ) {
			return DefaultTypeRegistry.Get(typ)
		}

		return typ
	}
	DefaultHashResolver HashResolver = func(element *Element) string {
		hash := element.NameCrc

		if DefaultHashRegistry.Has(hash) {
			return DefaultHashRegistry.Get(hash)
		}

		formattedHash := fmt.Sprintf("0x%08x", hash)

		return formattedHash
	}
)
