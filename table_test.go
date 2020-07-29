package partoo

import (
	"reflect"
	"testing"
)

type testModel struct {
	ID  string `sql:"id"`
	Foo string `sql:"foo"`
}

func (t testModel) TableName() string {
	return "test"
}

func (t *testModel) Columns() Cols {
	return Cols{
		&t.ID,
		&t.Foo,
	}
}

func TestColNames_Prefix(t *testing.T) {

	m := &testModel{}
	p := New(Postgres)
	cols := p.NamedFields(m)
	aliased := cols.Names().Prefix("alias")
	if reflect.DeepEqual(aliased, ColNames{"alias.id", "alias.foo"}) == false {
		t.Fatal("wrong aliases", aliased)
	}
}

func TestInsert(t *testing.T) {
	m := &testModel{}
	p := New(Postgres)

	sqlStr, args := p.Insert(m)
	if sqlStr != "INSERT INTO test (foo) VALUES ($1)" {
		t.Fatal(sqlStr)
	}
	if len(args) != 1 {
		t.Fatal(len(args))
	}
	t.Log(sqlStr)
}

func TestUpdate(t *testing.T) {
	m := &testModel{}
	p := New(Postgres)

	sqlStr, args := p.Update(m)
	if sqlStr != "UPDATE test SET foo = $1" {
		t.Fatal(sqlStr)
	}
	if len(args) != 1 {
		t.Fatal(len(args))
	}
	t.Log(sqlStr)
}

func TestUpdateOne(t *testing.T) {
	m := &testModel{}
	p := New(Postgres)

	sqlStr, args := p.UpdateOne(m)
	if sqlStr != "UPDATE test SET foo = $1 WHERE id = $2" {
		t.Fatal(sqlStr)
	}
	if len(args) != 2 {
		t.Fatal(len(args))
	}
	t.Log(sqlStr)
}


func TestPartoo_UpsertOne(t *testing.T) {
	m := &testModel{}
	p := New(Postgres)
	sqlStr, args := p.UpsertOne(m)
	if sqlStr != "INSERT INTO test (foo) VALUES ($1) ON CONFLICT (id) DO UPDATE SET foo = $2" {
		t.Fatal(sqlStr)
	}
	if len(args) != 2 {
		t.Fatal(len(args))
	}
}