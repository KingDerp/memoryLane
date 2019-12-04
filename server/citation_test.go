package server

import (
	"context"
	"testing"

	"github.com/KingDerp/memoryLane/common"
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

	//set required field to zero val
	c := defaultCitationRequest()
	c.Text = ""

	err := s.CitationServer.NewCitation(context.Background(), c)
	require.True(t, common.ValidationError.Has(err))
}

func TestValidateCitationRequestMissingText(t *testing.T) {
	s := newTest(t)
	defer s.tearDown()

	//set required field to zero val
	c := defaultCitationRequest()
	c.Text = ""

	err := ValidateCitationRequest(c)
	require.Error(t, err)
}
