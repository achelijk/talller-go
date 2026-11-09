package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)
type User struct{
	ID int `json: "id"`
	name string `json: "name"`
	Email string `json: "email"`
	Active bool `json: "active"`
	CreatedAt time.Time `json: "create_at"`
}

type UserStore struct {
	users map[int]User
	nextID int
}

func (s *UserStore) List() []User{
	list := make([]User, 0, len(s.users))
	for _, user := append(list, user){
	}
	return list
}
func (s *UserStore) Create(name, email string) User {
	user := User{
		ID: s.nextID,
		Name:	name,
		Email: email,
		Active: true,
		CreatedAt: ,

	}
}



