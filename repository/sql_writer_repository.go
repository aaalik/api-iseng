package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/RoseRocket/xerrs"
	"github.com/aaalik/api-iseng/model"
	"github.com/jmoiron/sqlx"
)

type SQLWriterRepository struct {
	dbw *sqlx.DB
}

func NewSQLWriterRepository(dbw *sqlx.DB) ISQLWriterRepository {
	return &SQLWriterRepository{dbw}
}

func (swr *SQLWriterRepository) CreateTransaction(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := swr.dbw.BeginTxx(ctx, nil)
	if err != nil {
		err = xerrs.Mask(err, model.ErrorSQLCreateTransaction)
		return nil, err
	}

	return tx, nil
}

func (swr *SQLWriterRepository) CommitTransaction(ctx context.Context, tx *sqlx.Tx) error {
	err := tx.Commit()
	if err != nil {
		err = xerrs.Mask(err, model.ErrorSQLCommitTransaction)
		return err
	}

	return nil
}

func (swr *SQLWriterRepository) RollbackTransaction(ctx context.Context, tx *sqlx.Tx) error {
	err := tx.Rollback()
	if err != sql.ErrTxDone {
		err = xerrs.Mask(err, model.ErrorSQLRollbackTransaction)
		return err
	}

	return nil
}

func (swr *SQLWriterRepository) CreateUser(ctx context.Context, tx *sqlx.Tx, user *model.User) error {
	query := fmt.Sprintf(`
		INSERT INTO
		%s
		(
			id, name, dob, created_at,
			updated_at
		)
		VALUES
		(
			:id, :name, :dob, :created_at,
			:updated_at
		)
	`, SQLUserTable)

	params := map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"dob":        user.Dob,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	var args []interface{}
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return err
	}
	query = swr.dbw.Rebind(query)

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		err = xerrs.Mask(err, model.ErrorSQLCreateUser)
		return err
	}

	return nil
}

func (swr *SQLWriterRepository) UpdateUser(ctx context.Context, tx *sqlx.Tx, user *model.User) error {
	query := fmt.Sprintf(`
		UPDATE
		%s
		SET
		name = :name,
		dob = :dob,
		updated_at = :updated_at
		WHERE id = :id
	`, SQLUserTable)

	params := map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"dob":        user.Dob,
		"updated_at": time.Now().Unix(),
	}

	var args []interface{}
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return err
	}
	query = swr.dbw.Rebind(query)

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		err = xerrs.Mask(err, model.ErrorSQLUpdateUser)
		return err
	}

	return nil
}

func (swr *SQLWriterRepository) DeleteUser(ctx context.Context, tx *sqlx.Tx, userID string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = :id
	`, SQLUserTable)

	params := map[string]interface{}{
		"id": userID,
	}

	var args []interface{}
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return err
	}
	query = swr.dbw.Rebind(query)

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		err = xerrs.Mask(err, model.ErrorSQLUpdateUser)
		return err
	}

	return nil
}
