package main

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/mmfshirokan/GoProject1/config"
	"github.com/mmfshirokan/GoProject1/handlers"
	"github.com/mmfshirokan/GoProject1/handlers/request"
	"github.com/mmfshirokan/GoProject1/repository"
	"github.com/mmfshirokan/GoProject1/service"
)

func main() {
	conf := config.Config{
		Database: "postgres",
	}

	repo := repository.NewRepository(conf)
	serv := service.NewUser(repo)

	pwrepo := repository.NewPasswordRepository(conf)
	pw := service.NewPassword(pwrepo)

	hand := handlers.NewHandler(serv, pw)

	e := echo.New()
	e.POST("/users/signin", hand.Register) // create changed to Register
	e.PUT("/users/login:id", hand.Login)
	g := e.Group("/users/auth")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(request.UserRequest)
		},
		SigningKey: []byte("secret"),
	}

	g.Use(echojwt.WithConfig(config))
	g.GET("/get", hand.GetUser)
	g.PUT("/update", hand.UpdateUser)
	g.DELETE("/delete", hand.DeleteUser)
	e.Logger.Fatal(e.Start(":8081"))
}
