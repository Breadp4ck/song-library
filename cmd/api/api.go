package api

import (
	"log/slog"

	"github.com/Breadp4ck/song-library/services/songs"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Breadp4ck/song-library/docs"
)

type APIServer struct {
	addr string
	db   *pgxpool.Pool
}

// Create new APIServer that binds to addr and uses db connection pool.
func NewAPIServer(addr string, db *pgxpool.Pool) *APIServer {
	return &APIServer{addr: addr, db: db}
}

// Run server and capture logs to provided logger.
func (s *APIServer) Run(logger *slog.Logger) {
	r := gin.New()

	// Use slog logger as middleware.
	r.Use(sloggin.New(logger))
	r.Use(gin.Recovery())

	// Setup API
	v1 := r.Group("/api/v1")
	{
		songsStore := songs.NewStore(s.db)
		songsHandler := songs.NewHandler(songsStore)
		songsHandler.RegisterRoutes(v1)
	}

	// Add swagger.
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(s.addr)
}
