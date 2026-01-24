package mids

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Yeremi528/itudy-back/kit/logger"
	"github.com/go-chi/cors"
)

type corsLogger struct {
	log *logger.Logger
}

func (cl *corsLogger) Printf(format string, params ...interface{}) {
	cl.log.Info(context.Background(), fmt.Sprintf(format, params...))
}

var opts = cors.Options{
	AllowedOrigins:   []string{"https://*", "http://*"},
	AllowedMethods:   []string{"GET", "POST"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: false,

	AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
}

func Cors(log *logger.Logger) func(http.Handler) http.Handler {
	corsLogger := &corsLogger{log: log}

	c := cors.New(opts)
	c.Log = corsLogger

	return c.Handler
}
