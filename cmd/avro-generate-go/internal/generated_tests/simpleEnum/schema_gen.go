// Code generated by avrogen. DO NOT EDIT.

package simpleEnum

import (
	"fmt"
	"github.com/heetch/avro/avrotypegen"
	"strconv"
)

type MyEnum int

const (
	MyEnumA MyEnum = iota
	MyEnumB
	MyEnumC
)

var _MyEnum_strings = []string{
	"a",
	"b",
	"c",
}

// String returns the textual representation of MyEnum.
func (e MyEnum) String() string {
	if e < 0 || int(e) >= len(_MyEnum_strings) {
		return "MyEnum(" + strconv.FormatInt(int64(e), 10) + ")"
	}
	return _MyEnum_strings[e]
}

// MarshalText implements encoding.TextMarshaler
// by returning the textual representation of MyEnum.
func (e MyEnum) MarshalText() ([]byte, error) {
	if e < 0 || int(e) >= len(_MyEnum_strings) {
		return nil, fmt.Errorf("MyEnum value %d is out of bounds", e)
	}
	return []byte(_MyEnum_strings[e]), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
// by expecting the textual representation of MyEnum.
func (e *MyEnum) UnmarshalText(data []byte) error {
	// Note for future: this could be more efficient.
	for i, s := range _MyEnum_strings {
		if string(data) == s {
			*e = MyEnum(i)
			return nil
		}
	}
	return fmt.Errorf("unknown value %q for MyEnum", data)
}

type R struct {
	E MyEnum
}

// AvroRecord implements the avro.AvroRecord interface.
func (R) AvroRecord() avrotypegen.RecordInfo {
	return avrotypegen.RecordInfo{
		Schema: `{"fields":[{"name":"E","type":{"name":"MyEnum","symbols":["a","b","c"],"type":"enum"}}],"name":"R","type":"record"}`,
		Required: []bool{
			0: true,
		},
	}
}
