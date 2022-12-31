package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/RoseRocket/xerrs"
	"github.com/aaalik/api-iseng/model"
	"github.com/jmoiron/sqlx"
)

type SQLReaderRepository struct {
	dbr *sqlx.DB
}

func NewSQLReaderRepository(dbr *sqlx.DB) ISQLReaderRepository {
	return &SQLReaderRepository{dbr}
}

func (srr *SQLReaderRepository) DetailUser(ctx context.Context, id string) (*model.User, error) {
	user := model.User{}
	query := fmt.Sprintf(`
		SELECT
		id, name, dob, created_at,
		updated_at
		FROM %s
		WHERE id = :id
	`, SQLUserTable)

	params := map[string]interface{}{
		"id": id,
	}

	var args []interface{}
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return nil, err
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	err = srr.dbr.QueryRowxContext(ctx, query, args...).StructScan(&user)
	if err == sql.ErrNoRows {
		return nil, model.ErrorUserNotFound
	} else if err != nil {
		return nil, xerrs.Mask(err, model.ErrorSQLRead)
	}

	return &user, nil
}

func (srr *SQLReaderRepository) ListUser(ctx context.Context, request *model.RequestListUser) ([]*model.User, error) {
	users := []*model.User{}
	selectQuery := fmt.Sprintf(`
		SELECT
		id, name, dob, created_at,
		updated_at
		FROM %s
	`, SQLUserTable)

	whereQuery, args := srr.buildUserFilter(request)
	orderQuery := srr.buildUserSort(request.Queries)

	query := strings.Join([]string{selectQuery, whereQuery, orderQuery}, " ")
	query = srr.dbr.Rebind(query)

	rows, err := srr.dbr.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, xerrs.Mask(err, model.ErrorSQLRead)
	}
	defer rows.Close()
	for rows.Next() {
		user := model.User{}
		err = rows.StructScan(&user)
		if err != nil {
			return nil, xerrs.Mask(err, model.ErrorSQLRead)
		}
		users = append(users, &user)
	}

	return users, nil
}

func (srr *SQLReaderRepository) CountUsers(ctx context.Context, request *model.RequestListUser) (int32, error) {
	invoiceCount := int32(0)
	selectQuery := fmt.Sprintf(`
		SELECT
		COUNT(u.id)
		FROM %s AS u 
	`, SQLUserTable)

	whereQuery, args := srr.buildUserFilter(request)

	query := strings.Join([]string{selectQuery, whereQuery}, " ")
	query = srr.dbr.Rebind(query)

	err := srr.dbr.QueryRowxContext(ctx, query, args...).Scan(&invoiceCount)
	if err != nil {
		return 0, xerrs.Mask(err, model.ErrorSQLRead)
	}

	return invoiceCount, nil
}

func (srr *SQLReaderRepository) buildUserFilter(req *model.RequestListUser) (string, []interface{}) {
	var whereQuery string
	var whereQueries []string
	var whereArgs []interface{}

	filter := req.Includes

	if ql, arguments, err := sqlx.In("id IN (?)", filter.IDs); err == nil {
		whereQueries = append(whereQueries, ql)
		whereArgs = append(whereArgs, arguments...)
	}

	if req.Queries.Keyword != "" {
		whereQueries = append(whereQueries, "name ~* ?")
		whereArgs = append(whereArgs, req.Queries.Keyword)
	}

	if ql, arguments, err := sqlx.In("dob IN (?)", filter.Dobs); err == nil {
		whereQueries = append(whereQueries, ql)
		whereArgs = append(whereArgs, arguments...)
	}

	if filter.CreatedAt.Min != 0 {
		whereQueries = append(whereQueries, "created_at >= ?")
		whereArgs = append(whereArgs, filter.CreatedAt.Min)
	}

	if filter.CreatedAt.Max != 0 {
		whereQueries = append(whereQueries, "created_at <= ?")
		whereArgs = append(whereArgs, filter.CreatedAt.Max)
	}

	if len(whereQueries) > 0 {
		whereQuery = "WHERE " + strings.Join(whereQueries, " AND ")
	}

	return whereQuery, whereArgs
}

func (srr *SQLReaderRepository) buildUserSort(queries *model.Queries) string {
	var query string
	var limit int32
	var offset int32

	orderBy := "ORDER BY created_at DESC"
	if queries.Sort != nil {
		if queries.Sort.Order == model.ORDER_ASC {
			orderBy = fmt.Sprintf("ORDER BY %s ASC", queries.Sort.Field)
		}
		if queries.Sort.Order == model.ORDER_DESC {
			orderBy = fmt.Sprintf("ORDER BY %s DESC", queries.Sort.Field)
		}
	}

	query += orderBy

	limit = queries.Rows
	if limit == 0 {
		limit = 10
	}

	offset = 0
	if queries.Page > 0 {
		offset = (queries.Page - 1) * limit
	}

	if queries.Page != 0 && queries.Rows != 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	return query
}
