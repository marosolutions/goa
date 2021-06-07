package goa

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// ErrorResponse ...
// Example: {errors: {name: ["can't be blank"], email: ["can't be blank", "already taken"]}}
type ErrorResponse struct {
	Errors map[string][]string `json:"errors"`
}

func WriteErrorResponse(w http.ResponseWriter, err error, status int) {
	log.Println("writeErrorResponse  : ", err)

	resp := ErrorResponse{
		Errors: map[string][]string{
			"base": {err.Error()},
		},
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println("ErrorResponse json.Marshal Error:", err)
	}

	w.WriteHeader(status)
	io.WriteString(w, string(jsonResp))
}
