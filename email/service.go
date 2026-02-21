package email

import (
	"context"
	"fmt"
	"strings"

	"github.com/resend/resend-go/v2"
)

type service struct {
	client *resend.Client
}

// Pasamos el cliente de Resend al constructor
func NewService(client *resend.Client) Service {
	return &service{
		client: client,
	}
}

const certEmailHTML = `
<!DOCTYPE html>
<html lang="es" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="color-scheme" content="light dark">
    <meta name="supported-color-schemes" content="light dark">
    <title>Certificación Itudy</title>
    
    <style>
        /* Reset básico para clientes de correo */
        body, table, td, p, h1, h2, h3 {
            margin: 0;
            padding: 0;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        }
        
        table {
            border-spacing: 0;
            border-collapse: collapse;
        }

        /* Reglas estrictas para cuando el usuario tenga activado el Modo Oscuro */
        @media (prefers-color-scheme: dark) {
            .email-bg {
                background-color: #0B0F19 !important;
            }
            .text-main {
                color: #E2E8F0 !important;
            }
            .info-box {
                background-color: #1E293B !important;
                border: 1px solid #334155 !important;
            }
            /* Forzamos el texto blanco dentro de la caja para evitar el gris oscuro de tu captura */
            .info-box td {
                color: #FFFFFF !important;
            }
        }
    </style>
</head>
<body style="margin: 0; padding: 0; background-color: #0B0F19; -webkit-text-size-adjust: 100%; text-size-adjust: 100%;">

    <table class="email-bg" width="100%" border="0" cellpadding="0" cellspacing="0" style="background-color: #0B0F19; padding: 40px 20px;">
        <tr>
            <td align="center">
                <table width="100%" max-width="600" border="0" cellpadding="0" cellspacing="0" style="max-width: 600px; width: 100%; background-color: transparent;">
                    
                    <tr>
                        <td align="center" style="padding-bottom: 20px;">
                            <h1 class="text-main" style="color: #FFFFFF; font-size: 28px; font-weight: bold; letter-spacing: 2px; margin-bottom: 5px;">ITUDY</h1>
                            <p style="color: #3B82F6; font-size: 14px; font-weight: bold; letter-spacing: 1px; text-transform: uppercase;">Ruta Backend</p>
                        </td>
                    </tr>

                    <tr>
                        <td align="left" style="padding-bottom: 20px;">
                            <h2 class="text-main" style="color: #FFFFFF; font-size: 20px; font-weight: bold; margin-bottom: 10px;">¡Hola, Gopher! 👋</h2>
                            <p class="text-main" style="color: #CBD5E1; font-size: 15px; line-height: 1.5;">
                                Tu certificación está lista. Has dado un gran paso en tu aprendizaje y es el momento de validar tus conocimientos. Aquí tienes los detalles de tu prueba agendada:
                            </p>
                        </td>
                    </tr>

                    <tr>
                        <td align="center" style="padding-bottom: 30px;">
                            <table class="info-box" width="100%" border="0" cellpadding="15" cellspacing="0" style="background-color: #1E293B; border-radius: 8px; border: 1px solid #334155;">
                                <tr>
                                    <td align="left" style="color: #FFFFFF; font-size: 14px; line-height: 1.8;">
                                        <strong>Prueba agendada:</strong> (Ej: Fundamentos de Go)<br>
                                        <strong>Nivel:</strong> Junior<br>
                                        <strong>Horario Agendado:</strong> 2026-03-01 a las 15:00<br>
                                        <strong>Duración Máxima:</strong> 90 Minutos<br>
                                        <strong>Requisito de Aprobación:</strong> 75% Mínimo
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>

                    <tr>
                        <td align="center">
                            <p class="text-main" style="color: #CBD5E1; font-size: 14px; margin-bottom: 15px;">Asegúrate de tener una conexión estable y un entorno tranquilo.</p>
                            <h2 class="text-main" style="color: #FFFFFF; font-size: 22px; font-weight: bold; margin-bottom: 30px; letter-spacing: 1px;">SUERTE EN TU EVALUACIÓN</h2>
                            
                            <p style="color: #64748B; font-size: 12px; line-height: 1.5;">
                                © 2026 Itudy. Todos los derechos reservados.<br>
                                Si tienes problemas, contáctanos a <a href="mailto:soporte@itudy.com" style="color: #3B82F6; text-decoration: none;">soporte@itudy.com</a>
                            </p>
                        </td>
                    </tr>

                </table>
            </td>
        </tr>
    </table>

</body>
</html>
`

func (s *service) SendEmail(ctx context.Context, date, email, nameTest string) error {

	// 2. Definimos las variables faltantes (Idealmente vienen de tu BD)
	difficultyLevel := "Junior"
	durationMinutes := "90"
	passingPercentage := "75"
	scheduledTime := "15:00" // Deberás parsear la hora de tu variable 'date'

	// 3. Reemplazamos las variables en el string de HTML
	htmlContent := certEmailHTML
	htmlContent = strings.ReplaceAll(htmlContent, "{{test_title}}", nameTest)
	htmlContent = strings.ReplaceAll(htmlContent, "{{difficulty_level}}", difficultyLevel)
	htmlContent = strings.ReplaceAll(htmlContent, "{{scheduled_date}}", date)
	htmlContent = strings.ReplaceAll(htmlContent, "{{scheduled_time}}", scheduledTime)
	htmlContent = strings.ReplaceAll(htmlContent, "{{duration_minutes}}", durationMinutes)
	htmlContent = strings.ReplaceAll(htmlContent, "{{passing_percentage}}", passingPercentage)

	params := &resend.SendEmailRequest{
		From:    "Itudy <no-reply@itudy.app>", // Asegúrate de haber verificado este dominio
		To:      []string{email},
		Subject: "Nueva tarea asignada: ",
		Html:    htmlContent,
	}

	sent, err := s.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("error al enviar email con Resend: %w", err)
	}

	fmt.Printf("Email enviado con ID: %s\n", sent.Id)
	return nil
}
