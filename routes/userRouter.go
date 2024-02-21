package routes

import (
	"github.com/Prasenjit43/golang-jwt-project/controllers"
	"github.com/Prasenjit43/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func UserAuth(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/getUser", controllers.GetUsers())
	incomingRoutes.GET("/user/:user_id", controllers.GetUser())

}
