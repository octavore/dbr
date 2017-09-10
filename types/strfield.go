package types

import "fmt"

type StrField struct {
	Field
}

func NewStrField(name string) *StrField {
	return &StrField{Field{name}}
}

func (f *StrField) Eq(val string) (string, string) {
	return fmt.Sprintf("%s = ?", f.name), val
}
