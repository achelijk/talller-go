package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type UserStore struct {
	users  map[int]User
	nextId int
}

func (s *UserStore) List() []User {
	list := make([]User, 0, len(s.users))
	for _, user := range s.users {
		list = append(list, user)
	}
	return list
}

func (s *UserStore) GetUserByID(id int) (User, error) {
	user, exist := s.users[id]
	if !exist {
		return User{}, errors.New("Usuario no encontrado")
	}
	return user, nil
}

func (s *UserStore) Create(name, email string) User {
	user := User{
		ID:        s.nextId,
		Name:      name,
		Email:     email,
		Active:    true,
		CreatedAt: time.Now(),
	}
	s.users[user.ID] = user
	s.nextId++
	return user
}

// PUT reemplaza los campos editables del usuario.
func (s *UserStore) Update(id int, name, email string, active bool) (User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return User{}, err
	}
	user.Name = name
	user.Email = email
	user.Active = active
	s.users[id] = user
	return user, nil
}

// PATCH modifica únicamente los campos enviados.
func (s *UserStore) Patch(id int, name, email *string, active *bool) (User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return User{}, err
	}
	if name != nil {
		user.Name = *name
	}
	if email != nil {
		user.Email = *email
	}
	if active != nil {
		user.Active = *active
	}
	s.users[id] = user
	return user, nil
}

func (s *UserStore) Delete(id int) error {
	if _, err := s.GetUserByID(id); err != nil {
		return err
	}
	delete(s.users, id)
	return nil
}

func NewStore() *UserStore {
	return &UserStore{
		users:  make(map[int]User),
		nextId: 1,
	}
}

// responder con json
func JsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// extraer ID del path
func extractPathId(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		return 0, errors.New("Ruta invalida")
	}
	return strconv.Atoi(parts[len(parts)-1])
}

type UserHandler struct {
	store *UserStore
}

// GET /api/users
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users := h.store.List()
	JsonResponse(w, http.StatusOK, users)
}

// GET /api/users/{id}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractPathId(r.URL.Path)
	if err != nil {
		http.Error(w, "ID no valido", http.StatusBadRequest)
		return
	}
	user, err := h.store.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	JsonResponse(w, http.StatusOK, user)
}

// POST /api/users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON no valido", http.StatusBadRequest)
		return
	}
	if input.Email == "" {
		http.Error(w, "Email es requerido", http.StatusBadRequest)
		return
	}
	if input.Name == "" {
		http.Error(w, "Name es requerido", http.StatusBadRequest)
		return
	}
	user := h.store.Create(input.Name, input.Email)
	JsonResponse(w, http.StatusCreated, user)
}

// PUT /api/users/{id}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := extractPathId(r.URL.Path)
	if err != nil {
		http.Error(w, "ID no valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var input struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Active bool   `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON no valido", http.StatusBadRequest)
		return
	}
	if input.Name == "" || input.Email == "" {
		http.Error(w, "Name y Email son requeridos", http.StatusBadRequest)
		return
	}
	user, err := h.store.Update(id, input.Name, input.Email, input.Active)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	JsonResponse(w, http.StatusOK, user)
}

// PATCH /api/users/{id}
func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id, err := extractPathId(r.URL.Path)
	if err != nil {
		http.Error(w, "ID no valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var input struct {
		Name   *string `json:"name"`
		Email  *string `json:"email"`
		Active *bool   `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON no valido", http.StatusBadRequest)
		return
	}
	if input.Name == nil && input.Email == nil && input.Active == nil {
		http.Error(w, "Debe enviar al menos un campo", http.StatusBadRequest)
		return
	}
	if input.Name != nil && *input.Name == "" {
		http.Error(w, "Name no puede estar vacio", http.StatusBadRequest)
		return
	}
	if input.Email != nil && *input.Email == "" {
		http.Error(w, "Email no puede estar vacio", http.StatusBadRequest)
		return
	}
	user, err := h.store.Patch(id, input.Name, input.Email, input.Active)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	JsonResponse(w, http.StatusOK, user)
}

// DELETE /api/users/{id}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := extractPathId(r.URL.Path)
	if err != nil {
		http.Error(w, "ID no valido", http.StatusBadRequest)
		return
	}
	if err := h.store.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	store := NewStore()
	handler := &UserHandler{store}

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.List(w, r)
		case http.MethodPost:
			handler.Create(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetUserByID(w, r)
		case http.MethodPut:
			handler.Update(w, r)
		case http.MethodPatch:
			handler.Patch(w, r)
		case http.MethodDelete:
			handler.Delete(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	fmt.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
