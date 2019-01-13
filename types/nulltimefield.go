package types

import (
	"time"
)

type NullTimeField struct {
	TimeField
}

func NewNullTimeField(table, name string) *NullTimeField {
	return &NullTimeField{TimeField{Field{table, name}}}
}

func (f *NullTimeField) As(val *time.Time) (string, *time.Time) {
	return f.name, val
}
