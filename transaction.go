package gosql

import (
	"errors"
)

func (c *Connection) Begin() error {
	if c.tx != nil {
		return errors.New("previous transaction not closed")
	}
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	c.tx = tx
	return nil
}

func (c *Connection) Commit() error {
	if c.tx == nil {
		return errors.New("not in transaction")
	}
	err := c.tx.Commit()
	if err == nil {
		return err
	}
	c.tx = nil
	return nil
}

func (c *Connection) Rollback() error {
	if c.tx == nil {
		return errors.New("not in transaction")
	}
	err := c.tx.Rollback()
	if err == nil {
		return err
	}
	c.tx = nil
	return nil
}

func (c *Connection) End() {
	if c.tx == nil {
		return
	}
	_ = c.Rollback()
}
