package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx"
	"github.com/labstack/echo"
)

var (
	conn       *pgx.Conn
	connConfig pgx.ConnConfig = pgx.ConnConfig{
		Host:     "project1-postgres-1",
		Port:     5432,
		Database: "echodb",
		User:     "echopguser",
		Password: "pgpw4echo",
	}
)

func main() {
	var err error
	conn, err = pgx.Connect(connConfig)

	if err != nil {
		fmt.Println("Error ocured while connecting to db: ", err)
		os.Exit(1)
	}

	defer conn.Close()

	e := echo.New()
	e.GET("/users:id", getUser)
	e.POST("/users:id", saveUser)
	e.PUT("/users:id", updateUser)
	e.DELETE("/users:id", deleteUser)
	e.Logger.Fatal(e.Start(":8080"))
}

func getUser(c echo.Context) error {
	id := c.Param("id")

	var (
		name string
		male bool
	)
	conn.QueryRow("SELECT name, male FROM entity WHERE id = "+id).Scan(&name, &male)
	return c.String(http.StatusOK, "User id: "+id+"\nUser name: "+name+"\nUser male: "+strconv.FormatBool(male)+"\n")
}

func saveUser(c echo.Context) error {
	var (
		id   int64
		err  error
		name string
		male bool
	)
	id, err = strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return err
	}
	name = c.FormValue("name")
	male, err = strconv.ParseBool(c.FormValue("male"))

	if err != nil {
		return err
	}
	_, err = conn.Exec("INSERT INTO entity VALUES ($1, $2, $3)", id, name, male) // обязательно ли id, name, male, должны быть int, string, bool,
	return err                                                                   // соответствеено или они могут быть string?
}

func updateUser(c echo.Context) error {
	var (
		id   int64
		err  error
		name string
		male bool
	)
	name = c.FormValue("name")
	male, err = strconv.ParseBool(c.FormValue("male"))

	if err != nil {
		return err
	}
	id, err = strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return err
	}
	_, err = conn.Exec("UPDATE entity SET name = $1, male = $2 WHERE id = $3", name, male, id)
	return err
}

func deleteUser(c echo.Context) error {
	var (
		id  int64
		err error
	)
	id, err = strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return err
	}
	_, err = conn.Exec("DELETE FROM entity WHERE id = $1", id)
	return err
}
