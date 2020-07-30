package partoo_test

import (
	"github.com/byrnedo/partoo"
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

func (t *baseModel) Columns() partoo.Cols {
	return partoo.Cols{
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
	p := partoo.New(partoo.Postgres)
	cols := p.NamedFields(m)
	aliased := cols.Names().Prefix("alias")
	if reflect.DeepEqual(aliased, partoo.ColNames{"alias.id", "alias.foo", "alias.pq_array", "alias.time"}) == false {
		t.Fatal("wrong aliases", aliased)
	}
}

func TestInsert(t *testing.T) {
	m := &baseModel{}
	p := partoo.New(partoo.Postgres)

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
	p := partoo.New(partoo.Postgres)

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
	p := partoo.New(partoo.Postgres)

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
	p := partoo.New(partoo.Postgres)
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
	p := partoo.New(partoo.Postgres)

	n := p.NamedFields(m)
	if len(n) != len(m.Columns()) {
		t.Fatal(len(n))
	}
}
