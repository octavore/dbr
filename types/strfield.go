package types

import "fmt"

type StrField struct {
	Field
}

func NewStrField(table, name string) *StrField {
	return &StrField{Field{table, name}}
}

func (f *StrField) As(val string) (string, string) {
	return f.Full(), val
}

func (f *StrField) Eq(val string) (string, string) {
	return fmt.Sprintf("%s = ?", f.Full()), val
}

func (f *StrField) In(val []string) (string, []string) {
	return fmt.Sprintf("%s IN ?", f.Full()), val
}
