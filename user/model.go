package user

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Yeremi528/itudy-back/kit/web"
)

type Stats struct {
	TotalXP       int    `json:"total_xp" bson:"total_xp"`
	StreakDays    int    `json:"streak_days" bson:"streak_days"`
	CurrentLeague string `json:"current_league" bson:"current_league"`
}

type User struct {
	Email          string    `json:"email" bson:"email"`
	Name           string    `json:"name" bson:"name"`
	Phone          string    `json:"phone,omitempty" bson:"phone,omitempty"`
	Country        string    `json:"country" bson:"country"`
	NativeLanguage string    `json:"native_language" bson:"native_language"`
	ImageURL       string    `json:"image_url,omitempty" bson:"image_url,omitempty"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	Stats          Stats     `json:"stats,omitempty" bson:"stats,omitempty"`
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
