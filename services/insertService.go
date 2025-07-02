package services

import (
	"github.com/elnerribeiro/go-ws-db-auth-v2/db"
	repo "github.com/elnerribeiro/go-ws-db-auth-v2/repositories"
	u "github.com/elnerribeiro/go-ws-db-auth-v2/utils"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

// ListInserts Lists all inserts for one batch by id
func ListInserts(logContext *u.LoggerContext, insert *repo.Insert) (*repo.Insert, error) {
	return insert.ListInserts(logContext)
}

// ClearInserts Clears tables
func ClearInserts(logContext *u.LoggerContext, insert *repo.Insert) error {
	logger := zerolog.Ctx(*logContext)
	tx, txContext, err := db.GetTransaction()
	if err != nil {
		logger.Error().Err(err).Msgf("[ClearInserts] Error starting transaction: %s", err)
		return err
	}
	defer db.Rollback(logContext, txContext, tx)
	if err2 := insert.ClearBatches(logContext, txContext, tx); err2 != nil {
		logger.Error().Err(err).Msgf("[ClearInserts] Error cleaning batches: %s", err2)
		return err2
	}
	db.Commit(logContext, txContext, tx)
	return nil
}

// InsertBatchSync Inserts a batch of given quantity synchronous
func InsertBatchSync(logContext *u.LoggerContext, insert *repo.Insert) (*repo.Insert, error) {
	logger := zerolog.Ctx(*logContext)
	tx, txContext, err := db.GetTransaction()
	if err != nil {
		logger.Error().Err(err).Msgf("[InsertBatchSync] Error starting transaction: %s", err)
		return nil, err
	}
	defer db.Rollback(logContext, txContext, tx)
	insert.Type = "sync"
	ins, err := insert.InsertID(logContext, txContext, tx)
	if err != nil {
		logger.Error().Err(err).Msgf("[InsertBatchSync] Error inserting a batch: %s", err)
		return nil, err
	}

	for i := 1; i <= insert.Quantity; i++ {
		insertBatch := &repo.InsertBatch{}
		insertBatch.ID_Ins_ID = ins.ID
		insertBatch.Pos = i
		if err := insertBatch.InsertOneBatch(logContext, txContext, tx); err != nil {
			logger.Error().Err(err).Msgf("[InsertBatchSync] Error inserting one item: %s", err)
		}
	}
	ins.Status = "Finished"
	_, err = ins.UpdateInsertID(logContext, txContext, tx)
	if err != nil {
		logger.Error().Err(err).Msgf("[InsertBatchSync] Error updating insert id: %s", err)
		return nil, err
	}
	db.Commit(logContext, txContext, tx)
	return ins.ListInserts(logContext)
}

// InsertBatchASync Inserts a batch of given quantity asynchronous
func InsertBatchASync(logContext *u.LoggerContext, insert *repo.Insert) (*repo.Insert, error) {
	logger := zerolog.Ctx(*logContext)
	tx, txContext, err := db.GetTransaction()
	if err != nil {
		logger.Error().Err(err).Msgf("[InsertBatchASync] Error starting transaction: %s", err)
		return nil, err
	}
	defer db.Rollback(logContext, txContext, tx)
	insert.Type = "async"
	ins, err := insert.InsertID(logContext, txContext, tx)
	ins.Quantity = insert.Quantity
	if err != nil {
		logger.Error().Err(err).Msgf("[InsertBatchASync] Error inserting a batch: %s", err)
		return nil, err
	}
	db.Commit(logContext, txContext, tx)
	go insertBatch(logContext, ins)
	return ins, nil
}

func insertBatch(logContext *u.LoggerContext, insert *repo.Insert) {
	logger := zerolog.Ctx(*logContext)
	tx, txContext, err := db.GetTransaction()
	if err != nil {
		logger.Error().Err(err).Msgf("[insertBatch] Error starting transaction: %s", err)
		return
	}
	defer db.Rollback(logContext, txContext, tx)
	for i := 1; i <= insert.Quantity; i++ {
		insertBatch := &repo.InsertBatch{}
		insertBatch.ID_Ins_ID = insert.ID
		insertBatch.Pos = i
		if err := insertBatch.InsertOneBatch(logContext, txContext, tx); err != nil {
			logger.Error().Err(err).Msgf("[insertBatch] Error inserting one item: %s", err)
			onError(logContext, txContext, tx, insert)
		}
	}
	insert.Status = "Finished"
	_, err = insert.UpdateInsertID(logContext, txContext, tx)
	if err != nil {
		logger.Error().Err(err).Msgf("[InsertBatchSync] Error updating insert id: %s", err)
	} else {
		db.Commit(logContext, txContext, tx)
	}
}

func onError(logContext *u.LoggerContext, txContext *db.DatabaseContext, tx *pgx.Tx, insert *repo.Insert) {
	logger := zerolog.Ctx(*logContext)
	db.Rollback(logContext, txContext, tx)
	newTx, newContext, err := db.GetTransaction()
	if err != nil {
		logger.Error().Msgf("[onError] Error starting transaction: %s", err)
		return
	}
	defer db.Rollback(logContext, newContext, newTx)
	insert.Status = "Error"
	_, err = insert.UpdateInsertID(logContext, newContext, newTx)
	if err != nil {
		logger.Error().Err(err).Msgf("[onError] Error updating insert id: %s", err)
	} else {
		db.Commit(logContext, newContext, newTx)
	}
}
