# Partoo

Very very very simple query builder for select, insert and update commands.
Just generates sql.

If you want to do something else, write the query by hand you degenerate!

## Usage

### Assumptions

- Your ID column is the first column in the Columns list
- You always update every field
- You always select every field

Implement the `Table` interface on your model type

```go
package mine
import (
    "github.com/byrnedo/partoo"
)

type MyModel struct {
    ID  string
    Foo string
}

func (t MyModel) TableName() string {
    return "some_table"
}

func (t *MyModel) Columns() Cols {
    return Cols{
        {"id", &t.ID},
        {"foo", &t.Foo},
    }
}
```

Then use partoo to build some queries:

```go
m := &MyModel{}
sqlStr, args := partoo.Insert(m)
// Return corresponds to:
// `INSERT INTO some_table ( foo ) VALUES ( $1 )`, []interface{}{&m.Foo}
```
