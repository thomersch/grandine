package fileformat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

var vtypers map[string]valueTyper

func init() {
	vtypers = map[string]valueTyper{
		"string": func(s interface{}) ([]byte, Tag_ValueType, error) {
			return []byte(s.(string)), Tag_STRING, nil
		},
		"int": func(s interface{}) ([]byte, Tag_ValueType, error) {
			var buf = make([]byte, 8)
			binary.LittleEndian.PutUint64(buf, uint64(s.(int)))
			return buf, Tag_INT, nil
		},
		"float64": func(f interface{}) ([]byte, Tag_ValueType, error) {
			var buf = make([]byte, 8)
			binary.LittleEndian.PutUint64(buf, math.Float64bits(f.(float64)))
			return buf, Tag_DOUBLE, nil
		},
		"<nil>": func(f interface{}) ([]byte, Tag_ValueType, error) {
			return []byte{}, Tag_STRING, nil
		},
	}
}

type valueTyper func(interface{}) ([]byte, Tag_ValueType, error)

func ValueType(i interface{}) ([]byte, Tag_ValueType, error) {
	t := fmt.Sprintf("%T", i)
	vt, ok := vtypers[t]
	if !ok {
		return nil, Tag_STRING, fmt.Errorf("unknown type: %s (value: %v)", t, i)
	}
	return vt(i)
}

// KeyValue retrieves key and value from a Tag.
func KeyValue(t *Tag) (string, interface{}, error) {
	switch t.GetType() {
	case Tag_STRING:
		return t.Key, string(t.GetValue()), nil
	case Tag_INT:
		return t.Key, int(binary.LittleEndian.Uint64(t.GetValue())), nil
	case Tag_DOUBLE:
		var (
			buf = bytes.NewBuffer(t.GetValue())
			f   float64
		)
		err := binary.Read(buf, binary.LittleEndian, &f)
		return t.Key, f, err
	default:
		// TODO
		return t.Key, nil, nil
	}
}
