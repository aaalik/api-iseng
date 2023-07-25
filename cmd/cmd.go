package cmd

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func StartServer(wg *sync.WaitGroup, logr *logrus.Logger, srv *http.Server) {
	defer wg.Done()

	logr.Infof("API serving at %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logr.Fatalf("Server error: %v", err)
	}
}

func StopGracefully(wg *sync.WaitGroup, logr *logrus.Logger, srv *http.Server, dbr, dbw *sqlx.DB) {
	defer wg.Done()

	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logr.Fatal(err)
		}
	}
	if dbr != nil {
		if err := dbr.Close(); err != nil {
			logr.Fatal(err)
		}
	}
	if dbw != nil {
		if err := dbw.Close(); err != nil {
			logr.Fatal(err)
		}
	}
	os.Exit(0)
}
