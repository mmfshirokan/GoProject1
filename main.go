package main

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/mmfshirokan/GoProject1/config"
	"github.com/mmfshirokan/GoProject1/handlers"
	"github.com/mmfshirokan/GoProject1/model"
	"github.com/mmfshirokan/GoProject1/repository"
	"github.com/mmfshirokan/GoProject1/service"
)

func main() {
	conf := config.Config{
		Database: "postgres",
	}

	repo := repository.NewRepository(conf)
	pw_repo := repository.NewPasswordRepository(conf)
	auth_repo := repository.NewAuthRpository()

	usr := service.NewUser(repo)
	pw := service.NewPassword(pw_repo)
	tok := service.NewToken(auth_repo)

	hand := handlers.NewHandler(usr, pw, tok)

	e := echo.New()
	e.POST("/users/signup", hand.SignUp)
	e.PUT("/users/signin", hand.SignIn)
	e.PUT("/users/refresh", hand.Refresh)
	g := e.Group("/users/auth")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(model.UserRequest)
		},
		SigningKey: []byte("secret"),
	}

	g.Use(echojwt.WithConfig(config))
	g.GET("/get", hand.GetUser)
	g.PUT("/update", hand.UpdateUser)
	g.DELETE("/delete", hand.DeleteUser)
	e.Logger.Fatal(e.Start(":8081"))
}
