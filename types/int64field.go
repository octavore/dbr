package types

import "fmt"

type Int64Field struct {
	Field
}

func NewInt64Field(name string) *Int64Field {
	return &Int64Field{Field{name}}
}

func (f *Int64Field) Eq(val int64) (string, int64) {
	return fmt.Sprintf("%s = ?", f.name), val
}
