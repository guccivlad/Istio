package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Status struct {
	Stat string `json:"status"`
}

func getEnv(key string) string {
	value := os.Getenv(key)
	return value
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	status := Status{
		Stat: "ok",
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(status)
}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	WELCOME_MESSAGE := getEnv("WELCOME_MESSAGE")
	fmt.Fprintf(w, WELCOME_MESSAGE)
}

func getLogHandler(w http.ResponseWriter, r *http.Request) {
	LOG_FILE := getEnv("LOG_FILE")

	if _, err := os.Stat(LOG_FILE); os.IsNotExist(err) {
		emptyFile, createErr := os.Create(LOG_FILE)

		if createErr != nil {
			http.Error(w, fmt.Sprintf("[ERROR] Could not create log file: %v", createErr), http.StatusInternalServerError)
			return
		}
		emptyFile.Close()
	}

	logs, err := os.ReadFile(LOG_FILE)
	if err != nil {
		http.Error(w, fmt.Sprintf("[ERROR] Read file error %v", err), http.StatusBadRequest)
		return
	}

	_, _ = w.Write(logs)
}

func postLogHandler(w http.ResponseWriter, r *http.Request) {
	LOG_FILE := getEnv("LOG_FILE")

	if _, err := os.Stat(LOG_FILE); os.IsNotExist(err) {
		emptyFile, createErr := os.Create(LOG_FILE)

		if createErr != nil {
			http.Error(w, fmt.Sprintf("[ERROR] Could not create log file: %v", createErr), http.StatusInternalServerError)
			return
		}
		emptyFile.Close()
	}

	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		http.Error(w, fmt.Sprintf("[ERROR] Read file error %v", err), http.StatusBadRequest)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	body, readErr := io.ReadAll(r.Body)

	if readErr != nil {
		http.Error(w, fmt.Sprintf("[ERROR] Read file error %v", err), http.StatusBadRequest)
		return
	}

	log.Print(string(body))
}

func main() {
	fmt.Println("Server start")

	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/log", postLogHandler)
	http.HandleFunc("/logs", getLogHandler)
	http.HandleFunc("/", simpleHandler)

	http.ListenAndServe("0.0.0.0:8081", nil)
}
