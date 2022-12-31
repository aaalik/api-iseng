package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func StopGracefully(logr *logrus.Logger, dbr, dbw *sqlx.DB) {
	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		sig := <-notify
		logr.Info("Caught signal ", sig, " Stop Gracefully")
		if dbr != nil {
			dbr.Close()
		}
		if dbw != nil {
			dbw.Close()
		}
		os.Exit(0)
	}()
}
