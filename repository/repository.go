package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mmfshirokan/GoProject1/model"

	"os"

	"context"
)

type Repository struct {
	conn *pgxpool.Pool
	err  error
}

func NewRepository() *Repository {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv(os.Getenv("DATABASE_URL")))
	return &Repository{
		conn: dbpool,
		err:  err,
	}
}

/*unc (rep *Repository) SetConnection() error {
	rep.ConnConfig =  pgx.ConnConfig{
		Host:     "project1-postgres-1",
		Port:     5432,
		Database: "echodb",
		User:     "echopguser",
		Password: "pgpw4echo",
	}
	rep.conn, rep.err = pgx.Connect(rep.ConnConfig)
	return rep.err
}*/

func (rep *Repository) GetUserTroughID(id string) (string, string, error) {
	usr := model.User{}
	rep.err = rep.conn.QueryRow(context.Background(), "SELECT name, male FROM entity WHERE id = "+id).Scan(&usr.Name, &usr.Male)
	return usr.Name, usr.Male, rep.err
}

func (rep *Repository) SaveUser(id string, name string, male string) error {
	_, rep.err = rep.conn.Exec(context.Background(), "INSERT INTO entity VALUES ($1, $2, $3)", id, name, male)
	return rep.err
}

func (rep *Repository) UpdateUser(id string, name string, male string) error {
	_, rep.err = rep.conn.Exec(context.Background(), "UPDATE entity SET name = $1, male = $2 WHERE id = $3", name, male, id)
	return rep.err
}

func (rep *Repository) DeleteUser(id string) error {
	_, rep.err = rep.conn.Exec(context.Background(), "DELETE FROM entity WHERE id = $1", id)
	return rep.err
}

func (rep *Repository) CreatEntity() error {
	_, rep.err = rep.conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS entity (id INT PRIMARY KEY, name CHARACTER VARYING(30) NOT NULL, male BOOLEAN NOT NULL)")
	return rep.err
}
