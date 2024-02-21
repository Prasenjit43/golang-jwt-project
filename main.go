package main

import (
	"fmt"
	"os"

	routes "github.com/Prasenjit43/golang-jwt-project/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Start of the golang-jwt-token project")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserAuth(router)

	router.Run(":" + port)

}
