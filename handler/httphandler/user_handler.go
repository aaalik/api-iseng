package httphandler

import (
	"net/http"
	"strconv"

	"github.com/aaalik/api-iseng/model"
	"github.com/aaalik/api-iseng/usecase"
	"github.com/aaalik/api-iseng/utils"
	"github.com/aaalik/api-iseng/utils/validator"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type UserHandler struct {
	uu usecase.IUserUsecase
	hr utils.IHTTPResponse
}

func NewUserHandler(
	uu usecase.IUserUsecase,
	hr utils.IHTTPResponse,
) IUserHandler {
	return &UserHandler{
		uu,
		hr,
	}
}

func (uh *UserHandler) DetailUser(w http.ResponseWriter, r *http.Request) {
	var (
		resp   *model.ResponseUser
		userID = chi.URLParam(r, "userID")
		ctx    = r.Context()
	)

	user, err := uh.uu.DetailUser(ctx, userID)
	if err != nil {
		uh.hr.Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	resp = &model.ResponseUser{User: user}
	uh.hr.Yay(w, r, http.StatusOK, resp)
}

func (uh *UserHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	var (
		req   *model.RequestListUser
		resp  *model.ResponseListUser
		ctx   = r.Context()
		query = r.URL.Query()
	)

	page, _ := strconv.Atoi(query.Get("page"))
	rows, _ := strconv.Atoi(query.Get("rows"))

	createdAtMin, _ := strconv.Atoi(query.Get("created_at[min]"))
	createdAtMax, _ := strconv.Atoi(query.Get("created_at[max]"))

	sortOrder, _ := strconv.Atoi(query.Get("sort[order]"))

	req = &model.RequestListUser{
		Includes: &model.RequestFilterUser{
			IDs:  query["ids[]"],
			Dobs: query["dobs[]"],
			CreatedAt: &model.Range{
				Min: int64(createdAtMin),
				Max: int64(createdAtMax),
			},
		},
		Queries: &model.Queries{
			Page: int32(page),
			Rows: int32(rows),
			Sort: &model.SortBy{
				Field: query.Get("sort[field]"),
				Order: model.Order(sortOrder),
			},
			Keyword: query.Get("keyword"),
		},
	}

	users, count, err := uh.uu.ListUser(ctx, req)
	if err != nil {
		uh.hr.Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	resp = &model.ResponseListUser{Users: users, Stats: &model.ResponseStats{Total: count}}
	uh.hr.Yay(w, r, http.StatusOK, resp)
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var (
		req  *model.RequestCreateUser
		resp *model.ResponseUser
		ctx  = r.Context()
	)

	if err := render.Decode(r, &req); err != nil {
		uh.hr.Nay(w, r, http.StatusBadRequest, err)
		return
	}

	if err := validator.ValidateByName(req, "create_user"); err != nil {
		uh.hr.Nay(w, r, http.StatusBadRequest, err)
		return
	}

	user, err := uh.uu.CreateUser(ctx, req)
	if err != nil {
		uh.hr.Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	resp = &model.ResponseUser{User: user}
	uh.hr.Yay(w, r, http.StatusCreated, resp)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		req  *model.RequestUpdateUser
		resp *model.ResponseUser
		ctx  = r.Context()
	)

	if err := render.Decode(r, &req); err != nil {
		uh.hr.Nay(w, r, http.StatusBadRequest, err)
		return
	}

	req.ID = chi.URLParam(r, "userID")

	if err := validator.ValidateByName(req, "update_user"); err != nil {
		uh.hr.Nay(w, r, http.StatusBadRequest, err)
		return
	}

	user, err := uh.uu.UpdateUser(ctx, req)
	if err != nil {
		uh.hr.Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	resp = &model.ResponseUser{User: user}
	uh.hr.Yay(w, r, http.StatusCreated, resp)
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	err := uh.uu.DeleteUser(ctx, chi.URLParam(r, "userID"))
	if err != nil {
		uh.hr.Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	uh.hr.Yay(w, r, http.StatusNoContent, nil)
}
