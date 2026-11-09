package models

import (
	"errors"
	"strings"
)

type Subject struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Credits int    `json:"credits"`
	Teacher string `json:"teacher"`
}

func (s *Subject) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("el nombre de la materia no puede estar vacío")
	}
	if s.Credits <= 0 {
		return errors.New("los créditos deben ser mayores a 0")
	}
	return nil
}