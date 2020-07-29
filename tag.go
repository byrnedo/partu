package partoo

import (
	"reflect"
)

// dumbStructToHeader converts a struct of string fields (with tag 'header') to a http.Header map
type namedField struct {
	Name  string
	Field interface{}
}

type namedFields []namedField

func (cm namedFields) Names() (ret ColNames) {
	idx := 0
	ret = make([]string, len(cm))
	for _, v := range cm {
		ret[idx] = v.Name
		idx++
	}
	return
}

func (cm namedFields) Fields() (ret []interface{}) {
	idx := 0
	ret = make([]interface{}, len(cm))
	for _, v := range cm {
		ret[idx] = v.Field
		idx++
	}
	return
}

func (p Builder) ColName(table Table, field interface{}) string {
	return findFieldTag(reflect.ValueOf(table), reflect.ValueOf(field))
}

func getColumnNames(table Table) (ret []namedField) {
	t := reflect.TypeOf(table)

	v := reflect.ValueOf(table)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for _, col := range table.Columns() {
		ret = append(ret, namedField{Name: findFieldTag(v, reflect.ValueOf(col)), Field: col})
	}
	return
}

func findFieldTag(structValue reflect.Value, fieldValue reflect.Value) string {
	sf := findStructField(structValue, fieldValue)

	colName := sf.Tag.Get("sql")
	if colName == "" {
		panic("struct field must have `sql` tag if included in Columns() output")
	}
	return colName
}

// Author: Copied from https://github.com/go-ozzo/ozzo-validation/
//
// findStructField looks for a field in the given struct.
// The field being looked for should be a pointer to the actual struct field.
// If found, the field info will be returned. Otherwise, nil will be returned.
func findStructField(structValue reflect.Value, fieldValue reflect.Value) *reflect.StructField {
	ptr := fieldValue.Pointer()
	for i := structValue.NumField() - 1; i >= 0; i-- {
		sf := structValue.Type().Field(i)
		if ptr == structValue.Field(i).UnsafeAddr() {
			// do additional type comparison because it's possible that the address of
			// an embedded struct is the same as the first field of the embedded struct
			if sf.Type == fieldValue.Elem().Type() {
				return &sf
			}
		}
		if sf.Anonymous {
			// delve into anonymous struct to look for the field
			fi := structValue.Field(i)
			if sf.Type.Kind() == reflect.Ptr {
				fi = fi.Elem()
			}
			if fi.Kind() == reflect.Struct {
				if f := findStructField(fi, fieldValue); f != nil {
					return f
				}
			}
		}
	}
	return nil
}
