package main

import (
	"log"
	"net/http"
	"os"
	"pvg/controllers"
	"pvg/infrastructure"
	"pvg/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func readEnvironmentFile() {
	//Environment file Load --------------------------------
	err := godotenv.Load(".pvg.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(3)
	}
}

func main() {
	readEnvironmentFile()

	repository.DB = infrastructure.OpenDbConnection()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/user", controllers.CreateUser())
	r.GET("/user/:userId", controllers.GetSpecificUser())
	r.GET("/users", controllers.GetAllUsers())
	r.PUT("/user/:userId", controllers.UpdateUser())
	r.DELETE("/user/:userId", controllers.DeleteUser())

	r.POST("/checkuserexist", controllers.CheckUserExist())
	r.POST("/forgotpassword", controllers.SendMailChangePassword())

	r.POST("/verify/user", controllers.VerifyUser())

	r.Run()
}
