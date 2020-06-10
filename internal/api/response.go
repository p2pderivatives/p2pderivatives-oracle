package api

import (
	"p2pderivatives-oracle/internal/database/entity"
	"p2pderivatives-oracle/internal/dlccrypto"
	"time"
)

// NewDLCDataReponse transforms a entity.DLCData to dlcData response
func NewDLCDataReponse(
	oraclePubKey *dlccrypto.PublicKey,
	dlcData *entity.DLCData) *DLCDataResponse {
	return &DLCDataResponse{
		OraclePublicKey: oraclePubKey.EncodeToString(),
		PublishedDate:   dlcData.PublishedDate,
		AssetID:         dlcData.AssetID,
		Rvalue:          dlcData.Rvalue,
		Signature:       dlcData.Signature,
		Value:           dlcData.Value,
	}
}

// DLCDataResponse represents the DLC data struct sent by AssetController
type DLCDataResponse struct {
	OraclePublicKey string    `json:"oraclePublicKey"`
	PublishedDate   time.Time `json:"publishDate"`
	AssetID         string    `json:"asset"`
	Rvalue          string    `json:"rvalue"`
	Signature       string    `json:"signature,omitempty"`
	Value           string    `json:"value,omitempty"`
}

// AssetConfigResponse represents the configuration of an asset api
type AssetConfigResponse struct {
	StartDate time.Time `json:"startDate"`
	Frequency string    `json:"frequency"`
	RangeD    string    `json:"range"`
}

// OraclePublicKeyResponse represents the public key of the oracle
type OraclePublicKeyResponse struct {
	PublicKey string `json:"publicKey"`
}
