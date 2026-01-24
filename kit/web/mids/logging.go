package mids

import (
	"bytes"
	"net/http"
	"time"

	"github.com/Yeremi528/itudy-back/kit/logger"
	"github.com/Yeremi528/itudy-back/kit/tracer"
	"github.com/go-chi/chi/v5/middleware"
)

func Logging(log *logger.Logger, skipPahts ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			fields := []any{
				"protocol", r.Proto,
				"method", r.Method,
				"URL", r.RequestURI,
				"userAgent", r.UserAgent(),
				"remoteAddr", r.RemoteAddr,
			}

			log.Info(ctx, "request started "+r.RequestURI, fields...)

			var res bytes.Buffer

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			ww.Tee(&res)

			next.ServeHTTP(ww, r.WithContext(r.Context()))

			fields = append(
				fields,
				"statusCode", ww.Status(),
				"response", res.String(),
				"since", time.Since(tracer.GetValues(ctx).Now).String(),
			)

			log.Info(ctx, "request completed "+r.RequestURI, fields...)

		}

		return http.HandlerFunc(h)
	}
}
