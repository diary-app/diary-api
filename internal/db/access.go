package db

import (
	"context"
	"diary-api/internal/auth"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
)

func CheckMyWriteAccessToDiary(ctx context.Context, tx TxOrDb, diaryID uuid.UUID) error {
	userID := auth.MustGetUserID(ctx)
	const checkAccessQuery = `SELECT EXISTS(
    	SELECT * FROM diaries
		WHERE id = $1 AND owner_id = $2)`

	var hasAccess bool
	if err := tx.QueryRowxContext(ctx, checkAccessQuery, diaryID, userID).Scan(&hasAccess); err != nil {
		return err
	}
	if !hasAccess {
		return &usecase.NoWriteAccessToDiaryError{DiaryID: diaryID}
	}
	return nil
}

func CheckMyWriteAccessToEntry(ctx context.Context, tx TxOrDb, entryID uuid.UUID) error {
	userID := auth.MustGetUserID(ctx)
	const checkAccessQuery = `SELECT EXISTS(
    	SELECT * FROM diary_entries e JOIN diaries d ON e.diary_id = d.id 
		WHERE e.id = $1 AND d.owner_id = $2)`
	var hasAccess bool
	if err := tx.QueryRowxContext(ctx, checkAccessQuery, entryID, userID).Scan(&hasAccess); err != nil {
		return err
	}
	if !hasAccess {
		return &usecase.NoWriteAccessToDiaryEntryError{EntryID: entryID}
	}

	return nil
}

func CheckMyReadAccessToEntry(ctx context.Context, tx TxOrDb, entryID uuid.UUID) error {
	userID := auth.MustGetUserID(ctx)
	return CheckUserReadAccessToEntry(ctx, tx, entryID, userID)
}

func CheckUserReadAccessToEntry(ctx context.Context, tx TxOrDb, entryID uuid.UUID, userID uuid.UUID) error {
	const checkAccessQuery = `SELECT EXISTS(
    	SELECT * FROM diary_entries e JOIN diaries d ON e.diary_id = d.id JOIN diary_keys dk ON d.id = dk.diary_id 
		WHERE e.id = $1 AND dk.user_id = $2)`
	var hasAccess bool
	if err := tx.QueryRowxContext(ctx, checkAccessQuery, entryID, userID).Scan(&hasAccess); err != nil {
		return err
	}
	if !hasAccess {
		return &usecase.NoReadAccessToDiaryEntryError{EntryID: entryID}
	}

	return nil
}
