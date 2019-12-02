package handlers

import (
	"encoding/json"
	"memoryLane/database"
	"memoryLane/server"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	http.Handler
}

func NewHandler(db *database.DB) *Handler {

	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin
		// hosts
		AllowedOrigins: []string{"*"}, //TODO(mac): before this roles out to prod this needs to include a config option which includes where this is publicly hosted
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true
		// },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	ss := server.NewScriptureServer(db)
	sh := newScriptureHandler(ss)

	r.Post("/api/scripture/new", http.HandlerFunc(sh.newScripture))

	return &Handler{Handler: r}
}

type scriptureHandler struct {
	scriptureServer *server.ScriptureServer
}

func newScriptureHandler(server *server.ScriptureServer) *scriptureHandler {
	return &scriptureHandler{scriptureServer: server}
}

func (ss *scriptureHandler) newScripture(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	decoder := json.NewDecoder(req.Body)
	var newScriptureRequest server.NewScriptureRequest
	err := decoder.Decode(&newScriptureRequest)
	if err != nil {
		logrus.Errorf("%+v", err)
		http.Error(w, "unable to parse json", http.StatusInternalServerError)
		return
	}

	err = ss.scriptureServer.NewScripture(ctx, &newScriptureRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Errorf("%+v", err)
		return
	}
}
