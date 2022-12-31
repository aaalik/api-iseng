package main

import (
	"net/http"

	"github.com/aaalik/api-iseng/cmd"
	"github.com/aaalik/api-iseng/config"
	"github.com/aaalik/api-iseng/handler/httphandler"
	apimiddleware "github.com/aaalik/api-iseng/middleware"
	"github.com/aaalik/api-iseng/repository"
	"github.com/aaalik/api-iseng/usecase"
	"github.com/aaalik/api-iseng/utils"
	"github.com/go-chi/chi"
	chim "github.com/go-chi/chi/middleware"
)

const ()

func main() {
	cf := config.NewConfig(true)

	logr := cf.NewLogrus()
	sqlRead, sqlWrite := cf.NewSQL()

	httpResp := utils.NewHTTPResponse(logr)
	randomizer := utils.NewRandomUtils()

	sqlReaderRepository := repository.NewSQLReaderRepository(sqlRead)
	sqlWriterRepository := repository.NewSQLWriterRepository(sqlWrite)

	userUsecase := usecase.NewUserUsecase(
		sqlReaderRepository,
		sqlWriterRepository,
		randomizer,
	)

	userHandler := httphandler.NewUserHandler(userUsecase, httpResp)

	router := chi.NewRouter()
	router.Use(
		chim.NoCache,
		chim.RedirectSlashes,
		chim.Heartbeat("/ping"),
		chim.RequestID,
		chim.Recoverer,
		chim.RealIP,
		apimiddleware.RequestLogger(logr),
	)

	router.Route("/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/", userHandler.ListUser)
			r.Get("/{userID}", userHandler.DetailUser)
			r.Put("/{userID}", userHandler.UpdateUser)
			r.Delete("/{userID}", userHandler.DeleteUser)
		})
	})

	if err := chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logr.Infof("%s %s", method, route)
		return nil
	}); err != nil {
		logr.Panicln(err)
	}

	server := http.Server{
		Addr:    cf.Host.Address,
		Handler: router,
	}

	// Graceful Stop handle
	cmd.StopGracefully(logr, sqlRead, sqlWrite)

	logr.Infof("API serving at %s", server.Addr)
	logr.Fatal(server.ListenAndServe())
}
