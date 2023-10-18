package main

import (
	"github.com/labstack/echo"

	"github.com/mmfshirokan/GoProject1/handlers"
	"github.com/mmfshirokan/GoProject1/repository"
	"github.com/mmfshirokan/GoProject1/service"
)

func main() {
	repo := repository.NewRepository()
	serv := service.NewUser(repo)
	hand := handlers.NewHandler(serv)

	e := echo.New()
	e.GET("/users:id", hand.GetUser)
	e.POST("/users:id", hand.SaveUser)
	e.PUT("/users:id", hand.UpdateUser)
	e.DELETE("/users:id", hand.DeleteUser)
	e.Logger.Fatal(e.Start(":8080"))
}
