package cryptocompare

// Config represents the crypto compare client configuration
type Config struct {
	APIBaseURL   string                   `configkey:"cryptoCompare.baseUrl" validate:"required"`
	APIKey       string                   `configkey:"cryptoCompare.apiKey"`
	AssetsConfig map[string]CCAssetConfig `configkey:"cryptoCompare.assetsConfig" validate:"required"`
}

// CCAssetConfig contains the request parameters to use for an asset
type CCAssetConfig struct {
	Fsym string `configkey:"fsym" validate:"required"`
	Tsym string `configkey:"tsym" validate:"required"`
}
