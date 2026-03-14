package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Yeremi528/itudy-back/kit/web"
)

type UserCoursesInfo struct {
	Active string           `json:"active" bson:"active,omitempty"`
	List   []EnrolledCourse `json:"list" bson:"list,omitempty"`
}

type EnrolledCourse struct {
	ID          string  `json:"ID" bson:"_id"`
	Name        string  `json:"name,omitempty" bson:"name,omitempty"`
	Progress    float64 `json:"progress,omitempty" bson:"progress,omitempty"`
	IsCompleted bool    `json:"is_completed,omitempty" bson:"is_completed,omitempty"`
}

type Stats struct {
	TotalXP       int       `json:"total_xp" bson:"total_xp,omitempty"`
	StreakDays    int       `json:"streak_days" bson:"streak_days,omitempty"`
	CurrentLeague string    `json:"current_league" bson:"current_league,omitempty"`
	LastLoginAt   time.Time `json:"last_login_at,omitempty" bson:"last_login_at,omitempty"`
}

type Achievement struct {
	ID       string    `json:"id" bson:"_id"`
	Title    string    `json:"title" bson:"title"`
	ExamName string    `json:"exam_name" bson:"exam_name"`
	EarnedAt time.Time `json:"earned_at" bson:"earned_at"`
}

type User struct {
	ID             string          `json:"id" bson:"_id"`
	Email          string          `json:"email" bson:"email,omitempty"`
	Name           string          `json:"name" bson:"name,omitempty"`
	Phone          string          `json:"phone,omitempty" bson:"phone,omitempty"`
	Country        string          `json:"country" bson:"country,omitempty"`
	NativeLanguage string          `json:"native_language" bson:"native_language,omitempty"`
	ImageURL       string          `json:"image_url,omitempty" bson:"image_url,omitempty"`
	CreatedAt      time.Time       `json:"created_at" bson:"created_at,omitempty"`
	CoursesInfo    UserCoursesInfo `json:"courses_info,omitempty" bson:"courses_info,omitempty"`
	Stats          Stats           `json:"stats,omitempty" bson:"stats,omitempty"`
	Achievements   []Achievement   `json:"achievements,omitempty" bson:"achievements,omitempty"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}

func errBadRequest(msg string) error {
	return web.NewRequestError(fmt.Errorf("%s", msg), http.StatusBadRequest)
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

	errorResponse := ErrorResponse{
		ErrorMessage: err.Error(),
	}
	errorResponseJSON, err := json.Marshal(errorResponse)
	if err != nil {
		return []byte(err.Error()), http.StatusInternalServerError
	}

	return errorResponseJSON, status
}

func validateRUT(rut string) bool {
	// Limpiar el RUT: quitar puntos, espacios y guiones
	rut = strings.ReplaceAll(rut, ".", "")
	rut = strings.ReplaceAll(rut, " ", "")
	rut = strings.ToUpper(strings.TrimSpace(rut))

	// Verificar formato básico
	if len(rut) < 3 || len(rut) > 10 || !strings.Contains(rut, "-") {
		return false
	}

	// Separar número y dígito verificador
	parts := strings.Split(rut, "-")
	if len(parts) != 2 {
		return false
	}

	numberStr, dv := parts[0], parts[1]

	// Validar que el número sea numérico y el DV sea válido (0-9 o K)
	if _, err := strconv.Atoi(numberStr); err != nil {
		return false
	}
	if len(dv) != 1 || !isValidDV(dv) {
		return false
	}

	// Calcular dígito verificador
	calculatedDV := calculateDV(numberStr)
	return calculatedDV == dv
}

// isValidDV verifica si el dígito verificador es válido (0-9 o K)
func isValidDV(dv string) bool {
	return (dv >= "0" && dv <= "9") || dv == "K"
}

// calculateDV calcula el dígito verificador para un número de RUT
func calculateDV(numberStr string) string {
	sum := 0
	multiplier := 2

	// Iterar desde el último dígito hacia el primero
	for i := len(numberStr) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(numberStr[i]))
		sum += digit * multiplier
		multiplier++
		if multiplier > 7 {
			multiplier = 2
		}
	}

	remainder := sum % 11
	dv := 11 - remainder

	switch dv {
	case 11:
		return "0"
	case 10:
		return "K"
	default:
		return strconv.Itoa(dv)
	}
}
