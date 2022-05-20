package db

import (
	"diary-api/internal/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDb(cfg config.DbConfig) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable",
		cfg.Host, cfg.Port, cfg.DbName, cfg.User, cfg.Password)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}