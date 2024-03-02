package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	env "github.com/Netflix/go-env"
	"github.com/julienschmidt/httprouter"
)

type Environment struct {
	PORT int `env:"PORT,default=3000"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
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

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	err = http.ListenAndServe(":"+strconv.Itoa(environment.PORT), router)
	if err != nil {
		return err
	}

	return nil
}
