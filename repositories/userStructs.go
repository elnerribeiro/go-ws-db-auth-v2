package repositories

import (
	"github.com/elnerribeiro/go-ws-db-auth-v2/db"
	u "github.com/elnerribeiro/go-ws-db-auth-v2/utils"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

// Token JWT Token
type Token struct {
	UserID int
	Role   string
	jwt.RegisteredClaims
}

// ContextKey Key to use on a context
type ContextKey string

// UserRepository Repository for table usuario
type UserRepository interface {
	ListUsers(logContext *u.LoggerContext) ([]User, error)
	Upsert(logContext *u.LoggerContext, txContext *db.DatabaseContext, tx *pgx.Tx) (*User, error)
	Delete(logContext *u.LoggerContext, txContext *db.DatabaseContext, tx *pgx.Tx)
	UserToData(data *db.SqlData) *db.SqlData
	GetUserByID(logContext *u.LoggerContext) (*User, error)
	GetUserByEmail(logContext *u.LoggerContext, password bool) (*User, error)
}

// User table usuario on database
type User struct {
	ID       int    `json:"id,omitempty" db:"id,omitempty"`
	Email    string `json:"email,omitempty" db:"email,omitempty"`
	Password string `json:"password,omitempty" db:"password,omitempty"`
	Token    string `json:"token,omitempty" db:"-"`
	Role     string `json:"role,omitempty" db:"role,omitempty"`
}
