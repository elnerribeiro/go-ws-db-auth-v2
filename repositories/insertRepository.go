package repositories

import (
	"github.com/elnerribeiro/go-ws-db-auth-v2/db"
	"github.com/elnerribeiro/go-ws-db-auth-v2/utils"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"time"
)

// ListInserts Retrieve one batch of inserts by id
func (insert *Insert) ListInserts(logContext *utils.LoggerContext) (*Insert, error) {
	logger := zerolog.Ctx(*logContext)
	var data db.SqlData
	paramsValue := db.SqlValue{Name: "id", Value: insert.ID}
	data = append(data, paramsValue)
	dbContext := db.GetDBContext()
	query := `
		select id, type, quantity, status, tstampinit, coalesce(tstampend,0) as tstampend from ins_id
		where id = $1
    `
	insert, err := db.SelectOne[Insert](logContext, dbContext, nil, query, data)
	if err != nil {
		logger.Error().Err(err).Msgf("[ListInserts] Error listing inserts: %s", err)
		return nil, err
	}
	queryInsertBatch := `
		select id, id_ins_id, pos from insert_batch
		where id_ins_id = $1
	`
	if insert != nil {
		var data db.SqlData
		paramsValue := db.SqlValue{Name: "id_ins_id", Value: insert.ID}
		data = append(data, paramsValue)

		val2, err := db.SelectAll[InsertBatch](logContext, dbContext, nil, queryInsertBatch, data)
		if err != nil {
			logger.Error().Err(err).Msgf("[ListInserts] Cannot find children: %s", err)
			return insert, nil
		}
		var result []InsertBatch
		for _, v := range val2 {
			result = append(result, v)
		}
		insert.ListVals = result
	}
	return insert, nil
}

// InsertOneBatch Inserts one item of the batch
func (insert *InsertBatch) InsertOneBatch(logContext *utils.LoggerContext, txContext *db.DatabaseContext, tx *pgx.Tx) error {
	logger := zerolog.Ctx(*logContext)

	var data db.SqlData
	paramsValue := db.SqlValue{Name: "id_ins_id", Value: insert.ID_Ins_ID}
	data = append(data, paramsValue)
	paramsValue2 := db.SqlValue{Name: "pos", Value: insert.Pos}
	data = append(data, paramsValue2)

	err := db.Insert(logContext, txContext, tx, "insert_batch", data)
	if err != nil {
		logger.Error().Err(err).Msgf("[InsertOneBatch] Cannot insert children for id %d: %s", insert.ID_Ins_ID, err)
		return err
	}
	return nil
}

// InsertID Inserts one batch of items
func (insert *Insert) InsertID(logContext *utils.LoggerContext, txContext *db.DatabaseContext, tx *pgx.Tx) (*Insert, error) {
	logger := zerolog.Ctx(*logContext)
	var data db.SqlData
	paramsValue := db.SqlValue{Name: "quantity", Value: insert.Quantity}
	data = append(data, paramsValue)
	paramsValue2 := db.SqlValue{Name: "status", Value: "Running"}
	data = append(data, paramsValue2)
	paramsValue3 := db.SqlValue{Name: "type", Value: insert.Type}
	data = append(data, paramsValue3)
	tstamp := time.Now().Unix()
	paramsValue4 := db.SqlValue{Name: "tstampinit", Value: tstamp}
	data = append(data, paramsValue4)
	insertReturn, err := db.InsertReturningPostgres[Insert](logContext, txContext, tx, "ins_id", data, "id")
	if err != nil || insertReturn == nil {
		logger.Error().Err(err).Msgf("[InsertID] Cannot insert new batch: %s", err)
		return nil, err
	}
	insertReturn.Quantity = insert.Quantity
	insertReturn.Status = "Running"
	insertReturn.Type = insert.Type
	insertReturn.Tstampinit = tstamp
	return insertReturn, nil
}

// UpdateInsertID Finishes batch insertion
func (insert *Insert) UpdateInsertID(logContext *utils.LoggerContext, txContext *db.DatabaseContext, tx *pgx.Tx) (int64, error) {
	logger := zerolog.Ctx(*logContext)
	var data db.SqlData
	paramsValue := db.SqlValue{Name: "status", Value: insert.Status}
	data = append(data, paramsValue)
	tstamp := time.Now().Unix()
	paramsValue2 := db.SqlValue{Name: "tstampend", Value: tstamp}
	data = append(data, paramsValue2)

	var filters db.SqlData
	filterValue := db.SqlValue{Name: "id", Value: insert.ID}
	filters = append(filters, filterValue)
	err := db.Update(logContext, txContext, tx, "ins_id", data, filters)
	if err != nil {
		logger.Error().Err(err).Msgf("[UpdateInsertID] Cannot update batch: %s", err)
		return 0, err
	}
	return tstamp, nil
}

// ClearBatches Removes all batches
func (insert *Insert) ClearBatches(logContext *utils.LoggerContext, txContext *db.DatabaseContext, tx *pgx.Tx) error {
	logger := zerolog.Ctx(*logContext)
	_, err := (*tx).Exec(*txContext, "delete from insert_batch where id > 0")
	if err != nil {
		logger.Error().Err(err).Msgf("[ClearBatches] Cannot remove children: %s", err)
		return err
	}
	_, err2 := (*tx).Exec(*txContext, "delete from ins_id where id > 0")
	if err2 != nil {
		logger.Error().Err(err).Msgf("[ClearBatches] Cannot remove batches: %s", err)
		return err2
	}
	return nil
}
