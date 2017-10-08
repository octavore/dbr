package types

import "fmt"

type Field struct {
	name string
}

func NewField(name string) *Field {
	return &Field{name}
}

func (f *Field) Eq(val interface{}) (string, interface{}) {
	return fmt.Sprintf("%s = ?", f.name), val
}

func (f *Field) NotNull() string {
	return fmt.Sprintf("%s is not null", f.name)
}

func (f *Field) Null() string {
	return fmt.Sprintf("%s is null", f.name)
}

func (f *Field) As(val interface{}) (string, interface{}) {
	return f.name, val
}

func (f *Field) String() string {
	return f.name
}

func (f *Field) Desc() (string, bool) {
	return f.name, false
}

func (f *Field) Asc() (string, bool) {
	return f.name, true
}
