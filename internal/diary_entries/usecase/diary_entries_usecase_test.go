package usecase

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		repo storage
	}
	tests := []struct {
		name string
		args args
		want *UseCase
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Create(t *testing.T) {
	type fields struct {
		repo storage
	}
	type args struct {
		ctx context.Context
		r   usecase.CreateDiaryEntryRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *usecase.DiaryEntry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &UseCase{
				repo: tt.fields.repo,
			}
			got, err := d.Create(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Delete(t *testing.T) {
	type fields struct {
		repo storage
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &UseCase{
				repo: tt.fields.repo,
			}
			err := d.Delete(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCase_GetByID(t *testing.T) {
	type fields struct {
		repo storage
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *usecase.DiaryEntry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &UseCase{
				repo: tt.fields.repo,
			}
			got, err := d.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetEntries(t *testing.T) {
	type fields struct {
		repo storage
	}
	type args struct {
		ctx     context.Context
		request usecase.GetDiaryEntriesParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []usecase.DiaryEntry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &UseCase{
				repo: tt.fields.repo,
			}
			got, err := d.GetEntries(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEntries() got = %v, want %v", got, tt.want)
			}
		})
	}
}
