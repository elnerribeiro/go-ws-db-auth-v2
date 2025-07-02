package controllers

import (
	"encoding/json"
	"net/http"

	repo "github.com/elnerribeiro/go-ws-db-auth-v2/repositories"
	"github.com/elnerribeiro/go-ws-db-auth-v2/services"
	u "github.com/elnerribeiro/go-ws-db-auth-v2/utils"
)

// Authenticate do user authentication
var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	logger, logContext := u.GetLoggerAndContext()
	account := &repo.User{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		logger.Error().Msgf("Invalid request: %s", err)
		u.Respond(logContext, w, u.Message(false, "Invalid request"))
		return
	}

	resp := services.Login(logContext, account, account.Password)
	u.Respond(logContext, w, resp)
}

// Validate do user validation - gets ID
var Validate = func(w http.ResponseWriter, r *http.Request) {
	_, logContext := u.GetLoggerAndContext()
	id := r.Context().Value(repo.ContextKey("user")).(int)
	role := r.Context().Value(repo.ContextKey("role")).(string)
	resp := u.Message(true, "success")
	resp["userId"] = id
	resp["role"] = role
	u.Respond(logContext, w, resp)
}
