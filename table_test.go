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

func TestSelect(t *testing.T) {

	m := &baseModel{}
	p := partu.New(partu.Postgres)

	sqlStr := p.Select(m)
	if sqlStr != "SELECT id,foo,pq_array,time" {
		t.Fatal(sqlStr)
	}

	sqlStr = p.SelectFrom(m)
	if sqlStr != "SELECT id,foo,pq_array,time FROM test" {
		t.Fatal(sqlStr)
	}

	sqlStr, _ = p.SelectOne(m)
	if sqlStr != "SELECT id,foo,pq_array,time FROM test WHERE id = $1" {
		t.Fatal(sqlStr)
	}
}

func TestInsert(t *testing.T) {
	m := &baseModel{}
	p := partu.New(partu.Postgres)

	sqlStr, args := p.Insert(m)
	if sqlStr != "INSERT INTO test (foo,pq_array,time) VALUES ($1,$2,$3)" {
		t.Fatal(sqlStr)
	}
	if len(args) != 3 {
		t.Fatal(len(args))
	}

	manualModel := &manualIDModel{}
	sqlStr, args = p.Insert(manualModel)
	if sqlStr != "INSERT INTO test (id,foo,pq_array,time) VALUES ($1,$2,$3,$4)" {
		t.Fatal(sqlStr)
	}
	if len(args) != 4 {
		t.Fatal(len(args))
	}
}

func TestUpdate(t *testing.T) {
	m := &baseModel{}
	p := partu.New(partu.Postgres)

	sqlStr, args := p.Update(m)
	if sqlStr != "UPDATE test SET foo = $1,pq_array = $2,time = $3" {
		t.Fatal(sqlStr)
	}
	if len(args) != 3 {
		t.Fatal(len(args))
	}
	t.Log(sqlStr)
}

func TestUpdateOne(t *testing.T) {
	m := &baseModel{}
	p := partu.New(partu.Postgres)

	sqlStr, args := p.UpdateOne(m)
	if sqlStr != "UPDATE test SET foo = $1,pq_array = $2,time = $3 WHERE id = $4" {
		t.Fatal(sqlStr)
	}
	if len(args) != 4 {
		t.Fatal(len(args))
	}
	t.Log(sqlStr)
}

func TestPartoo_UpsertOne(t *testing.T) {
	m := &baseModel{}
	p := partu.New(partu.Postgres)
	sqlStr, args := p.UpsertOne(m)
	if sqlStr != "INSERT INTO test (foo,pq_array,time) VALUES ($1,$2,$3) ON CONFLICT (id) DO UPDATE SET foo = $4,pq_array = $5,time = $6" {
		t.Fatal(sqlStr)
	}
	if len(args) != 6 {
		t.Fatal(len(args))
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
	if n != "time"{
		t.Fatal(n)
	}
}
