package dbrgen

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gocraft/dbr"

	"github.com/fatih/camelcase"
)

type Generator struct {
	// out           io.Writer
	baseDir       string
	models        []reflect.Type
	autoSnakeCase bool
}

func NewGenerator(baseDir string) *Generator {
	return &Generator{
		baseDir:       baseDir,
		models:        []reflect.Type{},
		autoSnakeCase: true,
	}
}

func (g *Generator) Register(i interface{}) {
	v := reflect.ValueOf(i).Type()
	if v.Kind() != reflect.Struct {
		panic("can only register struct types")
	}
	g.models = append(g.models, v)
}

var int64type = reflect.TypeOf(dbr.NullInt64{})
var strtype = reflect.TypeOf(dbr.NullString{})

func convertFieldType(t reflect.Type) string {
	switch t {
	case int64type:
		return "types.NewInt64Field"
	case strtype:
		return "types.NewStrField"
	}

	switch t.Kind() {
	case reflect.Ptr:
		return convertFieldType(t.Elem())
	case reflect.String:
		return "types.NewStrField"
	case reflect.Int64:
		return "types.NewInt64Field"
	default:
		return "types.NewField"
	}
}

func (g *Generator) Generate() error {
	for _, v := range g.models {
		buf := &bytes.Buffer{}
		pkgName := snakecase(v.Name())

		fmt.Fprintf(buf, "package %s\n", pkgName)
		fmt.Fprintf(buf, "import \"github.com/gocraft/dbr/types\"\n")
		fmt.Fprintf(buf, "var (")
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			name := snakecaseUpper(f.Name)
			fieldName := dbFieldname(f)
			fieldType := convertFieldType(f.Type)
			fmt.Fprintf(buf, "  %s = %s(%q)\n", name, fieldType, fieldName)
		}
		fmt.Fprintf(buf, ")")

		dir := filepath.Join(g.baseDir, pkgName)
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
		path := filepath.Join(dir, pkgName+".go")
		err = writeFormattedGo(path, buf)
		if err != nil {
			return err
		}
	}

	return nil
}

// extract the field name from the field. prefers db
// declared json name if it exists.
func dbFieldname(f reflect.StructField) string {
	dbTag, ok := f.Tag.Lookup("db")
	if ok {
		return strings.Split(dbTag, ",")[0]
	}
	return strings.ToLower(f.Name)
}

func snakecase(str string) string {
	return strings.ToLower(strings.Join(camelcase.Split(str), "_"))
}

func snakecaseUpper(str string) string {
	return strings.ToUpper(strings.Join(camelcase.Split(str), "_"))
}

// https://github.com/drone/sqlgen/blob/master/fmt.go
func writeFormattedGo(path string, in io.Reader) error {
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	gofmt := exec.Command("gofmt", "-s")
	gofmt.Stdin = in
	gofmt.Stdout = f
	gofmt.Stderr = os.Stderr
	return gofmt.Run()
}
