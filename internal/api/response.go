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
		DigitDecompositionDescriptor: DigitDecompositionDescriptor{
			Base:      eventData.Base,
			IsSigned:  eventData.IsSigned,
			Unit:      eventData.Unit,
			Precision: eventData.Precision,
			NbDigits:  len(eventData.Nonces),
		},
	}
	event := OracleEvent{
		Nonces:             eventData.Nonces,
		EventMaturityEpoch: eventData.PublishedDate.Unix(),
		EventDescriptor:    descriptor,
		EventID:            eventData.GetEventID(),
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

// DigitDecompositionDescriptor contains information about a numerical event.
type DigitDecompositionDescriptor struct {
	Base      int    `json:"base"`
	IsSigned  bool   `json:"isSigned"`
	Unit      string `json:"unit"`
	Precision int    `json:"precision"`
	NbDigits  int    `json:"nbDigits"`
}

// DecompositionDescriptor can contain information about either an enumerable event
// or a numerical event.
type DecompositionDescriptor struct {
	DigitDecompositionDescriptor DigitDecompositionDescriptor `json:"digitDecompositionEvent"`
}

// OracleEvent contains information about an event
type OracleEvent struct {
	Nonces             []string                `json:"oracleNonces"`
	EventMaturityEpoch int64                   `json:"eventMaturityEpoch"`
	EventDescriptor    DecompositionDescriptor `json:"eventDescriptor"`
	EventID            string                  `json:"eventId"`
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
