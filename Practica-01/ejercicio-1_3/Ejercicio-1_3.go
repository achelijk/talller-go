package main

import (
	"fmt"
	"time"
)

func main() {

	
	var timestamp int64 = 1734998400

	
	fecha := time.Unix(timestamp, 0).UTC()

	fechaString := fecha.Format("2006-01-02T15:04:05Z")

	var check rune = '\u2705'

	fmt.Printf("%c Conversión exitosa\n", check)
	fmt.Printf("Timestamp: %d (%T)\n", timestamp, timestamp)
	fmt.Printf("Fecha: %v (%T)\n", fecha, fecha)
	fmt.Printf("ISO 8601: %s (%T)\n", fechaString, fechaString)
}
