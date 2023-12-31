package main

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mmfshirokan/GoProject1/internal/config"
	"github.com/mmfshirokan/GoProject1/internal/handlers"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/mmfshirokan/GoProject1/internal/service"
)

func main() {
	conf := config.NewConfig()
	val := validator.New(validator.WithRequiredStructEnabled())

	if err := val.Struct(&conf); err != nil {
		fmt.Fprint(os.Stderr, "invalid config fild/s")
	}

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
	group.PUT("/uploadImage", hand.UploadImage)
	group.PUT("/downloadImage", hand.DownloadImage)
	echoServ.Logger.Fatal(echoServ.Start(":8081"))
}
