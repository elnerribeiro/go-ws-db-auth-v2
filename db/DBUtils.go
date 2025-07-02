package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/elnerribeiro/go-ws-db-auth-v2/utils"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"strconv"
)

type SqlValue struct {
	Name  string
	Value any
}
type SqlData []SqlValue

// GetTransaction returns a DB transaction
func GetTransaction() (*pgx.Tx, *DatabaseContext, error) {
	logger, logContext := utils.GetLoggerAndContext()

	database := GetDBConnection(logContext)
	if database == nil {
		logger.Error().Msgf("[GetTransaction] Database not initialized.")
		return nil, nil, errors.New("Database not initialized")
	}
	dbContext := GetDBContext()
	tx, err := database.Begin(*dbContext)
	if err != nil {
		logger.Error().Err(err).Msgf("[GetTransaction] Error starting transaction: %s", err)
		return nil, nil, err
	}
	transactionContext := context.Background()
	return &tx, (*DatabaseContext)(&transactionContext), nil
}

// Delete removes rows from a table, given filters
func Delete(logContext *utils.LoggerContext, txContext *DatabaseContext, transaction *pgx.Tx, table string, filters SqlData) error {
	logger := zerolog.Ctx(*logContext)
	if transaction == nil {
		logger.Error().Msgf("[Delete] Transaction not found.")
		return errors.New("not inside transaction")
	}
	_, err := (*transaction).Exec(*txContext, generateDeleteExpression(table, filters), generateArguments(filters)...)
	if err != nil {
		logger.Error().Err(err).Msgf("[Delete] Error executing delete query: %s", err)
		return err
	}
	return nil
}

func generateDeleteExpression(table string, filters SqlData) string {
	data := "delete from " + table + " where 1=1"
	if len(filters) > 0 {
		count := 0
		for _, v := range filters {
			count++
			data += " and " + v.Name + " = $" + strconv.Itoa(count)
		}
	}
	return data
}

func generateArguments(filters SqlData) []any {
	var arr []any
	for _, v := range filters {
		arr = append(arr, getValue(v))
	}
	return arr
}

func getValue(value SqlValue) any {
	if value.Value == nil {
		return nil
	}
	switch value.Value.(type) {
	default:
		return fmt.Sprintf("%s", value.Value)
	case int64:
		return value.Value.(int64)
	case int32:
		return value.Value.(int32)
	case int16:
		return value.Value.(int16)
	case int8:
		return value.Value.(int8)
	case int:
		x, _ := strconv.Atoi(fmt.Sprintf("%d", value.Value))
		return x
	case float64:
		return value.Value.(float64)
	case float32:
		return value.Value.(float32)
	}
}

func generateArgumentsForUpdate(filters SqlData, params SqlData) []any {
	var arr []any
	for _, v := range params {
		arr = append(arr, getValue(v))
	}
	for _, v := range filters {
		arr = append(arr, getValue(v))
	}
	return arr
}

// Update updates rows on a table, given filters
func Update(logContext *utils.LoggerContext, txContext *DatabaseContext, transaction *pgx.Tx, table string, params SqlData, filters SqlData) error {
	logger := zerolog.Ctx(*logContext)
	if transaction == nil {
		logger.Error().Msgf("[Update] Transaction not found.")
		return errors.New("not inside transaction")
	}
	update := generateUpdateExpression(table, filters, params)
	queryParams := generateArgumentsForUpdate(filters, params)

	_, err := (*transaction).Exec(*txContext, update, queryParams...)
	if err != nil {
		logger.Error().Err(err).Msgf("[Update] Error executing update query: %s", err)
		return err
	}
	return nil
}

func addCommaIfNeeded(count int) string {
	if count > 1 {
		return ", "
	} else {
		return " "
	}
}

func generateUpdateExpression(table string, filters SqlData, params SqlData) string {
	data := "update " + table + " set "
	count := 0
	if len(params) > 0 {
		for _, k := range params {
			count++
			data += addCommaIfNeeded(count) + k.Name + " = $" + strconv.Itoa(count)
		}
	}
	data = data + " where 1=1"
	if len(filters) > 0 {
		for _, k := range filters {
			count++
			data += " and " + k.Name + " = $" + strconv.Itoa(count)
		}
	}
	return data
}

// Insert inserts a row on the table
func Insert(logContext *utils.LoggerContext, txContext *DatabaseContext, transaction *pgx.Tx, table string, params SqlData) error {
	logger := zerolog.Ctx(*logContext)
	if transaction == nil {
		logger.Error().Msgf("[Insert] Transaction not found.")
		return errors.New("not inside transaction")
	}

	_, err := (*transaction).Exec(*txContext, generateInsertExpression(table, params), generateArguments(params)...)
	if err != nil {
		logger.Error().Err(err).Msgf("[Insert] Error executing insert query: %s", err)
		return err
	}
	return nil
}

func generateInsertExpression(table string, params SqlData) string {
	data := "insert into " + table + " ("
	if len(params) > 0 {
		count := 0
		for _, k := range params {
			count++
			data += addCommaIfNeeded(count) + k.Name
		}
	}
	data = data + ") values ("
	if len(params) > 0 {
		count := 0
		for range params {
			count++
			data += addCommaIfNeeded(count) + "$" + strconv.Itoa(count)
		}
	}
	data = data + ")"
	return data
}

// InsertReturningPostgres execs an insert and returns inserted id
func InsertReturningPostgres[T interface{}](logContext *utils.LoggerContext, txContext *DatabaseContext, transaction *pgx.Tx, table string, params SqlData, pk string) (*T, error) {
	logger := zerolog.Ctx(*logContext)
	if transaction == nil {
		logger.Error().Msgf("[Insert] Transaction not found.")
		return nil, errors.New("not inside transaction")
	}

	var id int32
	err := (*transaction).QueryRow(*txContext, generateInsertReturningExpression(table, params, pk), generateArguments(params)...).Scan(&id)
	if err != nil {
		logger.Error().Err(err).Msgf("[Insert] Error executing insert query: %s", err)
		return nil, err
	}
	paramsValue := SqlValue{Name: pk, Value: id}
	var paramsQuery SqlData
	paramsQuery = append(paramsQuery, paramsValue)
	return SelectOne[T](logContext, txContext, transaction, "select * from "+table+" where "+pk+" = $1", paramsQuery)
}

func generateInsertReturningExpression(table string, params SqlData, pk string) string {
	data := "insert into " + table + " ("
	if len(params) > 0 {
		count := 0
		for _, k := range params {
			count++
			data += addCommaIfNeeded(count) + k.Name
		}
	}
	data = data + ") values ("
	if len(params) > 0 {
		count := 0
		for range params {
			count++
			data += addCommaIfNeeded(count) + "$" + strconv.Itoa(count)
		}
	}
	data = data + ") RETURNING " + pk
	return data
}

// SelectAll retrieves all rows from a table, given a params filter
func SelectAll[T interface{}](logContext *utils.LoggerContext, txContext *DatabaseContext, transaction *pgx.Tx, query string, params SqlData) ([]T, error) {
	logger := zerolog.Ctx(*logContext)
	var rows pgx.Rows
	var err error
	if transaction == nil || txContext == nil {
		db := GetDBConnection(logContext)
		dbContext := GetDBContext()
		defer db.Release()
		rows, err = db.Query(*dbContext, query, generateArguments(params)...)
	} else {
		rows, err = (*transaction).Query(*txContext, query, generateArguments(params)...)
	}
	if err != nil {
		logger.Error().Err(err).Msgf("[SelectAll] Error executing query - %s: %s", query, err)
		return nil, err
	}

	if rows != nil {
		result, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	logger.Info().Msgf("[SelectAll] No rows found for query: %s", query)
	return nil, nil
}

// SelectOne retrieves the first result
func SelectOne[T interface{}](logContext *utils.LoggerContext, txContext *DatabaseContext, transaction *pgx.Tx, query string, params SqlData) (*T, error) {
	logger := zerolog.Ctx(*logContext)
	result, err := SelectAll[T](logContext, txContext, transaction, query, params)
	if err != nil {
		logger.Error().Err(err).Msgf("[SelectOne] Error executing query - %s: %s", query, err)
		return nil, err
	}
	if result != nil && len(result) > 0 {
		return &result[0], nil
	}
	logger.Info().Msgf("[SelectOne] No rows found for query: %s", query)
	return nil, nil
}

func Rollback(logContext *utils.LoggerContext, txContext *DatabaseContext, transaction *pgx.Tx) {
	logger := zerolog.Ctx(*logContext)
	err := (*transaction).Rollback(*txContext)
	if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		logger.Err(err).Msg("Error rolling back transaction")
	}
}

func Commit(logContext *utils.LoggerContext, txContext *DatabaseContext, transaction *pgx.Tx) {
	logger := zerolog.Ctx(*logContext)
	err := (*transaction).Commit(*txContext)
	if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		logger.Err(err).Msg("Error commiting transaction")
	}
}
