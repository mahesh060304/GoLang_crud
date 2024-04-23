package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mahesh060304/go-crud/controllers"
	"github.com/mahesh060304/go-crud/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()
	r.POST("/addusers",controllers.CreateNewUser)
        r.GET("/users", controllers.GetAllUsers)
	r.PUT("/updateuser/:id",controllers.UpdateUser)
	r.DELETE("deleteuser/:id",controllers.DeleteUser)

	r.Run()
}
