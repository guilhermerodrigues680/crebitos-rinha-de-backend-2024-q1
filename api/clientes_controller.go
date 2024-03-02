package api

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type TransacaoRequest struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type TransacaoResponse struct {
	Limite int `json:"limite"`
	Saldo  int `json:"saldo"`
}

type ExtratoResponse struct {
	Saldo             ExtratoSaldoResponse               `json:"saldo"`
	UltimasTransacoes []ExtratoUltimasTransacoesResponse `json:"ultimas_transacoes"`
}
type ExtratoSaldoResponse struct {
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int       `json:"limite"`
}
type ExtratoUltimasTransacoesResponse struct {
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

type ClientesController struct {
}

func NewClientesController() *ClientesController {
	return &ClientesController{}
}

func (*ClientesController) PostClientesIdTransacoes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clienteId := ps.ByName("id")

	var transacaoRequest TransacaoRequest
	err := parseJsonRequest(r, &transacaoRequest)
	if err != nil {
		sendErrorResponse(err, w)
		return
	}

	// TODO: chama o serviço de transações
	_ = clienteId

	res := TransacaoResponse{
		Limite: 1000,
		Saldo:  500,
	}

	sendOkJsonResponse(res, w)
}

func (*ClientesController) GetClientesIdExtrato(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clienteId := ps.ByName("id")

	// TODO: chama o serviço de extrato
	_ = clienteId

	res := ExtratoResponse{
		Saldo: ExtratoSaldoResponse{
			Total:       -9098,
			DataExtrato: time.Now(),
			Limite:      100000,
		},
		UltimasTransacoes: []ExtratoUltimasTransacoesResponse{
			{
				Valor:       10,
				Tipo:        "c",
				Descricao:   "descricao",
				RealizadaEm: time.Now(),
			},
			{
				Valor:       90000,
				Tipo:        "d",
				Descricao:   "descricao",
				RealizadaEm: time.Now(),
			},
		},
	}

	sendOkJsonResponse(res, w)
}
