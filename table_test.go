package partoo

import (
	"reflect"
	"testing"
)

type testModel struct {
	ID  string
	Foo string
}

func (t testModel) TableName() string {
	return "test"
}

func (t *testModel) Columns() Cols {
	return Cols{
		{"id", &t.ID},
		{"foo", &t.Foo},
	}
}

func TestColNames_Prefix(t *testing.T) {

	m := &testModel{}
	cols := m.Columns()
	aliased := cols.Names().Prefix("alias")

	if reflect.DeepEqual(aliased, ColNames{"alias.id", "alias.foo"}) == false {
		t.Fatal("wrong aliases", aliased)
	}
}

func TestInsert(t *testing.T) {
	m := &testModel{}

	sqlStr, args := Insert(m)
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

	sqlStr, args := Update(m)
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

	sqlStr, args := UpdateOne(m)
	if sqlStr != "UPDATE test SET foo = $1 WHERE id = $2" {
		t.Fatal(sqlStr)
	}
	if len(args) != 2 {
		t.Fatal(len(args))
	}
	t.Log(sqlStr)
}
