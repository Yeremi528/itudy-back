package certificates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Yeremi528/itudy-back/kit/web"
)

type Certificate struct {
	ID             string    `json:"id" bson:"_id"`
	UserID         string    `json:"user_id" bson:"user_id"`
	UserName       string    `json:"user_name" bson:"user_name"`
	ExamID         string    `json:"exam_id" bson:"exam_id"`
	ExamName       string    `json:"exam_name" bson:"exam_name"`
	Language       string    `json:"language" bson:"language"`
	Level          string    `json:"level" bson:"level"`
	ValidatedBy    string    `json:"validated_by" bson:"validated_by"`
	ExamStartedAt  time.Time `json:"exam_started_at" bson:"exam_started_at"`
	ApprovedAt     time.Time `json:"approved_at" bson:"approved_at"`
	CertificateURL string    `json:"certificate_url" bson:"certificate_url"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
}

// IssueCertificateRequest es el cuerpo del POST /certificates
type IssueCertificateRequest struct {
	UserID        string    `json:"user_id"`
	ExamID        string    `json:"exam_id"`
	WorkerName    string    `json:"worker_name"`
	ExamStartedAt time.Time `json:"exam_started_at"`
	ApprovedAt    time.Time `json:"approved_at"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}

func newError(err error) ([]byte, int) {
	var status int
	switch {
	case web.IsRequestError(err):
		errReq := web.GetRequestError(err)
		status = errReq.Status
	default:
		status = http.StatusInternalServerError
	}

	errorResponse := ErrorResponse{ErrorMessage: err.Error()}
	errorResponseJSON, err := json.Marshal(errorResponse)
	if err != nil {
		return []byte(err.Error()), http.StatusInternalServerError
	}
	return errorResponseJSON, status
}

func errBadRequest(msg string) error {
	return web.NewRequestError(fmt.Errorf("%s", msg), http.StatusBadRequest)
}
