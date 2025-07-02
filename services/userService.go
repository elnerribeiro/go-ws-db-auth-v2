package services

import (
	"database/sql"
	"errors"
	"github.com/elnerribeiro/go-ws-db-auth-v2/db"
	repo "github.com/elnerribeiro/go-ws-db-auth-v2/repositories"
	u "github.com/elnerribeiro/go-ws-db-auth-v2/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"time"
)

// ListUsers Lists all users
func ListUsers(logContext *u.LoggerContext, user *repo.User) ([]repo.User, error) {
	return user.ListUsers(logContext)
}

// Login Authenticates an user
func Login(logContext *u.LoggerContext, user *repo.User, password string) map[string]interface{} {
	logger := zerolog.Ctx(*logContext)
	account, err := user.GetUserByEmail(logContext, true)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error().Err(err).Msgf("[Login] Email not found.")
			return u.Message(false, "Email address not found")
		} else {
			logger.Error().Err(err).Msgf("[Login] Connection error. Please retry: %s", err)
			return u.Message(false, "Connection error. Please retry")
		}
	}
	if account.Password != password { //Password does not match!
		logger.Error().Msgf("[Login] Invalid login credentials. Please try again: %s != %s", account.Password, password)
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	expirationTime := time.Now().Add(12 * time.Hour)
	numericDate := jwt.NumericDate{expirationTime}
	// Create the JWT claims, which includes the username and expiry time
	claims := &repo.Token{
		UserID: account.ID,
		Role:   account.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: &numericDate,
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("JWTpassword123@"))
	account.Token = tokenString //Store the token in the response

	logger.Info().Msgf("[Login] User Logged In: %d", account.ID)

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

// GetUserByID Gets an user by ID
func GetUserByID(logContext *u.LoggerContext, user *repo.User) (*repo.User, error) {
	logger := zerolog.Ctx(*logContext)
	val, err := user.GetUserByID(logContext)
	if err != nil {
		logger.Error().Err(err).Msgf("[GetUserByID] Error : %s", err)
		return nil, err
	}
	return val, nil
}

// Upsert Inserts or updates an user
func Upsert(logContext *u.LoggerContext, user *repo.User) (*repo.User, error) {
	logger := zerolog.Ctx(*logContext)
	tx, txContext, err := db.GetTransaction()
	if err != nil {
		logger.Error().Err(err).Msgf("[Upsert] Error starting transaction: %s", err)
		return nil, err
	}
	defer db.Rollback(logContext, txContext, tx)
	val, err := user.Upsert(logContext, txContext, tx)
	if err != nil {
		logger.Error().Err(err).Msgf("[Upsert] Error executing upsert: %s", err)
		return nil, err
	}
	db.Commit(logContext, txContext, tx)
	return val, nil
}

// Delete Deletes an user
func Delete(logContext *u.LoggerContext, user *repo.User) error {
	logger := zerolog.Ctx(*logContext)
	tx, txContext, err := db.GetTransaction()
	if err != nil {
		logger.Error().Err(err).Msgf("[Delete] Error starting transaction: %s", err)
		return err
	}
	defer db.Rollback(logContext, txContext, tx)
	if err2 := user.Delete(logContext, txContext, tx); err2 != nil {
		logger.Error().Err(err2).Msgf("[Delete] Error executing delete: %s", err)
		return err2
	}
	db.Commit(logContext, txContext, tx)
	return nil
}
