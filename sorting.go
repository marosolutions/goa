package goa

import (
	"net/http"
	"strings"
)

// Sorting ...
func (config *Config) Sorting(req *http.Request) (sortOrder string, sortBy string) {
	urlQuery := req.URL.Query()

	validSortOrders := map[string]bool{
		"desc": true,
		"asc":  true,
	}
	sortOrder = strings.ToLower(urlQuery.Get("sort_order"))

	// Set default sort order and sort by
	if !validSortOrders[sortOrder] {
		sortOrder = config.Sort.DefaultSortOrder
	}

	sortBy = strings.ToLower(urlQuery.Get("sort_by"))
	if !SliceContains(config.Sort.ValidSortBy, sortBy) {
		sortBy = config.Sort.DefaultSortBy
	}

	return
}
