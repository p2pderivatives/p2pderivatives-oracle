package datafeed

import "time"

// DataFeed interface represents a datafeed with any sorts of data
type DataFeed interface {
	AssetPriceFeed
}

// AssetPriceFeed interface represents a datafeed which implemented price related services
type AssetPriceFeed interface {
	FindCurrentAssetPrice(assetID string, currency string) (*float64, error)
	FindPastAssetPrice(assetID string, currency string, date time.Time) (*float64, error)
}
