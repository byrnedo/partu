package partu

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

func (p *Builder) SetTag(tag string) *Builder {
	p.tagName = tag
	return p
}

func (p Builder) Tag() string {
	if p.tagName == "" {
		return "sql"
	}
	return p.tagName
}

func (p Builder) placeholder(i int) string {
	switch p.dialect {
	case Mysql:
		return fmt.Sprintf("?")
	default:
		return fmt.Sprintf("$%d", i)
	}
}

func (p Builder) placeholders(low, high int) string {
	parts := make([]string, high-low)
	for i := low; i < high; i++ {
		parts[i-low] = p.placeholder(i)
	}
	return strings.Join(parts, ", ")
}

func (p Builder) Select(t Table) string {
	cols := p.NamedFields(t)
	return fmt.Sprintf(
		"SELECT %s",
		cols.Names().String(),
	)
}

func (p Builder) SelectFrom(t Table) string {
	return fmt.Sprintf(
		"%s\nFROM %s",
		p.Select(t),
		t.TableName(),
	)
}

func (p Builder) SelectOne(t Table) (string, interface{}) {
	cols := p.NamedFields(t)
	first := cols[0]
	return fmt.Sprintf(
		"%s\nWHERE %s = %s",
		p.SelectFrom(t),
		first.Name,
		p.placeholder(1),
	), first.Field
}

func (p Builder) Insert(t Table) (string, []interface{}) {
	cols := p.NamedFields(t)

	colsToInsert := cols.Names()[1:]
	args := cols.Fields()[1:]
	autoID, ok := t.(AutoID)
	if ok && !autoID.AutoID() {
		colsToInsert = cols.Names()
		args = cols.Fields()
	}

	return fmt.Sprintf(
		"INSERT INTO %s (%s)\nVALUES (%s)",
		t.TableName(),
		colsToInsert.String(),
		p.placeholders(1, len(colsToInsert)+1),
	), args
}

func (p Builder) Update(t Table) (string, []interface{}) {
	namedFields := p.NamedFields(t)

	setPlaceholders := p.AssignmentString(namedFields, 1)

	return fmt.Sprintf(
		"UPDATE %s\nSET %s",
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
	return fmt.Sprintf("%s\nWHERE %s = %s", upd, names[0], p.placeholder(len(fields))), args
}

func (p Builder) UpsertOne(t Table) (string, []interface{}) {

	cols := p.NamedFields(t)
	args := append(cols.Fields()[1:], cols.Fields()[1:]...)
	colsToInsert := cols.Names()[1:]
	autoID, ok := t.(AutoID)
	if ok && !autoID.AutoID() {
		colsToInsert = cols.Names()
		args = append(cols.Fields(), cols.Fields()[1:]...)
	}
	insertP := p.placeholders(1, len(colsToInsert) + 1)
	updateP := p.AssignmentString(cols, len(colsToInsert)+1)
	iNames := colsToInsert.String()

	if p.dialect == Mysql {
		return p.upsertMysql(t, cols, iNames, insertP, updateP, args)
	}
	return p.upsertPostgres(t, cols, iNames, insertP, updateP, args)
}

func (p Builder) AssignmentString(cols namedFields, startIndex int) string {

	names := cols.Names()[1:]
	setPlaceholders := make([]string, len(names))
	for i, n := range names {
		setPlaceholders[i] = fmt.Sprintf("%s = %s", n, p.placeholder(i+startIndex))
	}

	return strings.Join(setPlaceholders, ", ")
}

func (p Builder) upsertMysql(t Table, cols namedFields,iNames, iPlaceholders string, uPlaceholders string, args []interface{}) (string, []interface{}) {

	return fmt.Sprintf(
		"INSERT INTO %s (%s)\nVALUES (%s)\nON DUPLICATE KEY UPDATE\n%s",
		t.TableName(),
		iNames,
		iPlaceholders,
		uPlaceholders,
	), args
}

func (p Builder) upsertPostgres(t Table, cols namedFields, iNames, iPlaceholders string, uPlaceholders string, args []interface{}) (string, []interface{}) {

	return fmt.Sprintf(
		"INSERT INTO %s (%s)\nVALUES (%s)\nON CONFLICT (%s) DO UPDATE\nSET %s",
		t.TableName(),
		iNames,
		iPlaceholders,
		cols[0].Name,
		uPlaceholders,
	), args
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

func (c ColNames) String() string {
	return strings.Join(c, ", ")
}

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
