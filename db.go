package goa

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// newDBClient ...
func (config *Config) newDBClient() *pgxpool.Pool {
	// NOTE: Configure database
	// https://pkg.go.dev/github.com/jackc/pgx/v4@v4.11.0/pgxpool#ParseConfig

	// Example URL
	// "postgres://username:password@localhost:5432/database_name"

	// Example DSN
	// "user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10"

	var connString string

	if os.Getenv("DATABASE_URL") != "" {
		connString = os.Getenv("DATABASE_URL")
	} else {
		connString = fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable pool_max_conns=5",
			config.Database.Username,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Database,
		)
	}

	// NOTE: Using github.com/jackc/pgx/v4
	// https://pkg.go.dev/github.com/jackc/pgx/v4@v4.11.0/pgxpool#pkg-overview

	dbConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse database config: %v\n", err)
		os.Exit(1)
	}

	// dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	// 	// do something with every new connection
	// }

	dbpool, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}
