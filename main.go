package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/mmfshirokan/GoProject1/config"
	"github.com/mmfshirokan/GoProject1/handlers"
	"github.com/mmfshirokan/GoProject1/passwordRepository"
	"github.com/mmfshirokan/GoProject1/passwordService"
	"github.com/mmfshirokan/GoProject1/repository"
	"github.com/mmfshirokan/GoProject1/service"
)

func main() {
	conf := config.Config{
		Database: "mongodb",
	}

	repo := repository.NewRepository(conf)
	serv := service.NewUser(repo)

	pwrepo := passwordRepository.NewPasswordRepository(conf)
	pw := passwordService.NewPassword(pwrepo)

	hand := handlers.NewHandler(serv, pw)

	e := echo.New()
	e.Use(middleware.BasicAuth(hand.Login))

	e.GET("/users:id", hand.GetUser)
	e.POST("/users:id", hand.Register) // create changed to Register
	e.PUT("/users:id", hand.UpdateUser)
	e.DELETE("/users:id", hand.DeleteUser)
	e.Logger.Fatal(e.Start(":8080"))
}
