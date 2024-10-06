package songs

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	_ "github.com/Breadp4ck/song-library/docs"
	"github.com/Breadp4ck/song-library/types"
)

const PAGE_SIZE_DEAFULT = uint(10)
const PAGE_SIZE_MAX = uint(50)
const VERSE_COUNT_DEFAULT = uint(20)
const VERSE_COUNT_MAX = uint(50)

type Handler struct {
	store *Store
}

func NewHandler(s *Store) *Handler {
	return &Handler{s}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/info")
	{
		g.POST("", h.handleCreateSong)
		g.DELETE("/:song_id", h.handleRemoveSong)
		g.PUT("/:song_id", h.handleUpdateSong)
		g.GET("/:song_id", h.handleGetSong)
		g.GET("", h.handleGetSongs)
		g.GET("/:song_id/lyrcs", h.handleGetLyrcs)
	}
}

type BindUUID struct {
	UUID uuid.UUID
}

func (u *BindUUID) UnmarshalParam(param string) error {
	parsedUUID, err := uuid.Parse(param)
	if err != nil {
		return err
	}
	u.UUID = parsedUUID
	return nil
}

type BindDate struct {
	date time.Time
}

func (rd *BindDate) UnmarshalParam(param string) error {
	const DATE_FORMAT = "02.01.2006"
	date, err := time.Parse(DATE_FORMAT, param)
	if err != nil {
		return err
	}
	rd.date = date
	return nil
}

func (rd *BindDate) Inner() *time.Time {
	return &rd.date
}

type CreateSongRequest struct {
	GroupName   string    `json:"group_name" binding:"required" example:"Jamiroquai"`
	SongName    string    `json:"song_name" binding:"required" example:"Virtual Insanity"`
	SongText    *string   `json:"song_text" example:"Oh yeah, aw\nWhat we're livin' in?\nLet me tell ya\n\nYeah, it's a wonder man can eat at all\nWhen things are big that should be small\nWho can tell what magic spells we'll be doin' for us?\nAnd I'm givin' all my love to this world"`
	Link        *string   `json:"link" example:"https://www.youtube.com/watch?v=4JkIs37a2JE"`
	ReleaseDate *BindDate `json:"release_date"`
}

// CreateSong godoc
//
//	@Summary	Create song
//	@Accept		json
//	@Produce	json
//	@Param		request	body	CreateSongRequest	true	"Create Song Request"
//	@Success	201		{string}		ok
//	@Failure	400		{object}		ServiceError
//	@Router		/info [post]
func (h *Handler) handleCreateSong(c *gin.Context) {
	var request CreateSongRequest

	if err := c.Bind(&request); err != nil {
		c.JSON(WrongParametersError())
		return
	}

	var releaseDate *time.Time
	if request.ReleaseDate != nil {
		releaseDate = request.ReleaseDate.Inner()
	}

	err := h.store.CreateSong(Song{
		GroupName:   &request.GroupName,
		SongName:    &request.SongName,
		SongText:    request.SongText,
		Link:        request.Link,
		ReleaseDate: releaseDate,
	})
	if err != nil {
		c.JSON(WrongParametersError())
		slog.Info(err.Error())
		return
	}

	c.JSON(http.StatusCreated, types.ResponseOK)
}

type RemoveSongRequest struct {
	SongID BindUUID `uri:"song_id" binding:"required"`
}

// RemoveSong godoc
//
//	@Summary	Remove a song
//	@Accept		json
//	@Produce	json
//	@Param		song_id	path	string	true	"Song ID"
//	@Success	200		{string}	ok
//	@Failure	400		{object}	ServiceError
//	@Failure	404		{object}	ServiceError
//	@Router		/info/{song_id} [delete]
func (h *Handler) handleRemoveSong(c *gin.Context) {
	var request RemoveSongRequest

	if err := c.BindUri(&request); err != nil {
		c.JSON(WrongParametersError())
		return
	}

	err := h.store.RemoveSongByID(request.SongID.UUID)
	if err != nil {
		c.JSON(SongNotFoundError(request.SongID.UUID))
		return
	}

	c.JSON(http.StatusOK, types.ResponseOK)
}

type UpdateSongRequest struct {
	SongID      BindUUID  `uri:"song_id" binding:"required"`
	SongName    *string   `json:"song_name" example:"Virtual Insanity"`
	SongText    *string   `json:"song_text" example:"Oh yeah, aw\nWhat we're livin' in?\nLet me tell ya\n\nYeah, it's a wonder man can eat at all\nWhen things are big that should be small\nWho can tell what magic spells we'll be doin' for us?\nAnd I'm givin' all my love to this world"`
	GroupName   *string   `json:"group_name" example:"Jamiroquai"`
	Link        *string   `json:"link" example:"https://www.youtube.com/watch?v=4JkIs37a2JE"`
	ReleaseDate *BindDate `json:"release_date"`
}

// UpdateSong godoc
//
//	@Summary	Update a song
//	@Accept		json
//	@Produce	json
//	@Param		song_id	path string	true	"Song ID"
//	@Param		request	body UpdateSongRequest	true	"Update Song Request"
//	@Success	200		{string}	ok
//	@Failure	400		{object}	ServiceError
//	@Failure	404		{object}	ServiceError
//	@Router		/info/{song_id} [put]
func (h *Handler) handleUpdateSong(c *gin.Context) {
	var request UpdateSongRequest

	if err := c.BindUri(&request); err != nil {
		c.JSON(WrongParametersError())
		return
	}

	if err := c.Bind(&request); err != nil {
		c.JSON(WrongParametersError())
		return
	}

	var releaseDate *time.Time
	if request.ReleaseDate != nil {
		releaseDate = request.ReleaseDate.Inner()
	}

	err := h.store.UpdateSongByID(request.SongID.UUID, &Song{
		SongName:    request.SongName,
		SongText:    request.SongText,
		GroupName:   request.GroupName,
		Link:        request.Link,
		ReleaseDate: releaseDate,
	})
	if err != nil {
		slog.Info(err.Error())
		c.JSON(SongNotFoundError(request.SongID.UUID))
		return
	}

	c.JSON(http.StatusOK, types.ResponseOK)
}

type GetSongRequest struct {
	SongID BindUUID `uri:"song_id" binding:"required"`
}

// GetSong godoc
//
//	@Summary	Get a song
//	@Accept		json
//	@Produce	json
//	@Param		song_id	path string	true	"Song ID"
//	@Success	200		{object}	Song
//	@Failure	400		{object}	ServiceError
//	@Failure	404		{object}	ServiceError
//	@Router		/info/{song_id} [get]
func (h *Handler) handleGetSong(c *gin.Context) {
	var request GetSongRequest

	if err := c.BindUri(&request); err != nil {
		c.JSON(WrongParametersError())
		return
	}

	song, err := h.store.GetSongByID(request.SongID.UUID)
	if err != nil {
		c.JSON(SongNotFoundError(request.SongID.UUID))
		return
	}

	c.JSON(http.StatusOK, types.Response{Message: song})
}

type GetSongsRequest struct {
	PageCurrent uint  `form:"page_current"`
	PageSize    *uint `form:"page_size"`

	SongName    *string   `form:"song_name" example:"Absolute Territory"`
	GroupName   *string   `form:"group_name" example:"Ken Ashcorp"`
	ReleaseDate *BindDate `form:"release_date" example:"09.03.2013"`
}

// GetSongs godoc
//
//	@Summary	Get songs
//	@Accept		json
//	@Produce	json
//	@Param		song_id			path	string				true	"Song ID"
//	@Param		page_current	query	int					false	"Current songs page"
//	@Param		page_size		query	int					false	"Songs per request"
//	@Param		page_size		body	GetSongsRequest		false	"Songs per request"
//	@Success	200		{object}	[]Song
//	@Failure	400		{object}	ServiceError
//	@Router		/info [get]
func (h *Handler) handleGetSongs(c *gin.Context) {
	var request GetSongsRequest

	if err := c.Bind(&request); err != nil {
		c.JSON(WrongParametersError())
		return
	}

	pageSize := PAGE_SIZE_DEAFULT
	if request.PageSize != nil {
		pageSize = *request.PageSize
	}
	if pageSize > PAGE_SIZE_MAX {
		c.JSON(BadPageSizeError(uint(pageSize)))
		return
	}

	var filter SongsFilter
	filter.SongName = request.SongName
	filter.GroupName = request.GroupName
	if request.ReleaseDate != nil {
		var releaseDate *time.Time
		if request.ReleaseDate != nil {
			releaseDate = request.ReleaseDate.Inner()
		}
		filter.ReleaseDate = releaseDate
	}

	songs, err := h.store.GetSongsFiltered(request.PageCurrent, pageSize, &filter)
	if err != nil {
		c.JSON(WrongParametersError())
		return
	}

	c.JSON(http.StatusOK, types.Response{Message: songs})
}

type GetLyrcsRequest struct {
	SongID BindUUID `uri:"song_id" binding:"required"`

	VerseCurrent uint  `form:"verse_current" example:"20"`
	VerseCount   *uint `form:"verse_count" example:"5"`
}

// GetLyrcs godoc
//
//	@Summary	Get song's verses as array
//	@Accept		json
//	@Produce	json
//	@Param		song_id	path	string	true	"Song ID"
//	@Param		verse_current	query	int		false	"Current verse"
//	@Param		verse_count		query	int		false	"Verse per request"
//	@Success	200		{object}	[]string
//	@Failure	400		{object}	ServiceError
//	@Router		/info/{song_id}/lyrcs [get]
func (h *Handler) handleGetLyrcs(c *gin.Context) {
	var request GetLyrcsRequest

	if err := c.BindUri(&request); err != nil {
		c.JSON(WrongParametersError())
		return
	}

	if err := c.Bind(&request); err != nil {
		c.JSON(WrongParametersError())
		return
	}

	verseCount := VERSE_COUNT_DEFAULT
	if request.VerseCount != nil {
		verseCount = *request.VerseCount
	}
	if verseCount > VERSE_COUNT_DEFAULT {
		c.JSON(BadPageSizeError(uint(verseCount)))
		return
	}

	lyrcs, err := h.store.GetLyrcsBySongID(request.SongID.UUID, request.VerseCurrent, verseCount)
	if err != nil {
		c.JSON(WrongParametersError())
		return
	}

	c.JSON(http.StatusOK, types.Response{Message: lyrcs})
}
