package db

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type TxOrDb interface {
	Rebind(query string) string
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Get(dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	Preparex(query string) (*sqlx.Stmt, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Prepare(query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRow(query string, args ...interface{}) *sql.Row
}

type Db interface {
	TxOrDb
	MustBegin() *sqlx.Tx
	Beginx() (*sqlx.Tx, error)
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
}

type Tx interface {
	TxOrDb
	Commit() error
	Rollback() error
	NamedStmt(stmt *sqlx.NamedStmt) *sqlx.NamedStmt
	StmtxContext(ctx context.Context, stmt interface{}) *sqlx.Stmt
	NamedStmtContext(ctx context.Context, stmt *sqlx.NamedStmt) *sqlx.NamedStmt
	StmtContext(ctx context.Context, stmt *sql.Stmt) *sql.Stmt
	Stmt(stmt *sql.Stmt) *sql.Stmt
}
