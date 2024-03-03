package rinha2024q1crebito

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidParameter é retornado quando um parâmetro é inválido
	ErrInvalidParameter = errors.New("invalid parameter")

	// ErrNotFound é retornado quando um recurso não é encontrado
	ErrNotFound = errors.New("not found")

	// ErrUnprocessable é retornado quando um recurso não pode ser processado
	ErrUnprocessable = errors.New("unprocessable")

	// XXX: Eu realmente preciso de um erro interno? Um erro sem tipo não seria suficiente?
	// ErrInternal é retornado quando um erro interno ocorre
	ErrInternal = errors.New("internal error")

	// Erros específicos

	ErrSaldoInsuficiente = fmt.Errorf("saldo insuficiente: (%w)", ErrUnprocessable)
)

func NewErrInternal(message string, causes ...error) error {
	if len(causes) == 0 {
		return fmt.Errorf("%s: (%w)", message, ErrInternal)
	}
	return fmt.Errorf("%s: %w (%w)", message, errors.Join(causes...), ErrInternal)
}
