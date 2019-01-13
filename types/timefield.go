package types

import (
	"fmt"
	"time"
)

type TimeField struct {
	Field
}

func NewTimeField(table, name string) *TimeField {
	return &TimeField{Field{table, name}}
}

func (f *TimeField) As(val time.Time) (string, time.Time) {
	return f.name, val
}

func (f *TimeField) Eq(val time.Time) (string, time.Time) {
	return fmt.Sprintf("%s = ?", f.Full()), val
}

func (f *TimeField) Gt(val time.Time) (string, time.Time) {
	return fmt.Sprintf("%s > ?", f.Full()), val
}

func (f *TimeField) Lt(val time.Time) (string, time.Time) {
	return fmt.Sprintf("%s < ?", f.Full()), val
}
