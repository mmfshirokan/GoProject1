package main

import (
	"fmt"
	"handlers"
	"os"
	"repository"
	"service"
)

func main() {
	err := repository.OpenConnection()
	defer repository.CloseConnection()

	if err != nil {
		fmt.Println("Error ocured while opening connection: ", err)
		os.Exit(1)
	}

	h := handlers.NewUserHandler()

	service.StartServer(h.GetUser, handlers.SaveUser, handlers.UpdateUser, handlers.DeleteUser)
}
