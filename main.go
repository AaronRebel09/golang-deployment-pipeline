package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

const version string = "2.0.1"

// VersionHandler handles incoming requests to /version
// and just returns a simple version number
func versionHandler(w http.ResponseWriter, r *http.Request) {
	logFile := fileLog()
	defer logFile.Close()
	// Set log out put and enjoy :)
	log.SetOutput(logFile)
	// optional: log date-time, filename, and line number
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	log.Println("log*** versionHandler", version)
	io.WriteString(w, version)
}

func fileLog() *os.File {
	// log to custom file
	LOG_FILE := "/tmp/app_log"
	// open log file
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	return logFile
}

func main() {

	logFile := fileLog()
	defer logFile.Close()
	// Set log out put and enjoy :)
	log.SetOutput(logFile)
	// optional: log date-time, filename, and line number
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	log.Println("Logging to custom file")

	log.Println("This is a log from GOLANG")
	log.Println("Listening on port 8000...")
	http.HandleFunc("/version", versionHandler)
	http.ListenAndServe(":8000", nil)
}