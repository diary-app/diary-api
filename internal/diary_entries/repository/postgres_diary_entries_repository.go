package repository

import (
	"context"
	"diary-api/internal/auth"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type pgRepo struct {
	db *sqlx.DB
}

type diaryEntry struct {
	Id      uuid.UUID `db:"id"`
	DiaryId uuid.UUID `db:"diary_id"`
	Name    string    `db:"name"`
	Date    time.Time `db:"date"`
}

type diaryEntryContent struct {
	DiaryEntryId uuid.UUID   `db:"diary_entry_id"`
	Value        interface{} `db:"value"`
}

type diaryEntryWithContent struct {
	diaryEntry
	Contents []diaryEntryContent
}

func (p *pgRepo) GetEntries(ctx context.Context, r usecase.GetDiaryEntriesParams) ([]usecase.DiaryEntry, error) {
	userId := auth.MustGetUserId(ctx)
	namedArgs := map[string]interface{}{"user_id": userId}
	query := `
		SELECT * FROM diary_entries 
        WHERE diary_id IN (
        	SELECT id FROM diaries d JOIN diary_keys dk ON d.id = dk.diary_id WHERE dk.user_id = :user_id)`
	if r.DiaryId != nil {
		query += `AND WHERE diary_id = :diary_id`
		namedArgs["diary_id"] = r.DiaryId
	}
	if r.Date != nil {
		query += `AND WHERE date = :date`
		namedArgs["date"] = r.Date
	}

	query, args, err := p.db.BindNamed(query, namedArgs)
	if err != nil {
		return nil, err
	}
	var entries []diaryEntry
	if err = p.db.SelectContext(ctx, &entries, query, args...); err != nil {
		return nil, err
	}

	ucEntries := make([]usecase.DiaryEntry, len(entries))
	for i, e := range entries {
		ucEntries[i] = mapDiaryEntryToUc(e)
	}
	return ucEntries, nil
}

func (p *pgRepo) GetById(ctx context.Context, id uuid.UUID) (*usecase.DiaryEntry, error) {
	const entryQuery = `SELECT id, diary_id, name, date FROM diary_entries WHERE id = $1`
	entry := &diaryEntryWithContent{}
	if err := p.db.GetContext(ctx, entry, entryQuery, id); err != nil {
		return nil, err
	}

	userId := auth.MustGetUserId(ctx)
	const checkAccessQuery = `SELECT EXISTS(
    	SELECT * FROM diaries d JOIN diary_keys dk ON d.id = dk.diary_id 
		WHERE d.id = $1 AND dk.user_id = $2)`
	var hasAccess bool
	if err := p.db.QueryRowxContext(ctx, checkAccessQuery, entry.DiaryId, userId).Scan(&hasAccess); err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, usecase.ErrNoAccessToDiary
	}

	const contentsQuery = `SELECT diary_entry_id, value FROM diary_entry_contents WHERE diary_entry_id = $1`
	if err := p.db.SelectContext(ctx, entry.Contents, contentsQuery, id); err != nil {
		return nil, err
	}

	result := &usecase.DiaryEntry{
		Id:       entry.Id,
		DiaryId:  entry.DiaryId,
		Name:     entry.Name,
		Date:     entry.Date,
		Contents: make([]interface{}, len(entry.Contents)),
	}
	for i, c := range entry.Contents {
		result.Contents[i] = c.Value
	}
	return result, nil
}

func (p *pgRepo) UpdateContents(ctx context.Context, contentsChanges usecase.DiaryEntryContentsChangeList) {
	//TODO implement me
	panic("implement me")
}

func (p *pgRepo) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func New(db *sqlx.DB) usecase.DiaryEntriesRepository {
	return &pgRepo{
		db: db,
	}
}

func (p *pgRepo) Create(ctx context.Context, entry *usecase.DiaryEntry) (*usecase.DiaryEntry, error) {
	const query = `INSERT INTO diary_entries (id, diary_id, name, date) VALUES (:id, :diary_id, :name, :date)`
	if _, err := p.db.NamedExecContext(ctx, query, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func mapDiaryEntryToUc(e diaryEntry) usecase.DiaryEntry {
	return usecase.DiaryEntry{
		Id:       e.Id,
		DiaryId:  e.DiaryId,
		Name:     e.Name,
		Date:     e.Date,
		Contents: nil,
	}
}
