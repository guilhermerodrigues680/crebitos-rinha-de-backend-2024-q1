package postgres

import (
	"context"
	"log"
	"rinha2024q1crebito"

	"github.com/jackc/pgx/v5"
)

type PostgresClient struct {
	conn *pgx.Conn
}

func NewPostgresClient(ctx context.Context, dbHostname string) (*PostgresClient, error) {
	// TODO: Trocar para pgxpool
	// DB_INITIAL_POOL_SIZE
	// DB_MAX_POOL_SIZE

	connConf, err := pgx.ParseConfig("postgres://admin:123@" + dbHostname + ":5432/rinha")
	if err != nil {
		return nil, rinha2024q1crebito.NewErrInternal("erro no parse da configuração de conexão", err)
	}

	conn, err := pgx.ConnectConfig(ctx, connConf)
	if err != nil {
		return nil, rinha2024q1crebito.NewErrInternal("erro na conexão com banco de dados", err)
	}
	log.Println("Conexão com banco de dados estabelecida")

	err = conn.Ping(ctx)
	if err != nil {
		return nil, rinha2024q1crebito.NewErrInternal("erro no ping do banco de dados", err)
	}
	log.Println("Ping do banco de dados realizado")

	return &PostgresClient{
		conn: conn,
	}, nil
}
