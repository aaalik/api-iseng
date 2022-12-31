package httphandler

import "net/http"

type IUserHandler interface {
	DetailUser(w http.ResponseWriter, r *http.Request)
	ListUser(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}
