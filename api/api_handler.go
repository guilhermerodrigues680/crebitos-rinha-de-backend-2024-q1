package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ApiHandler struct {
	router *httprouter.Router
}

var (
	// ApiServer implements http.Handler
	_ http.Handler = (*ApiHandler)(nil)
)

func NewApiHandler() *ApiHandler {

	apiHandler := &ApiHandler{
		router: httprouter.New(),
	}

	// Endpoints na raiz "/"
	apiHandler.router.GET("/", apiHandler.Index)
	apiHandler.router.GET("/hello/:name", apiHandler.Hello)

	// Endpoints para clientes "/clientes/"
	clientesController := NewClientesController()
	apiHandler.router.POST("/clientes/:id/transacoes", clientesController.PostClientesIdTransacoes)
	apiHandler.router.GET("/clientes/:id/extrato", clientesController.GetClientesIdExtrato)

	return apiHandler
}

func (ah *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ah.router.ServeHTTP(w, r)
}

func (*ApiHandler) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(
		w,
		"Bem-vindo! API HTTP github.com/guilhermerodrigues680/rinha-2024q1-crebito para participação na 2ª edição da Rinha de Backend",
	)
}

func (*ApiHandler) Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
