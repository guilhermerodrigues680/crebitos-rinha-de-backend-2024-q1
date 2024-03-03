package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"rinha2024q1crebito"
	"rinha2024q1crebito/api"
	"rinha2024q1crebito/infra/postgres"
	"strconv"

	env "github.com/Netflix/go-env"
)

type Environment struct {
	Port         int    `env:"PORT,default=3000"`
	DbHostname   string `env:"DB_HOSTNAME,required=true"`
	ApiUseDbFunc bool   `env:"API_USE_DB_FUNC,default=false"`
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
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

	// Postgres

	postgresClient, err := postgres.NewPostgresClient(context.Background(), environment.DbHostname)
	if err != nil {
		return err
	}
	crebitoPostgresRepo := postgres.NewCrebitoPostgresRepository(postgresClient, environment.ApiUseDbFunc)

	// Services

	crebitoService := rinha2024q1crebito.NewCrebitoServiceImpl(crebitoPostgresRepo)

	// Server

	apiHandler := api.NewApiHandler(crebitoService)
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(environment.Port),
		Handler: apiHandler,
	}
	log.Println("Servidor HTTP ouvindo no endereço:", server.Addr)
	return server.ListenAndServe()
}
