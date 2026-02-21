package payin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Yeremi528/itudy-back/email"
	"github.com/Yeremi528/itudy-back/movements"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

// type service is the implementation of Service interface containing all the business logic
// and dependencies required to complete the given tasks without exposing the implementation.
type service struct {
	repository          Repository
	mpPreferenceClient  preference.Client
	mpPaymentClient     payment.Client
	repositoryMovement  movements.Repository
	rAssignmentFlexible repositoryAssignments
	emailService        email.Service
}

type MercadoPagoConfig struct {
	AccessToken string
}

type Config struct {
	MercadoPago MercadoPagoConfig
}

func NewService(r Repository, rMovements movements.Repository, rAssignments repositoryAssignments, cfg Config, emailService email.Service) (*service, error) {
	cfgMercadoPago, err := config.New(cfg.MercadoPago.AccessToken)
	if err != nil {
		return &service{}, fmt.Errorf("error configurando Mercado Pago: %w", err)
	}

	paymentMercadoPago := payment.NewClient(cfgMercadoPago)
	preferenceMercadoPago := preference.NewClient(cfgMercadoPago)
	return &service{
		repository:          r,
		mpPreferenceClient:  preferenceMercadoPago,
		mpPaymentClient:     paymentMercadoPago,
		repositoryMovement:  rMovements,
		rAssignmentFlexible: rAssignments,
		emailService:        emailService,
	}, nil
}

// Tener una segunda opción, ya que esto no era buena practica JAJ LOL XD
func (s *service) RechargeLink(ctx context.Context, email, IDAssignment, examName string, amount int) (string, error) {

	request := preference.Request{
		Items: []preference.ItemRequest{
			{
				Title:      examName,
				UnitPrice:  float64(amount),
				Quantity:   1,
				CurrencyID: "CLP",
			},
		},
		Metadata: map[string]any{
			"email":      email,
			"assignment": IDAssignment,
		},
		BackURLs: &preference.BackURLsRequest{
			Success: "https://colectigo-back-uwu-579390796383.southamerica-west1.run.app",
			Failure: "https://colectigo-back-uwu-579390796383.southamerica-west1.run.app",
			Pending: "https://colectigo-back-uwu-579390796383.southamerica-west1.run.app",
		},
		// Usar webhooks V1 suele ser la recomendación actual de MP
		NotificationURL: "https://itudy-947017986235.southamerica-west1.run.app/payin",
		AutoReturn:      "approved",
	}

	// CORRECCIÓN: Usar mpPreferenceClient en lugar de mpPaymentClient
	// Y pasar el ctx de la función en lugar de context.Background()
	resource, err := s.mpPreferenceClient.Create(ctx, request)
	if err != nil {
		return "", fmt.Errorf("error creando preferencia: %w", err)
	}

	return resource.InitPoint, nil
}

func (s *service) WebHook(ctx context.Context, ID, topic string) error {
	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return err
	}

	// TO DO: Auditar el correo
	mercadoPago := MercadoPago{ID: ID, Topic: topic}
	if err := s.repository.Audit(ctx, mercadoPago); err != nil {
		return fmt.Errorf("payin.webhook: repository.audit: %w", err)
	}

	switch topic {
	case "payment":
		pay, err := s.mpPaymentClient.Get(ctx, int(id))
		if err != nil {
			return fmt.Errorf("error obteniendo detalles del pago: %w", err)
		}

		// Verificamos que el pago realmente esté aprobado antes de liberar el horario
		if pay.Status != "approved" {
			fmt.Printf("El pago %d está en estado: %s. No se libera el horario aún.\n", id, pay.Status)
			return nil
		}

		fmt.Println("Entramos en payments, pago aprobado!")

		// OJO: Asegúrate de que createMovement reciba los datos correctos.
		// Si necesitas el email, lo puedes sacar de pay.Metadata["email"]
		user := createMovement(pay.Payer.Identification.Number, pay.TransactionDetails.TotalPaidAmount)
		if err := s.Payin(ctx, user); err != nil {
			return fmt.Errorf("payin.webhook: s.Payin: %w", err)
		}

		assignmentID := pay.Metadata["assignment"]
		if err := s.rAssignmentFlexible.UpdateAssignment(ctx, fmt.Sprintf("%v", assignmentID)); err != nil {
			return fmt.Errorf("payin.Webhook: assignmentsFlexible.Update: %w", err)
		}

		if err := s.emailService.SendEmail(ctx, "2025-03", "geremiararaya@gmail.com", "panchito"); err != nil {
			fmt.Println("error al enviar el correo", err)
		}

		return nil

	case "merchant_order", "merchant":
		// CORRECCIÓN: Ignoramos la orden comercial pero devolvemos nil para que el handler HTTP envíe un 200 OK
		fmt.Printf("Notificación de merchant_order recibida (ID: %s). Ignorando...\n", ID)
		return nil

	default:
		return nil
	}

}

func (s *service) Payin(ctx context.Context, movement movements.Movement) error {

	err := s.repository.Update(ctx, Pay{
		Email:  movement.Email,
		Amount: movement.Amount,
	})
	if err != nil {
		return fmt.Errorf("payin.Payin: repository.Update: %w", err)
	}

	err = s.repositoryMovement.Insert(ctx, movement)
	if err != nil {
		return fmt.Errorf("payin.Payin: repositoryMovement.Insert: %w", err)
	}

	return nil
}
