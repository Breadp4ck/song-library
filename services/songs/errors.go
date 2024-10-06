package songs

import (
	"fmt"
	"net/http"

	"github.com/Breadp4ck/song-library/types"
	"github.com/google/uuid"
)

type ServiceError struct {
	Type   string `json:"type"`   // Unique identifier of the error
	Detail string `json:"detail"` // Human-readable detail of the error
}

func SongNotFoundError(songID uuid.UUID) (int, *types.Response) {
	return http.StatusNotFound, &types.Response{Error: ServiceError{
		Type:   "SongNotFound",
		Detail: fmt.Sprintf("Song with id %s is not found.", songID.String()),
	}}
}

func WrongParametersError() (int, *types.Response) {
	return http.StatusBadRequest, &types.Response{Error: ServiceError{
		Type:   "WrongParameters",
		Detail: "Wrong parameters for endpoint. Consider reading documentation.",
	}}
}

func BadPageSizeError(suppliedSize uint) (int, *types.Response) {
	return http.StatusBadRequest, &types.Response{Error: ServiceError{
		Type:   "BadPageSize",
		Detail: fmt.Sprintf("Page size more than %d is not allowed. Yours is %d.", PAGE_SIZE_MAX, suppliedSize),
	}}
}

func BadVerseCountError(suppliedCount uint) (int, *types.Response) {
	return http.StatusBadRequest, &types.Response{Error: ServiceError{
		Type:   "BadVerseCount",
		Detail: fmt.Sprintf("Verse count more than %d is not allowed. Yours is %d.", VERSE_COUNT_MAX, suppliedCount),
	}}
}
