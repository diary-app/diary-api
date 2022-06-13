package repository

import (
	"context"
	"database/sql"
	"diary-api/internal/db"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
	"time"
)

func checkForBlocksFromOtherEntries(
	ctx context.Context,
	tx db.TxOrDb,
	entryID uuid.UUID,
	blocksIds []uuid.UUID,
) error {
	getBadBlocksIdsQuery := `SELECT id FROM diary_entry_blocks WHERE id IN (?) AND diary_entry_id != ?`
	query, args, err := sqlx.In(getBadBlocksIdsQuery, blocksIds, entryID)
	if err != nil {
		return err
	}
	query = tx.Rebind(query)

	var badIds []uuid.UUID
	err = tx.SelectContext(ctx, &badIds, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows || len(badIds) == 0 {
		return nil
	}

	return &usecase.AlienEntryBlocksError{AlienBlocksIds: badIds}
}

func checkForAlienBlocks(ctx context.Context, tx db.TxOrDb, entryID uuid.UUID, blocksToDelete []uuid.UUID) error {
	const query = `SELECT id FROM diary_entry_blocks WHERE diary_entry_id = $1`
	var existingIds []uuid.UUID
	if err := tx.SelectContext(ctx, &existingIds, query, entryID); err != nil {
		return err
	}
	existingIdsMap := map[uuid.UUID]struct{}{}
	for _, id := range existingIds {
		existingIdsMap[id] = struct{}{}
	}
	var alienBlocks []uuid.UUID
	for _, id := range blocksToDelete {
		if _, ok := existingIdsMap[id]; !ok {
			alienBlocks = append(alienBlocks, id)
		}
	}
	if len(alienBlocks) > 0 {
		return &usecase.AlienEntryBlocksError{AlienBlocksIds: alienBlocks}
	}
	return nil
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
		if rbErr := tx.Rollback(); rbErr != nil {
			return multierror.Append(err, rbErr)
		}
		return err
	}
	return db.ShouldCommitOrRollback(tx)
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

func validateUpdateRequest(ctx context.Context, tx db.TxOrDb, entryID uuid.UUID, r *usecase.UpdateDiaryEntryRequest) error {
	err := db.CheckMyWriteAccessToEntry(ctx, tx, entryID)
	if err != nil {
		return err
	}
	if len(r.BlocksToUpsert) > 0 {
		blocksIds := getBlocksIds(r.BlocksToUpsert)
		if err = checkForBlocksFromOtherEntries(ctx, tx, entryID, blocksIds); err != nil {
			return err
		}
	}
	if len(r.BlocksToDelete) > 0 {
		if err = checkForAlienBlocks(ctx, tx, entryID, r.BlocksToDelete); err != nil {
			return err
		}
	}
	return nil
}

func updateEntry(ctx context.Context, tx db.TxOrDb, id uuid.UUID, req *usecase.UpdateDiaryEntryRequest) error {
	if req.Name == nil && req.Date == nil && req.DiaryId == nil && req.Value == nil {
		return nil
	}
	const getQuery = `SELECT id, diary_id, name, date, value FROM diary_entries WHERE id = $1`
	entry := &diaryEntry{}
	if err := tx.GetContext(ctx, entry, getQuery, id); err != nil {
		return err
	}
	if req.DiaryId != nil {
		if err := db.CheckMyWriteAccessToDiary(ctx, tx, *req.DiaryId); err != nil {
			return err
		}
		entry.DiaryID = *req.DiaryId
	}
	if req.Name != nil {
		entry.Name = *req.Name
	}
	if req.Date != nil {
		entry.Date = time.Time(*req.Date)
	}
	if req.Value != nil {
		entry.Value = req.Value
	}
	const updateQuery = `UPDATE diary_entries SET name = :name, date = :date, diary_id = :diary_id, value = :value WHERE id = :id`
	if _, err := tx.NamedExecContext(ctx, updateQuery, entry); err != nil {
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
	query = tx.Rebind(query)
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
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
ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value`

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
