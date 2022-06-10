package db

import (
	"context"
	"diary-api/internal/auth"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
)

func CheckMyAccessToDiaryContext(ctx context.Context, tx TxOrDb, diaryID uuid.UUID) error {
	userID := auth.MustGetUserID(ctx)
	return CheckUserAccessToDiary(ctx, tx, diaryID, userID)
}

func CheckUserAccessToDiary(ctx context.Context, tx TxOrDb, diaryID uuid.UUID, userID uuid.UUID) error {
	const checkAccessQuery = `SELECT EXISTS(
    	SELECT * FROM diaries d JOIN diary_keys dk ON d.id = dk.diary_id 
		WHERE d.id = $1 AND dk.user_id = $2)`

	var hasAccess bool
	if err := tx.QueryRowxContext(ctx, checkAccessQuery, diaryID, userID).Scan(&hasAccess); err != nil {
		return err
	}
	if !hasAccess {
		return &usecase.NoAccessToDiaryError{DiaryID: diaryID}
	}
	return nil
}

func CheckMyAccessToEntry(ctx context.Context, tx TxOrDb, entryID uuid.UUID) error {
	userID := auth.MustGetUserID(ctx)
	return CheckUserAccessToEntry(ctx, tx, entryID, userID)
}

func CheckUserAccessToEntry(ctx context.Context, tx TxOrDb, entryID uuid.UUID, userID uuid.UUID) error {
	const checkAccessQuery = `SELECT EXISTS(
    	SELECT * FROM diary_entries e JOIN diaries d ON e.diary_id = d.id JOIN diary_keys dk ON d.id = dk.diary_id 
		WHERE e.id = $1 AND dk.user_id = $2)`
	var hasAccess bool
	if err := tx.QueryRowxContext(ctx, checkAccessQuery, entryID, userID).Scan(&hasAccess); err != nil {
		return err
	}
	if !hasAccess {
		return &usecase.NoAccessToDiaryEntryError{EntryID: entryID}
	}

	return nil
}
