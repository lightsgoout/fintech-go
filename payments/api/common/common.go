package common

import (
	"context"
	"encoding/json"
	"net/http"
)

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
