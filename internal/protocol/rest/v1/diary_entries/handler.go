package diary_entries

import (
	"diary-api/internal/protocol/rest/common"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetEntries() gin.HandlerFunc
	GetByID() gin.HandlerFunc
	Create() gin.HandlerFunc
	PatchEntry() gin.HandlerFunc
	Delete() gin.HandlerFunc
}

func New(uc usecase.DiaryEntriesUseCase) Handler {
	return &handler{
		uc: uc,
	}
}

type handler struct {
	uc usecase.DiaryEntriesUseCase
}

func mapToEntryResponse(e *usecase.DiaryEntry) usecase.DiaryEntryResponse {
	blocks := mapToBlocksResponse(e)
	return usecase.DiaryEntryResponse{
		ID:      e.ID,
		DiaryID: e.DiaryID,
		Name:    e.Name,
		Date:    common.DateOnly(e.Date),
		Value:   e.Value,
		Blocks:  blocks,
	}
}

func mapToShortEntry(e *usecase.DiaryEntry) usecase.ShortDiaryEntryResponse {
	return usecase.ShortDiaryEntryResponse{
		ID:      e.ID,
		DiaryID: e.DiaryID,
		Name:    e.Name,
		Date:    common.DateOnly(e.Date),
	}
}

func mapToBlocksResponse(e *usecase.DiaryEntry) []usecase.DiaryEntryBlockResponse {
	if e.Blocks != nil && len(e.Blocks) > 0 {
		blocks := make([]usecase.DiaryEntryBlockResponse, len(e.Blocks))
		for i, b := range e.Blocks {
			blocks[i] = usecase.DiaryEntryBlockResponse{
				ID:    b.ID,
				Value: b.Value,
			}
		}
		return blocks
	} else {
		return make([]usecase.DiaryEntryBlockResponse, 0)
	}
}
