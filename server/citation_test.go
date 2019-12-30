package server

import (
	"context"
	"testing"
	"time"

	"github.com/KingDerp/memoryLane/database"
	wu "github.com/KingDerp/memoryLane/webutil"
	uuid "github.com/gofrs/uuid"

	"github.com/stretchr/testify/require"
)

func TestCitationsInTheLastWeek(t *testing.T) {
	s := newTest(t)
	defer s.tearDown()

}

func TestCitationCreate(t *testing.T) {
	s := newTest(t)
	defer s.tearDown()

	err := s.CitationServer.NewCitation(context.Background(), defaultCitationRequest())
	require.NoError(t, err)
}

func TestCitationMissingTextField(t *testing.T) {
	s := newTest(t)
	defer s.tearDown()

	//set required field set to zero val
	c := defaultCitationRequest()
	c.Text = ""

	err := s.CitationServer.NewCitation(context.Background(), c)
	require.True(t, wu.HasValidationError(err))
}

func TestValidateCitationRequestMissingText(t *testing.T) {
	s := newTest(t)
	defer s.tearDown()

	//set required field set to zero val
	c := defaultCitationRequest()
	c.Text = ""

	err := ValidateCitationReq(c)
	require.Error(t, err)
}

//========================================== UTILITY ==============================================

func (st *serverTest) defaultCitationsForLastNDays(ctx context.Context, n int) (err error) {

	for i := 0; i < n; i++ {
		var newId uuid.UUID
		newId, err = uuid.NewV4()
		if err != nil {
			return err
		}

		defaultCitation := defaultCitationRequest()

		err = st.db.WithTx(ctx, func(ctx context.Context, tx *database.Tx) error {

			return tx.CreateNoReturn_Citation(ctx,
				database.Citation_MemDate(time.Now().AddDate(0, 0, (-1*i))),
				database.Citation_Id(newId.String()),
				database.Citation_Text(defaultCitation.Text),
				database.Citation_Create_Fields{
					Reference: database.Citation_Reference(defaultCitation.Reference),
					Author:    database.Citation_Author(defaultCitation.Author),
					Book:      database.Citation_Book(defaultCitation.Book),
					Hint:      database.Citation_Hint(defaultCitation.Hint),
					Year:      database.Citation_Year(defaultCitation.Year),
				},
			)

		})
	}

	return err
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
