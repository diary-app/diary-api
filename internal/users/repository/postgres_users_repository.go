package repository

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type postgresUsersRepository struct {
	db *sqlx.DB
}

func (p *postgresUsersRepository) CreateUser(ctx context.Context, user *usecase.FullUser) (*usecase.FullUser, error) {
	const insertQuery = `
INSERT INTO users(username, password_hash, salt_for_keys, public_key_for_sharing, encrypted_private_key_for_sharing) 
VALUES(:username,:password_hash,:salt_for_keys,:public_key_for_sharing,:encrypted_private_key_for_sharing) 
RETURNING id`
	query, args, err := p.db.BindNamed(insertQuery, user)
	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	err = p.db.GetContext(ctx, &id, query, args...)
	if err != nil {
		return nil, err
	}

	user.Id = id
	return user, nil
}

func (p *postgresUsersRepository) GetUserById(ctx context.Context, id uuid.UUID) (*usecase.FullUser, error) {
	const query = `SELECT * FROM users WHERE id = $1`

	user := &usecase.FullUser{}
	if err := p.db.GetContext(ctx, user, query, id); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *postgresUsersRepository) GetUserByName(ctx context.Context, username string) (*usecase.FullUser, error) {
	const query = `SELECT * FROM users WHERE username = $1`
	user := &usecase.FullUser{}
	if err := p.db.GetContext(ctx, user, query, username); err != nil {
		return nil, err
	}
	return user, nil
}

func New(db *sqlx.DB) usecase.UsersRepository {
	return &postgresUsersRepository{
		db: db,
	}
}
