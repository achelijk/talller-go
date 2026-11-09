package service

import (
	"errors"
	"sistema-control-escolar/models"
	"sistema-control-escolar/storage"
	"strings"
	"sync"
	"time"
)

type SchoolService struct {
	Students []models.Student
	Subjects []models.Subject
	Grades   []models.Grade
	mu       sync.RWMutex
}

func NewSchoolService() *SchoolService {
	s := &SchoolService{}
	s.LoadAll()
	return s
}

func (s *SchoolService) LoadAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = storage.LoadData("students.json", &s.Students)
	_ = storage.LoadData("subjects.json", &s.Subjects)
	_ = storage.LoadData("grades.json", &s.Grades)
}

func (s *SchoolService) SaveAll() {
	_ = storage.SaveData("students.json", s.Students)
	_ = storage.SaveData("subjects.json", s.Subjects)
	_ = storage.SaveData("grades.json", s.Grades)
}

// 1. Alta Alumno
func (s *SchoolService) AddStudent(id, firstName, lastName, email string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, student := range s.Students {
		if student.ID == id {
			return errors.New("ya existe un alumno con ese ID")
		}
		if strings.EqualFold(student.Email, email) {
			return errors.New("el email ya está registrado")
		}
	}

	student := models.Student{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		CreatedAt: time.Now(),
	}

	if err := student.Validate(); err != nil {
		return err
	}

	s.Students = append(s.Students, student)
	s.SaveAll()
	return nil
}

// 2. Alta Materia
func (s *SchoolService) AddSubject(id, name string, credits int, teacher string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sub := range s.Subjects {
		if sub.ID == id {
			return errors.New("ya existe una materia con ese ID")
		}
		if strings.EqualFold(sub.Name, name) {
			return errors.New("ya existe una materia con ese nombre")
		}
	}

	subject := models.Subject{
		ID:      id,
		Name:    name,
		Credits: credits,
		Teacher: teacher,
	}

	if err := subject.Validate(); err != nil {
		return err
	}

	s.Subjects = append(s.Subjects, subject)
	s.SaveAll()
	return nil
}

// 3. Alta/Actualización Calificación
func (s *SchoolService) AddOrUpdateGrade(studentID, subjectID string, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	studentExists := false
	for _, st := range s.Students {
		if st.ID == studentID {
			studentExists = true
			break
		}
	}
	if !studentExists {
		return errors.New("el alumno especificado no existe")
	}

	subjectExists := false
	for _, sub := range s.Subjects {
		if sub.ID == subjectID {
			subjectExists = true
			break
		}
	}
	if !subjectExists {
		return errors.New("la materia especificada no existe")
	}

	grade := models.Grade{
		StudentID: studentID,
		SubjectID: subjectID,
		Value:     value,
		Date:      time.Now(),
	}

	if err := grade.Validate(); err != nil {
		return err
	}

	for i, g := range s.Grades {
		if g.StudentID == studentID && g.SubjectID == subjectID {
			s.Grades[i].Value = value
			s.Grades[i].Date = time.Now()
			s.SaveAll()
			return nil
		}
	}

	s.Grades = append(s.Grades, grade)
	s.SaveAll()
	return nil
}

// 4. Consultas y Reportes
func (s *SchoolService) CalculateStudentAverage(studentID string) float64 {
	var total float64
	var count int

	for _, g := range s.Grades {
		if g.StudentID == studentID {
			total += g.Value
			count++
		}
	}

	if count == 0 {
		return 0.0
	}
	return total / float64(count)
}

// Reporte Concurrente usando Goroutines y Channels
type GeneralReport struct {
	TotalStudents int
	TotalSubjects int
	GlobalAverage float64
}

func (s *SchoolService) GenerateGeneralReportConcurrently() GeneralReport {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var wg sync.WaitGroup
	report := GeneralReport{
		TotalStudents: len(s.Students),
		TotalSubjects: len(s.Subjects),
	}

	avgChan := make(chan float64, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		if len(s.Grades) == 0 {
			avgChan <- 0.0
			return
		}
		var sum float64
		for _, g := range s.Grades {
			sum += g.Value
		}
		avgChan <- sum / float64(len(s.Grades))
	}()

	wg.Wait()
	close(avgChan)
	report.GlobalAverage = <-avgChan

	return report
}