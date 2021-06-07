package goapi

import "net/http"

// SetHeaders ...
func (config Config) SetHeaders(w http.ResponseWriter, req *http.Request) {
	// Set Content Type header to respond with a json response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Set CORS Headers
	w.Header().Set("Access-Control-Allow-Origin", config.Cors.AllowOrigins)
	w.Header().Set("Access-Control-Allow-Methods", config.Service.Method+", options")
	w.Header().Set("Access-Control-Allow-Headers", config.Cors.AllowHeaders)
}
