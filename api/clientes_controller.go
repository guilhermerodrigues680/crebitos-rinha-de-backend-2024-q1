package api

import (
	"net/http"
	"rinha2024q1crebito"
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
	Saldo             *ExtratoSaldoResponse               `json:"saldo"`
	UltimasTransacoes []*ExtratoUltimasTransacoesResponse `json:"ultimas_transacoes"`
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
	crebitoService rinha2024q1crebito.CrebitoService
}

func NewClientesController(crebitoService rinha2024q1crebito.CrebitoService) *ClientesController {
	return &ClientesController{
		crebitoService: crebitoService,
	}
}

func (cc *ClientesController) PostClientesIdTransacoes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clientIdNum, err := parseInt(ps.ByName("id"))
	if err != nil {
		sendErrorResponse(err, w)
		return
	}

	clientId, err := rinha2024q1crebito.NewClientID(clientIdNum)
	if err != nil {
		sendErrorResponse(err, w)
		return
	}
	var transacaoRequest TransacaoRequest
	err = parseJsonRequest(r, &transacaoRequest)
	if err != nil {
		sendErrorResponse(err, w)
		return
	}

	tReq, err := rinha2024q1crebito.NewTransactionRequest(
		clientId,
		transacaoRequest.Valor,
		transacaoRequest.Tipo,
		transacaoRequest.Descricao,
	)
	if err != nil {
		sendErrorResponse(err, w)
		return
	}

	balanceUpdate, err := cc.crebitoService.DoTransaction(tReq)
	if err != nil {
		sendErrorResponse(err, w)
		return
	}

	res := TransacaoResponse{
		Limite: balanceUpdate.Limite,
		Saldo:  balanceUpdate.Saldo,
	}

	sendOkJsonResponse(res, w)
}

func (cc *ClientesController) GetClientesIdExtrato(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clientIdNum, err := parseInt(ps.ByName("id"))
	if err != nil {
		sendErrorResponse(err, w)
		return
	}

	clientId, err := rinha2024q1crebito.NewClientID(clientIdNum)
	if err != nil {
		sendErrorResponse(err, w)
		return
	}

	extratoCliente, err := cc.crebitoService.GetExtratoCliente(clientId)
	if err != nil {
		sendErrorResponse(err, w)
		return
	}

	ultimasTransacoes := make([]*ExtratoUltimasTransacoesResponse, 0, len(extratoCliente.Transacoes))
	for _, t := range extratoCliente.Transacoes {
		ultimasTransacoes = append(ultimasTransacoes, &ExtratoUltimasTransacoesResponse{
			Valor:       t.Valor,
			Tipo:        t.Tipo,
			Descricao:   t.Descricao,
			RealizadaEm: t.RealizadaEm,
		})
	}

	res := ExtratoResponse{
		Saldo: &ExtratoSaldoResponse{
			Total:       extratoCliente.Saldo.Total,
			DataExtrato: extratoCliente.Saldo.DataExtrato,
			Limite:      extratoCliente.Saldo.Limite,
		},
		UltimasTransacoes: ultimasTransacoes,
	}

	sendOkJsonResponse(res, w)
}
