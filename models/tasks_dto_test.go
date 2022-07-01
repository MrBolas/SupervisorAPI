package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeFormatValidation(t *testing.T) {

	taskr := TaskRequest{
		Summary: "mock_summary",
		Date:    "2006-01-02 03:04:05PM",
	}

	err := taskr.Validate()
	assert.Nil(t, err)

	taskr = TaskRequest{
		Summary: "mock_summary",
		Date:    "2006-01-02 03:04:05",
	}

	err = taskr.Validate()
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("invalid date format, use yyyy-mm-dd hh:mm:ssPM"), err)
	}
}

func TestSummaryConstraintValidation(t *testing.T) {

	taskr := TaskRequest{
		Summary: "mock_summary",
		Date:    "2006-01-02 03:04:05PM",
	}

	err := taskr.Validate()
	assert.Nil(t, err)

	taskr = TaskRequest{
		Summary: string(make([]byte, 2501)),
		Date:    "2006-01-02 03:04:05PM",
	}

	err = taskr.Validate()
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("summary max size is 2500 characters"), err)
	}
}

func TestTaskRequestToTask(t *testing.T) {

	tr := TaskRequest{
		Summary: "mock_request",
		Date:    "2006-01-02 03:04:05PM",
	}

	task, err := tr.ToTask("mocked_worker_id")
	assert.Nil(t, err)

	assert.Equal(t, task.Summary, tr.Summary)
	assert.Equal(t, task.WorkerId, "mocked_worker_id")
	assert.NotNil(t, task.Date.Time)
	assert.NotNil(t, task.Id)
}

func TestToTaskResponseList(t *testing.T) {

	tr := TaskRequest{
		Summary: "mock_request",
		Date:    "2006-01-02 03:04:05PM",
	}

	tasks := make([]Task, 0)

	for i := 0; i < 5; i++ {
		task, err := tr.ToTask("mocked_worker_id")
		assert.Nil(t, err)

		tasks = append(tasks, task)
	}

	tRespList := ToListResponse(tasks, 1, 10)

	assert.Equal(t, len(tRespList.Data), 5)
	assert.Equal(t, tRespList.Metadata.Page, 1)
	assert.Equal(t, tRespList.Metadata.PageSize, 10)
}
