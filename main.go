package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/mmfshirokan/GoProject1/config"
	"github.com/mmfshirokan/GoProject1/handlers"
	"github.com/mmfshirokan/GoProject1/repository"
	"github.com/mmfshirokan/GoProject1/service"
)

func main() {
	conf := config.Config{
		Database: "mongodb",
	}

	repo := repository.NewRepository(conf)
	serv := service.NewUser(repo)

	pwrepo := repository.NewPasswordRepository(conf)
	pw := service.NewPassword(pwrepo)

	hand := handlers.NewHandler(serv, pw)

	e := echo.New()
	e.POST("/users", hand.Register) // create changed to Register
	g := e.Group("/users")

	g.Use(middleware.BasicAuth(hand.Login))
	g.GET("/auth:id", hand.GetUser)
	g.PUT("/auth:id", hand.UpdateUser)
	g.DELETE("/auth:id", hand.DeleteUser)
	e.Logger.Fatal(e.Start(":8080"))
}
