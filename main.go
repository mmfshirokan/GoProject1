package main

import (
	"context"
	//"fmt"
	//"os"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/mmfshirokan/GoProject1/docs"
	"github.com/mmfshirokan/GoProject1/internal/config"
	"github.com/mmfshirokan/GoProject1/internal/consumer"
	"github.com/mmfshirokan/GoProject1/internal/handlers"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/mmfshirokan/GoProject1/internal/service"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		log.Error("invalid config fild/s")
	}

	var (
		repo     repository.Interface
		pwRepo   repository.PwRepositoryInterface
		authRepo repository.AuthRepositoryInterface
	)

	if conf.Database == "mongodb" {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.MongoURI))
		if err != nil {
			log.Error("invalid config fild/s %w", err)
		}

		repo = repository.NewMongoRepository(client.Database("users").Collection("entity"))
		pwRepo = repository.NewMongoPasswordRepository(client.Database("users").Collection("passwords"))

	} else if conf.Database == "postgresql" {

		dbpool, err := pgxpool.New(ctx, conf.PostgresURI)
		if err != nil {
			dbpool.Close()
			log.Error("can't connect to the pgxpool: %w", err)
		}

		repo = repository.NewPostgresRepository(dbpool)
		pwRepo = repository.NewPostgresPasswordRepository(dbpool)
		authRepo = repository.NewAuthRpository(dbpool)

	} else {
		log.Fatal("unexpected error occurred (wrong config)")
		repo = nil
		pwRepo = nil
		authRepo = nil
	}

	redisClient := repository.NewCLient(conf)
	redisUsr := repository.NewUserRedisRepository(redisClient)
	redisTok := repository.NewRftRedisRepository(redisClient)

	userMap := make(map[string]*model.User)
	rftMap := make(map[string][]*model.RefreshToken)

	userMapConn := repository.NewUserMap(userMap)
	rftMapConn := repository.NewRftMap(rftMap)

	usr := service.NewUser(repo, redisUsr, userMapConn)
	pw := service.NewPassword(pwRepo)
	tok := service.NewToken(authRepo, redisTok, rftMapConn)

	cons := consumer.NewConsumer(
		redisClient,
		redisUsr,
		redisTok,
		userMapConn,
		rftMapConn,
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
