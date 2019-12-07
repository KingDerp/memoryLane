package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/KingDerp/memoryLane/server"
	wu "github.com/KingDerp/memoryLane/webutil"
)

type citationHandler struct {
	citationServer *server.CitationServer
}

func newCitationHandler(server *server.CitationServer) *citationHandler {
	return &citationHandler{citationServer: server}
}

func (ss *citationHandler) newCitation(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	newCitationRequest, err := decodeCitationReq(req.Body)
	if err != nil {
		wu.HandleError(w, err)
		return
	}

	err = ss.citationServer.NewCitation(ctx, newCitationRequest)
	if err != nil {
		wu.HandleError(w, err)
		return
	}

	wu.RenderJSON(w, &server.NewCitationResponse{
		Message: "citation was succesfully received and stored",
	})
}

func decodeCitationReq(io io.ReadCloser) (*server.CitationRequest, error) {
	c := new(server.CitationRequest)

	if err := json.NewDecoder(io).Decode(c); err != nil {
		return nil, err
	}

	return c, nil
}
