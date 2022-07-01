package repositories

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

const DEFAULT_PAGE = 1
const DEFAULT_PAGE_SIZE = 20
const MAX_PAGE_SIZE = 40

type ListQuery struct {
	Filters    map[string]interface{}
	Sort       Sort
	Pagination Pagination
}

type Sort struct {
	By    string
	Order string
}

type Pagination struct {
	Page     int
	PageSize int
}

func NewListQuery() ListQuery {
	return ListQuery{
		Filters: make(map[string]interface{}),
		Sort: Sort{
			By:    "",
			Order: "",
		},
		Pagination: Pagination{
			Page:     0,
			PageSize: 0,
		},
	}
}

func (lq *ListQuery) AddListTaskFilters(queryParamaters url.Values, isManager bool) error {

	for key, val := range queryParamaters {

		// the value is an array of data, only use if something is there
		if len(val) == 0 {
			continue
		}

		if isManager {
			if key == "worker_name" {
				lq.Filters["worker_name"] = val[0]
			}
		}

		// validate date
		if key == "before" {
			lq.Filters["before"] = val[0]
		}

		// validate date
		if key == "after" {
			lq.Filters["after"] = val[0]
		}
	}

	return nil
}

func (lq *ListQuery) AddPageAndPageSize(pageParam string, pageSizeParam string) error {

	page := DEFAULT_PAGE
	pageSize := DEFAULT_PAGE_SIZE
	var err error

	if pageParam != "" {
		page, err = strconv.Atoi(pageParam)
		if err != nil || page <= 0 {
			return errors.New("page must be bigger than 0")
		}
	}

	if pageSizeParam != "" {
		pageSize, err = strconv.Atoi(pageSizeParam)
		if err != nil || pageSize < 1 {
			return errors.New("page_size must be bigger than 0")
		}
	}

	if pageSize > MAX_PAGE_SIZE {
		return fmt.Errorf("page_size must be less or equal than %d", MAX_PAGE_SIZE)
	}

	lq.Pagination.Page = page
	lq.Pagination.PageSize = pageSize

	return nil
}

func (lq *ListQuery) AddSorting(sortBy string, order string) error {

	if order != "asc" && order != "desc" && order != "" {
		return errors.New("order must be desc or asc")
	}

	dbOrder := "desc"
	dbSortBy := "date"

	if order == "asc" {
		dbOrder = "asc"
	}

	if sortBy == "name" {
		dbSortBy = "worker_name"
	}

	lq.Sort.By = dbSortBy
	lq.Sort.Order = dbOrder

	return nil
}

func (lq *ListQuery) GetOffsetLimit() (int, int) {
	return (lq.Pagination.Page - 1) * lq.Pagination.PageSize, lq.Pagination.PageSize + 1 // this plus one is used for pagination
}
