package gosql

import (
	"database/sql"
)

// Connection can not be safely shared between goroutines (if a transaction is in progress)
// Connection.Clone() should be called if it is shared to other goroutines
type Connection struct {
	db *sql.DB
	tx *sql.Tx
}

func Connect(driverName, dataSourceName string) (*Connection, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Connection{db: db}, nil
}

func (c *Connection) Close() error {
	if c.tx != nil {
		_ = c.tx.Rollback()
	}
	return c.db.Close()
}

func (c *Connection) Clone() *Connection {
	return &Connection{c.db, nil}
}
