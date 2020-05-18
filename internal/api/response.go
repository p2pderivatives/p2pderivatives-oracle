package api

import (
	"p2pderivatives-oracle/internal/database/entity"
	"p2pderivatives-oracle/internal/utils/iso8601"
	"time"
)

// DLCDataResponse represents the DLC data struct sent by AssetController
type DLCDataResponse struct {
	PublishedDate time.Time `json:"publish_date"`
	AssetID       string    `json:"asset"`
	Rvalue        string    `json:"rvalue"`
	Signature     string    `json:"signature,omitempty"`
	Value         string    `json:"value,omitempty"`
}

// AssetConfigResponse represents the configuration of an asset api
type AssetConfigResponse struct {
	Frequency string `json:"frequency"`
	RangeD    string `json:"range"`
}

// OraclePublicKeyResponse represents the public key of the oracle
type OraclePublicKeyResponse struct {
	PublicKey string `json:"public_key"`
}

// NewResponse returns a new response from a matching response object
func NewResponse(o interface{}) (res interface{}) {
	switch o.(type) {
	case *entity.DLCData:
		t := o.(*entity.DLCData)
		res = &DLCDataResponse{
			PublishedDate: t.PublishedDate.UTC(),
			AssetID:       t.AssetID,
			Rvalue:        t.Rvalue,
			Signature:     t.Signature,
			Value:         t.Value,
		}
	case AssetConfig:
		t := o.(AssetConfig)
		res = &AssetConfigResponse{
			Frequency: iso8601.EncodeDuration(t.Frequency),
			RangeD:    iso8601.EncodeDuration(t.RangeD),
		}

	default:
		return o
	}
	return res
}
