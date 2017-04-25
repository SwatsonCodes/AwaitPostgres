package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func getEnvVarAsInt(varName string, varDefault string) (*int64, error) {
	varStr := os.Getenv(varName)
	if varStr == "" {
		varStr = varDefault
	}
	varInt, err := strconv.ParseInt(varStr, 10, 0)
	if err != nil {
		return nil, err
	}
	return &varInt, nil
}

func main() {
	var numAttempts int64
	var connectionErr error

	url := os.Getenv("POSTGRES_URL")
	if url == "" {
		log.Fatal("Must supply postgres URL to connect to as environment variable POSTGRES_URL")
	}
	numRetries, err := getEnvVarAsInt("RETRIES", "10")
	if err != nil {
		log.Fatal(err)
	}
	pause, err := getEnvVarAsInt("WAIT_SECS", "2")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	connectionErr = db.Ping()
	for connectionErr != nil && numAttempts < *numRetries {
		log.Printf("Unable to connect: %s\n", connectionErr.Error())
		numAttempts++
		time.Sleep(time.Duration(*pause) * time.Second)
		connectionErr = db.Ping()
	}
	if numAttempts >= *numRetries {
		log.Fatal("Unable to establish connection to postgres. Aborting.")
	}
	log.Println("Connected to postgres!")
}
