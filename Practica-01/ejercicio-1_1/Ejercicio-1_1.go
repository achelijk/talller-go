package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {

	const DB_HOST = "localhost"
	const DB_PORT = "5432"
	const DB_MAX_CONNS = "25"
	const DB_USE_SSL = "true"
	const DB_TIMEOUT_MS = "5000"

	host := DB_HOST

	port, err := strconv.Atoi(DB_PORT)
	if err != nil {
		fmt.Println("Error en DB_PORT:", err)
		port = 5432
	}

	maxConns, err := strconv.Atoi(DB_MAX_CONNS)
	if err != nil {
		fmt.Println("Error en DB_MAX_CONNS:", err)
		maxConns = 10
	}

	useSSL, err := strconv.ParseBool(DB_USE_SSL)
	if err != nil {
		fmt.Println("Error en DB_USE_SSL:", err)
		useSSL = false
	}

	timeoutMS, err := strconv.ParseInt(DB_TIMEOUT_MS, 10, 64)
	if err != nil {
		fmt.Println("Error en DB_TIMEOUT_MS:", err)
		timeoutMS = 5000
	}

	timeout := time.Duration(timeoutMS) * time.Millisecond

	fmt.Printf("DB_HOST: %v (%T)\n", host, host)
	fmt.Printf("DB_PORT: %v (%T)\n", port, port)
	fmt.Printf("DB_MAX_CONNS: %v (%T)\n", maxConns, maxConns)
	fmt.Printf("DB_USE_SSL: %v (%T)\n", useSSL, useSSL)
	fmt.Printf("DB_TIMEOUT_MS: %v (%T)\n", timeout, timeout)
}
