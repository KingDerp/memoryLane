package webutil

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeebo/errs"
)

func TestStatusFromErrors(t *testing.T) {
	var statusFromErrorTests = []struct {
		err      error
		expected int
	}{
		{ServerError.New(""), http.StatusInternalServerError},
		{ValidationError.New(""), http.StatusBadRequest},
		{errs.New(""), http.StatusInternalServerError},
		{nil, http.StatusInternalServerError},
	}

	for _, v := range statusFromErrorTests {
		assert.Equal(t, v.expected, statusFromError(v.err))
	}
}
