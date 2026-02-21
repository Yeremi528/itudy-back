package movements

import "time"

type Movement struct {
	ID            string
	Email         string
	Amount        float64
	TestID        string
	State         int
	TransactionID string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type response struct {
	Data any `json:"data"`
}

func responseMobile(data any) response {
	return response{
		Data: data,
	}
}
