package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

const MAX_SUMMARY_CHARS = 2500

type TaskRequest struct {
	Summary string `json:"summary"`
	Date    string `json:"date"`
}

func (tr *TaskRequest) ToTask(workerId string) (Task, error) {

	t, err := time.Parse(
		"2006-01-02 03:04:05PM",
		tr.Date)
	if err != nil {
		return Task{}, errors.New("invalid date format, use yyyy-mm-dd hh:mm:ssPM")
	}

	genUuid, err := uuid.NewV4()
	if err != nil {
		return Task{}, errors.New("Task Id generation error")
	}

	return Task{
		Id:       genUuid,
		Summary:  tr.Summary,
		WorkerId: workerId,
		Date: sql.NullTime{
			Valid: true,
			Time:  t,
		},
	}, nil
}

func (tr *TaskRequest) Validate() error {

	// validate date
	_, err := time.Parse(
		"2006-01-02 03:04:05PM",
		tr.Date)
	if err != nil {
		return errors.New("invalid date format, use yyyy-mm-dd hh:mm:ssPM")
	}

	// validate number of charecters of summary
	if len(tr.Summary) > MAX_SUMMARY_CHARS {
		return errors.New("summary max size is 2500 characters")
	}

	return nil
}

type TaskResponse struct {
	Id       uuid.UUID `json:"id"`
	WorkerId string    `json:"worker_name"`
	Summary  string    `json:"summary"`
	Date     string    `json:"date"`
}

type TaskListResponse struct {
	Data     []TaskResponse `json:"data"`
	Metadata Metadata       `json:"metadata"`
}

type Metadata struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func ToListResponse(tasks []Task, page int, pageSize int) TaskListResponse {

	tasksResponse := make([]TaskResponse, 0)

	if len(tasks) > pageSize {
		tasks = tasks[:len(tasks)-1]
	}

	for _, t := range tasks {
		tasksResponse = append(tasksResponse, t.ToResponse())
	}

	return TaskListResponse{
		Data: tasksResponse,
		Metadata: Metadata{
			Page:     page,
			PageSize: pageSize,
		},
	}
}
