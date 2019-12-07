package server

import (
	"testing"

	"github.com/KingDerp/memoryLane/database"

	"github.com/stretchr/testify/require"
)

type serverTest struct {
	t              *testing.T
	db             *database.DB
	CitationServer *CitationServer
}

func newTest(t *testing.T) *serverTest {

	db, err := database.Open("sqlite3", "file:memdb1?mode=memory&cache=shared")
	require.NoError(t, err)

	//initialize database with schema
	_, err = db.Exec(db.Schema())
	require.NoError(t, err)

	citationServer := NewCitationServer(db)

	return &serverTest{
		t:              t,
		db:             db,
		CitationServer: citationServer,
	}
}

func (h *serverTest) tearDown() {
	require.NoError(h.t, h.db.Close())
}

func defaultCitationRequest() *CitationRequest {
	return &CitationRequest{
		Reference: "page 1 paragraph 1",
		Author:    "Charel Dickens",
		Text:      `It was the best of times, it was the worst of times, it was the age of wisdom, it was the age of foolishness, it was the epoch of belief, it was the epoch of incredulity...`,
		Book:      "A Tale of Two Cities",
		Hint:      "Best and Worst",
		Year:      2019,
	}
}
