package common

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeebo/errs"
)

var statusFromErrorTests = []struct {
	err      error
	expected int
}{
	{ServerError.New("server"), http.StatusInternalServerError},
	{ValidationError.New("validation"), http.StatusBadRequest},
	{MarshalError.New("marshal"), http.StatusInternalServerError},
	{UnmarshalError.New("unmarshal"), http.StatusBadRequest},
	{errs.New("unknown"), http.StatusInternalServerError},
	{nil, http.StatusInternalServerError},
}

func TestStatusFromErrors(t *testing.T) {
	for _, v := range statusFromErrorTests {
		require.Equal(t, v.expected, statusFromError(v.err))
	}
}
