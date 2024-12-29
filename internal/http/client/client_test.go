package client

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trb1maker/self-sign-cert-example/internal/http/dto"
)

func TestUnmarshal(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	require.NoError(t, json.NewEncoder(buf).Encode(dto.NewTimeResponse()))

	ts := new(dto.TimeResponse)
	require.NoError(t, json.NewDecoder(buf).Decode(ts))
	require.NotEmpty(t, ts.TS)
}
