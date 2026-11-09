package models

import (
	"errors"
	"time"
)

type Grade struct {
	StudentID string    `json:"student_id"`
	SubjectID string    `json:"subject_id"`
	Value     float64   `json:"value"`
	Date      time.Time `json:"date"`
}

func (g *Grade) Validate() error {
	if g.Value < 0 || g.Value > 10 {
		return errors.New("la calificación debe estar entre 0 y 10")
	}
	return nil
}