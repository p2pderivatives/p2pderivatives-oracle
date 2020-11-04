package api

import (
	"p2pderivatives-oracle/internal/database/entity"
	"p2pderivatives-oracle/internal/dlccrypto"
	"time"
)

// NewOracleAnnouncement converts a DLCData structure to an oracle announcement
func NewOracleAnnouncement(
	oraclePubKey *dlccrypto.SchnorrPublicKey,
	eventData *entity.EventData) *OracleAnnouncement {
	descriptor := DecompositionDescriptor{
		Base:      eventData.Base,
		IsSigned:  eventData.IsSigned,
		Unit:      eventData.Unit,
		Precision: eventData.Precision,
	}
	event := OracleEvent{
		Nonces:          eventData.Nonces,
		EventMaturity:   eventData.PublishedDate,
		EventDescriptor: descriptor,
		EventID:         eventData.GetEventID(),
	}
	announcement := OracleAnnouncement{
		AnnouncementSignature: eventData.AnnouncementSignature,
		OraclePublicKey:       oraclePubKey.EncodeToString(),
		OracleEvent:           event,
	}

	return &announcement
}

// NewOracleAttestation creates a new OracleAttestation structure from the given eventData
func NewOracleAttestation(eventData *entity.EventData) *OracleAttestation {
	return &OracleAttestation{
		EventID:    eventData.GetEventID(),
		Signatures: eventData.Signatures,
		Values:     eventData.Values,
	}
}

// DecompositionDescriptor contains information to about an event using
// numerical decomposition
type DecompositionDescriptor struct {
	Base      int    `json:"base"`
	IsSigned  bool   `json:"isSigned"`
	Unit      string `json:"unit"`
	Precision int    `json:"precision"`
}

// OracleEvent contains information about an event
type OracleEvent struct {
	Nonces          []string                `json:"nonces"`
	EventMaturity   time.Time               `json:"eventMaturity"`
	EventDescriptor DecompositionDescriptor `json:"eventDescriptor"`
	EventID         string                  `json:"eventId"`
}

// OracleAnnouncement contains information about an event and a signature over
// the OracleEvent structure
type OracleAnnouncement struct {
	AnnouncementSignature string      `json:"announcementSignature"`
	OraclePublicKey       string      `json:"oraclePublicKey"`
	OracleEvent           OracleEvent `json:"oracleEvent"`
}

// OracleAttestation contains information about the outcome of an event
type OracleAttestation struct {
	EventID    string   `json:"eventId"`
	Signatures []string `json:"signatures"`
	Values     []string `json:"values"`
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
