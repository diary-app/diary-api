package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"testing"
)

func Test_postgresDiaryRepository_GetDiariesByUser(t *testing.T) {
	userID := uuid.New()
	diaryIDs := make([]uuid.UUID, 3)
	for i := range diaryIDs {
		diaryIDs[i] = uuid.New()
	}

	query, args, err := sqlx.In("SELECT * FROM diary_keys WHERE user_id = ? AND diary_id IN (?)", userID, diaryIDs)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Print(query)
	fmt.Print(args)
}
