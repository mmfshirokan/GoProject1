package main

import (
	"context"
	"net"

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
	"github.com/mmfshirokan/GoProject1/internal/server"
	"github.com/mmfshirokan/GoProject1/internal/service"
	"github.com/mmfshirokan/GoProject1/proto/pb"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
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
		repo     repository.RepositoryInterface
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

	go rpcServerStart(repo, pwRepo, authRepo)

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

func rpcServerStart(
	repo repository.RepositoryInterface,
	pwRepo repository.PwRepositoryInterface,
	authRepo repository.AuthRepositoryInterface,
) {
	lis, err := net.Listen("tcp", "localhost:9091")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryServerInterceptor))

	rpcUser := server.NewUserServer(repo)
	rpsToken := server.NewTokenServer(repo, pwRepo, authRepo)
	rpsImage := server.NewImageServer()

	pb.RegisterUserServer(grpcServer, rpcUser)
	pb.RegisterTokenServer(grpcServer, rpsToken)
	pb.RegisterImageServer(grpcServer, rpsImage)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("rpc fatal error")
	}
}

func unaryServerInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	// if info.FullMethod != "/pb.Token/SignUp" && info.FullMethod != "/pb.Token/SignIn" && info.FullMethod != "/pb.Token/Refresh" {
	// 	incomingMetadata, ok := metadata.FromIncomingContext(ctx)
	// 	if !ok {
	// 		err := errors.New("missing methadata")
	// 		log.Error("warning! metadata missing in main", err)
	// 		return nil, err
	// 	}

	// 	val, ok := incomingMetadata["authorization"]
	// 	if !ok || len(val) != 1 {
	// 		err := errors.New("missingAuth")
	// 		log.Error("warning! auth missing in metadata", err)
	// 		return nil, err
	// 	}

	// 	_, err := jwt.Parse(val[0], func(t *jwt.Token) (interface{}, error) {
	// 		return []byte("secret"), nil
	// 	})
	// 	if err != nil {
	// 		log.Error("jwt token parse failed in main: ", err)
	// 		return nil, err
	// 	}
	// }

	result, err := handler(ctx, req)
	if err != nil {
		log.Error("rpc failed with error: ", err)
		return nil, err
	}

	return result, nil
}
