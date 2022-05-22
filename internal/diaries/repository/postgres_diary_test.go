package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"testing"
)

func Test_postgresDiaryRepository_GetDiariesByUser(t *testing.T) {
	userId := uuid.New()
	diaryIds := make([]uuid.UUID, 3)
	for i := range diaryIds {
		diaryIds[i] = uuid.New()
	}

	query, args, err := sqlx.In("SELECT * FROM diary_keys WHERE user_id = ? AND diary_id IN (?)", userId, diaryIds)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Print(query)
	fmt.Print(args)
}
