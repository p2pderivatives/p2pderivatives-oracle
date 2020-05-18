package cryptocompare

// Config represents the crypto compare client configuration
type Config struct {
	APIBaseURL string `configkey:"cryptoCompare.baseUrl" validate:"required"`
	APIKey     string `configkey:"cryptoCompare.apiKey" validate:"required"`
}
