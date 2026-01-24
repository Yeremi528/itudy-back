package tracer

import (
	"context"
	"time"
)

// Values struct represents the state for each request.
type Values struct {
	TraceID       string
	CorrelationID string
	Rut           string
	Now           time.Time
}

type contextKey int

const ctxKey contextKey = 1
const defaultTraceId = "00000000-0000-000000000000"

// GetValues returns the values from the context.
func GetValues(ctx context.Context) *Values {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return &Values{
			TraceID: defaultTraceId,
			Now:     time.Now().UTC(),
		}
	}

	return v
}
