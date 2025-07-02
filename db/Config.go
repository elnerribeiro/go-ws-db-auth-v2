package db

import (
	"context"
	"github.com/elnerribeiro/go-ws-db-auth-v2/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magiconair/properties"
	"github.com/rs/zerolog"
	"time"
)

func Config() *pgxpool.Config {
	logger, _ := utils.GetLoggerAndContext()
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	fileNameStr := "application.properties"
	propertyURLStr := "db.url"
	fileName := &fileNameStr
	propertyURL := &propertyURLStr
	p := properties.MustLoadFile(*fileName, properties.UTF8)
	dbURL := p.GetString(*propertyURL, "")

	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create a config, error: ")
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		logger.Info().Msg("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		logger.Info().Msg("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		logger.Info().Msg("Closed the connection pool to the database!!")
	}

	return dbConfig
}

func FinalizeDB(logger *zerolog.Logger) {
	logger.Info().Msg("DB finalizer executed.")
	pool := GetDBConnPool()
	if pool != nil {
		logger.Info().Msg("DB connection pool closed.")
		pool.Close()
	}
}
