package api

import "time"

// Config contains the API configuration
type Config struct {
	AssetConfigs map[string]AssetConfig `configkey:"api.assets" validate:"required"`
}

// SigningConfig contains parameters for the oracle to sign event outcomes
type SigningConfig struct {
	Base      int  `configkey:"base" validate:"required"`
	NbDigits  int  `configkey:"nbDigits" validate:"required"`
	IsSigned  bool `configkey:"isSigned"`
	Precision int  `configkey:"precision"`
}

// AssetConfig represents one asset configuration delivered by the oracle
type AssetConfig struct {
	StartDate  time.Time     `configkey:"startDate" validate:"required"`
	Frequency  time.Duration `configkey:"frequency,duration,iso8601" validate:"required"`
	RangeD     time.Duration `configkey:"range,duration,iso8601" validate:"required"`
	SignConfig SigningConfig `configkey:"signconfig" validate:"required"`
	Unit       string        `configkey:"unit" validate:"required"`
}
