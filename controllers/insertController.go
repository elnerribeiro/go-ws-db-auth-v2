package controllers

import (
	"net/http"

	"strconv"

	repo "github.com/elnerribeiro/go-ws-db-auth-v2/repositories"
	"github.com/elnerribeiro/go-ws-db-auth-v2/services"
	u "github.com/elnerribeiro/go-ws-db-auth-v2/utils"
	"github.com/gorilla/mux"
)

// ClearInserts Clears the database
var ClearInserts = func(w http.ResponseWriter, r *http.Request) {
	logger, logContext := u.GetLoggerAndContext()
	insert := &repo.Insert{}
	err := services.ClearInserts(logContext, insert)
	if err != nil {
		logger.Error().Msgf("[ClearInserts] Error while cleaning inserts: %s", err)
		resp := u.Message(false, "Error while cleaning batches")
		u.Respond(logContext, w, resp)
		return
	}
	u.Respond(logContext, w, u.Message(true, "success"))
}

// ListInsert Lists one insert batch
var ListInsert = func(w http.ResponseWriter, r *http.Request) {
	logger, logContext := u.GetLoggerAndContext()
	vars := mux.Vars(r)
	uID, _ := strconv.Atoi(vars["id"])
	insert := &repo.Insert{}
	insert.ID = uID
	data, err := services.ListInserts(logContext, insert)
	if err != nil {
		logger.Error().Msgf("[ListInsert] Error listing inserts: %s", err)
		resp := u.Message(false, "Error querying batch")
		u.Respond(logContext, w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(logContext, w, resp)
}

// InsertSync Inserts a batch of given quantity sync
var InsertSync = func(w http.ResponseWriter, r *http.Request) {
	logger, logContext := u.GetLoggerAndContext()
	vars := mux.Vars(r)
	qty, _ := strconv.Atoi(vars["qty"])
	insert := &repo.Insert{}
	insert.Quantity = qty
	data, err := services.InsertBatchSync(logContext, insert)
	if err != nil {
		logger.Error().Msgf("[InsertSync] Error inserting batch synchronous: %s", err)
		resp := u.Message(false, "Error inserting batch")
		u.Respond(logContext, w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(logContext, w, resp)
}

// InsertASync Inserts a batch of given quantity async
var InsertASync = func(w http.ResponseWriter, r *http.Request) {
	logger, logContext := u.GetLoggerAndContext()
	vars := mux.Vars(r)
	qty, _ := strconv.Atoi(vars["qty"])
	insert := &repo.Insert{}
	insert.Quantity = qty
	data, err := services.InsertBatchASync(logContext, insert)
	if err != nil {
		logger.Error().Msgf("[InsertASync] Error inserting batch asynchronous: %s", err)
		resp := u.Message(false, "Error inserting batch async")
		u.Respond(logContext, w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(logContext, w, resp)
}
