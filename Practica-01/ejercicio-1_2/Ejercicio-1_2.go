package main

import (
	"fmt"
	"math"
)

func main() {

	
	precio1 := 19.99
	precio2 := 5.50
	precio3 := 12.00

	
	centavos1 := int64(math.Round(precio1 * 100))
	centavos2 := int64(math.Round(precio2 * 100))
	centavos3 := int64(math.Round(precio3 * 100))

	cantidad1 := int64(2)
	cantidad2 := int64(1)
	cantidad3 := int64(3)

	totalCentavos := centavos1*cantidad1 +
		centavos2*cantidad2 +
		centavos3*cantidad3

	total := float64(totalCentavos) / 100

	
	iva := total * 0.16
	totalFinal := total + iva

	fmt.Printf("Total sin IVA: $%.2f\n", total)
	fmt.Printf("IVA (16%%): $%.2f\n", iva)
	fmt.Printf("Total final: $%.2f\n", totalFinal)
}
