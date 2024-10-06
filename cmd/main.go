package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/Breadp4ck/song-library/cmd/api"
	"github.com/Breadp4ck/song-library/configs"
	"github.com/Breadp4ck/song-library/db"

	_ "github.com/Breadp4ck/song-library/docs"
)

//	@title		Music Info
//	@version	0.0.1

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	// 1. Setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 2. Load config
	slog.Info("Load configuration file.")
	configs.Setup()

	// 3. Init db
	slog.Info("Connect to database.")
	db, err := db.NewPGSQLPool()
	if err != nil {
		log.Panic(err)
	}

	// 4. Run server
	slog.Info("Run server.")
	server := api.NewAPIServer("localhost:8080", db)
	server.Run(logger)
}
