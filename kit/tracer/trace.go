package tracer

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func NewTraceID(r *http.Request) context.Context {
	var (
		ctx           = r.Context()
		now           = time.Now().UTC()
		traceID       = uuid.NewString()
		correlationID = r.Header.Get("correlationID")
	)

	if correlationID == "" {
		correlationID = uuid.NewString()
	}

	v := Values{
		Now:           now,
		TraceID:       traceID,
		CorrelationID: correlationID,
	}

	return context.WithValue(ctx, ctxKey, &v)
}
