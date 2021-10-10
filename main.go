package main

import (
	"io"
	"log"
	"log/syslog"
	"net/http"
)

const version string = "2.0.1"

// VersionHandler handles incoming requests to /version
// and just returns a simple version number
func versionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("log*** versionHandler", version)
	io.WriteString(w, version)
}

func main() {
	// Log to syslog
	logWriter, err := syslog.New(syslog.LOG_SYSLOG, "My Awesome App")
	if err != nil {
		log.Fatalln("Unable to set logfile:", err.Error())
	}
	// + set log flag
	defer log.SetFlags(log.Lshortfile)
	// set the log output
	defer log.SetOutput(logWriter)

	log.Println("This is a log from GOLANG")
	log.Println("Listening on port 8000...")
	http.HandleFunc("/version", versionHandler)
	http.ListenAndServe(":8000", nil)
}