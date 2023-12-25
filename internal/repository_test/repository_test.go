package repository_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmfshirokan/GoProject1/internal/config"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"
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

	pgResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:   "postgres_test",
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

	rdResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Hostname:   "redis_test",
		Repository: "redis",
		Tag:        "latest",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	postgresHostAndPort := pgResource.GetHostPort("5432/tcp")
	postgresUrl := fmt.Sprintf("postgres://echopguser:pgpw4echo@%s/echodb?sslmode=disable", postgresHostAndPort)
	redisHostAndPort := rdResource.GetHostPort("6379/tcp")
	redisUrl := fmt.Sprintf("redis://%s/0?protocol=3", redisHostAndPort)

	log.Println("Connecting to database on url: ", postgresUrl)
	log.Println("Connecting to redis on url: ", redisUrl)

	var dbpool *pgxpool.Pool
	if err = pool.Retry(func() error { // remove retry? (not nessesary)
		dbpool, err = pgxpool.New(ctx, postgresUrl)
		if err != nil {
			dbpool.Close()
			log.Error("can't connect to the pgxpool: %w", err)
		}
		return dbpool.Ping(ctx)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	redisConf := config.Config{RedisURI: redisUrl}

	var client *redis.Client
	if err = pool.Retry(func() error {
		client = repository.NewCLient(redisConf)
		return client.Ping(ctx).Err()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	commandArr := []string{
		"-url=jdbc:postgresql://" + postgresHostAndPort + "/echodb",
		"-user=echopguser",
		"-password=pgpw4echo",
		"-locations=filesystem:../../migrations/sql",
		"-schemas=apps", //remove?
		"-connectRetries=60",
		"migrate",
	}
	cmd := exec.Command("flyway", commandArr[:]...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err = cmd.Run()
	if err != nil {
		log.Error(fmt.Printf("error: %s", err))
	}
	log.Info(fmt.Printf("out: %s%s", outb.String(), errb.String()))

	pool.MaxWait = 120 * time.Second
	conn = repository.NewPostgresRepository(dbpool)
	pwConn = repository.NewPostgresPasswordRepository(dbpool)
	authConn = repository.NewAuthRpository(dbpool)

	redisUsrConn = repository.NewUserRedisRepository(client)
	redisRftConn = repository.NewRftRedisRepository(client)

	mapUsr = make(map[string]*model.User)
	mapRft = make(map[string][]*model.RefreshToken)
	mapUsrConn = repository.NewUserMap(mapUsr)
	mapRftConn = repository.NewRftMap(mapRft)

	code := m.Run()

	if err := pool.Purge(pgResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestCreate(t *testing.T) { // TODO add error comparison
	type testCase struct { // TODO add more cases (negative number id?)
		name     string
		input    model.User
		hasError bool
	}
	testTable := []testCase{
		{
			name: "standart input with ID=110",
			input: model.User{
				ID:       110,
				Name:     "Mark",
				Male:     true,
				Password: "",
			},
			hasError: false,
		},
		{
			name: "repeated intput",
			input: model.User{
				ID:       110,
				Name:     "Mark",
				Male:     true,
				Password: "",
			},
			hasError: true,
		},
		{
			name: "ID bigger than allowed",
			input: model.User{
				ID:       1000000001,
				Name:     "Jane",
				Male:     false,
				Password: "",
			},
			hasError: true,
		},
		{
			name: "name bigger than allowed",
			input: model.User{
				ID:       111,
				Name:     "12345678901234567890123456789012345678901",
				Male:     false,
				Password: "",
			},
			hasError: true,
		},
		{
			name: "standart input with ID=113",
			input: model.User{
				ID:       113,
				Name:     "Jhon",
				Male:     true,
				Password: "",
			},
			hasError: false,
		},
	}

	for _, test := range testTable {
		err := conn.Create(context.Background(), test.input)
		if test.hasError {
			assert.Error(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
	fmt.Print("TestCreate Finished!\n")
}

func TestGetTroughID(t *testing.T) {
	type testCase struct {
		name     string
		input    int
		output   *model.User
		hasError bool
	}
	testTable := []testCase{
		{
			name:  "standart input with ID=110",
			input: 110,
			output: &model.User{
				ID:       110,
				Name:     "Mark",
				Male:     true,
				Password: "",
			},
			hasError: false,
		},
		{
			name:  "standart input with ID=113",
			input: 113,
			output: &model.User{
				ID:       113,
				Name:     "Jhon",
				Male:     true,
				Password: "",
			},
			hasError: false,
		},
		{
			name:     "ID that does not exist",
			input:    2288,
			output:   nil,
			hasError: true,
		},
		{
			name:     "ID bigger than allowed therfore it does not exist",
			input:    1000000001,
			output:   nil,
			hasError: true,
		},
	}

	for _, test := range testTable {
		usrForTesting, err := conn.GetTroughID(context.Background(), test.input)
		if test.hasError {
			assert.Error(t, err, "this test case should have error: ", test.name)
			assert.Nil(t, test.output, "otput must be nil: ", test.name)
		} else {
			assert.Equal(t, test.output, usrForTesting)
			assert.Nil(t, err, "err must be nil: ", test.name)
		}
	}
	fmt.Print("TestGetTroughID Finished!\n")
}

func TestUpdate(t *testing.T) {
	type testCase struct {
		name     string
		input    model.User
		hasError bool
	}
	testTable := []testCase{
		{
			name: "standart input with ID= 110",
			input: model.User{
				ID:   110,
				Name: "Lilu",
				Male: false,
			},
			hasError: false,
		},
		{
			name: "standart input with ID=113",
			input: model.User{
				ID:   113,
				Name: "Jane",
				Male: false,
			},
			hasError: false,
		},
		{ // TODO add hasError=true chage logic
			name: "ID does not exist",
			input: model.User{
				ID:   1234567,
				Name: "Jane",
				Male: false,
			},
			hasError: false,
		},
		{
			name: "name bigger than allowed",
			input: model.User{
				ID:       113,
				Name:     "12345678901234567890123456789012345678901",
				Male:     false,
				Password: "",
			},
			hasError: true,
		},
	}

	for _, test := range testTable {
		err := conn.Update(context.Background(), test.input)
		if test.hasError {
			assert.Error(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
	fmt.Print("TestUpdate Finished!\n")
}

func TestDelete(t *testing.T) {
	type testCase struct {
		name     string
		input    int
		hasError bool
	}
	testTable := []testCase{
		{
			name:     "standart input with ID=110",
			input:    110,
			hasError: false,
		},
		{
			name:     "standart input with ID=113",
			input:    113,
			hasError: false,
		},
		{ // TODO add hasError=true chage logic
			name:     "Id that does not exist",
			input:    12345,
			hasError: false,
		},
	}

	for _, test := range testTable {
		err := conn.Delete(context.Background(), test.input)
		if test.hasError {
			assert.Error(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
	fmt.Print("TestDelete Finished!\n")
}
