package main

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mmfshirokan/GoProject1/internal/config"
	"github.com/mmfshirokan/GoProject1/internal/handlers"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/mmfshirokan/GoProject1/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Out = os.Stdout

	conf := config.NewConfig()

	repo := repository.NewRepository(conf)
	pwRepo := repository.NewPasswordRepository(conf)
	authRepo := repository.NewAuthRpository(conf)

	usr := service.NewUser(repo)
	pw := service.NewPassword(pwRepo)
	tok := service.NewToken(authRepo)

	hand := handlers.NewHandler(usr, pw, tok)

	echoServ := echo.New()
	echoServ.POST("/users/signup", hand.SignUp)
	echoServ.PUT("/users/signin", hand.SignIn)
	echoServ.PUT("/users/refresh", hand.Refresh)
	group := echoServ.Group("/users/auth")

	echoConf := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(model.UserRequest)
		},
		SigningKey: []byte("secret"),
	}

	group.Use(echojwt.WithConfig(echoConf))
	group.GET("/get", hand.GetUser)
	group.PUT("/update", hand.UpdateUser)
	group.DELETE("/delete", hand.DeleteUser)
	echoServ.Logger.Fatal(echoServ.Start(":8081"))
}
