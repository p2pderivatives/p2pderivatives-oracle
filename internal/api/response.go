package api

import (
	"p2pderivatives-oracle/internal/database/entity"
	"p2pderivatives-oracle/internal/dlccrypto"
	"time"
)

// NewDLCDataResponse transforms a entity.DLCData to dlcData response
func NewDLCDataResponse(
	oraclePubKey *dlccrypto.SchnorrPublicKey,
	dlcData *entity.DLCData) *DLCDataResponse {
	return &DLCDataResponse{
		OraclePublicKey: oraclePubKey.EncodeToString(),
		PublishedDate:   dlcData.PublishedDate,
		AssetID:         dlcData.AssetID,
		Rvalues:         dlcData.Rvalues,
		Signatures:      dlcData.Signatures,
		Values:          dlcData.Values,
	}
}

// DLCDataResponse represents the DLC data struct sent by AssetController
type DLCDataResponse struct {
	OraclePublicKey string    `json:"oraclePublicKey"`
	PublishedDate   time.Time `json:"publishDate"`
	AssetID         string    `json:"asset"`
	Rvalues         []string  `json:"rvalues"`
	Signatures      []string  `json:"signatures,omitempty"`
	Values          []string  `json:"values,omitempty"`
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
