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

type Partoo struct {
	dialect Dialect
}

func New(dialect Dialect) *Partoo {
	return &Partoo{dialect: dialect}
}

func (p Partoo) placeholder(i int) string {
	switch p.dialect {
	case Mysql:
		return fmt.Sprintf("$k")
	default:
		return fmt.Sprintf("$%d", i)
	}
}

func (p Partoo) placeholders(low, high int) string {
	parts := make([]string, high-low)
	for i := low; i < high; i++ {
		parts[i-low] = p.placeholder(i)
	}
	return strings.Join(parts, ",")
}

type Cols []interface{}

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

func (p Partoo) Select(t Table) string {
	cols := p.NamedFields(t)
	return fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(cols.Names().Strings(), ","),
		t.TableName(),
	)
}

func (p Partoo) SelectOne(t Table) (string, interface{}) {
	cols := p.NamedFields(t)
	first := cols[0]
	return fmt.Sprintf(
		"%s WHERE %s = %s",
		p.Select(t),
		first.Name,
		p.placeholder(1),
	), first.Field
}

func (p Partoo) NamedFields(t Table) NamedFields {
	return getColumnNames(t)
}

func (p Partoo) Insert(t Table) (string, []interface{}) {
	cols := p.NamedFields(t)

	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		t.TableName(),
		strings.Join(cols.Names()[1:].Strings(), ","),
		p.placeholders(1, len(cols)),
	), cols.Fields()[1:]
}

func (p Partoo) Update(t Table) (string, []interface{}) {
	namedFields := p.NamedFields(t)

	names := namedFields.Names()[1:]
	setPlaceholders := make([]string, len(names))
	for i, n := range names {
		setPlaceholders[i] = fmt.Sprintf("%s = %s", n, p.placeholder(i+1))
	}

	return fmt.Sprintf(
		"UPDATE %s SET %s",
		t.TableName(),
		strings.Join(setPlaceholders, ","),
	), namedFields.Fields()[1:]
}

func (p Partoo) UpdateOne(t Table) (string, []interface{}) {
	upd, args := p.Update(t)
	namedFields := p.NamedFields(t)
	fields := namedFields.Fields()
	names := namedFields.Names()
	args = append(args, fields[0])
	return fmt.Sprintf("%s WHERE %s = %s", upd, names[0], p.placeholder(len(fields))), args
}
