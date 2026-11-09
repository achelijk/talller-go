package main

import (
	"bufio"
	"fmt"
	"os"
	"sistema-control-escolar/service"
	"strconv"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)

func main() {
	svc := service.NewSchoolService()

	for {
		fmt.Println("\n=================================")
		fmt.Println("    SISTEMA DE CONTROL ESCOLAR   ")
		fmt.Println("=================================")
		fmt.Println("1. Registrar Alumno")
		fmt.Println("2. Registrar Materia")
		fmt.Println("3. Registrar/Actualizar Calificación")
		fmt.Println("4. Consultar Boleta de Alumno")
		fmt.Println("5. Listar Alumnos por Promedio Mínimo")
		fmt.Println("6. Reporte por Materia")
		fmt.Println("7. Reporte General (Concurrente)")
		fmt.Println("8. Salir")
		fmt.Print("Seleccione una opción: ")

		scanner.Scan()
		option := strings.TrimSpace(scanner.Text())

		switch option {
		case "1":
			registerStudent(svc)
		case "2":
			registerSubject(svc)
		case "3":
			registerGrade(svc)
		case "4":
			showReportCard(svc)
		case "5":
			listStudentsByMinAverage(svc)
		case "6":
			showSubjectReport(svc)
		case "7":
			showGeneralReport(svc)
		case "8":
			fmt.Println("¡Hasta luego!")
			return
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
		}
	}
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func registerStudent(svc *service.SchoolService) {
	fmt.Println("\n--- Registrar Alumno ---")
	id := readInput("ID: ")
	name := readInput("Nombre: ")
	lastName := readInput("Apellido: ")
	email := readInput("Email: ")

	err := svc.AddStudent(id, name, lastName, email)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("✅ Alumno registrado correctamente.")
}

func registerSubject(svc *service.SchoolService) {
	fmt.Println("\n--- Registrar Materia ---")
	id := readInput("ID: ")
	name := readInput("Nombre de Materia: ")
	creditsStr := readInput("Créditos (> 0): ")
	teacher := readInput("Profesor: ")

	credits, err := strconv.Atoi(creditsStr)
	if err != nil {
		fmt.Println("Error: los créditos deben ser un número entero.")
		return
	}

	err = svc.AddSubject(id, name, credits, teacher)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("✅ Materia registrada correctamente.")
}

func registerGrade(svc *service.SchoolService) {
	fmt.Println("\n--- Registrar/Actualizar Calificación ---")
	studentID := readInput("ID Alumno: ")
	subjectID := readInput("ID Materia: ")
	valStr := readInput("Calificación (0 - 10): ")

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		fmt.Println("Error: ingrese un número válido.")
		return
	}

	err = svc.AddOrUpdateGrade(studentID, subjectID, val)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("✅ Calificación registrada/actualizada exitosamente.")
}

func showReportCard(svc *service.SchoolService) {
	studentID := readInput("ID Alumno: ")

	var studentName string
	found := false
	for _, st := range svc.Students {
		if st.ID == studentID {
			studentName = st.FirstName + " " + st.LastName
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Alumno no encontrado.")
		return
	}

	fmt.Printf("\n--- Boleta de %s (ID: %s) ---\n", studentName, studentID)
	var total float64
	var count int

	for _, g := range svc.Grades {
		if g.StudentID == studentID {
			materiaName := g.SubjectID
			for _, sub := range svc.Subjects {
				if sub.ID == g.SubjectID {
					materiaName = sub.Name
					break
				}
			}
			fmt.Printf("- %s: %.2f\n", materiaName, g.Value)
			total += g.Value
			count++
		}
	}

	if count == 0 {
		fmt.Println("El alumno no tiene calificaciones registradas.")
	} else {
		fmt.Printf("Promedio General: %.2f\n", total/float64(count))
	}
}

func listStudentsByMinAverage(svc *service.SchoolService) {
	minAvgStr := readInput("Ingrese el promedio mínimo a filtrar: ")
	minAvg, err := strconv.ParseFloat(minAvgStr, 64)
	if err != nil {
		fmt.Println("Error: valor numérico no válido.")
		return
	}

	fmt.Printf("\n--- Alumnos con promedio >= %.2f ---\n", minAvg)
	for _, st := range svc.Students {
		avg := svc.CalculateStudentAverage(st.ID)
		if avg >= minAvg {
			fmt.Printf("- [%s] %s %s | Promedio: %.2f\n", st.ID, st.FirstName, st.LastName, avg)
		}
	}
}

func showSubjectReport(svc *service.SchoolService) {
	subjectID := readInput("ID Materia: ")

	var subjectName string
	found := false
	for _, sub := range svc.Subjects {
		if sub.ID == subjectID {
			subjectName = sub.Name
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Materia no encontrada.")
		return
	}

	var total float64
	var count int
	var maxVal float64 = -1
	var minVal float64 = 11

	for _, g := range svc.Grades {
		if g.SubjectID == subjectID {
			total += g.Value
			count++
			if g.Value > maxVal {
				maxVal = g.Value
			}
			if g.Value < minVal {
				minVal = g.Value
			}
		}
	}

	fmt.Printf("\n--- Reporte de Materia: %s ---\n", subjectName)
	fmt.Printf("Alumnos inscritos/evaluados: %d\n", count)
	if count > 0 {
		fmt.Printf("Calificación más alta: %.2f\n", maxVal)
		fmt.Printf("Calificación más baja: %.2f\n", minVal)
		fmt.Printf("Promedio de la materia: %.2f\n", total/float64(count))
	} else {
		fmt.Println("Sin calificaciones registradas para esta materia.")
	}
}

func showGeneralReport(svc *service.SchoolService) {
	report := svc.GenerateGeneralReportConcurrently()
	fmt.Println("\n--- Reporte General del Sistema ---")
	fmt.Printf("Total de Alumnos: %d\n", report.TotalStudents)
	fmt.Printf("Total de Materias: %d\n", report.TotalSubjects)
	fmt.Printf("Promedio Global: %.2f\n", report.GlobalAverage)
}	