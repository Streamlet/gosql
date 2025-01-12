package gosql

import (
	"database/sql"
)

func (c *Connection) exec(sql string, bind ...interface{}) (r sql.Result, err error) {
	if c.tx != nil {
		r, err = c.tx.Exec(sql, bind...)
	} else {
		r, err = c.db.Exec(sql, bind...)
	}
	return
}

func (c *Connection) Insert(sql string, bind ...interface{}) (int64, error) {
	r, err := c.exec(sql, bind...)
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

func (c *Connection) Update(sql string, bind ...interface{}) (int64, error) {
	r, err := c.exec(sql, bind...)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}
