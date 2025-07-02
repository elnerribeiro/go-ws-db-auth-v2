package controllers

import (
	"encoding/json"
	"net/http"

	"strconv"

	repo "github.com/elnerribeiro/go-ws-db-auth-v2/repositories"
	"github.com/elnerribeiro/go-ws-db-auth-v2/services"
	u "github.com/elnerribeiro/go-ws-db-auth-v2/utils"
	"github.com/gorilla/mux"
)

// ListUsers Lists all users
var ListUsers = func(w http.ResponseWriter, r *http.Request) {
	_, logContext := u.GetLoggerAndContext()
	role := r.Context().Value(repo.ContextKey("role")).(string)
	if role != "admin" {
		resp := u.Message(false, "Unauthorized user")
		u.Respond(logContext, w, resp)
		return
	}
	account := &repo.User{}
	data, err := services.ListUsers(logContext, account)
	if err != nil {
		resp := u.Message(false, "Error searching users")
		u.Respond(logContext, w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(logContext, w, resp)
}

// GetUserByID Get an user by ID
var GetUserByID = func(w http.ResponseWriter, r *http.Request) {
	_, logContext := u.GetLoggerAndContext()
	role := r.Context().Value(repo.ContextKey("role")).(string)
	if role != "admin" {
		resp := u.Message(false, "Unauthorized user")
		u.Respond(logContext, w, resp)
		return
	}
	vars := mux.Vars(r)
	uID, _ := strconv.Atoi(vars["id"])
	account := &repo.User{}
	account.ID = uID
	data, err := services.GetUserByID(logContext, account)
	if err != nil {
		resp := u.Message(false, "Error searching user")
		u.Respond(logContext, w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(logContext, w, resp)
}

// Upsert Inserts or updates an user
var Upsert = func(w http.ResponseWriter, r *http.Request) {
	_, logContext := u.GetLoggerAndContext()
	role := r.Context().Value(repo.ContextKey("role")).(string)
	if role != "admin" {
		resp := u.Message(false, "Unauthorized user")
		u.Respond(logContext, w, resp)
		return
	}
	account := &repo.User{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(logContext, w, u.Message(false, "Invalid request"))
		return
	}

	us, err2 := services.Upsert(logContext, account)
	if err2 != nil {
		resp := u.Message(false, "Error updating user")
		u.Respond(logContext, w, resp)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = us
	u.Respond(logContext, w, resp)
}

// Delete Deletes an user by ID
var Delete = func(w http.ResponseWriter, r *http.Request) {
	_, logContext := u.GetLoggerAndContext()
	role := r.Context().Value(repo.ContextKey("role")).(string)
	if role != "admin" {
		resp := u.Message(false, "Unauthorized user")
		u.Respond(logContext, w, resp)
		return
	}
	vars := mux.Vars(r)
	uID, _ := strconv.Atoi(vars["id"])
	account := &repo.User{}
	account.ID = uID
	if err := services.Delete(logContext, account); err != nil {
		resp := u.Message(false, "Error deleting user")
		u.Respond(logContext, w, resp)
		return
	}
	u.Respond(logContext, w, u.Message(true, "success"))
}
