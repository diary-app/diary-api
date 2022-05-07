package users

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
}

type Repo interface {
	GetUser(username string)
}

func NewRepo(conn *sqlx.Conn) Repo {
	return &postgresRepo{
		conn: conn,
	}
}

type postgresRepo struct {
	conn *sqlx.Conn
}

func (p *postgresRepo) GetUser(username string) {
	return func(ctx *gin.Context) {
		//TODO implement me
		panic("implement me")
	}
}
