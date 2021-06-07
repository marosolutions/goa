package goa

import (
	"io"
	"net/http"
)

// Router ...
func (config *Config) Router(w http.ResponseWriter, req *http.Request) {
	// log.Printf("%v %v\n", req.Method, req.URL.Path)

	if req.URL.Path != config.Service.Path {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{\"errors\": {\"base\": [\"invalid request path\"]}}")
		return
	}

	// Set Content Type and CORS headers
	config.SetHeaders(w, req)

	// NOTE: req.Method and config.Service.Method should be converted to same case to match
	switch req.Method {
	case config.Service.Method:
		// Valid request to be processed by the Controller method
		service.Controller(config, w, req)
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	default:
		// Return a bad request error when the request method is invalid
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{\"errors\": {\"base\": [\"invalid request method\"]}}")
	}
}
