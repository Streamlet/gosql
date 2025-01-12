# gosql

A simple SQL helper for golang.

## Connect

```go
// Similar with sql.Open(driverName, dataSourceName string)
connection, err := gosql.Connect(driverName, dataSourceName string)
```

## Transaction

```go
err := connection.Begin()
defer connection.End()

// Do work with transaction
rows, err := connection.Update(...)

// Save return when error occurs
if err != nil {
return false
}

connection.Commit()

// defer connection.End() takes no effect after Commit/Rollback
return true
```

Re-enter guard:

```go
func outer() {
err := connection.Begin()
defer connection.End()
// ...
inner()
}

func inner(connection *Connection) {
err := connection.Begin() // will get an error, it is by design
// ...
}
```

We strongly disagree with the view that DB interface should be designed to be transaction-insensitive.

On the contrary, we believe that everyone should clearly know whether the current code will be executed within or
outside the transaction when writing the code.
Only in this way can the scope of the transaction be well controlled.

# Update and insert

```go
rows, err := connection.Update("UPDATE ...", ...)
// rows will be the number of affected rows
```

```go
id, err := connection.Update("INSERT ...", ...)
// id will be last insert id
```

But in fact, we won't parse SQL, so we can't identify whether the SQL statement is update or insert.
If using the Update method to execute the insert statement, it is feasible, but you will only get the number of affected
rows, not the last insert id.
Using the Insert method to execute the update statement is also possible, but you can't get the number of affected rows.

# Select

Ordinary select:

```go
rs, err := connection.Select("SELECT * FROM table", ...)
// rs will be sql.Rows
```

Structured select:

```go
// Define a struct that represent a record
type Item struct {
    Id      int64     `db:"id"`
    Name    string    `db:"name"` // if non-pointer member received null value, it will be set to zero value for the type
    Value   *string   `db:"value"` // if pointer member receives null value, it will be set to nil
    Time    time.Time `db:"time"`
}
items, err := connection.Select[Item]("SELECT id, name, value, time FROM table", ...)
// items will be slice of Item
```
