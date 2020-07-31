package partu_test

import (
	"github.com/byrnedo/partu"
	"github.com/lib/pq"
	"reflect"
	"testing"
	"time"
)

type baseModel struct {
	ID      string                `sql:"id"`
	Foo     string                `sql:"foo,somescrap"`
	PQArray []string              `sql:"pq_array"`
	Time    time.Time             `sql:"time"`
	Omitted *struct{ Foo string } `sql:"random_struct"`
}

func (t baseModel) TableName() string {
	return "test"
}

func (t *baseModel) Columns() partu.Cols {
	return partu.Cols{
		&t.ID,
		&t.Foo,
		pq.Array(&t.PQArray),
		&t.Time,
	}
}

type manualIDModel struct {
	baseModel
}

func (m manualIDModel) AutoID() bool {
	return false
}

func TestColNames_Prefix(t *testing.T) {

	m := &baseModel{}
	p := partu.New(partu.Postgres)
	cols := p.NamedFields(m)
	aliased := cols.Names().Prefix("alias")
	if reflect.DeepEqual(aliased, partu.ColNames{"alias.id", "alias.foo", "alias.pq_array", "alias.time"}) == false {
		t.Fatal("wrong aliases", aliased)
	}
}

func TestBuilders(t *testing.T) {

	type test struct {
		f          func() (string, []interface{})
		outputSQL  string
		outputArgs []interface{}
	}

	baseTable := &baseModel{}
	manualIDTable := &manualIDModel{}
	pPostgres := partu.New(partu.Postgres)
	pMysql := partu.New(partu.Mysql)

	tests := []test{
		{
			func() (string, []interface{}) {
				return pPostgres.Select(baseTable), nil
			},
			"SELECT id, foo, pq_array, time", nil,
		},
		{
			func() (string, []interface{}) {
				return pMysql.Select(baseTable), nil
			},
			"SELECT id, foo, pq_array, time", nil,
		},
		{
			func() (string, []interface{}) {
				return pPostgres.SelectFrom(baseTable), nil
			},
			"SELECT id, foo, pq_array, time\nFROM test", nil,
		},
		{
			func() (string, []interface{}) {
				return pMysql.SelectFrom(baseTable), nil
			},
			"SELECT id, foo, pq_array, time\nFROM test", nil,
		},
		{
			func() (string, []interface{}) {
				s, a := pPostgres.SelectOne(baseTable)
				return s, []interface{}{a}
			},
			"SELECT id, foo, pq_array, time\nFROM test\nWHERE id = $1", []interface{}{&baseTable.ID},
		},
		{
			func() (string, []interface{}) {
				s, a := pMysql.SelectOne(baseTable)
				return s, []interface{}{a}
			},
			"SELECT id, foo, pq_array, time\nFROM test\nWHERE id = ?", baseTable.Columns()[0:1],
		},
		{
			func() (string, []interface{}) {
				return pPostgres.Insert(baseTable)
			},
			"INSERT INTO test (foo, pq_array, time)\nVALUES ($1, $2, $3)", baseTable.Columns()[1:],
		},
		{
			func() (string, []interface{}) {
				return pMysql.Insert(baseTable)
			},
			"INSERT INTO test (foo, pq_array, time)\nVALUES (?, ?, ?)", baseTable.Columns()[1:],
		},
		{
			func() (string, []interface{}) {
				return pPostgres.Insert(manualIDTable)
			},
			"INSERT INTO test (id, foo, pq_array, time)\nVALUES ($1, $2, $3, $4)", baseTable.Columns(),
		},
		{
			func() (string, []interface{}) {
				return pMysql.Insert(manualIDTable)
			},
			"INSERT INTO test (id, foo, pq_array, time)\nVALUES (?, ?, ?, ?)", baseTable.Columns(),
		},
		{
			func() (string, []interface{}) {
				return pPostgres.Update(baseTable)
			},
			"UPDATE test\nSET foo = $1, pq_array = $2, time = $3", baseTable.Columns()[1:],
		},
		{
			func() (string, []interface{}) {
				return pMysql.Update(baseTable)
			},
			"UPDATE test\nSET foo = ?, pq_array = ?, time = ?", baseTable.Columns()[1:],
		},
		{
			func() (string, []interface{}) {
				return pPostgres.UpdateOne(baseTable)
			},
			"UPDATE test\nSET foo = $1, pq_array = $2, time = $3\nWHERE id = $4", append(baseTable.Columns()[1:], baseTable.Columns()[0]),
		},
		{
			func() (string, []interface{}) {
				return pMysql.UpdateOne(baseTable)
			},
			"UPDATE test\nSET foo = ?, pq_array = ?, time = ?\nWHERE id = ?", append(baseTable.Columns()[1:], baseTable.Columns()[0]),
		},
		{
			func() (string, []interface{}) {
				return pPostgres.UpsertOne(baseTable)
			},
			"INSERT INTO test (foo, pq_array, time)\nVALUES ($1, $2, $3)\nON CONFLICT (id) DO UPDATE\nSET foo = $4, pq_array = $5, time = $6", append(baseTable.Columns()[1:], baseTable.Columns()[1:]...),
		},
		{
			func() (string, []interface{}) {
				return pMysql.UpsertOne(baseTable)
			},
			"INSERT INTO test (foo, pq_array, time)\nVALUES (?, ?, ?)\nON DUPLICATE KEY UPDATE\nfoo = ?, pq_array = ?, time = ?", append(baseTable.Columns()[1:], baseTable.Columns()[1:]...),
		},
		{
			func() (string, []interface{}) {
				return pPostgres.UpsertOne(manualIDTable)
			},
			"INSERT INTO test (id, foo, pq_array, time)\nVALUES ($1, $2, $3, $4)\nON CONFLICT (id) DO UPDATE\nSET foo = $5, pq_array = $6, time = $7", append(baseTable.Columns(), baseTable.Columns()[1:]...),
		},
		{
			func() (string, []interface{}) {
				return pMysql.UpsertOne(manualIDTable)
			},
			"INSERT INTO test (id, foo, pq_array, time)\nVALUES (?, ?, ?, ?)\nON DUPLICATE KEY UPDATE\nfoo = ?, pq_array = ?, time = ?", append(baseTable.Columns(), baseTable.Columns()[1:]...),
		},

	}

	for i, toRun := range tests {
		now := time.Now()
		sqlStr, args := toRun.f()
		diff := time.Now().Sub(now)
		t.Logf("time for test %d: %s", i, diff)
		if sqlStr != toRun.outputSQL {
			t.Fatalf("expected SQL \n`%s`,\n\tgot SQL \n`%s`", toRun.outputSQL, sqlStr)
		}
		if !reflect.DeepEqual(args, toRun.outputArgs) {
			t.Fatalf("expected ARGS `%s`,\n\tgot ARGS `%s`", toRun.outputArgs, args)
		}
	}
}


func TestNamedFields(t *testing.T) {
	m := &baseModel{}
	p := partu.New(partu.Postgres)

	n := p.NamedFields(m)
	if len(n) != len(m.Columns()) {
		t.Fatal(len(n))
	}
}

func TestBuilder_ColName(t *testing.T) {
	m := &baseModel{}
	p := partu.New(partu.Postgres)
	n := p.ColName(m, &m.Time)
	if n != "time" {
		t.Fatal(n)
	}
}
