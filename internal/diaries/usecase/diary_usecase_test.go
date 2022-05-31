package usecase_test

import (
	"context"
	"diary-api/internal/auth"
	diaryUsecase "diary-api/internal/diaries/usecase"
	"diary-api/internal/usecase"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type getDiariesByUserTestCase struct {
	name       string
	context    func() context.Context
	storage    func(ctrl *gomock.Controller) *Mockstorage
	assertFunc func(t *testing.T, got []usecase.Diary, err error)
}

func TestUseCase_GetDiariesByUser(t *testing.T) {
	existingUserID := uuid.MustParse("c3790d41-3099-4bef-a611-44b772115990")
	ctxWithExistingUserID := context.WithValue(context.Background(), auth.UserIDKey, existingUserID)
	correctDiaries := []usecase.Diary{
		{ID: uuid.New()},
	}
	tests := []getDiariesByUserTestCase{
		{
			name: "extracts user ID from context and returns result from user storage by this ID",
			storage: func(ctrl *gomock.Controller) *Mockstorage {
				mock := NewMockstorage(ctrl)
				mock.EXPECT().GetDiariesByUser(gomock.Any(), existingUserID).Return(correctDiaries, nil)
				return mock
			},
			context: func() context.Context {
				return ctxWithExistingUserID
			},
			assertFunc: func(t *testing.T, gotDiaries []usecase.Diary, err error) {
				assert.Equal(t, gotDiaries, correctDiaries)
				assert.Nil(t, err)
			},
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := diaryUsecase.New(tc.storage(ctrl))
			d, err := u.GetDiariesByUser(tc.context())
			tc.assertFunc(t, d, err)
		})
	}
}
