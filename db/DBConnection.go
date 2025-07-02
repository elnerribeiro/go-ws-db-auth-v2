package db

import (
	"context"
	"github.com/elnerribeiro/go-ws-db-auth-v2/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type DatabaseContext context.Context

var connPool *pgxpool.Pool
var dbContext *DatabaseContext

func GetDBConnection(logContext *utils.LoggerContext) *pgxpool.Conn {
	logger := zerolog.Ctx(*logContext)
	// Create database connection
	var err error
	connection, err := connPool.Acquire(*dbContext)
	if err != nil {
		logger.Error().Err(err).Msg("Error while acquiring connection from the database pool!!")
		panic("Error while acquiring connection from the database pool!!")
	}

	err = connection.Ping(*dbContext)
	if err != nil {
		logger.Error().Err(err).Msg("Could not ping database")
	}

	logger.Info().Msg("Connected to the database!!")

	return connection

}

func GetDBContext() *DatabaseContext {
	return dbContext
}

func GetDBConnPool() *pgxpool.Pool {
	return connPool
}

func init() {
	initDBContext := context.Background()
	dbContext = (*DatabaseContext)(&initDBContext)
	logger, _ := utils.GetLoggerAndContext()
	initPool, err := pgxpool.NewWithConfig(*dbContext, Config())
	if err != nil {
		logger.Error().Err(err).Msg("Error while creating connection pool to the database!!")
		panic("Error while creating connection pool to the database!!")
	}
	connPool = initPool
}
