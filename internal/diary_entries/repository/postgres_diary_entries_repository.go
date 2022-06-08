package repository

import (
	"context"
	"database/sql"
	"diary-api/internal/auth"
	"diary-api/internal/db"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
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
	Value   string    `db:"value"`
}

type diaryEntryBlock struct {
	ID           uuid.UUID `db:"id"`
	DiaryEntryID uuid.UUID `db:"diary_entry_id"`
	Value        string    `db:"value"`
}

func (p *pgRepo) GetEntries(ctx context.Context, r usecase.GetDiaryEntriesParams) ([]usecase.DiaryEntry, error) {
	userID := auth.MustGetUserID(ctx)
	namedArgs := map[string]interface{}{"user_id": userID}
	query := `
		SELECT * FROM diary_entries 
        WHERE diary_id IN (
        	SELECT id FROM diaries d JOIN diary_keys dk ON d.id = dk.diary_id WHERE dk.user_id = :user_id)`
	if r.DiaryID != nil {
		query += `AND WHERE diary_id = :diary_id`
		namedArgs["diary_id"] = r.DiaryID
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
		ucEntries[i] = mapDiaryEntry(e)
	}
	return ucEntries, nil
}

func checkAccess(ctx context.Context, tx db.TxOrDb, entryID uuid.UUID) error {
	const checkAccessQuery = `SELECT EXISTS(
    	SELECT * FROM diary_entries e JOIN diaries d ON e.diary_id = d.id JOIN diary_keys dk ON d.id = dk.diary_id 
		WHERE e.id = $1 AND dk.user_id = $2)`

	userID := auth.MustGetUserID(ctx)
	var hasAccess bool
	if err := tx.QueryRowxContext(ctx, checkAccessQuery, entryID, userID).Scan(&hasAccess); err != nil {
		return err
	}
	if !hasAccess {
		return &usecase.NoAccessToDiaryEntryError{DiaryId: entryID}
	}

	return nil
}

func (p *pgRepo) GetByID(ctx context.Context, id uuid.UUID) (*usecase.DiaryEntry, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}
	err = checkAccess(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	entry, err := getEntry(ctx, id, tx)
	if err != nil {
		return nil, err
	}

	const blocksQuery = `SELECT id, diary_entry_id, value FROM diary_entry_blocks WHERE diary_entry_id = $1`
	var blocks []diaryEntryBlock
	if err := p.db.SelectContext(ctx, blocks, blocksQuery, id); err != nil {
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
	const entryQuery = `SELECT id, diary_id, name, date FROM diary_entries WHERE id = $1`
	entry := &diaryEntry{}
	if err := tx.GetContext(ctx, entry, entryQuery, id); err != nil {
		return nil, err
	}
	return entry, nil
}

func (p *pgRepo) Update(ctx context.Context, id uuid.UUID, r *usecase.UpdateDiaryEntryRequest) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}
	if err = validateUpdateRequest(ctx, tx, id, r); err != nil {
		return err
	}
	if err = update(ctx, tx, id, r); err != nil {
		if rbErr := tx.Rollback(); err != nil {
			return multierror.Append(err, rbErr)
		}
		return err
	}
	return nil
}

func validateUpdateRequest(ctx context.Context, tx db.TxOrDb, id uuid.UUID, r *usecase.UpdateDiaryEntryRequest) error {
	err := checkAccess(ctx, tx, id)
	if err != nil {
		return err
	}
	if r.BlocksToUpsert != nil {
		blocksIds := getBlocksIds(r.BlocksToUpsert)
		if err := validateBlocksIds(ctx, tx, id, blocksIds); err != nil {
			return err
		}
	}
	if r.BlocksToDelete != nil {
		if err := validateBlocksIds(ctx, tx, id, r.BlocksToDelete); err != nil {
			return err
		}
	}
	return nil
}

func update(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, r *usecase.UpdateDiaryEntryRequest) error {
	if err := updateEntry(ctx, tx, id, r); err != nil {
		return err
	}
	if err := upsertBlocks(ctx, tx, id, r.BlocksToUpsert); err != nil {
		return err
	}
	if err := deleteBlocks(ctx, tx, r.BlocksToDelete); err != nil {
		return err
	}
	return nil
}

func updateEntry(ctx context.Context, tx db.TxOrDb, id uuid.UUID, req *usecase.UpdateDiaryEntryRequest) error {
	if req.Name == nil && req.Date == nil && req.DiaryId == nil {
		return nil
	}
	const getQuery = `SELECT id, diary_id, name, date, value FROM diary_entries WHERE id = $1`
	entry := &diaryEntry{}
	if err := tx.GetContext(ctx, entry, getQuery, id); err != nil {
		return err
	}
	if req.DiaryId != nil {
		if err := checkAccess(ctx, tx, *req.DiaryId); err != nil {
			return err
		}
		entry.DiaryID = *req.DiaryId
	}
	if req.Name != nil {
		entry.Name = *req.Name
	}
	if req.Date != nil {
		entry.Date = *req.Date
	}
	if req.Value != nil {
		entry.Value = *req.Value
	}
	const updateQuery = `UPDATE diary_entries SET name = :name, date = :date, diary_id = :diary_id, value = :value WHERE id = :id`
	if _, err := tx.NamedExecContext(ctx, updateQuery, entry); err != nil {
		return err
	}
	return nil
}

func upsertBlocks(ctx context.Context, tx db.TxOrDb, id uuid.UUID, blocks []usecase.DiaryEntryBlockDto) error {
	if blocks == nil || len(blocks) == 0 {
		return nil
	}

	const query = `
INSERT INTO diary_entry_blocks (id, diary_entry_id, value) 
VALUES (:id, :diary_entry_id, :value)
ON CONFLICT DO UPDATE SET value = :value`

	data := make([]diaryEntryBlock, len(blocks))
	for i, b := range blocks {
		data[i] = diaryEntryBlock{
			ID:           b.ID,
			DiaryEntryID: id,
			Value:        b.Value,
		}
	}
	if _, err := tx.NamedExecContext(ctx, query, data); err != nil {
		return err
	}
	return nil
}

func deleteBlocks(ctx context.Context, tx db.TxOrDb, ids []uuid.UUID) error {
	if ids == nil || len(ids) == 0 {
		return nil
	}
	query, args, err := sqlx.In(`DELETE FROM diary_entry_blocks WHERE id IN (?)`, ids)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, query, args); err != nil {
		return err
	}
	return nil
}

func getBlocksIds(blocks []usecase.DiaryEntryBlockDto) []uuid.UUID {
	res := make([]uuid.UUID, len(blocks))
	for i, b := range blocks {
		res[i] = b.ID
	}
	return res
}

func validateBlocksIds(
	ctx context.Context,
	txOrDb db.TxOrDb,
	entryId uuid.UUID,
	blocksIds []uuid.UUID,
) error {
	getBadBlocksIdsQuery := `SELECT id FROM diary_entry_blocks WHERE id IN (?) AND diary_entry_id != ?`
	query, args, err := sqlx.In(getBadBlocksIdsQuery, blocksIds, entryId)
	if err != nil {
		return err
	}

	var badIds []uuid.UUID
	err = txOrDb.SelectContext(ctx, badIds, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows || len(badIds) == 0 {
		return nil
	}

	return &usecase.AlienEntryBlocksError{DiaryEntryId: entryId, AlienBlocksIds: badIds}
}

func (p *pgRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}
	if err := checkAccess(ctx, tx, id); err != nil {
		return err
	}
	const query = `DELETE FROM diary_entries WHERE id = $1`
	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}

func New(db *sqlx.DB) usecase.DiaryEntriesRepository {
	return &pgRepo{
		db: db,
	}
}

func (p *pgRepo) Create(ctx context.Context, entry *usecase.DiaryEntry) (*usecase.DiaryEntry, error) {
	const insertEntryQuery = `INSERT INTO diary_entries (id, diary_id, name, date) VALUES (:id, :diary_id, :name, :date)`
	if _, err := p.db.NamedExecContext(ctx, insertEntryQuery, entry); err != nil {
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
