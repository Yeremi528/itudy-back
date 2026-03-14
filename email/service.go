package email

import (
	"context"
	"fmt"
	"strings"

	"github.com/Yeremi528/itudy-back/exam"
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
        /* Reset básico */
        body, table, td, p, h1, h2, h3 {
            margin: 0;
            padding: 0;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
        }
        
        table {
            border-spacing: 0;
            border-collapse: collapse;
        }

        /* Utilidades responsivas */
        .wrapper {
            width: 100%;
            table-layout: fixed;
            background-color: #F8FAFC;
            padding-bottom: 40px;
        }
        
        .main {
            background-color: #FFFFFF;
            margin: 0 auto;
            width: 100%;
            max-width: 600px;
            border-radius: 12px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
            overflow: hidden;
        }

        /* Modo Oscuro */
        @media (prefers-color-scheme: dark) {
            .wrapper {
                background-color: #0B0F19 !important;
            }
            .main {
                background-color: #151E2E !important;
                box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.3) !important;
            }
            .text-title {
                color: #FFFFFF !important;
            }
            .text-body {
                color: #CBD5E1 !important;
            }
            .info-box {
                background-color: #1E293B !important;
                border-left: 4px solid #3B82F6 !important;
            }
            .info-text {
                color: #F8FAFC !important;
            }
            .footer-text {
                color: #64748B !important;
            }
        }

        /* Ajustes para móviles */
        @media screen and (max-width: 600px) {
            .main {
                border-radius: 0 !important;
            }
            .pad-mobile {
                padding: 30px 20px !important;
            }
        }
    </style>
</head>
<body style="margin: 0; padding: 0; background-color: #F8FAFC; -webkit-text-size-adjust: 100%; text-size-adjust: 100%;">

    <center class="wrapper" style="width: 100%; background-color: #F8FAFC; padding: 40px 0;">
        <table class="main" width="100%" border="0" cellpadding="0" cellspacing="0" style="max-width: 600px; background-color: #FFFFFF; border-radius: 12px; margin: 0 auto;">
            
            <tr>
                <td style="height: 6px; background: linear-gradient(90deg, #3B82F6 0%, #06B6D4 100%);"></td>
            </tr>

            <tr>
                <td class="pad-mobile" style="padding: 40px 40px 30px 40px;">
                    
                    <table width="100%" border="0" cellpadding="0" cellspacing="0">
                        <tr>
                            <td align="center" style="padding-bottom: 30px;">
                                <h1 class="text-title" style="color: #0F172A; font-size: 28px; font-weight: 800; letter-spacing: 2px; margin: 0;">ITUDY</h1>
                                <p style="color: #3B82F6; font-size: 13px; font-weight: 700; letter-spacing: 1.5px; text-transform: uppercase; margin-top: 5px;">Ruta Backend</p>
                            </td>
                        </tr>
                    </table>

                    <table width="100%" border="0" cellpadding="0" cellspacing="0">
                        <tr>
                            <td align="left" style="padding-bottom: 25px;">
                                <h2 class="text-title" style="color: #1E293B; font-size: 22px; font-weight: 700; margin-bottom: 12px;">¡Hola, Gopher! 👋</h2>
                                <p class="text-body" style="color: #475569; font-size: 16px; line-height: 1.6; margin: 0;">
                                    Tu certificación está lista. Has dado un gran paso en tu aprendizaje y es el momento de validar tus conocimientos. Aquí tienes los detalles de tu prueba agendada:
                                </p>
                            </td>
                        </tr>
                    </table>

                    <table width="100%" border="0" cellpadding="0" cellspacing="0">
                        <tr>
                            <td align="center" style="padding-bottom: 30px;">
                                <table class="info-box" width="100%" border="0" cellpadding="20" cellspacing="0" style="background-color: #F1F5F9; border-radius: 8px; border-left: 4px solid #3B82F6;">
                                    <tr>
                                        <td align="left">
                                            <p class="info-text" style="color: #1E293B; font-size: 15px; line-height: 2; margin: 0;">
                                                <strong style="color: #3B82F6;">Prueba agendada:</strong> Fundamentos de Go<br>
                                                <strong style="color: #3B82F6;">Nivel:</strong> Junior<br>
                                                <strong style="color: #3B82F6;">Horario:</strong> 2026-03-01 a las 15:00<br>
                                                <strong style="color: #3B82F6;">Duración Máxima:</strong> 90 Minutos<br>
                                                <strong style="color: #3B82F6;">Aprobación:</strong> 75% Mínimo
                                            </p>
                                        </td>
                                    </tr>
                                </table>
                            </td>
                        </tr>
                    </table>

                    <table width="100%" border="0" cellpadding="0" cellspacing="0">
                        <tr>
                            <td align="center">
                                <p class="text-body" style="color: #475569; font-size: 15px; margin-bottom: 20px;">Asegúrate de tener una conexión estable y un entorno tranquilo.</p>
                                <h2 class="text-title" style="color: #0F172A; font-size: 20px; font-weight: 800; letter-spacing: 0.5px; margin: 0;">¡SUERTE EN TU EVALUACIÓN!</h2>
                            </td>
                        </tr>
                    </table>

                </td>
            </tr>
            
            <tr>
                <td style="background-color: #F8FAFC; padding: 20px 40px; border-top: 1px solid #E2E8F0;">
                    <table width="100%" border="0" cellpadding="0" cellspacing="0">
                        <tr>
                            <td align="center">
                                <p class="footer-text" style="color: #94A3B8; font-size: 13px; line-height: 1.6; margin: 0;">
                                    © 2026 Itudy. Todos los derechos reservados.<br>
                                    Si tienes problemas, contáctanos a <a href="mailto:soporte@itudy.com" style="color: #3B82F6; text-decoration: none; font-weight: 500;">soporte@itudy.com</a>
                                </p>
                            </td>
                        </tr>
                    </table>
                </td>
            </tr>

        </table>
    </center>

</body>
</html>
`

func (s *service) SendEmail(ctx context.Context, examInfo exam.Exam, date, email string) error {

	htmlContent := certEmailHTML
	htmlContent = strings.ReplaceAll(htmlContent, "{{test_title}}", examInfo.Title)
	htmlContent = strings.ReplaceAll(htmlContent, "{{difficulty_level}}", examInfo.DifficultyLevel)
	htmlContent = strings.ReplaceAll(htmlContent, "{{scheduled_date}}", date)
	htmlContent = strings.ReplaceAll(htmlContent, "{{duration_minutes}}", fmt.Sprintf("%d", examInfo.DurationMinutes))
	htmlContent = strings.ReplaceAll(htmlContent, "{{passing_percentage}}", fmt.Sprintf("%+v", examInfo.PassingPercentage))

	fmt.Println("llegamos a exam")
	fmt.Printf("%+v", examInfo)
	fmt.Println(email, "email")
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
