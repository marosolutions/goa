package goapi

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// newDBClient ...
func (config *Config) newDBClient() *pgxpool.Pool {
	// Example dbUrl
	// "postgres://username:password@localhost:5432/database_name"

	var dbUrl string
	if os.Getenv("DATABASE_URL") != "" {
		dbUrl = os.Getenv("DATABASE_URL")
	} else {
		dbUrl = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
			config.Database.Username,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Database,
		)
	}

	dbpool, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}
