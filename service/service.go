package service

import "github.com/labstack/echo"

func StartServer(getUser func(echo.Context) error, saveUser func(echo.Context) error, updateUser func(echo.Context) error, deleteUser func(echo.Context) error) {
	e := echo.New()
	e.GET("/users:id", getUser)
	e.POST("/users:id", saveUser)
	e.PUT("/users:id", updateUser)
	e.DELETE("/users:id", deleteUser)
	e.Logger.Fatal(e.Start(":8080"))
}
