package repository_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

var conn repository.Interface

func TestMain(m *testing.M) {
	ctx, _ := context.WithCancel(context.Background())

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		//Hostname:   "localhost:5432",
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_PASSWORD=pgpw4echo",
			"POSTGRES_USER=echopguser",
			"POSTGRES_DB=echodb",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://echopguser:pgpw4echo@%s/echodb?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120)

	dbpool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		dbpool.Close()
		log.Error("can't connect to the pgxpool: %w", err)
	}

	commandArr := []string{
		"echo",
		"k0ntraltOO",
		"|",
		"sudo",
		"-S",
		"flyway",
		"-url=jdbc:postgresql://postgres/echodb",
		"-user=echopguser",
		"-password=pgpw4echo",
		"-schemas=\"apps\"", // remove?
		//"-initSql=\"CREATE TABLE IF NOT EXISTS entity (id INT PRIMARY KEY, name TEXT NOT NULL, male BOOLEAN NOT NULL)\"",
		"-connectRetries=60",
		"migrate",
		"wget",
	}
	cmd := exec.Command(commandArr[0], commandArr[1:]...)
	cmd.Env = []string{
		"FLYWAY_CONNECT_RETRIES=60",
		"FLYWAY_LOCATIONS=filesystem:/home/andreishyrakanau/projects/project1/GoProject1/migrations/sql", //app/sql
		"FLYWAY_SCHEMAS=apps",
	}

	b, err := cmd.CombinedOutput()

	pool.MaxWait = 120 * time.Second
	conn = repository.NewPostgresRepository(dbpool)

	//go cmd.Run()

	if err != nil {
		log.Fatal("following error occured while migrating in repository_test:", string(b), err)
	} else {
		log.Info(string(b))
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestCreate(t *testing.T) {
	type testCase struct {
		name     string
		input    model.User
		hasError bool
	}
	testTable := []testCase{
		{
			name: "standart input",
			input: model.User{
				ID:       110,
				Name:     "Mark",
				Male:     true,
				Password: "Abcd",
			},
			hasError: false,
		},
		{
			name: "repeated intput",
			input: model.User{
				ID:       110,
				Name:     "Mark",
				Male:     true,
				Password: "Abcd",
			},
			hasError: true,
		},
		{
			name: "ID bigger than allowed",
			input: model.User{
				ID:       1000000001,
				Name:     "Jane",
				Male:     false,
				Password: "ldrd",
			},
			hasError: true,
		},
		{
			name: "Name bigger than allowed",
			input: model.User{
				ID:       111,
				Name:     "12345678901234567890123456789012345678901",
				Male:     false,
				Password: "ldrd",
			},
			hasError: true,
		},
		{
			name: "Password bigger than allowed",
			input: model.User{
				ID:       112,
				Name:     "Alice",
				Male:     false,
				Password: "12345678901234567890123456789012345678901",
			},
			hasError: true,
		},
		{
			name: "standart input",
			input: model.User{
				ID:       113,
				Name:     "Jhon",
				Male:     true,
				Password: "4040",
			},
			hasError: false,
		},
	}

	for _, test := range testTable {
		err := conn.Create(context.Background(), test.input)
		if test.hasError {
			assert.NotNil(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
}
