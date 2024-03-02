package rinha2024q1crebito

import "time"

type AtualizacaoSaldo struct {
	// deve ser o limite cadastrado do cliente.
	Limite int

	// deve ser o novo saldo após a conclusão da transação.
	Saldo int
}

type Extrato struct {
	Saldo      *ExtratoSaldo
	Transacoes []*ExtratoTransacao
}

type ExtratoSaldo struct {
	Total       int
	DataExtrato time.Time
	Limite      int
}

type ExtratoTransacao struct {
	Valor       int
	Tipo        string
	Descricao   string
	RealizadaEm time.Time
}
