package db

import (
	"context"
	"diary-api/internal/config"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

func InitDb(cfg config.DbConfig) (*sqlx.DB, error) {
	var connStr string
	if os.Getenv("APP_ENV") == "production" {
		connStr = os.Getenv("DATABASE_URL")
	} else {
		connStr = fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable",
			cfg.Host, cfg.Port, cfg.DbName, cfg.User, cfg.Password)
	}
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetNamedContext(tx *sqlx.Tx, ctx context.Context, dest interface{}, query string, arg interface{}) error {
	query, args, err := tx.BindNamed(query, arg)
	if err != nil {
		return err
	}
	if err = tx.GetContext(ctx, dest, query, args...); err != nil {
		return err
	}
	return nil
}

func ShouldCommitOrRollback(tx Tx) error {
	if commitErr := tx.Commit(); commitErr != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return multierror.Append(commitErr, rbErr)
		}
		return commitErr
	}
	return nil
}
