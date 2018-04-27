package db

import (
	"context"
	"database/sql"

	//pq is the sql driver for database/sql
	_ "github.com/lib/pq"
)

type contextKey int

var contextKeyTx = contextKey(0)

//NewPostgresDB gives a new instance of DB
func NewPostgresDB(sqlURL string, noOfConns int) (*DB, error) {
	db, err := sql.Open("postgres", sqlURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(noOfConns)
	db.SetMaxOpenConns(noOfConns)

	return &DB{db}, nil
}

//TransactionalStore gives the signature for models to be able to
//expose this method
type TransactionalStore interface {
	WithTxContext(ctx context.Context, f func(context.Context) error) error
}

//DB is a light holder of *sql.DB
type DB struct {
	*sql.DB
}

//WithTxContext wraps the contextcallers for ease of use
func (t *DB) WithTxContext(ctx context.Context, f func(context.Context) error) error {

	tx, err := t.Begin()
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, contextKeyTx, tx)
	err = f(ctx)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

//WithTx executes f within a transaction
func (t *DB) WithTx(f func(tx *sql.Tx) error) error {
	tx, err := t.Begin()
	if err != nil {
		return err
	}
	err = f(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

//ContextDB gives definitions common between sql.DB and sql.Tx
type ContextDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

//GetContextDB gets the transaction in the context or returns a new one
func (t *DB) GetContextDB(ctx context.Context) (ContextDB, error) {
	txValue := ctx.Value(contextKeyTx)
	switch contextTx := txValue.(type) {
	case *sql.Tx:
		return contextTx, nil
	default:
		return t.DB, nil
	}
}
