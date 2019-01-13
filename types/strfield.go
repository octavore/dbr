package types

import "fmt"

type StrField struct {
	Field
}

func NewStrField(table, name string) *StrField {
	return &StrField{Field{table, name}}
}

func (f *StrField) As(val string) (string, string) {
	return f.name, val
}

func (f *StrField) Eq(val string) (string, string) {
	return fmt.Sprintf("%s = ?", f.Full()), val
}

func (f *StrField) Neq(val string) (string, string) {
	return fmt.Sprintf("%s != ?", f.Full()), val
}

func (f *StrField) In(val []string) (string, []string) {
	return fmt.Sprintf("%s IN ?", f.Full()), val
}

func (f *StrField) NotIn(val []string) (string, []string) {
	return fmt.Sprintf("%s NOT IN ?", f.Full()), val
}
