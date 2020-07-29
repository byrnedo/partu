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

// Note you must have the tag right now, there is no default
type MyModel struct {
    ID  string `sql:"id"`
    Foo string `sql:"foo"`
}

func (t MyModel) TableName() string {
    return "some_table"
}

func (t *MyModel) Columns() partoo.Cols {
    return partoo.Cols{
        t.ID,
        &t.Foo,
    }
}
```

Then use partoo to build some queries:

```go
m := &MyModel{}
p := partoo.New(partoo.Postgres)
sqlStr, args := p.Insert(m)
// Return corresponds to:
// `INSERT INTO some_table ( foo ) VALUES ( $1 )`, []interface{}{&m.Foo}
```
