package repositories

import (
	"database/sql"
	"github.com/elnerribeiro/go-ws-db-auth-v2/db"
	u "github.com/elnerribeiro/go-ws-db-auth-v2/utils"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

// ListUsers Lists all users
func (user *User) ListUsers(logContext *u.LoggerContext) ([]User, error) {
	logger := zerolog.Ctx(*logContext)
	var data db.SqlData
	val, err := db.SelectAll[User](logContext, nil, nil, "select id, email, password as password, role from user_db", data)
	if err != nil {
		logger.Error().Err(err).Msgf("[ListUsers] Error listing users: %s", err)
		return nil, err
	}
	var account []User
	for _, v := range val {
		account = append(account, v)
	}
	if account != nil {
		for i, user := range account {
			user.Password = ""
			account[i] = user
		}
		return account, nil
	}
	logger.Error().Msgf("[ListUsers] No user found.")
	return account, nil
}

// Upsert Inserts or updates an user
func (user *User) Upsert(logContext *u.LoggerContext, txContext *db.DatabaseContext, tx *pgx.Tx) (*User, error) {
	logger := zerolog.Ctx(*logContext)

	var data *db.SqlData
	data = user.UserToData()
	if user.ID == 0 {
		res, err := db.InsertReturningPostgres[User](logContext, txContext, tx, "user_db", *data, "id")
		if err != nil {
			logger.Error().Err(err).Msgf("[Upsert] Error inserting user: %s", err)
			return nil, err
		}
		user.ID = res.ID
		user.Password = ""
		return user, nil
	}

	var filter db.SqlData
	paramsValueFilter := db.SqlValue{Name: "id", Value: user.ID}
	filter = append(filter, paramsValueFilter)
	err2 := db.Update(logContext, txContext, tx, "user_db", *data, filter)
	if err2 != nil {
		logger.Error().Err(err2).Msgf("[Upsert] Error updating user: %s", err2)
		return nil, err2
	}
	user.Password = ""
	return user, nil
}

// Delete Deletes an user
func (user *User) Delete(logContext *u.LoggerContext, txContext *db.DatabaseContext, tx *pgx.Tx) error {
	var filter db.SqlData
	paramsValue := db.SqlValue{Name: "id", Value: user.ID}
	filter = append(filter, paramsValue)
	err := db.Delete(logContext, txContext, tx, "user_db", filter)
	return err
}

// UserToDados Fills a map containing the struct User
func (user *User) UserToData() *db.SqlData {
	var paramsQuery db.SqlData

	if user.ID != 0 {
		paramsValue := db.SqlValue{Name: "id", Value: user.ID}
		paramsQuery = append(paramsQuery, paramsValue)
	}
	if user.Email != "" {
		paramsValue := db.SqlValue{Name: "email", Value: user.Email}
		paramsQuery = append(paramsQuery, paramsValue)
	}
	if user.Role != "" {
		paramsValue := db.SqlValue{Name: "role", Value: user.Role}
		paramsQuery = append(paramsQuery, paramsValue)
	}
	if user.Password != "" {
		paramsValue := db.SqlValue{Name: "password", Value: user.Password}
		paramsQuery = append(paramsQuery, paramsValue)
	}
	return &paramsQuery
}

// GetUserByID Get an user by ID
func (user *User) GetUserByID(logContext *u.LoggerContext) (*User, error) {
	logger := zerolog.Ctx(*logContext)
	var data db.SqlData
	paramsValue := db.SqlValue{Name: "id", Value: user.ID}
	data = append(data, paramsValue)
	query := `
		select id, email, password as password, role from user_db
		where id = $1
	`
	val, err := getUser(logContext, data, query)
	if err != nil {
		logger.Error().Err(err).Msgf("[GetUserByID] Error retrieving user: %s", err)
		return nil, err
	}
	account := val.(*User)
	if account.Email == "" { //User not found!
		return nil, sql.ErrNoRows
	}
	account.Password = ""
	return account, nil
}

// GetUserByEmail Get an user by email
func (user *User) GetUserByEmail(logContext *u.LoggerContext, password bool) (*User, error) {
	logger := zerolog.Ctx(*logContext)
	var data db.SqlData
	paramsValue := db.SqlValue{Name: "email", Value: user.Email}
	data = append(data, paramsValue)
	query := `
		select id, email, password as password, role from user_db
		where email = lower($1)
	`
	val, err := getUser(logContext, data, query)
	if err != nil {
		logger.Error().Err(err).Msgf("[GetUserByEmail] Error retrieving user: %s", err)
		return nil, err
	}
	account := val.(*User)
	if account.Email == "" { //User not found!
		logger.Error().Msgf("[GetUserByEmail] Error retrieving user: %s", user.Email)
		return nil, sql.ErrNoRows
	}

	if !password {
		account.Password = ""
	}
	return account, nil
}

func getUser(logContext *u.LoggerContext, data db.SqlData, query string) (interface{}, error) {
	logger := zerolog.Ctx(*logContext)
	val, err := db.SelectOne[User](logContext, nil, nil, query, data)
	if err != nil {
		logger.Error().Err(err).Msgf("[getUser] Error selecting user: %s", err)
		return nil, err
	}
	return val, nil
}
