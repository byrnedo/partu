# Partu

Opinionated and simple query builder for select, insert and update commands.
Just generates sql.

Supports Mysql and Postgres

```go
m := &MyModel{}
p := partu.New(partu.Postgres)
sqlStr, args := p.Insert(m)
// Return corresponds to:
//
//  "INSERT INTO some_table (id, foo)
//   VALUES ($1,$2)",
//  []interface{}{&m.ID, &m.Foo}

// Or if AutoID isn't fulfilled
//
// "INSERT INTO some_table (foo)
//  VALUES ($1)",
// []interface{}{&m.Foo}
```


Implement the `Table` interface on your model type
```go
package mine
import (
    "github.com/byrnedo/partu"
)

// Note you must have the tag right now, there is no default, but you can override it with `SetTag`
type MyModel struct {
    ID  string `sql:"id"`
    Foo string `sql:"foo"`
}

func (t MyModel) TableName() string {
    return "some_table"
}

func (t *MyModel) Columns() partu.Cols {
    return partu.Cols{
        &t.ID,
        &t.Foo,
    }
}

// OPTIONAL: if you want to manually create your ids, return false
func (t *MyModel) AutoID() bool {
    return false
}
```



### Assumptions

- Your ID column is the first column in the Columns list
- You always update every field except the ID 
- You always select every field


### Available SQL generating methods:

- `Select(t Table)`     
    - `SELECT [cols]`
- `SelectFrom(t Table)`
    - `SELECT [cols] FROM [table]`
- `SelectOne(t Table)`
    - `SELECT [cols] FROM [table] WHERE [id] = [placeholder]`
- `Insert(t Table)`
    - `INSERT INTO [table] ([cols]) VALUES ([placeholders])`
- `Update(t Table)`
    - `UPDATE [table] SET [cols = placeholders]`
- `UpdateOne(t Table)`
    - `UPDATE [table] SET [cols = placeholders] WHERE [id] = [placeholder]`
- `UpsertOne(t Table)`
    - `INSERT INTO [table] ([cols]) VALUES ([placeholders]) ON CONFLICT ([id]) UPDATE SET [cols = placeholders]`

