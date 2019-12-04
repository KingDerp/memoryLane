package server

import (
	"context"

	"github.com/KingDerp/memoryLane/common"
	"github.com/KingDerp/memoryLane/database"

	uuid "github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type CitationServer struct {
	db *database.DB
}

func NewCitationServer(db *database.DB) *CitationServer {
	return &CitationServer{db: db}
}

type CitationRequest struct {
	Reference string `json:"reference"` //helps to locate item within a book. Ex: page number etc.
	Author    string `json:"author"`
	Text      string `json:"text"`
	Book      string `json:"book"`
	Hint      string `json:"hint"`
	Year      int64  `json:"year"`
}

type NewCitationResponse struct {
	Message string `json:"message"`
}

func (ss *CitationServer) NewCitation(ctx context.Context, req *CitationRequest) (err error) {

	logrus.Info("entered citation server")
	err = ValidateCitationRequest(req)
	if err != nil {
		return err
	}

	newId, err := uuid.NewV4()
	if err != nil {
		return err
	}

	err = ss.db.WithTx(ctx, func(ctx context.Context, tx *database.Tx) error {

		return tx.CreateNoReturn_Citation(ctx,
			database.Citation_Id(newId.String()),
			database.Citation_Text(req.Text),
			database.Citation_Create_Fields{
				Reference: database.Citation_Reference(req.Reference),
				Author:    database.Citation_Author(req.Author),
				Book:      database.Citation_Book(req.Book),
				Hint:      database.Citation_Hint(req.Hint),
				Year:      database.Citation_Year(req.Year),
			},
		)
	})

	if err != nil {
		return common.ServerError.Wrap(err)
	}

	return nil
}

func ValidateCitationRequest(c *CitationRequest) error {
	if c.Text == "" {
		return common.ValidationError.New("text cannot be empty. You must have something to memorize.")
	}

	return nil
}
