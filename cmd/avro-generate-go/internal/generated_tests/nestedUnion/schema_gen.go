// Code generated by avrogen. DO NOT EDIT.

package nestedUnion

import (
	"github.com/heetch/avro/avrotypegen"
)

type R struct {
	F []interface{}
}

// AvroRecord implements the avro.AvroRecord interface.
func (R) AvroRecord() avrotypegen.RecordInfo {
	return avrotypegen.RecordInfo{
		Schema: `{"fields":[{"name":"F","type":{"items":["int",{"items":["null","string"],"type":"array"}],"type":"array"}}],"name":"R","type":"record"}`,
		Required: []bool{
			0: true,
		},
		Unions: []avrotypegen.UnionInfo{
			0: {
				Type: new([]interface{}),
				Union: []avrotypegen.UnionInfo{{
					Type: new(int),
				}, {
					Type: new([]*string),
					Union: []avrotypegen.UnionInfo{{
						Type: nil,
					}, {
						Type: new(string),
					}},
				}},
			},
		},
	}
}

// TODO implement MarshalBinary and UnmarshalBinary methods?