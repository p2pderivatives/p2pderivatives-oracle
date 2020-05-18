package datafeed

import (
	"math/rand"
	"time"
)

// NewDummyDataFeed returns a dummy datafeed !
func NewDummyDataFeed() DataFeed {
	return &dummyDataFeed{}
}

type dummyDataFeed struct{}

func (d *dummyDataFeed) FindCurrentAssetPrice(assetID string, currency string) (*float64, error) {
	f := randomfloat64()
	return &f, nil
}

func (d *dummyDataFeed) FindPastAssetPrice(assetID string, currency string, date time.Time) (*float64, error) {
	f := randomfloat64()
	return &f, nil
}

func randomfloat64() float64 {
	return rand.Float64()
}
