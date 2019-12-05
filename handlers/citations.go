package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	ml "github.com/KingDerp/memoryLane"
	"github.com/KingDerp/memoryLane/server"

	"github.com/sirupsen/logrus"
)

type citationHandler struct {
	citationServer *server.CitationServer
}

func newCitationHandler(server *server.CitationServer) *citationHandler {
	return &citationHandler{citationServer: server}
}

func (ss *citationHandler) newCitation(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	newCitationRequest, err := unmarshalCitationReq(req.Body)
	if err != nil {
		ml.HandleUnmarshalError(w, err)
		return
	}

	err = ss.citationServer.NewCitation(ctx, newCitationRequest)
	if err != nil {
		ml.HandleCitationError(w, err)
		return
	}

	b, err := marshalCitationResp()
	if err != nil {
		ml.HandleMarshalError(w, err)
		return
	}

	renderCitationResp(w, b)
}

func renderCitationResp(w http.ResponseWriter, b []byte) {
	h := w.Header()
	h.Set("Content-Type", "application/json")
	w.Write(b)
}

func unmarshalCitationReq(io io.ReadCloser) (*server.CitationRequest, error) {
	var c *server.CitationRequest

	if err := json.NewDecoder(io).Decode(c); err != nil {
		logrus.Errorf("error decoding request: %+v", err)
		return nil, err
	}

	return c, nil
}

func marshalCitationResp() (b []byte, err error) {
	return json.Marshal(&server.NewCitationResponse{
		Message: "citation was succesfully received and stored",
	})
}
