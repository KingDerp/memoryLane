package handlers

import (
	"net/http"

	"github.com/KingDerp/memoryLane/database"
	"github.com/KingDerp/memoryLane/server"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type Handler struct {
	http.Handler
}

func NewHandler(db *database.DB) *Handler {

	r := chi.NewRouter()

	r.Use(newCORS().Handler)

	ss := server.NewCitationServer(db)
	sh := newCitationHandler(ss)

	r.Post("/api/citation/new", http.HandlerFunc(sh.newCitation))

	return &Handler{Handler: r}
}

func newCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, //TODO(mac): before this roles out to prod this needs to include a config option which includes where this is publicly hosted
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
