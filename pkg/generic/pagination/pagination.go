package pagination

import (
	"fmt"
	"math"
)

type Links struct {
	First    string `json:"first"`
	Previous string `json:"previous"`
	Self     string `json:"self"`
	Next     string `json:"next"`
	Last     string `json:"last"`
}

type Pagination struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
	PageCount  int    `json:"pageCount"`
	TotalCount int    `json:"totalCount"`
	Sort       string `json:"sort"`
	Links      Links  `json:"links"`
}

func GetPaginationMetaData(totalCount int, currentPage int, perPage int, sortValue string, searchValue string) Pagination {
	pageCount := getLastPageCount(totalCount, perPage)
	searchParam := fmt.Sprintf("searchParam=%s", searchValue)
	perPageParam := fmt.Sprintf("&perPage=%d", perPage)
	sortParam := fmt.Sprintf("&sort=%s", getSortDirection(sortValue))
	links := getLinks(searchParam, currentPage, pageCount, perPageParam, sortParam)

	return Pagination{currentPage, perPage, pageCount, totalCount, sortValue, links}
}

func getLinks(searchParam string, currentPage int, pageCount int, perPageParam string, sortParam string) Links {
	previousPage := currentPage - 1
	if currentPage == 1 {
		previousPage = currentPage
	}
	const stringFragment = "?%s&page=%d%s%s"
	first := fmt.Sprintf(stringFragment, searchParam, 1, perPageParam, sortParam)
	previous := fmt.Sprintf(stringFragment, searchParam, previousPage, perPageParam, sortParam)
	self := fmt.Sprintf(stringFragment, searchParam, currentPage, perPageParam, sortParam)
	next := fmt.Sprintf(stringFragment, searchParam, currentPage+1, perPageParam, sortParam)
	last := fmt.Sprintf(stringFragment, searchParam, pageCount, perPageParam, sortParam)

	return Links{first, previous, self, next, last}
}

type SortDirection string

const (
	Ascending  SortDirection = "asc"
	Descending SortDirection = "desc"
)

func getSortDirection(sort string) string {
	if SortDirection(sort) == Descending {
		return "desc"
	}
	return "asc"
}

func GetSortOptionType(sort string) int {
	if SortDirection(sort) == Descending {
		return -1
	}
	return 1
}

func getLastPageCount(totalCount int, maxPerPage int) int {
	if totalCount <= GetMaxPerPage(maxPerPage) {
		return 1
	}
	return int(math.Ceil(float64(totalCount) / float64(GetMaxPerPage(maxPerPage))))
}

func GetCurrentPage(page int) int {
	if page < 1 {
		return 1
	}
	return page
}

func GetMaxPerPage(perPage int) int {
	if perPage <= 0 {
		return 100
	} else if perPage > 250 {
		return 250
	} else {
		return perPage
	}
}
