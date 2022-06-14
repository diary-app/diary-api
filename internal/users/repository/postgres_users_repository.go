package repository

import (
	"context"
	"database/sql"
	"diary-api/internal/db"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type postgresUsersRepository struct {
	db *sqlx.DB
}

func (p *postgresUsersRepository) CreateUser(
	ctx context.Context, user *usecase.FullUser, diary *usecase.Diary) (*usecase.FullUser, *usecase.Diary, error) {
	const insertUserQuery = `
INSERT INTO users(username, password_hash, salt_for_keys, public_key_for_sharing, encrypted_private_key_for_sharing) 
VALUES(:username,:password_hash,:salt_for_keys,:public_key_for_sharing,:encrypted_private_key_for_sharing) 
RETURNING id`
	const insertDiaryQuery = `INSERT INTO diaries(name, owner_id) VALUES (:name, :owner_id) RETURNING id`
	const insertDiaryKeyQuery = `
INSERT INTO diary_keys(diary_id, user_id, encrypted_key) 
VALUES ($1, $2, $3)`

	tx, err := p.db.Beginx()
	if err != nil {
		return nil, nil, err
	}
	var userID uuid.UUID
	if err = db.GetNamedContext(tx, ctx, &userID, insertUserQuery, user); err != nil {
		return nil, nil, err
	}

	diary.OwnerID = userID
	var diaryID uuid.UUID
	if err = db.GetNamedContext(tx, ctx, &diaryID, insertDiaryQuery, diary); err != nil {
		return nil, nil, err
	}

	keyBytes := diary.Keys[0].EncryptedKey
	if _, err = tx.ExecContext(ctx, insertDiaryKeyQuery, diaryID, userID, keyBytes); err != nil {
		return nil, nil, err
	}

	if err = db.ShouldCommitOrRollback(tx); err != nil {
		return nil, nil, err
	}

	user.ID = userID
	diary.ID = diaryID
	return user, diary, nil
}

func (p *postgresUsersRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*usecase.FullUser, error) {
	const query = `SELECT * FROM users WHERE id = $1`

	user := &usecase.FullUser{}
	if err := p.db.GetContext(ctx, user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, usecase.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (p *postgresUsersRepository) GetUserByName(ctx context.Context, username string) (*usecase.FullUser, error) {
	const query = `SELECT * FROM users WHERE username = $1`
	user := &usecase.FullUser{}
	if err := p.db.GetContext(ctx, user, query, username); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func New(db *sqlx.DB) usecase.UsersRepository {
	return &postgresUsersRepository{
		db: db,
	}
}
