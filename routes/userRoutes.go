package routes

import (
	controller "github.com/AaronRebel09/golang-deployment-pipeline/controllers"
	"github.com/AaronRebel09/golang-deployment-pipeline/middleware"
	"github.com/gin-gonic/gin"
)

// UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controller.SignUp())
	incomingRoutes.POST("/users/signin", controller.Login())
	//Se agrega middleware Authentication solo a este endpoint
	incomingRoutes.GET("/users/info", middleware.Authentication(), controller.SimpleEndpoint())
}
