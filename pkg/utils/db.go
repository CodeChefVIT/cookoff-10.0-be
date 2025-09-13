package utils

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/jackc/pgx/v5/pgtype"
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

// InterfaceToNumeric converts an interface{} to pgtype.Numeric
// Supports string, float64, int, int64, and float32 types
func InterfaceToNumeric(val interface{}) (pgtype.Numeric, error) {
	numeric := pgtype.Numeric{}

	switch v := val.(type) {
	case string:
		// Try to parse as float first
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			err := numeric.Scan(f)
			return numeric, err
		}
		// If not a valid float, try as int
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			err := numeric.Scan(i)
			return numeric, err
		}
		return numeric, fmt.Errorf("could not convert string to numeric: %v", v)

	case float64:
		err := numeric.Scan(v)
		return numeric, err

	case float32:
		err := numeric.Scan(float64(v))
		return numeric, err

	case int:
		err := numeric.Scan(int64(v))
		return numeric, err

	case int64:
		err := numeric.Scan(v)
		return numeric, err

	case uint:
		err := numeric.Scan(int64(v))
		return numeric, err

	case uint64:
		// Handle potential overflow
		if v > 1<<63-1 {
			return numeric, fmt.Errorf("value too large for int64: %v", v)
		}
		err := numeric.Scan(int64(v))
		return numeric, err

	case int32:
		err := numeric.Scan(int64(v))
		return numeric, err

	case uint32:
		err := numeric.Scan(int64(v))
		return numeric, err

	case int16:
		err := numeric.Scan(int64(v))
		return numeric, err

	case uint16:
		err := numeric.Scan(int64(v))
		return numeric, err

	case int8:
		err := numeric.Scan(int64(v))
		return numeric, err

	case uint8:
		err := numeric.Scan(int64(v))
		return numeric, err

	default:
		return numeric, fmt.Errorf("unsupported type for numeric conversion: %T", v)
	}
}
