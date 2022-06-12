package repository

import (
	"context"
	"diary-api/internal/auth"
	"diary-api/internal/db"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type pgRepo struct {
	db db.Db
}

type diaryEntry struct {
	ID      uuid.UUID `db:"id"`
	DiaryID uuid.UUID `db:"diary_id"`
	Name    string    `db:"name"`
	Date    time.Time `db:"date"`
	Value   []byte    `db:"value"`
}

type diaryEntryBlock struct {
	ID           uuid.UUID `db:"id"`
	DiaryEntryID uuid.UUID `db:"diary_entry_id"`
	Value        []byte    `db:"value"`
}

func (p *pgRepo) GetEntries(ctx context.Context, params usecase.GetDiaryEntriesParams) ([]usecase.DiaryEntry, error) {
	userID := auth.MustGetUserID(ctx)
	namedArgs := map[string]interface{}{"user_id": userID}
	query := `
		SELECT * FROM diary_entries 
        WHERE diary_id IN (
        	SELECT id FROM diaries d JOIN diary_keys dk ON d.id = dk.diary_id WHERE dk.user_id = :user_id)`
	if params.DiaryID != nil {
		query += ` AND diary_id = :diary_id`
		namedArgs["diary_id"] = *params.DiaryID
	}
	if params.Date != nil {
		query += ` AND date = :date`
		namedArgs["date"] = time.Time(*params.Date)
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
		ucEntries[i] = mapDiaryEntry(e)
	}
	return ucEntries, nil
}

func (p *pgRepo) GetByID(ctx context.Context, id uuid.UUID) (*usecase.DiaryEntry, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}
	err = db.CheckMyAccessToEntry(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	entry, err := getEntry(ctx, id, tx)
	if err != nil {
		return nil, err
	}

	const blocksQuery = `SELECT id, diary_entry_id, value FROM diary_entry_blocks WHERE diary_entry_id = $1`
	var blocks []diaryEntryBlock
	if err := p.db.SelectContext(ctx, &blocks, blocksQuery, id); err != nil {
		return nil, err
	}

	result := mapToUseCaseDiaryEntry(entry, blocks)
	return result, nil
}

func mapToUseCaseDiaryEntry(entry *diaryEntry, blocks []diaryEntryBlock) *usecase.DiaryEntry {
	result := &usecase.DiaryEntry{
		ID:      entry.ID,
		DiaryID: entry.DiaryID,
		Name:    entry.Name,
		Date:    entry.Date,
		Value:   entry.Value,
		Blocks:  make([]usecase.DiaryEntryBlock, len(blocks)),
	}
	for i, b := range blocks {
		result.Blocks[i] = usecase.DiaryEntryBlock{
			ID:    b.ID,
			Value: b.Value,
		}
	}
	return result
}

func getEntry(ctx context.Context, id uuid.UUID, tx db.Tx) (*diaryEntry, error) {
	const entryQuery = `SELECT id, diary_id, name, date, value FROM diary_entries WHERE id = $1`
	entry := &diaryEntry{}
	if err := tx.GetContext(ctx, entry, entryQuery, id); err != nil {
		return nil, err
	}
	return entry, nil
}

func getBlocksIds(blocks []usecase.DiaryEntryBlockDto) []uuid.UUID {
	res := make([]uuid.UUID, len(blocks))
	for i, b := range blocks {
		res[i] = b.ID
	}
	return res
}

func (p *pgRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}
	if err := db.CheckMyAccessToEntry(ctx, tx, id); err != nil {
		return err
	}
	const query = `DELETE FROM diary_entries WHERE id = $1`
	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return db.ShouldCommitOrRollback(tx)
}

func New(db *sqlx.DB) usecase.DiaryEntriesRepository {
	return &pgRepo{
		db: db,
	}
}

func (p *pgRepo) Create(ctx context.Context, entry *usecase.DiaryEntry) (*usecase.DiaryEntry, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}
	if err = db.CheckMyAccessToDiary(ctx, tx, entry.DiaryID); err != nil {
		return nil, err
	}

	const insertEntryQuery = `
INSERT INTO diary_entries (id, diary_id, name, date, value) VALUES (:id, :diary_id, :name, :date, :value)`

	if _, err = tx.NamedExecContext(ctx, insertEntryQuery, entry); err != nil {
		return nil, err
	}

	if err = db.ShouldCommitOrRollback(tx); err != nil {
		return nil, err
	}
	return entry, nil
}

func mapDiaryEntry(e diaryEntry) usecase.DiaryEntry {
	return usecase.DiaryEntry{
		ID:      e.ID,
		DiaryID: e.DiaryID,
		Name:    e.Name,
		Date:    e.Date,
		Blocks:  nil,
	}
}
