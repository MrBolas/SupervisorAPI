package repositories

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDefaultQuery(t *testing.T) {
	query := NewListQuery()
	err := query.AddPageAndPageSize("", "")
	assert.Nil(t, err)
	err = query.AddSorting("", "")
	assert.Nil(t, err)

	assert.Equal(t, query.Filters, make(map[string]interface{}))
	assert.Equal(t, query.IntervalFilters, make(map[string]interface{}))
	assert.Equal(t, query.Pagination.Page, 1)
	assert.Equal(t, query.Pagination.PageSize, 20)
	assert.Equal(t, query.Sort.By, "date")
	assert.Equal(t, query.Sort.Order, "desc")
}

func TestPaginationErrorHandling(t *testing.T) {
	query := NewListQuery()
	err := query.AddPageAndPageSize("", "")
	assert.Nil(t, err)

	err = query.AddPageAndPageSize("0", "")
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("page must be bigger than 0"), err)
	}

	err = query.AddPageAndPageSize("", "0")
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("page_size must be bigger than 0"), err)
	}

	err = query.AddPageAndPageSize("", "41")
	if assert.Error(t, err) {
		assert.Equal(t, fmt.Errorf("page_size must be less or equal than %d", MAX_PAGE_SIZE), err)
	}
}

func TestSortingOptions(t *testing.T) {
	query := NewListQuery()
	err := query.AddPageAndPageSize("", "")
	assert.Nil(t, err)
	err = query.AddSorting("", "")
	assert.Nil(t, err)
	assert.Equal(t, query.Sort.By, "date")
	assert.Equal(t, query.Sort.Order, "desc")

	err = query.AddSorting("name", "asc")
	assert.Nil(t, err)
	assert.Equal(t, query.Sort.By, "worker_name")
	assert.Equal(t, query.Sort.Order, "asc")
}

func TestGetOffsetLimit(t *testing.T) {
	query := NewListQuery()
	err := query.AddPageAndPageSize("1", "20")
	assert.Nil(t, err)

	offset, limit := query.GetOffsetLimit()

	assert.Equal(t, query.Pagination.Page, 1)
	assert.Equal(t, query.Pagination.PageSize, 20)
	assert.Equal(t, offset, 0)
	assert.Equal(t, limit, 21)
}
