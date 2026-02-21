package mercadopago

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	UserID       int    `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	PublicKey    string `json:"public_key"`
	LiveMode     bool   `json:"live_mode"`
}
