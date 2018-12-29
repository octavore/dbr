package dbr

import (
	"database/sql/driver"
	"reflect"
	"strings"
)

var NameMapping = camelCaseToSnakeCase

func isUpper(b byte) bool {
	return 'A' <= b && b <= 'Z'
}

func isLower(b byte) bool {
	return 'a' <= b && b <= 'z'
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func toLower(b byte) byte {
	if isUpper(b) {
		return b - 'A' + 'a'
	}
	return b
}

func camelCaseToSnakeCase(name string) string {
	var buf strings.Builder
	buf.Grow(len(name) * 2)

	for i := 0; i < len(name); i++ {
		buf.WriteByte(toLower(name[i]))
		if i != len(name)-1 && isUpper(name[i+1]) &&
			(isLower(name[i]) || isDigit(name[i]) ||
				(i != len(name)-2 && isLower(name[i+2]))) {
			buf.WriteByte('_')
		}
	}

	return buf.String()
}

var (
	typeValuer = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
)

type tagInfo struct {
	tag  string
	name string
}

type tagStore struct {
	m map[reflect.Type][]tagInfo
}

func newTagStore() *tagStore {
	return &tagStore{
		m: make(map[reflect.Type][]tagInfo),
	}
}

func (s *tagStore) get(t reflect.Type) []tagInfo {
	if t.Kind() != reflect.Struct {
		return nil
	}
	if _, ok := s.m[t]; !ok {
		l := make([]tagInfo, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" && !field.Anonymous {
				// unexported
				continue
			}
			tag := field.Tag.Get("db")
			fieldName := tag
			if tag == "-" {
				// ignore
				continue
			}
			if fieldName == "" {
				// no fieldName from tag, but we can record the field name
				fieldName = NameMapping(field.Name)
			}
			l[i] = tagInfo{tag: tag, name: fieldName}
		}
		s.m[t] = l
	}
	return s.m[t]
}

func (s *tagStore) findPtr(value reflect.Value, name []string, ptr []interface{}) error {
	if value.CanAddr() && value.Addr().Type().Implements(typeScanner) {
		ptr[0] = value.Addr().Interface()
		return nil
	}
	switch value.Kind() {
	case reflect.Struct:
		s.findValueByName(value, name, ptr, true, nil)
		return nil
	case reflect.Ptr:
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		return s.findPtr(value.Elem(), name, ptr)
	default:
		ptr[0] = value.Addr().Interface()
		return nil
	}
}

func (s *tagStore) findValueByName(value reflect.Value, name []string, ret []interface{}, retPtr bool, prefix *string) {
	if value.Type().Implements(typeValuer) {
		return
	}
	switch value.Kind() {
	case reflect.Ptr:
		// embedded ptr
		if value.IsNil() {
			return
		}
		s.findValueByName(value.Elem(), name, ret, retPtr, prefix)
	case reflect.Struct:
		// embedded type
		l := s.get(value.Type())
		for i := 0; i < value.NumField(); i++ {
			tagInfo := l[i]
			if tagInfo.name == "" {
				continue
			}
			fieldValue := value.Field(i)

			queryColName := tagInfo.name
			if prefix != nil {
				queryColName = *prefix + "___" + tagInfo.name
			}
			for i, want := range name {
				if want != queryColName {
					continue
				}
				if ret[i] == nil {
					if retPtr {
						ret[i] = fieldValue.Addr().Interface()
					} else {
						ret[i] = fieldValue
					}
				}
			}
			var newPrefix *string
			if tagInfo.tag != "" {
				newPrefix = &tagInfo.tag
			}
			s.findValueByName(fieldValue, name, ret, retPtr, newPrefix)
		}
	}
}
