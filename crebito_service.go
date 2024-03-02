package rinha2024q1crebito

import (
	"log"
)

type CrebitoRepository interface {
	Creditar(clientId, valor int, descricao string) (*AtualizacaoSaldo, error)
	Debitar(clientId, valor int, descricao string) (*AtualizacaoSaldo, error)
}

type CrebitoService interface {
	DoTransaction(transactionReq *TransactionRequest) (*AtualizacaoSaldo, error)
	GetExtratoCliente(clientId ClientID) (*Extrato, error)
}

type CrebitoServiceImpl struct {
	crebitoRepo CrebitoRepository
}

var (
	// CrebitoServiceImpl implementa CrebitoService
	_ CrebitoService = (*CrebitoServiceImpl)(nil)
)

func NewCrebitoServiceImpl(crebitoRepo CrebitoRepository) *CrebitoServiceImpl {
	return &CrebitoServiceImpl{
		crebitoRepo: crebitoRepo,
	}
}

func (cs *CrebitoServiceImpl) DoTransaction(transactionReq *TransactionRequest) (*AtualizacaoSaldo, error) {
	err := transactionReq.validar()
	if err != nil {
		return nil, err
	}

	if transactionReq.Tipo == "c" {
		return cs.crebitoRepo.Creditar(
			transactionReq.ClientId.Value(),
			transactionReq.Valor,
			transactionReq.Descricao)
	} else if transactionReq.Tipo == "d" {
		return cs.crebitoRepo.Debitar(
			transactionReq.ClientId.Value(),
			transactionReq.Valor,
			transactionReq.Descricao)
	} else {
		log.Println("Erro inesperado, tipo de transação não reconhecido.")
		return nil, ErrInternal
	}
}

func (cs *CrebitoServiceImpl) GetExtratoCliente(clientId ClientID) (*Extrato, error) {
	panic("not implemented")
}
