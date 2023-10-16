package repository

import (
	"github.com/jackc/pgx"
)

var (
	conn       *pgx.Conn
	ConnConfig pgx.ConnConfig = pgx.ConnConfig{
		Host:     "project1-postgres-1",
		Port:     5432,
		Database: "echodb",
		User:     "echopguser",
		Password: "pgpw4echo",
	}
)

func GetUserTroughID(id string) (string, string, error) {
	var (
		name string
		male string
	)
	err := conn.QueryRow("SELECT name, male FROM entity WHERE id = "+id).Scan(&name, &male)
	return name, male, err
}

func SaveUser(id string, name string, male string) error {
	_, err := conn.Exec("INSERT INTO entity VALUES ($1, $2, $3)", id, name, male)
	return err
}

func UpdateUser(id string, name string, male string) error {
	_, err := conn.Exec("UPDATE entity SET name = $1, male = $2 WHERE id = $3", name, male, id)
	return err
}

func DeleteUser(id string) error {
	_, err := conn.Exec("DELETE FROM entity WHERE id = $1", id)
	return err
}

func CreatEntity() error {
	_, err := conn.Exec("CREATE TABLE IF NOT EXISTS entity (id INT PRIMARY KEY, name CHARACTER VARYING(30) NOT NULL, male BOOLEAN NOT NULL)")
	return err
}

func OpenConnection() error {
	var err error
	conn, err = pgx.Connect(ConnConfig)
	return err
}

func CloseConnection() error {
	err := conn.Close()
	return err
}
