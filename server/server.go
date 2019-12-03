package server

import (
	"context"

	"github.com/KingDerp/memoryLane/common"
	"github.com/KingDerp/memoryLane/database"

	uuid "github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
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

type NewScriptureResponse struct {
	Message string `json:"message"`
}

func (ss *ScriptureServer) NewScripture(ctx context.Context, req *NewScriptureRequest) (err error) {

	logrus.Info("entered scripture server")
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
		return common.ServerError.Wrap(err)
	}

	return nil
}

func ValidateScriptureRequest(sr *NewScriptureRequest) error {
	if sr.Book == "" {
		return common.ValidationError.New("book cannot be empty")
	}

	if sr.VerseNumber <= 0 {
		return common.ValidationError.New("verseNumber must be greater than 0")
	}

	if sr.VerseText == "" {
		return common.ValidationError.New("verseText cannot be empty")
	}

	if sr.Chapter <= 0 {
		return common.ValidationError.New("chapter must be greater than 0")
	}

	return nil
}
