package main

import (
	middleware "github.com/AaronRebel09/golang-deployment-pipeline/middleware"
	routes "github.com/AaronRebel09/golang-deployment-pipeline/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

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
	log.Println("Iniciando web service ...")
	//port := os.Getenv("PORT")
	//if port == "" {
	//	port = "8000"
	//}
	port := "8588"
	log.Println("port : ", port)
	router := gin.New()
	router.Use(gin.Logger())
	//Se agrega Cross Origin Resource Sharing
	router.Use(middleware.CORSMiddleware())
	routes.UserRoutes(router)
	// API-1
	router.GET("/api/v1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api v1"})
	})

	// API-2
	router.GET("/api/v2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api v2"})
	})
	router.Run(":" + port)

}
