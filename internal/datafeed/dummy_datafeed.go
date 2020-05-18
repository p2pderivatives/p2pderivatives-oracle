package datafeed

import (
	"encoding/binary"
	"math/rand"
	"time"
)

// NewDummyDlcOracle returns a dummy DLCOracle, does not implement any crypto !
func NewDummyDataFeed() DataFeed {
	return &dummyDataFeed{}
}

type dummyDataFeed struct{}

func (d *dummyDataFeed) FindCurrentAssetPrice(assetID string, currency string) (uint64, error) {
	return randomUint64(), nil
}

func (d *dummyDataFeed) FindPastAssetPrice(assetID string, currency string, date time.Time) (uint64, error){
	return randomUint64(), nil
}

func randomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf) // random so bytes order irrelevant
}