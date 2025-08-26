package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type PostgresConfig struct {
	Uri string
}

func SetupPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	config := &PostgresConfig{}

	err := viper.UnmarshalKey("postgres.database", config)
	if err != nil {
		return nil, err
	}

	client, err := pgxpool.New(ctx, config.Uri)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
