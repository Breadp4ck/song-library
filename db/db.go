package db

import (
	"context"

	"github.com/Breadp4ck/song-library/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPGSQLPool() (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), configs.Envs().DBUrl())
}
