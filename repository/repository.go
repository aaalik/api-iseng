package repository

import (
	"context"

	"github.com/aaalik/api-iseng/model"
	"github.com/jmoiron/sqlx"
)

const (
	SQLUserTable = "users"
)

type (
	ErrorResponse struct {
		Code    int32             `json:"code"`
		Message string            `json:"message"`
		Reasons map[string]string `json:"reasons"`
	}
	Response struct {
		RequestID string                 `json:"request_id"`
		Status    int                    `json:"status"`
		Content   map[string]interface{} `json:"content,omitempty"`
		Error     *ErrorResponse         `json:"error,omitempty"`
	}
)

type ISQLWriterRepository interface {
	CreateTransaction(ctx context.Context) (*sqlx.Tx, error)
	CommitTransaction(ctx context.Context, tx *sqlx.Tx) error
	RollbackTransaction(ctx context.Context, tx *sqlx.Tx) error

	CreateUser(ctx context.Context, tx *sqlx.Tx, user *model.User) error
	UpdateUser(ctx context.Context, tx *sqlx.Tx, user *model.User) error
	DeleteUser(ctx context.Context, tx *sqlx.Tx, userID string) error
}

type ISQLReaderRepository interface {
	DetailUser(ctx context.Context, id string) (*model.User, error)
	ListUser(ctx context.Context, request *model.RequestListUser) ([]*model.User, error)
	CountUsers(ctx context.Context, request *model.RequestListUser) (int32, error)
}
