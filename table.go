package partoo

//type Dialect string
//
//const (
//	Mysql Dialect = "mysql"
//	Postgres Dialect = "postgres"
//)
//
//var dialect Dialect
//
//func SetDialect(d Dialect) {
//	dialect = d
//}
//

type Table interface {
	// The name of the ... table
	TableName() string
	// The columns mapped to the fields in your struct, by ref please
	Columns() map[string]interface{}
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

func ColumnNames(t Table) (ret []string) {
	idx := 0
	cols := t.Columns()
	ret = make([]string, len(cols))
	for k, _ := range cols {
		ret[idx] = k
		idx ++
	}
	return
}

func ColumnFields(t Table) (ret []interface{}) {

	idx := 0
	cols := t.Columns()
	ret = make([]interface{}, len(cols))
	for _, v := range cols {
		ret[idx] = v
		idx ++
	}
	return
}
