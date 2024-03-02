package rinha2024q1crebito

import "errors"

var (
	// ErrInvalidParameter é retornado quando um parâmetro é inválido
	ErrInvalidParameter = errors.New("invalid parameter")

	// ErrNotFound é retornado quando um recurso não é encontrado
	ErrNotFound = errors.New("not found")

	// ErrUnprocessable é retornado quando um recurso não pode ser processado
	ErrUnprocessable = errors.New("unprocessable")

	// ErrInternal é retornado quando um erro interno ocorre
	ErrInternal = errors.New("internal error")
)
