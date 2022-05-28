package diaries_test

import (
	"context"
	"diary-api/internal/diaries"
	"diary-api/internal/usecase"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUseCase_CreateDiary(t *testing.T) {
	tests := []struct {
		name       string
		storage    func(ctrl *gomock.Controller) *Mockstorage
		assertFunc func(t *testing.T, actual *usecase.Diary, err error)
	}{
		{
			name: "returns error",
			storage: func(ctrl *gomock.Controller) *Mockstorage {
				mock := NewMockstorage(ctrl)
				mock.EXPECT().CreateDiary(context.Background(), gomock.Any()).Return(nil, errors.New(""))
				return mock
			},
			assertFunc: func(t *testing.T, actual *usecase.Diary, err error) {
				assert.Nil(t, actual)
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := diaries.NewDiaryUseCase(tc.storage(ctrl))
			d, err := u.CreateDiary(context.Background(), uuid.New(), &usecase.CreateDiaryRequest{})
			tc.assertFunc(t, d, err)
		})
	}
}
