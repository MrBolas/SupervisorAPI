package models

import (
	"database/sql"

	"github.com/gofrs/uuid"
)

type Task struct {
	Id       uuid.UUID    `gorm:"primary_key;default:UUID()"`
	WorkerId string       `gorm:"column:workerId"`
	Summary  string       `gorm:"column:summary"`
	Date     sql.NullTime `gorm:"column:date"`
}

func (t *Task) ToResponse() TaskResponse {
	return TaskResponse{
		Id:       t.Id,
		WorkerId: t.WorkerId,
		Summary:  t.Summary,
		Date:     t.Date.Time.String(),
	}
}
