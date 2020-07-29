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

type Cols []interface{}

type Table interface {
	// The name of the ... table
	TableName() string
	// The columns mapped to the fields in your struct, by ref please
	Columns() Cols
}

// Implement if you want to control inserting the id
type AutoID interface {
	AutoID() bool
}

type Builder struct {
	dialect Dialect
	tagName string
}

func New(dialect Dialect) *Builder {
	return &Builder{dialect: dialect}
}

func (p *Builder) SetTag(tag string) Builder {
	p.tagName = tag
	return p
}

func  (p Builder) Tag() string {
	if p.tagName == "" {
		return "sql"
	}
	return p.tagName
}

func (p Builder) placeholder(i int) string {
	switch p.dialect {
	case Mysql:
		return fmt.Sprintf("$k")
	default:
		return fmt.Sprintf("$%d", i)
	}
}

func (p Builder) placeholders(low, high int) string {
	parts := make([]string, high-low)
	for i := low; i < high; i++ {
		parts[i-low] = p.placeholder(i)
	}
	return strings.Join(parts, ",")
}

func (p Builder) Select(t Table) string {
	cols := p.NamedFields(t)
	return fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(cols.Names().Strings(), ","),
		t.TableName(),
	)
}

func (p Builder) SelectOne(t Table) (string, interface{}) {
	cols := p.NamedFields(t)
	first := cols[0]
	return fmt.Sprintf(
		"%s WHERE %s = %s",
		p.Select(t),
		first.Name,
		p.placeholder(1),
	), first.Field
}

func (p Builder) NamedFields(t Table) namedFields {
	return getColumnNames(p.Tag(), t)
}

func (p Builder) Insert(t Table) (string, []interface{}) {
	cols := p.NamedFields(t)

	autoID, ok := t.(AutoID)
	colsToInsert := cols.Names()[1:].Strings()
	args := cols.Fields()[1:]
	if ok && !autoID.AutoID() {
		colsToInsert = cols.Names().Strings()
		args = cols.Fields()
	}

	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		t.TableName(),
		strings.Join(colsToInsert, ","),
		p.placeholders(1, len(colsToInsert)+1),
	), args
}

func (p Builder) Update(t Table) (string, []interface{}) {
	namedFields := p.NamedFields(t)

	setPlaceholders := p.generateUpdatePlaceholders(namedFields, 1)

	return fmt.Sprintf(
		"UPDATE %s SET %s",
		t.TableName(),
		setPlaceholders,
	), namedFields.Fields()[1:]
}

func (p Builder) UpdateOne(t Table) (string, []interface{}) {
	upd, args := p.Update(t)
	namedFields := p.NamedFields(t)
	fields := namedFields.Fields()
	names := namedFields.Names()
	args = append(args, fields[0])
	return fmt.Sprintf("%s WHERE %s = %s", upd, names[0], p.placeholder(len(fields))), args
}

func (p Builder) UpsertOne(t Table) (string, []interface{}) {
	if p.dialect == Mysql {
		return p.upsertMysql(t)
	}
	return p.upsertPostgres(t)
}

func (p Builder) generateUpdatePlaceholders(cols namedFields, startIndex int) string {

	names := cols.Names()[1:]
	setPlaceholders := make([]string, len(names))
	for i, n := range names {
		setPlaceholders[i] = fmt.Sprintf("%s = %s", n, p.placeholder(i+startIndex))
	}

	return strings.Join(setPlaceholders, ",")
}

func (p Builder) upsertMysql(t Table) (string, []interface{}) {

	cols := p.NamedFields(t)

	setPlaceholders := p.generateUpdatePlaceholders(cols, len(cols))

	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		t.TableName(),
		strings.Join(cols.Names()[1:].Strings(), ","),
		p.placeholders(2, len(cols)),
		setPlaceholders,
	), append(cols.Fields()[1:], cols.Fields()[1:]...)
}

func (p Builder) upsertPostgres(t Table) (string, []interface{}) {

	cols := p.NamedFields(t)

	setPlaceholders := p.generateUpdatePlaceholders(cols, len(cols))

	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (%s) DO UPDATE SET %s",
		t.TableName(),
		strings.Join(cols.Names()[1:].Strings(), ","),
		p.placeholders(1, len(cols)),
		cols[0].Name,
		setPlaceholders,
	), append(cols.Fields()[1:], cols.Fields()[1:]...)
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
