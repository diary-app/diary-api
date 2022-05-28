package db

import (
	"diary-api/internal/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

func InitDb(cfg config.DbConfig) (*sqlx.DB, error) {
	var connStr string
	if os.Getenv("APP_ENV") == "production" {
		connStr = os.Getenv("DATABASE_URL")
	} else {
		connStr = fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=%v",
			cfg.Host, cfg.Port, cfg.DbName, cfg.User, cfg.Password, cfg.SslMode)
	}
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
