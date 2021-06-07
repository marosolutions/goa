package goa

import (
	"net/http"
)

// Router ...
func (config *Config) Router(w http.ResponseWriter, req *http.Request) {
	// log.Printf("%v %v\n", req.Method, req.URL.Path)

	if req.URL.Path != config.Service.Path {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Set Content Type and CORS headers
	config.SetHeaders(w, req)

	switch req.Method {
	case config.Service.Method:
		// Valid request to be processed by the Controller method
		service.Controller(config, w, req)
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	default:
		// Return a bad request error when the request method is invalid
		w.WriteHeader(http.StatusBadRequest)
	}
}
