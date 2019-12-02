package server

import (
	"context"
	"fmt"
	"memoryLane/database"

	uuid "github.com/satori/go.uuid"
	"github.com/zeebo/errs"
	"gopkg.in/spacemonkeygo/dbx.v1/prettyprint"
)

type ScriptureServer struct {
	db *database.DB
}

func NewScriptureServer(db *database.DB) *ScriptureServer {
	return &ScriptureServer{db: db}
}

type NewScriptureRequest struct {
	Book        string `json:"book"`
	VerseNumber int    `json:"verseNumber"`
	VerseText   string `json:"verseText"`
	Chapter     int    `json:"chapter"`
	Hint        string `json:"hint"`
}

func (ss *ScriptureServer) NewScripture(ctx context.Context, req *NewScriptureRequest) (err error) {

	fmt.Println("entered scripture server")
	prettyprint.Println(req)
	err = ValidateScriptureRequest(req)
	if err != nil {
		return err
	}

	newId, err := uuid.NewV4()
	if err != nil {
		return err
	}

	err = ss.db.WithTx(ctx, func(ctx context.Context, tx *database.Tx) error {

		return tx.CreateNoReturn_Scripture(ctx,
			database.Scripture_Id(newId.String()),
			database.Scripture_Chapter(int64(req.Chapter)),
			database.Scripture_Book(req.Book),
			database.Scripture_VerseNumber(int64(req.VerseNumber)),
			database.Scripture_VerseText(req.VerseText),
			database.Scripture_Hint(req.Hint),
		)
	})

	if err != nil {
		return err
	}

	return nil
}

func ValidateScriptureRequest(sr *NewScriptureRequest) error {
	if sr.Book == "" {
		return errs.New("Book cannot be empty")
	}

	if sr.VerseNumber <= 0 {
		return errs.New("Verse Number must be greater than 0")
	}

	if sr.VerseText == "" {
		return errs.New("Verse Text cannot be empty")
	}

	if sr.Chapter <= 0 {
		return errs.New("Chapter must be greater than or equal to 0")
	}

	if sr.Hint == "" {
		return errs.New("Hint cannot be empty")
	}

	return nil
}
