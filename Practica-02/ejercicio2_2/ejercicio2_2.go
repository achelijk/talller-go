package main

import (
	"fmt"
)

func clasificarHTTP(status int) {
	switch status / 100 {
	case 2:
		fmt.Printf("HTTP %d → Éxito\n", status)
	case 3:
		fmt.Printf("HTTP %d → Redirección\n", status)
	case 4:
		fmt.Printf("HTTP %d → Error del cliente\n", status)
	case 5:
		fmt.Printf("HTTP %d → Error del servidor\n", status)
	default:
		fmt.Printf("HTTP %d → Código desconocido\n", status)
	}
}

func main() {
	var code int

	// Pedir el código HTTP al usuario
	fmt.Print("Ingresa un código de estado HTTP (ej. 200, 404, 500): ")
	_, err := fmt.Scanln(&code)

	if err != nil {
		fmt.Println("❌ Error: Debes ingresar un número entero válido.")
		return
	}

	// Mostrar la clasificación
	fmt.Println()
	clasificarHTTP(code)
}
