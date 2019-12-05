package server

import (
	"context"
	"testing"

	ml "github.com/KingDerp/memoryLane"

	"github.com/stretchr/testify/require"
)

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
	require.True(t, ml.HasValidationError(err))
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
