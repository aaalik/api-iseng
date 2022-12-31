package model

import (
	"github.com/RoseRocket/xerrs"
)

var (
	// 404 Errors
	ErrorUserNotFound = xerrs.New("User tidak ditemukan")

	// 422 Errors
	ErrorInvalidUser = xerrs.New("User invalid")

	// 500 Errors
	ErrorSQLCreateTransaction   = xerrs.New("Failed to Write Data")
	ErrorSQLCommitTransaction   = xerrs.New("Failed to Commit Data")
	ErrorSQLRollbackTransaction = xerrs.New("Failed to Rollback Data")
	ErrorSQLRead                = xerrs.New("Failed to Retrieve Data")
	ErrFileProcessing           = xerrs.New("Failed to Process File")
	ErrInvalidContentType       = xerrs.New("Invalid Content Type")

	ErrorSQLCreateUser = xerrs.New("Gagal membuat User")
	ErrorSQLUpdateUser = xerrs.New("Gagal update User")
)

func init() {
	xerrs.SetData(ErrorUserNotFound, "status_code", 404)

	xerrs.SetData(ErrorInvalidUser, "status_code", 422)
}
