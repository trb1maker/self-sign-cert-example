package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trb1maker/self-sign-cert-example/internal/http/dto"
)

func TestHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	timeHandler(w, r)

	ts := new(dto.TimeResponse)
	require.NoError(t, json.NewDecoder(w.Body).Decode(ts))
	require.NotEmpty(t, ts.Time)

	require.Equal(t, ts.Time().Year(), time.Now().Year())
}
