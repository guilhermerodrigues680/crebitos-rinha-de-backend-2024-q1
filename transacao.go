package rinha2024q1crebito

import "fmt"

type TransactionRequest struct {
	ClientId ClientID

	// deve ser um número inteiro positivo que representa centavos
	// (não vamos trabalhar com frações de centavos). Por exemplo,
	// R$ 10 são 1000 centavos.
	Valor int

	// deve ser apenas c para crédito ou d para débito.
	Tipo string // XXX: usar rune?

	// deve ser uma string de 1 a 10 caracteres.
	Descricao string
}

func NewTransactionRequest(clientId ClientID, valor int, tipo, descricao string) (*TransactionRequest, error) {
	t := &TransactionRequest{
		ClientId:  clientId,
		Valor:     valor,
		Tipo:      tipo,
		Descricao: descricao,
	}
	if err := t.validar(); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *TransactionRequest) validar() error {
	if t.Valor <= 0 {
		return fmt.Errorf(
			"valor da transação deve ser positivo: %w",
			ErrInvalidParameter)
	}

	if t.Tipo != "c" && t.Tipo != "d" {
		return fmt.Errorf(
			"tipo da transação deve ser c ou d: %w",
			ErrInvalidParameter)
	}

	if len(t.Descricao) < 1 || len(t.Descricao) > 10 {
		return fmt.Errorf(
			"descrição da transação deve ter entre 1 e 10 caracteres: %w",
			ErrInvalidParameter)
	}

	return nil
}
