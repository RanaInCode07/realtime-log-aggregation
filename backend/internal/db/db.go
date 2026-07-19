package db

import (
	"context"
	"fmt"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schema string

func InitDB(connString string) (*pgxpool.Pool, error){
	config, parseErr := pgxpool.ParseConfig(connString)
	if parseErr != nil{
		return nil, fmt.Errorf("Error during parsing connection string %v", parseErr)
	}
	config.MaxConns = 20;
	config.MinConns = 5;

	pool, connErr := pgxpool.NewWithConfig(context.Background(), config)
	if connErr != nil{
		return nil, fmt.Errorf("Error during connecting through config %v", connErr)
	}
	if pingErr := pool.Ping(context.Background()); pingErr != nil {
		return nil, fmt.Errorf("Database is not reachable %v", pingErr)
	}
	_, execErr := pool.Exec(context.Background(), string(schema))
	if execErr != nil {
		return nil, fmt.Errorf("Error during exec %v", execErr)
	}

	return pool, nil
}