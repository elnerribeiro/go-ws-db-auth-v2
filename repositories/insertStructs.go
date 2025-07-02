package repositories

import (
	"github.com/elnerribeiro/go-ws-db-auth-v2/db"
	"github.com/elnerribeiro/go-ws-db-auth-v2/utils"
	"github.com/jackc/pgx/v5"
)

// InsertBatch table insert_batch on database
type InsertBatch struct {
	ID        int `json:"id,omitempty" db:"id,omitempty"`
	ID_Ins_ID int `json:"id_ins_id,omitempty" db:"id_ins_id,omitempty"`
	Pos       int `json:"pos,omitempty" db:"pos,omitempty"`
}

// Insert table ins_id on database
type Insert struct {
	ID         int           `json:"id,omitempty" db:"id,omitempty"`
	Type       string        `json:"type,omitempty" db:"type,omitempty"`
	Quantity   int           `json:"quantity,omitempty" db:"quantity,omitempty"`
	Status     string        `json:"status,omitempty" db:"status,omitempty"`
	Tstampinit int64         `json:"tstampinit,omitempty" db:"tstampinit,omitempty"`
	Tstampend  *int64        `json:"tstampend,omitempty" db:"tstampend,omitempty"`
	ListVals   []InsertBatch `json:"list,omitempty" db:"-"`
}

// InsertInterface interface for batch insert tables
type InsertInterface interface {
	UpdateInsertID(logContext *utils.LoggerContext, txContext *db.DatabaseContext, tx *pgx.Tx) (int64, error)
	ClearBatches(logContext *utils.LoggerContext, txContext *db.DatabaseContext, tx *pgx.Tx) error
	InsertID(logContext *utils.LoggerContext, txContext *db.DatabaseContext, tx *pgx.Tx) (*Insert, error)
	InsertOneBatch(logContext *utils.LoggerContext, txContext *db.DatabaseContext, tx *pgx.Tx) error
	ListInserts(logContext *utils.LoggerContext) (*Insert, error)
}
