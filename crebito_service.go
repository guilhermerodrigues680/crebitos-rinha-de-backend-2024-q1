package rinha2024q1crebito

import (
	"context"
	"log"
)

type CrebitoRepository interface {
	Creditar(ctx context.Context, clientId, valor int, descricao string) (*AtualizacaoSaldo, error)
	Debitar(ctx context.Context, clientId, valor int, descricao string) (*AtualizacaoSaldo, error)
	GetExtratoCliente(clientId int) (*Extrato, error)
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

	ctx := context.Background()

	if transactionReq.Tipo == "c" {
		return cs.crebitoRepo.Creditar(
			ctx,
			transactionReq.ClientId.Value(),
			transactionReq.Valor,
			transactionReq.Descricao)
	} else if transactionReq.Tipo == "d" {
		return cs.crebitoRepo.Debitar(
			ctx,
			transactionReq.ClientId.Value(),
			transactionReq.Valor,
			transactionReq.Descricao)
	} else {
		log.Println("Erro inesperado, tipo de transação não reconhecido.")
		return nil, ErrInternal
	}
}

func (cs *CrebitoServiceImpl) GetExtratoCliente(clientId ClientID) (*Extrato, error) {
	err := clientId.validate()
	if err != nil {
		return nil, err
	}

	return cs.crebitoRepo.GetExtratoCliente(clientId.Value())
}
