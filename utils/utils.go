package utils

import (
	"net/http"
)

type IRandomUtils interface {
	GenerateUniqueID() string
}

type IHTTPResponse interface {
	Nay(w http.ResponseWriter, r *http.Request, status int, err error)
	Yay(w http.ResponseWriter, r *http.Request, status int, content interface{})
	HTMLYay(w http.ResponseWriter, r *http.Request, status int, content string)
	DataYay(w http.ResponseWriter, r *http.Request, filename string, content []byte)
}
