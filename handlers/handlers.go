package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/KingDerp/memoryLane/common"
	"github.com/KingDerp/memoryLane/database"
	"github.com/KingDerp/memoryLane/server"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	http.Handler
}

func NewHandler(db *database.DB) *Handler {

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, //TODO(mac): before this roles out to prod this needs to include a config option which includes where this is publicly hosted
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
		logrus.Errorf("error decoding request: %+v", err)
		http.Error(w, "unable to parse json", http.StatusBadRequest)
		return
	}

	err = ss.scriptureServer.NewScripture(ctx, &newScriptureRequest)

	if err != nil {
		if common.ServerError.Has(err) {
			logrus.Errorf("server error: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if common.ValidationError.Has(err) {
			logrus.Errorf("validation error: %+v", err)
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		logrus.Errorf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(&server.NewScriptureResponse{
		Message: "Scripture was succesfully received and stored",
	})
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	h := w.Header()
	h.Set("Content-Type", "application/json")
	w.Write(b)
}
