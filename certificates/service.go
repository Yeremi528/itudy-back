package certificates

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/google/uuid"
)

type service struct {
	repo     Repository
	storage  Storage
	examSvc  examService
	userRepo userRepository
}

func NewService(repo Repository, storage Storage, examSvc examService, userRepo userRepository) Service {
	return &service{
		repo:    repo,
		storage: storage,
		examSvc: examSvc,
		userRepo: userRepo,
	}
}

func (s *service) IssueCertificate(ctx context.Context, req IssueCertificateRequest) (Certificate, error) {
	// Obtener datos del usuario
	u, err := s.userRepo.GetUser(ctx, req.UserID)
	if err != nil {
		return Certificate{}, fmt.Errorf("certificates.IssueCertificate: getUser: %w", err)
	}
	if u.ID == "" {
		return Certificate{}, fmt.Errorf("certificates.IssueCertificate: usuario no encontrado")
	}

	// Obtener datos del examen
	exam := s.examSvc.ExambyID(ctx, req.ExamID)
	if exam.ID == "" {
		return Certificate{}, fmt.Errorf("certificates.IssueCertificate: examen no encontrado")
	}

	approvedAt := req.ApprovedAt
	if approvedAt.IsZero() {
		approvedAt = time.Now().UTC()
	}

	certID := uuid.New().String()

	// Generar PDF
	pdfBytes, err := generatePDF(certID, u.Name, exam.Title, exam.Language, exam.DifficultyLevel, req.WorkerName, req.ExamStartedAt, approvedAt)
	if err != nil {
		return Certificate{}, fmt.Errorf("certificates.IssueCertificate: generatePDF: %w", err)
	}

	// Subir al bucket
	objectName := fmt.Sprintf("certificates/%s/%s.pdf", req.UserID, certID)
	url, err := s.storage.UploadPublic(ctx, objectName, bytes.NewReader(pdfBytes), "application/pdf")
	if err != nil {
		return Certificate{}, fmt.Errorf("certificates.IssueCertificate: uploadPDF: %w", err)
	}

	cert := Certificate{
		ID:             certID,
		UserID:         req.UserID,
		UserName:       u.Name,
		ExamID:         req.ExamID,
		ExamName:       exam.Title,
		Language:       exam.Language,
		Level:          exam.DifficultyLevel,
		ValidatedBy:    req.WorkerName,
		ExamStartedAt:  req.ExamStartedAt,
		ApprovedAt:     approvedAt,
		CertificateURL: url,
		CreatedAt:      time.Now().UTC(),
	}

	if err := s.repo.Save(ctx, cert); err != nil {
		return Certificate{}, fmt.Errorf("certificates.IssueCertificate: save: %w", err)
	}

	return cert, nil
}

func (s *service) GetByUserID(ctx context.Context, userID string) ([]Certificate, error) {
	certs, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("certificates.GetByUserID: %w", err)
	}
	return certs, nil
}

// generatePDF construye el certificado en formato A4 apaisado.
func generatePDF(certID, userName, examName, language, level, validatedBy string, startedAt, approvedAt time.Time) ([]byte, error) {
	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(0, 0, 0)
	pdf.AddPage()

	pageW := 297.0
	pageH := 210.0

	// ── Fondo oscuro (bordes / marco) ──────────────────────────────────────
	pdf.SetFillColor(15, 23, 42) // slate-900
	pdf.Rect(0, 0, pageW, pageH, "F")

	// Franja dorada superior
	pdf.SetFillColor(203, 161, 53) // dorado
	pdf.Rect(0, 0, pageW, 4, "F")
	// Franja dorada inferior
	pdf.Rect(0, pageH-4, pageW, 4, "F")
	// Franja izquierda y derecha
	pdf.Rect(0, 0, 4, pageH, "F")
	pdf.Rect(pageW-4, 0, 4, pageH, "F")

	// ── Área blanca interior ────────────────────────────────────────────────
	pdf.SetFillColor(255, 255, 255)
	pdf.RoundedRect(14, 14, pageW-28, pageH-28, 4, "1234", "F")

	// ── Banda de color superior interna ────────────────────────────────────
	pdf.SetFillColor(15, 23, 42)
	pdf.RoundedRect(14, 14, pageW-28, 32, 4, "1200", "F")
	pdf.Rect(14, 30, pageW-28, 16, "F") // hace cuadrado el borde inferior de la banda

	// ── Logo / Título de la plataforma ─────────────────────────────────────
	pdf.SetTextColor(203, 161, 53)
	pdf.SetFont("Helvetica", "B", 22)
	pdf.SetXY(14, 18)
	pdf.CellFormat(pageW-28, 10, "ITUDY", "", 0, "C", false, 0, "")

	pdf.SetTextColor(200, 200, 200)
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetXY(14, 28)
	pdf.CellFormat(pageW-28, 6, "Plataforma de aprendizaje tecnológico", "", 0, "C", false, 0, "")

	// ── Título del certificado ─────────────────────────────────────────────
	pdf.SetTextColor(15, 23, 42)
	pdf.SetFont("Helvetica", "B", 20)
	pdf.SetXY(14, 52)
	pdf.CellFormat(pageW-28, 10, "CERTIFICADO DE APROBACIÓN", "", 0, "C", false, 0, "")

	// Línea decorativa bajo el título
	pdf.SetDrawColor(203, 161, 53)
	pdf.SetLineWidth(0.6)
	centerX := pageW / 2
	pdf.Line(centerX-50, 64, centerX+50, 64)

	// ── Cuerpo: texto principal ────────────────────────────────────────────
	pdf.SetTextColor(80, 80, 80)
	pdf.SetFont("Helvetica", "", 11)
	pdf.SetXY(14, 68)
	pdf.CellFormat(pageW-28, 8, "Se certifica que:", "", 0, "C", false, 0, "")

	// Nombre del usuario (grande y destacado)
	pdf.SetTextColor(15, 23, 42)
	pdf.SetFont("Helvetica", "B", 26)
	pdf.SetXY(14, 76)
	pdf.CellFormat(pageW-28, 14, userName, "", 0, "C", false, 0, "")

	pdf.SetTextColor(80, 80, 80)
	pdf.SetFont("Helvetica", "", 11)
	pdf.SetXY(14, 91)
	pdf.CellFormat(pageW-28, 8, "ha aprobado exitosamente el examen:", "", 0, "C", false, 0, "")

	// Nombre del examen
	pdf.SetTextColor(15, 23, 42)
	pdf.SetFont("Helvetica", "B", 16)
	pdf.SetXY(14, 99)
	pdf.CellFormat(pageW-28, 10, examName, "", 0, "C", false, 0, "")

	// ── Detalles en dos columnas ───────────────────────────────────────────
	pdf.SetLineWidth(0.3)
	pdf.SetDrawColor(220, 220, 220)
	pdf.Line(14, 116, pageW-14, 116)

	colW := (pageW - 28) / 4
	startY := 120.0
	lineH := 8.0

	details := []struct{ label, value string }{
		{"Lenguaje", strings.Title(strings.ToLower(language))},
		{"Nivel", strings.Title(strings.ToLower(level))},
		{"Inicio del examen", startedAt.In(time.FixedZone("CLT", -4*3600)).Format("02/01/2006 15:04")},
		{"Fecha de aprobación", approvedAt.In(time.FixedZone("CLT", -4*3600)).Format("02/01/2006 15:04")},
	}

	for i, d := range details {
		x := 14 + float64(i)*colW
		// Etiqueta
		pdf.SetTextColor(150, 150, 150)
		pdf.SetFont("Helvetica", "", 8)
		pdf.SetXY(x, startY)
		pdf.CellFormat(colW, lineH, d.label, "", 0, "C", false, 0, "")
		// Valor
		pdf.SetTextColor(15, 23, 42)
		pdf.SetFont("Helvetica", "B", 10)
		pdf.SetXY(x, startY+lineH)
		pdf.CellFormat(colW, lineH, d.value, "", 0, "C", false, 0, "")
	}

	// ── Validado por ───────────────────────────────────────────────────────
	pdf.SetLineWidth(0.3)
	pdf.SetDrawColor(220, 220, 220)
	pdf.Line(14, 142, pageW-14, 142)

	pdf.SetTextColor(150, 150, 150)
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetXY(14, 145)
	pdf.CellFormat(pageW-28, lineH, "Validado por", "", 0, "C", false, 0, "")

	pdf.SetTextColor(15, 23, 42)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetXY(14, 153)
	pdf.CellFormat(pageW-28, lineH, validatedBy, "", 0, "C", false, 0, "")

	// Línea de firma
	pdf.SetDrawColor(15, 23, 42)
	pdf.SetLineWidth(0.4)
	pdf.Line(centerX-40, 166, centerX+40, 166)
	pdf.SetTextColor(100, 100, 100)
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetXY(14, 167)
	pdf.CellFormat(pageW-28, 6, "Instructor Evaluador — Itudy", "", 0, "C", false, 0, "")

	// ── Footer: ID del certificado ─────────────────────────────────────────
	pdf.SetTextColor(160, 160, 160)
	pdf.SetFont("Helvetica", "", 7)
	pdf.SetXY(14, pageH-20)
	pdf.CellFormat(pageW-28, 6, fmt.Sprintf("ID de Certificado: %s", certID), "", 0, "C", false, 0, "")
	pdf.SetXY(14, pageH-15)
	pdf.CellFormat(pageW-28, 6, "Certificado verificable en itudy.app", "", 0, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
