package mercadopago

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Yeremi528/itudy-back/kit/httpcaller"
	"github.com/Yeremi528/itudy-back/kit/logger"
	"github.com/Yeremi528/itudy-back/kit/mask"
)

type ServiceConfig struct {
	Base  string
	Token string
}

type Config struct {
	ClientID string
	Password string
}

type Core struct {
	Cfg         Config
	logger      *logger.Logger
	Credentials Token
}

func New(cfg Config, logger *logger.Logger) *Core {
	c := &Core{
		Cfg:    cfg,
		logger: logger,
	}

	credentials, err := c.Token()
	if err != nil {
		panic(err)
	}

	c.Credentials = credentials

	return c

}

func (c *Core) Token() (Token, error) {
	ctx := context.Background()

	data := map[string]string{
		"client_id":     c.Cfg.ClientID,
		"client_secret": c.Cfg.Password,
		"grant_type":    "client_credentials",
	}

	var headers = map[string]string{
		"Content-Type": "application/json",
	}

	req := httpcaller.RequestParams{
		Client:     &http.Client{},
		URL:        "https://api.mercadopago.com/oauth/token",
		Headers:    headers,
		Body:       nil,
		QueryParam: nil,
		Urlencoded: data,
	}
	var resp Token
	since, res, status, err := httpcaller.POST(ctx, req, &resp)
	if err != nil {
		c.logUnsuccessfulCall(ctx, "POST - TOKEN - FAILED", req.URL, nil, res, status, since)
		return Token{}, fmt.Errorf("POST: %w", err)
	}

	switch status {
	case http.StatusOK:
		return resp, nil
	default:
		return Token{}, err
	}

}

func (c *Core) logUnsuccessfulCall(ctx context.Context, message, URL string, req []byte, res []byte, status int, since time.Duration) {
	maskedResponse, _ := mask.JSONBytes(res, "access_token", "id_token")

	c.logger.Errorc(ctx, 4, message, "URL", URL, "status", status, "request", string(req), "response", string(maskedResponse), "since", since.String())
}
