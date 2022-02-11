package store

import (
	"database/sql"
	"fmt"
)

type sqlEngine struct {
	db *sql.DB
}

type sqlTxn struct {
	*sql.Tx
	db *sql.DB
}

type sqlSnapshot struct {
}

func newSQLEngine(dir string) (Engine, error) {
	return &sqlEngine{}, nil
}

func (e *sqlEngine) Name() string {
	return "sql-engine"
}

func (e *sqlEngine) Txn(fn func(Txn) error) error {
	return nil
}

func (e *sqlEngine) View(fn func(Txn) error) error {
	return e.Txn(fn)
}

func (e *sqlEngine) Update(fn func(Txn) error) error {
	return e.Txn(fn)
}

func (e *sqlEngine) Snapshot() (Snapshot, error) {
	return newSQLSnapshot(), nil
}

func (e *sqlEngine) Close() error {
	return nil
}

func newSQLSnapshot() *sqlSnapshot {
	return &sqlSnapshot{}
}

func (s *sqlSnapshot) Name() string {
	return "sql-snapshot"
}

func (t *sqlTxn) Get(key interface{}) (interface{}, error) {
	sql := fmt.Sprintf("select * from %s where")
	t.Tx.Query(sql)
}

func getInode() {

}
