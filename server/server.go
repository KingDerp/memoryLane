package server

import (
	"context"
	"database/sql"
	"time"

	"github.com/KingDerp/memoryLane/database"
	wu "github.com/KingDerp/memoryLane/webutil"

	uuid "github.com/gofrs/uuid"
	"github.com/k0kubun/pp"
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

type CitationResponse struct {
	Reference  string    `json:"reference"` //helps to locate item within a book. Ex: page number etc.
	Author     string    `json:"author"`
	Text       string    `json:"text"`
	Book       string    `json:"book"`
	Hint       string    `json:"hint"`
	Year       int64     `json:"year"`
	MemoryDate time.Time `json:"memoryDate"`
}

type NewCitationResponse struct {
	Message string `json:"message"`
}

func (ss *CitationServer) NewCitation(ctx context.Context, req *CitationRequest) (err error) {

	logrus.Info("Creating New Citation")
	err = ValidateCitationReq(req)
	if err != nil {
		return err
	}

	newId, err := uuid.NewV4()
	if err != nil {
		return err
	}

	err = ss.db.WithTx(ctx, func(ctx context.Context, tx *database.Tx) error {

		return tx.CreateNoReturn_Citation(ctx,
			database.Citation_MemDate(time.Now()),
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
		return wu.ServerError.Wrap(err)
	}

	return nil
}

func (ss *CitationServer) GetDailyCitations(ctx context.Context, req *CitationRequest) (
	c []CitationResponse, err error) {
	logrus.Info("Getting today's citations")

	rows, err := ss.db.GetDailyScriptures()
	if err != nil {
		return nil, wu.ServerError.Wrap(err)
	}
	defer rows.Close()

	return rowsToCitationResponse(rows)
}

func rowsToCitationResponse(rows *sql.Rows) (c []CitationResponse, err error) {
	var crs []CitationResponse
	for rows.Next() {
		var reference, author, text, book, hint string
		var year int64
		var memDate time.Time
		err = rows.Scan(&reference, &author, &text, &book, &hint, &year, &memDate)
		if err != nil {
			return nil, wu.ServerError.Wrap(err)
		}

		crs = append(crs, CitationResponse{
			Reference:  reference,
			Author:     author,
			Text:       text,
			Book:       book,
			Hint:       hint,
			Year:       year,
			MemoryDate: memDate,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, wu.ServerError.Wrap(err)
	}

	pp.Print(crs)
	return crs, err
}

func ValidateCitationReq(c *CitationRequest) error {
	if c.Text == "" {
		return wu.ValidationError.New("text cannot be empty. You must have something to memorize.")
	}

	return nil
}
