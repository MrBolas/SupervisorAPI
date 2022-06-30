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
		time.RFC3339,
		tr.Date)
	if err != nil {
		return Task{}, errors.New("invalid date format, use yyyy-mm-ddThh:mm:ss as per rfc3339")
	}

	return Task{
		Summary:  tr.Summary,
		WorkerId: workerId,
		Date: sql.NullTime{
			Valid: true,
			Time:  t,
		},
	}, nil
}

type TaskResponse struct {
	Id       uuid.UUID `json:"id"`
	WorkerId string    `json:"worker_id"`
	Summary  string    `json:"summary"`
	Date     string    `json:"date"`
}
