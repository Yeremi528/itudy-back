package certificates

import (
	"context"
	"io"

	"github.com/Yeremi528/itudy-back/exam"
	"github.com/Yeremi528/itudy-back/user"
)

// Service define el contrato de negocio para certificados.
type Service interface {
	IssueCertificate(ctx context.Context, req IssueCertificateRequest) (Certificate, error)
	GetByUserID(ctx context.Context, userID string) ([]Certificate, error)
}

// Repository define el contrato de persistencia.
type Repository interface {
	Save(ctx context.Context, cert Certificate) error
	GetByUserID(ctx context.Context, userID string) ([]Certificate, error)
}

// Storage define el contrato para subir archivos al bucket.
type Storage interface {
	UploadPublic(ctx context.Context, objectName string, data io.Reader, contentType string) (string, error)
}

// examService es el subconjunto de exam.Service que necesitamos.
type examService interface {
	ExambyID(ctx context.Context, ID string) exam.Exam
}

// userRepository es el subconjunto de user.Repository que necesitamos.
type userRepository interface {
	GetUser(ctx context.Context, idOrEmail string) (user.User, error)
}
