// Code generated by avrogen. DO NOT EDIT.

package sharedUnion

import "github.com/rogpeppe/avro"

type R struct {
	A interface{}
	B interface{}
}

// AvroRecord implements the avro.AvroRecord interface.
func (R) AvroRecord() avro.RecordInfo {
	return avro.RecordInfo{
		Schema: `{"fields":[{"name":"A","type":["int","string","float"]},{"name":"B","type":["int","string","float"]}],"name":"R","type":"record"}`,
		Fields: []avro.FieldInfo{
			0: {
				Info: &avro.TypeInfo{
					Type: new(interface{}),
					Union: []avro.TypeInfo{{
						Type: new(int),
					}, {
						Type: new(string),
					}, {
						Type: new(float32),
					}},
				},
			},
			1: {
				Info: &avro.TypeInfo{
					Type: new(interface{}),
					Union: []avro.TypeInfo{{
						Type: new(int),
					}, {
						Type: new(string),
					}, {
						Type: new(float32),
					}},
				},
			},
		},
	}
}

// TODO implement MarshalBinary and UnmarshalBinary methods?
