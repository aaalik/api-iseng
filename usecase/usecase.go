package usecase

import (
	"context"

	"github.com/aaalik/api-iseng/model"
)

type IUserUsecase interface {
	DetailUser(ctx context.Context, userID string) (*model.User, error)
	ListUser(ctx context.Context, request *model.RequestListUser) ([]*model.User, int32, error)
	CreateUser(ctx context.Context, request *model.RequestCreateUser) (*model.User, error)
	UpdateUser(ctx context.Context, request *model.RequestUpdateUser) (*model.User, error)
	DeleteUser(ctx context.Context, userID string) error
}
