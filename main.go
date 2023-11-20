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
	pwRepo := repository.NewPasswordRepository(conf)
	authRepo := repository.NewAuthRpository()

	usr := service.NewUser(repo)
	pw := service.NewPassword(pwRepo)
	tok := service.NewToken(authRepo)

	hand := handlers.NewHandler(usr, pw, tok)

	echoServ := echo.New()
	echoServ.POST("/users/signup", hand.SignUp)
	echoServ.PUT("/users/signin", hand.SignIn)
	echoServ.PUT("/users/refresh", hand.Refresh)
	group := echoServ.Group("/users/auth")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(model.UserRequest)
		},
		SigningKey: []byte("secret"),
	}

	group.Use(echojwt.WithConfig(config))
	group.GET("/get", hand.GetUser)
	group.PUT("/update", hand.UpdateUser)
	group.DELETE("/delete", hand.DeleteUser)
	echoServ.Logger.Fatal(echoServ.Start(":8081"))
}
