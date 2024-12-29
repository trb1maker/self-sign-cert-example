package server

import (
	"encoding/json"
	"net/http"

	"github.com/trb1maker/self-sign-cert-example/internal/http/dto"
)

func timeHandler(w http.ResponseWriter, _ *http.Request) {
	ts := dto.NewTimeResponse()

	if err := json.NewEncoder(w).Encode(ts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
