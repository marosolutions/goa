package goa

import (
	"strconv"
)

var MaxPerPage = 50
var PerPage = 10

// Pagination ...
func Pagination(paramsPage string, paramsLimit string) (limit int, offset int) {
	// Per page records limit
	limit = PerPage

	if paramsLimit != "" {
		limit, _ = strconv.Atoi(paramsLimit)
	}

	// Set maximum per page limit
	if limit > MaxPerPage {
		limit = MaxPerPage
	}

	// Page number requested
	pageNumber := 1

	if paramsPage != "" {
		pageNumber, _ = strconv.Atoi(paramsPage)
	}

	// Set pageNumber to first page by default
	if pageNumber <= 0 {
		pageNumber = 1
	}

	offset = (pageNumber - 1) * limit
	return
}
