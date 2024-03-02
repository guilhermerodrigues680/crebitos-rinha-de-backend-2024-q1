package main

import (
	"log"
	"net/http"
	"os"
	"rinha2024q1crebito/api"
	"strconv"

	env "github.com/Netflix/go-env"
)

type Environment struct {
	PORT int `env:"PORT,default=3000"`
}

func main() {
	log.SetOutput(os.Stdout)
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		return err
	}

	apiHandler := api.NewApiHandler()
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(environment.PORT),
		Handler: apiHandler,
	}
	log.Println("Servidor HTTP ouvindo no endereço:", server.Addr)
	return server.ListenAndServe()
}
