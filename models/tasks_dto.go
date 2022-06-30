package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

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

	// validate number of charecters of summary?

	return nil
}

type TaskResponse struct {
	Id       uuid.UUID `json:"id"`
	WorkerId string    `json:"worker_id"`
	Summary  string    `json:"summary"`
	Date     string    `json:"date"`
}
