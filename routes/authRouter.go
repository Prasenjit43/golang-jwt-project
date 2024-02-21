package routes

import (
	"fmt"

	controllers "github.com/Prasenjit43/golang-jwt-project/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {

	fmt.Println("111111")

	incomingRoutes.POST("/users/signup", controllers.Signup())
	incomingRoutes.POST("/users/login", controllers.Login())

}
