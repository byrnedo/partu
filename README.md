# Partoo

Very very very simple query builder for select, insert and update commands.
Just generates sql.

If you want to do something else, write the query by hand you degenerate!

## Usage

### Assumptions

- Your ID column is the first column in the Columns list
- You always update every field except the ID
- You always select every field

Implement the `Table` interface on your model type

```go
package mine
import (
    "github.com/byrnedo/partoo"
)

// Note you must have the tag right now, there is no default, but you can override it with `SetTag`
type MyModel struct {
    ID  string `sql:"id"`
    Foo string `sql:"foo"`
}

func (t MyModel) TableName() string {
    return "some_table"
}

func (t *MyModel) Columns() partoo.Cols {
    return partoo.Cols{
        &t.ID,
        &t.Foo,
    }
}

// OPTIONAL: if you want to manually create your ids, return false
func (t *MyModel) AutoID() bool {
    return false
}
```

Then use partoo to build some queries:

```go
m := &MyModel{}
p := partoo.New(partoo.Postgres)
sqlStr, args := p.Insert(m)
// Return corresponds to:
// `INSERT INTO some_table (id,foo) VALUES ($1,$2)`, []interface{}{&m.ID, &m.Foo}
// Or if AutoID isn't fulfilled
// `INSERT INTO some_table (foo) VALUES ($1)`, []interface{}{&m.Foo}
```

##