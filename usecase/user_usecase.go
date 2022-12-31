package usecase

import (
	"context"
	"time"

	"github.com/aaalik/api-iseng/model"
	"github.com/aaalik/api-iseng/repository"
	"github.com/aaalik/api-iseng/utils"
)

type UserUsecase struct {
	srr repository.ISQLReaderRepository
	swr repository.ISQLWriterRepository
	ru  utils.IRandomUtils
}

func NewUserUsecase(
	srr repository.ISQLReaderRepository,
	swr repository.ISQLWriterRepository,
	ru utils.IRandomUtils,
) IUserUsecase {
	return &UserUsecase{
		srr,
		swr,
		ru,
	}
}

func (uu *UserUsecase) DetailUser(ctx context.Context, userID string) (*model.User, error) {
	invoice, err := uu.srr.DetailUser(ctx, userID)
	return invoice, err
}

func (uu *UserUsecase) ListUser(ctx context.Context, request *model.RequestListUser) ([]*model.User, int32, error) {
	invoices, err := uu.srr.ListUser(ctx, request)
	if err != nil {
		return nil, 0, err
	}

	count, err := uu.srr.CountUsers(ctx, request)
	if err != nil {
		return nil, 0, err
	}

	return invoices, count, nil
}

func (uu *UserUsecase) CreateUser(ctx context.Context, request *model.RequestCreateUser) (*model.User, error) {
	user := model.User{
		ID:        uu.ru.GenerateUniqueID(),
		Name:      request.Name,
		Dob:       request.Dob,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	tx, err := uu.swr.CreateTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer uu.swr.RollbackTransaction(ctx, tx)

	err = uu.swr.CreateUser(ctx, tx, &user)
	if err != nil {
		return nil, err
	}

	err = uu.swr.CommitTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (uu *UserUsecase) UpdateUser(ctx context.Context, request *model.RequestUpdateUser) (*model.User, error) {
	user, err := uu.DetailUser(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	reqUser := model.User{
		ID:        user.ID,
		Name:      request.Name,
		Dob:       request.Dob,
		CreatedAt: user.CreatedAt,
		UpdatedAt: time.Now().Unix(),
	}

	tx, err := uu.swr.CreateTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer uu.swr.RollbackTransaction(ctx, tx)

	err = uu.swr.UpdateUser(ctx, tx, &reqUser)
	if err != nil {
		return nil, err
	}

	err = uu.swr.CommitTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	return &reqUser, nil
}

func (uu *UserUsecase) DeleteUser(ctx context.Context, userID string) error {
	tx, err := uu.swr.CreateTransaction(ctx)
	if err != nil {
		return err
	}
	defer uu.swr.RollbackTransaction(ctx, tx)

	err = uu.swr.DeleteUser(ctx, tx, userID)
	if err != nil {
		return err
	}

	err = uu.swr.CommitTransaction(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}
