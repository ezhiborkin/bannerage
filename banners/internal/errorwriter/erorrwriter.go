package errorwriter

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONError struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, error string, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(JSONError{Error: error})
	if err != nil {
		fmt.Printf("%s", err)
	}
}
