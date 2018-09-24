package dbr

import (
	"reflect"
)

type Table struct {
	Name       string
	Columns    []string
	PrimaryKey string
}

func MustNewTable(table string, t interface{}) *Table {
	tbl, err := NewTable(table, t)
	if err != nil {
		panic(err)
	}
	return tbl
}

func NewTable(table string, t interface{}) (*Table, error) {
	tbl := &Table{Name: table, PrimaryKey: "id"}
	v := reflect.ValueOf(t)
	if v.Kind() != reflect.Struct {
		return nil, ErrNotSupported
	}

	for k, _ := range colNames(t) {
		tbl.Columns = append(tbl.Columns, k)
	}
	return tbl, nil
}

func (t *Table) colsWithName() (out []string) {
	for _, col := range t.Columns {
		out = append(out, t.Name+"."+col)
	}
	return
}

func (t *Table) Select(session SessionRunner) *SelectBuilder {
	return session.Select(t.colsWithName()...).From(t.Name)
}

func (t *Table) Update(session SessionRunner) *UpdateBuilder {
	return session.Update(t.Name)
}

func (t *Table) Insert(session SessionRunner) *InsertBuilder {
	colsWithPKey := []string{}
	for _, col := range t.Columns {
		if col != t.PrimaryKey {
			colsWithPKey = append(colsWithPKey, col)
		}
	}
	return session.InsertInto(t.Name).Columns(colsWithPKey...)
}

func (t *Table) UpdateRecord(session SessionRunner, record interface{}) *UpdateBuilder {
	data := colNames(record)
	updateBuilder := session.Update(t.Name)
	for k, v := range data {
		updateBuilder = updateBuilder.Set(k, v.Interface())
	}
	return updateBuilder
}

func colNames(i interface{}) map[string]reflect.Value {
	cols := map[string]reflect.Value{}
	v := reflect.ValueOf(i)
	t := reflect.TypeOf(i)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		col := field.Tag.Get("db")
		if col == "-" {
			continue
		}
		if col == "" {
			col = camelCaseToSnakeCase(field.Name)
		}
		cols[col] = v.Field(i)
	}
	return cols
}
