package partoo

import (
	"fmt"
	"strings"
)

type Dialect string

const (
	Mysql    Dialect = "mysql"
	Postgres Dialect = "postgres"
)

var (
	dialect = Postgres
)

func SetDialect(d Dialect) {
	dialect = d
}

func placeholder(i int) string {
	switch dialect {
	case Mysql:
		return fmt.Sprintf("$k")
	default:
		return fmt.Sprintf("$%d", i)
	}
}

func placeholders(low, high int) string {
	parts := make([]string, high-low)
	for i := low; i < high; i++ {
		parts[i-low] = placeholder(i)
	}
	return strings.Join(parts, ",")
}

type Col struct {
	Name  string
	Field interface{}
}
type Cols []Col

type Table interface {
	// The name of the ... table
	TableName() string
	// The columns mapped to the fields in your struct, by ref please
	Columns() Cols
}

type ColNames []string

func (c ColNames) Prefix(alias string) ColNames {

	if alias != "" {
		alias += "."
	}
	for i, colName := range c {
		c[i] = alias + colName
	}
	return c
}

func (c ColNames) Strings() []string {
	return c
}

func (cm Cols) Names() (ret ColNames) {
	idx := 0
	ret = make([]string, len(cm))
	for _, v := range cm {
		ret[idx] = v.Name
		idx++
	}
	return
}

func (cm Cols) Fields() (ret []interface{}) {

	idx := 0
	ret = make([]interface{}, len(cm))
	for _, v := range cm {
		ret[idx] = v.Field
		idx++
	}
	return
}

func Select(t Table) string {
	cols := t.Columns()
	return fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(cols.Names().Strings(), ","),
		t.TableName(),
	)
}

func SelectOne(t Table) (string, interface{}) {
	cols := t.Columns()
	return fmt.Sprintf(
		"%s WHERE %s = %s",
		Select(t),
		cols.Names()[0],
		placeholder(1),
	), cols.Fields()[0]
}

func Insert(t Table) (string, []interface{}) {
	cols := t.Columns()

	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		t.TableName(),
		strings.Join(cols.Names()[1:].Strings(), ","),
		placeholders(1, len(cols)),
	), cols.Fields()[1:]
}

func Update(t Table) (string, []interface{}) {
	cols := t.Columns()

	names := cols.Names()[1:]
	setPlaceholders := make([]string, len(names))
	for i, n := range names {
		setPlaceholders[i] = fmt.Sprintf("%s = %s", n, placeholder(i+1))
	}

	return fmt.Sprintf(
		"UPDATE %s SET %s",
		t.TableName(),
		strings.Join(setPlaceholders, ","),
	), cols.Fields()[1:]
}

func UpdateOne(t Table) (string, []interface{}) {
	upd, args := Update(t)
	cols := t.Columns()
	fields := cols.Fields()
	names := cols.Names()
	args = append(args, fields[0])
	return fmt.Sprintf("%s WHERE %s = %s", upd, names[0], placeholder(len(fields))), args
}
