package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/swatsoncodes/AwaitPostgres/reachability"
)

type connector func() error

func attemptConnectionRepeatedly(connect connector, errorMsg string, numRetries int64, pause time.Duration) error {
	var numAttempts int64

	for attempt := connect(); attempt != nil && numAttempts < numRetries; {
		log.Printf("%s: %s \n", errorMsg, attempt.Error())
		numAttempts++
		time.Sleep(pause)
	}
	if numAttempts >= numRetries {
		return fmt.Errorf("%s after %d attempts. Aborting", errorMsg, numAttempts)
	}
	return nil
}

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
	pgURLStr := os.Getenv("POSTGRES_URL")
	if pgURLStr == "" {
		log.Fatal("Must supply postgres URL to connect to as environment variable POSTGRES_URL")
	}
	pgURL, err := url.Parse(pgURLStr)
	if err != nil {
		log.Fatal(err)
	}
	numRetries, err := getEnvVarAsInt("RETRIES", "10")
	if err != nil {
		log.Fatal(err)
	}
	pause, err := getEnvVarAsInt("WAIT_SECS", "2")
	if err != nil {
		log.Fatal(err)
	}
	pauseDuration := time.Duration(*pause) * time.Second
	hostTimeout := reachability.HostTimeout{Host: pgURL.Host, Timeout: time.Duration(2) * time.Second}

	err = attemptConnectionRepeatedly(hostTimeout.IsHostReachable, "Unable to reach host", *numRetries, pauseDuration)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", pgURLStr)
	if err != nil {
		log.Fatal(err)
	}

	err = attemptConnectionRepeatedly(db.Ping, "Unable to connect to postgres", *numRetries, pauseDuration)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to postgres!")
}
