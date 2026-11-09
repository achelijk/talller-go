package main

import (
	"fmt"
	"strings"
)

func main() {
	var username, email, password string
	var age int

	// Pedir datos al usuario
	fmt.Print("Ingresa tu username: ")
	fmt.Scanln(&username)

	fmt.Print("Ingresa tu email: ")
	fmt.Scanln(&email)

	fmt.Print("Ingresa tu password: ")
	fmt.Scanln(&password)

	fmt.Print("Ingresa tu edad: ")
	fmt.Scanln(&age)

	var errors []string

	// 1. Validar username
	if len(username) < 4 || len(username) > 20 {
		errors = append(errors, "username debe tener entre 4 y 20 caracteres")
	}

	// 2. Validar email
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		errors = append(errors, "email debe contener '@' y '.'")
	}

	// 3. Validar password
	if len(password) < 8 {
		errors = append(errors, "password debe tener al menos 8 caracteres")
	}

	// 4. Validar age
	if age < 18 {
		errors = append(errors, "debes ser mayor o igual a 18 años")
	}

	// Mostrar resultados
	fmt.Println()
	if len(errors) == 0 {
		fmt.Println("✅ Registro exitoso")
	} else {
		fmt.Println("❌ Error:")
		for _, err := range errors {
			fmt.Printf(" - %s\n", err)
		}
	}
}
