package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/mmfshirokan/GoProject1/docs"
	"github.com/mmfshirokan/GoProject1/internal/config"
	"github.com/mmfshirokan/GoProject1/internal/consumer"
	"github.com/mmfshirokan/GoProject1/internal/handlers"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/mmfshirokan/GoProject1/internal/service"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Echo Serevr
// @version 1.0
// @description This is a server for using JWT, Swagger, and exetra with Echo.

// @securityDefinitions.apiKey JWT
// @in       header
// @name      token

// @host localhost:8081
// @BasePath /users
// @schemes http
func main() {
	conf := config.NewConfig()
	val := validator.New(validator.WithRequiredStructEnabled())
	ctx, _ := context.WithCancel(context.Background())

	if err := val.Struct(&conf); err != nil {
		fmt.Fprint(os.Stderr, "invalid config fild/s")
	}

	repo := repository.NewRepository(conf)
	pwRepo := repository.NewPasswordRepository(conf)
	authRepo := repository.NewAuthRpository(conf)

	redisClient := repository.NewCLient(conf)
	redisUsr := repository.NewUserRedisRepository(redisClient)
	srcUserMap := repository.NewUserMap()
	redisTok := repository.NewRftRedisRepository(redisClient)
	srcRftMap := repository.NewRftMap()

	usr := service.NewUser(repo, redisUsr, srcUserMap)
	pw := service.NewPassword(pwRepo)
	tok := service.NewToken(authRepo, redisTok, srcRftMap)

	cons := consumer.NewConsumer(
		redisClient,
		redisUsr,
		redisTok,
		srcUserMap,
		srcRftMap,
	)

	go cons.Consume(ctx)

	hand := handlers.NewHandler(usr, pw, tok)

	echoServ := echo.New()
	echoServ.POST("/users/signup", hand.SignUp)
	echoServ.PUT("/users/signin", hand.SignIn)
	echoServ.PUT("/users/refresh", hand.Refresh)
	echoServ.GET("/swagger/*", echoSwagger.WrapHandler)
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
