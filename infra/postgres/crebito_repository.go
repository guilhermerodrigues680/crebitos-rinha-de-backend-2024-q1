package postgres

import (
	"context"
	"errors"
	"fmt"
	"rinha2024q1crebito"
	"strings"
	"time"
)

type AtualizacaoSaldoRow struct {
	// novo_saldo INT,
	// possui_erro BOOL,
	// mensagem VARCHAR(20))

	NovoSaldo  int
	PossuiErro bool
	Mensagem   string
}

type CrebitoPostgresRepository struct {
	client       *PostgresClient
	apiUseDbFunc bool
}

var (
	_ rinha2024q1crebito.CrebitoRepository = (*CrebitoPostgresRepository)(nil)
)

func NewCrebitoPostgresRepository(client *PostgresClient, apiUseDbFunc bool) *CrebitoPostgresRepository {

	// TODO: Implementar sem uso de função de banco de dados
	if !apiUseDbFunc {
		panic("apiUseDbFunc must be true")
	}

	return &CrebitoPostgresRepository{
		client:       client,
		apiUseDbFunc: apiUseDbFunc,
	}
}

func (cpr *CrebitoPostgresRepository) Creditar(ctx context.Context, clientId int, valor int, descricao string) (*rinha2024q1crebito.AtualizacaoSaldo, error) {
	atualizacaoSaldoRow := AtualizacaoSaldoRow{}
	err := cpr.client.Conn().
		QueryRow(
			ctx,
			"SELECT * from creditar($1, $2, $3)",
			clientId,
			valor,
			descricao,
		).
		Scan(
			&atualizacaoSaldoRow.NovoSaldo,
			&atualizacaoSaldoRow.PossuiErro,
			&atualizacaoSaldoRow.Mensagem,
		)
	if err != nil {
		return nil, rinha2024q1crebito.NewErrInternal("erro ao creditar", err)
	}

	if atualizacaoSaldoRow.PossuiErro {
		return nil, rinha2024q1crebito.NewErrInternal("erro ao creditar", errors.New(atualizacaoSaldoRow.Mensagem))
	}

	return &rinha2024q1crebito.AtualizacaoSaldo{
		Limite: 0, // FIXME: Implementar limite
		Saldo:  atualizacaoSaldoRow.NovoSaldo,
	}, nil

}

func (cpr *CrebitoPostgresRepository) Debitar(ctx context.Context, clientId int, valor int, descricao string) (*rinha2024q1crebito.AtualizacaoSaldo, error) {
	atualizacaoSaldoRow := AtualizacaoSaldoRow{}
	err := cpr.client.Conn().
		QueryRow(
			ctx,
			"SELECT * from debitar($1, $2, $3)",
			clientId,
			valor,
			descricao,
		).
		Scan(
			&atualizacaoSaldoRow.NovoSaldo,
			&atualizacaoSaldoRow.PossuiErro,
			&atualizacaoSaldoRow.Mensagem,
		)
	if err != nil {
		return nil, rinha2024q1crebito.NewErrInternal("erro ao debitar", err)
	}

	if atualizacaoSaldoRow.PossuiErro {
		if strings.Contains(atualizacaoSaldoRow.Mensagem, "saldo insuficente") {
			return nil, fmt.Errorf(
				"não é possível debitar: %s (%w)",
				atualizacaoSaldoRow.Mensagem,
				rinha2024q1crebito.ErrSaldoInsuficiente)
		}

		return nil, rinha2024q1crebito.NewErrInternal("erro ao debitar", errors.New(atualizacaoSaldoRow.Mensagem))
	}

	return &rinha2024q1crebito.AtualizacaoSaldo{
		Limite: 0, // FIXME: Implementar limite
		Saldo:  atualizacaoSaldoRow.NovoSaldo,
	}, nil
}

func (cpr *CrebitoPostgresRepository) GetExtratoCliente(clientId int) (*rinha2024q1crebito.Extrato, error) {

	// Faz tudo em uma única transação para garantir consistência
	tx, err := cpr.client.Conn().Begin(context.TODO())
	if err != nil {
		return nil, rinha2024q1crebito.NewErrInternal("erro ao buscar extrato, iniciar transação", err)
	}
	defer tx.Rollback(context.TODO())

	type ExtratoRow struct {
		Valor       int
		DataExtrato time.Time
	}

	type TransacaoRow struct {
		Valor       int
		Tipo        string
		Descricao   string
		RealizadaEm time.Time
	}

	var clientLimit int
	err = tx.
		QueryRow(
			context.Background(),
			"SELECT limite FROM clientes WHERE id = $1",
			clientId,
		).
		Scan(&clientLimit)
	if err != nil {
		// Checa se o erro é de cliente não encontrado
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("cliente não encontrado: (%w)", rinha2024q1crebito.ErrNotFound)
		}

		return nil, rinha2024q1crebito.NewErrInternal("erro ao buscar extrato, buscar limite", err)
	}

	extratoRow := ExtratoRow{}
	err = tx.
		QueryRow(
			context.Background(),
			"SELECT valor, NOW() FROM saldos WHERE cliente_id = $1",
			clientId,
		).
		Scan(
			&extratoRow.Valor,
			&extratoRow.DataExtrato,
		)
	if err != nil {
		return nil, rinha2024q1crebito.NewErrInternal("erro ao buscar extrato, buscar saldo", err)
	}

	transacoes := []*rinha2024q1crebito.ExtratoTransacao{}
	rows, err := tx.
		Query(
			context.Background(),
			"SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE cliente_id = $1",
			clientId,
		)
	if err != nil {
		return nil, rinha2024q1crebito.NewErrInternal("erro ao buscar extrato, buscar transações", err)
	}
	defer rows.Close()
	for rows.Next() {
		transacaoRow := TransacaoRow{}
		err := rows.Scan(
			&transacaoRow.Valor,
			&transacaoRow.Tipo,
			&transacaoRow.Descricao,
			&transacaoRow.RealizadaEm,
		)
		if err != nil {
			return nil, rinha2024q1crebito.NewErrInternal("erro ao buscar extrato, ler transação", err)
		}

		transacoes = append(transacoes, &rinha2024q1crebito.ExtratoTransacao{
			Valor:       transacaoRow.Valor,
			Tipo:        transacaoRow.Tipo,
			Descricao:   transacaoRow.Descricao,
			RealizadaEm: transacaoRow.RealizadaEm,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, rinha2024q1crebito.NewErrInternal("erro ao buscar extrato, ler transações", err)
	}

	extrato := &rinha2024q1crebito.Extrato{
		Saldo: &rinha2024q1crebito.ExtratoSaldo{
			Total:       extratoRow.Valor,
			DataExtrato: extratoRow.DataExtrato,
			Limite:      clientLimit,
		},
		Transacoes: transacoes,
	}

	return extrato, nil
}
