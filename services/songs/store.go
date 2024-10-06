package songs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Song struct {
	SongID      uuid.UUID  `json:"song_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	SongName    *string    `json:"song_name" example:"Absolute Territory"`
	SongText    *string    `json:"song_text" example:"She's got a fetish for fine art\nA pair of knee-socks and an oversized sweatshirt\nShe goes right to my heart\nShe comes a'knocking with her stocking and I get hurt\n\nI get the feeling I'm in deep\nTroubled waters, but they're only thigh-high\nThis kind of girl don't get no sleep\nDon't wake your father, skip the starters, strap those garters up\nOh my my!"`
	GroupName   *string    `json:"group_name" example:"Ken Ashcorp"`
	Link        *string    `json:"link" example:"https://www.youtube.com/watch?v=kFZKgf5WG0g"`
	ReleaseDate *time.Time `json:"release_date" example:"09.03.2013"`
}

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) CreateSong(song Song) error {
	_, err := s.db.Exec(
		context.Background(),
		"INSERT INTO songs (song_name, song_text, group_name, link, release_date) VALUES ($1, $2, $3, $4, $5);",
		song.SongName,
		song.SongText,
		song.GroupName,
		song.Link,
		song.ReleaseDate,
	)

	return err
}

func (s *Store) RemoveSongByID(songID uuid.UUID) error {
	_, err := s.db.Exec(
		context.Background(),
		"DELETE FROM songs WHERE song_id=$1;",
		songID,
	)
	return err
}

func (s *Store) UpdateSongByID(songID uuid.UUID, song *Song) error {
	query, args := buildUpdateQuery(songID, song)
	_, err := s.db.Exec(context.Background(), query, args...)
	return err
}

// Returns pgx query string and array of values in exact order.
func buildUpdateQuery(songID uuid.UUID, song *Song) (string, []any) {
	var setClauses []string
	var args []any

	if song.SongName != nil {
		args = append(args, song.SongName)
		setClauses = append(setClauses, fmt.Sprintf("song_name = $%d", len(args)))
	}
	if song.SongText != nil {
		args = append(args, song.SongText)
		setClauses = append(setClauses, fmt.Sprintf("song_text = $%d", len(args)))
	}
	if song.GroupName != nil {
		args = append(args, song.GroupName)
		setClauses = append(setClauses, fmt.Sprintf("group_name = $%d", len(args)))
	}
	if song.Link != nil {
		args = append(args, song.Link)
		setClauses = append(setClauses, fmt.Sprintf("link = $%d", len(args)))
	}
	if song.ReleaseDate != nil {
		args = append(args, song.ReleaseDate)
		setClauses = append(setClauses, fmt.Sprintf("release_date = $%d", len(args)))
	}

	args = append(args, songID)
	query := fmt.Sprintf("UPDATE songs SET %s WHERE song_id = $%d", strings.Join(setClauses, ", "), len(args))

	return query, args
}

func (s *Store) GetSongByID(songID uuid.UUID) (Song, error) {
	var song Song
	err := s.db.QueryRow(
		context.Background(),
		"SELECT song_id, song_name, song_text, group_name, link, release_date FROM songs WHERE song_id=$1;",
		songID,
	).Scan(&song.SongID, &song.SongName, &song.SongText, &song.GroupName, &song.Link, &song.ReleaseDate)
	return song, err
}

type SongsFilter struct {
	SongName    *string
	ReleaseDate *time.Time
	GroupName   *string
}

func (s *Store) GetSongsFiltered(pageCurrent uint, pageSize uint, filter *SongsFilter) ([]Song, error) {
	query, args := buildFilterQuery(pageCurrent, pageSize, filter)

	rows, err := s.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := make([]Song, 0)
	for rows.Next() {
		var song Song
		err := rows.Scan(&song.SongID, &song.SongName, &song.SongText, &song.GroupName, &song.Link, &song.ReleaseDate)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

// Returns pgx query string and array of values in exact order.
func buildFilterQuery(pageCurrent uint, pageSize uint, filter *SongsFilter) (string, []any) {
	var whereClauses []string
	var args []any

	if filter.SongName != nil {
		args = append(args, filter.SongName)
		whereClauses = append(whereClauses, fmt.Sprintf("song_name LIKE $%d", len(args)))
	}

	if filter.ReleaseDate != nil {
		args = append(args, filter.ReleaseDate)
		whereClauses = append(whereClauses, fmt.Sprintf("release_date == $%d", len(args)))
	}

	if filter.GroupName != nil {
		args = append(args, filter.GroupName)
		whereClauses = append(whereClauses, fmt.Sprintf("group_name LIKE $%d", len(args)))
	}

	where := ""
	if len(args) != 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(whereClauses, " AND "))
	}

	query := fmt.Sprintf(
		"SELECT song_id, song_name, song_text, group_name, link, release_date FROM songs %s LIMIT %d OFFSET %d;",
		where,
		pageSize,
		pageCurrent,
	)

	return query, args
}

func (s *Store) GetLyrcsBySongID(songID uuid.UUID, verseCurrent, verseCount uint) ([]string, error) {
	song, err := s.GetSongByID(songID)

	if song.SongText != nil {
		return extractVerses(*song.SongText, verseCurrent, verseCount), nil
	}

	return nil, err
}

func extractVerses(songText string, verseCurrent, verseCount uint) []string {
	verses := strings.Split(songText, "\n\n")
	versesLength := uint(len(verses))

	if verseCurrent >= versesLength {
		return make([]string, 0)
	}

	versesEnd := verseCurrent + verseCount
	if versesEnd > versesLength {
		versesEnd = uint(len(verses))
	}
	return verses[verseCurrent:versesEnd]
}
