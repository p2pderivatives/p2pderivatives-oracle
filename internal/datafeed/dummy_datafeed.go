package datafeed

import (
	"time"
)

// NewDummyDataFeed returns a dummy datafeed !
func NewDummyDataFeed(config *DummyConfig) DataFeed {
	return &dummyDataFeed{
		config: config,
	}
}

type dummyDataFeed struct {
	config *DummyConfig
}

func (d *dummyDataFeed) FindCurrentAssetPrice(assetID string) (*float64, error) {
	f := d.config.ReturnValue
	return &f, nil
}

func (d *dummyDataFeed) FindPastAssetPrice(assetID string, date time.Time) (*float64, error) {
	f := d.config.ReturnValue
	return &f, nil
}

// DummyConfig configuration for the dummy Datafeed
type DummyConfig struct {
	ReturnValue float64 `configkey:"dummy.returnValue" validate:"required"`
}
