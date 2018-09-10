package types

import "fmt"

type Field struct {
	table string
	name  string
}

type GetFielder interface {
	GetField() *Field
}

func NewField(table, name string) *Field {
	return &Field{table, name}
}

func (f *Field) Eq(val interface{}) (string, interface{}) {
	return fmt.Sprintf("%s = ?", f.Full()), val
}

func (f *Field) NotNull() string {
	return fmt.Sprintf("%s is not null", f.Full())
}

func (f *Field) Null() string {
	return fmt.Sprintf("%s is null", f.Full())
}

func (f *Field) As(val interface{}) (string, interface{}) {
	return f.name, val
}

func (f *Field) String() string {
	return f.name
}

func (f *Field) Alias(alias string) string {
	return fmt.Sprintf("%s.%s", alias, f.name)
}

func (f *Field) Full() string {
	return f.Alias(f.table)
}

func (f *Field) Desc() (string, bool) {
	return f.Full(), false
}

func (f *Field) Asc() (string, bool) {
	return f.Full(), true
}

func (f *Field) On(g GetFielder) (string, string) {
	return f.table, fmt.Sprintf("%s = %s", f.Full(), g.GetField().Full())
}

func (f *Field) GetField() *Field {
	return f
}
