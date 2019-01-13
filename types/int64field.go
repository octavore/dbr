package types

import "fmt"

type Int64Field struct {
	Field
}

func NewInt64Field(table, name string) *Int64Field {
	return &Int64Field{Field{table, name}}
}

func (f *Int64Field) As(val int64) (string, int64) {
	return f.name, val
}

func (f *Int64Field) Eq(val int64) (string, int64) {
	return fmt.Sprintf("%s = ?", f.Full()), val
}

func (f *Int64Field) Neq(val int64) (string, int64) {
	return fmt.Sprintf("%s != ?", f.Full()), val
}

func (f *Int64Field) In(val []int64) (string, []int64) {
	return fmt.Sprintf("%s IN ?", f.Full()), val
}

func (f *Int64Field) NOtIn(val []int64) (string, []int64) {
	return fmt.Sprintf("%s NOT IN ?", f.Full()), val
}
