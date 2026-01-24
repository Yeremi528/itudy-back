package mids

import (
	"bytes"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// Middleware para capturar errores y panics
func ErrorRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var res bytes.Buffer

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		ww.Tee(&res)

		next.ServeHTTP(ww, r)
		switch ww.Status() {

		}
	})
}
