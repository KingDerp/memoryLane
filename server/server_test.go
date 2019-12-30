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
