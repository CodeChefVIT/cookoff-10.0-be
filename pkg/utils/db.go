package utils

import (
	"context"
	"fmt"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Queries *db.Queries

func InitDB() {
	dbHost := Config.PostgresHost
	dbUser := Config.PostgresUser
	dbPassword := Config.PostgresPassword
	dbName := Config.PostgresDB
	dbPort := Config.PostgresPort

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
		logger.Errorf("Database connection parameters are not set")
		return
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
	)

	var err error
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Errorf(err.Error())
		panic(err)
	}

	logger.Infof("Connected to the postgres successfully")
	Queries = db.New(pool)
	Ping(pool)
}

func Ping(pool *pgxpool.Pool) {
	if pool == nil {
		logger.Errorf("Postgres connection is not initialized")
		return
	}

	ctx := context.Background()
	err := pool.Ping(ctx)
	if err != nil {
		logger.Errorf("Unable to ping the postgres: %v", err)
		return
	}

	logger.Infof("Postgres ping successful")
}
