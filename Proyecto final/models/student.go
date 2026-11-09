package models

import (
	"errors"
	"strings"
	"time"
)

type Student struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Student) Validate() error {
	if strings.TrimSpace(s.FirstName) == "" || strings.TrimSpace(s.LastName) == "" {
		return errors.New("el nombre y el apellido no pueden estar vacíos")
	}
	if !strings.Contains(s.Email, "@") || !strings.Contains(s.Email, ".") {
		return errors.New("el formato del email no es válido")
	}
	return nil
}