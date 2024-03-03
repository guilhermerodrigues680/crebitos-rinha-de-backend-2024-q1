package postgres

import (
	"context"
	"log"
	"rinha2024q1crebito"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	pool *pgxpool.Pool
}

func NewPostgresClient(ctx context.Context, dbHostname string, DbInitialPoolSize int, DbMaxPoolSize int) (*PostgresClient, error) {
	poolConf, err := pgxpool.ParseConfig("postgres://admin:123@" + dbHostname + ":5432/rinha")
	if err != nil {
		return nil, rinha2024q1crebito.NewErrInternal("erro no parse da configuração de conexão", err)
	}

	poolConf.MinConns = int32(DbInitialPoolSize)
	poolConf.MaxConns = int32(DbMaxPoolSize)

	pool, err := pgxpool.NewWithConfig(ctx, poolConf)
	if err != nil {
		return nil, rinha2024q1crebito.NewErrInternal("erro na conexão com banco de dados", err)
	}
	log.Println("Conexão com banco de dados estabelecida")

	err = pool.Ping(ctx)
	if err != nil {
		return nil, rinha2024q1crebito.NewErrInternal("erro no ping do banco de dados", err)
	}
	log.Println("Ping do banco de dados realizado")

	return &PostgresClient{
		pool: pool,
	}, nil
}

func (pc *PostgresClient) Conn() *pgxpool.Pool {
	return pc.pool
}
