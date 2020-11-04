package entity

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// EventData represents the db model of the oracle data rvalue/signature relative an asset
type EventData struct {
	Timestamp
	PublishedDate         time.Time   `gorm:"primary_key"`
	AssetID               string      `gorm:"primary_key"`
	Nonces                StringArray `gorm:"not null"`
	Signatures            StringArray
	Values                StringArray
	Asset                 Asset `gorm:"foreignkey:AssetID" json:"-"`
	AnnouncementSignature string
	Base                  int
	Unit                  string
	IsSigned              bool
	Precision             int

	// TODO should be stored somewhere secure
	Kvalues StringArray `gorm:"not null" json:"-"`
}

// GetEventID returns the event ID for the given eventData structure
func (eventData *EventData) GetEventID() string {
	return eventData.AssetID + strconv.FormatInt(eventData.PublishedDate.Unix(), 10)
}

// StringArray is an alias type for an array of strings
type StringArray []string

// Scan implements the Scanner interface for gorm custom types
func (s *StringArray) Scan(value interface{}) error {
	csv, ok := value.(string)
	if !ok {
		return fmt.Errorf("Failed to unmarshal string value: %v", value)
	}

	*s = strings.Split(csv, ",")

	return nil
}

// Value implements the Valuer interface for gorm custom types
func (s StringArray) Value() (driver.Value, error) {
	if s == nil || len(s) == 0 {
		return nil, nil
	}
	return strings.Join(s, ","), nil
}

// GormDataType implements the GormDataTypeInterface for gorm custom types
func (StringArray) GormDataType() string {
	return "text"
}

// HasSignature returns true if the Signature is set
func (eventData *EventData) HasSignature() bool {
	return len(eventData.Signatures) > 0
}

// CreateEventData will try to create a DLCData with a new Rvalue corresponding to an asset and publishDate
// if already in db, it will return the value found with no error
func CreateEventData(db *gorm.DB, assetID string, publishDate time.Time, signingks []string, rvalues []string, base int) (*EventData, error) {
	tx := db.Begin()

	newDLCData := &EventData{
		PublishedDate: publishDate,
		AssetID:       assetID,
		Kvalues:       signingks,
		Nonces:        rvalues,
		Base:          base,
	}

	tx = tx.Create(newDLCData)

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return newDLCData, nil
}

// FindDLCDataPublishedNear will try to retrieve the oldest dlcData which has been published between nearTime and rangeD
// limit
func FindDLCDataPublishedNear(db *gorm.DB, assetID string, nearTime time.Time, rangeD time.Duration) (*EventData, error) {
	dlcData := &EventData{}
	filterCondition := &EventData{
		AssetID: assetID,
	}
	req := db.Where(filterCondition)
	limit := nearTime.Add(rangeD)
	// look for the oldest
	req = req.Where("published_date BETWEEN ? AND ?", nearTime, limit)
	req = req.Order("published_date ASC")
	err := req.First(dlcData).Error
	if err != nil {
		return nil, err
	}
	return dlcData, nil
}

// FindDLCDataPublishedBefore will try to retrieve the most recent dlcData which has been published BEFORE a specific
// date
func FindDLCDataPublishedBefore(db *gorm.DB, assetID string, nearTime time.Time) (*EventData, error) {
	dlcData := &EventData{}
	filterCondition := &EventData{
		AssetID: assetID,
	}
	req := db.Where(filterCondition)
	// look for the most recent
	req = req.Where("published_date < ?", nearTime)
	req = req.Order("published_date DESC")
	err := req.First(dlcData).Error
	if err != nil {
		return nil, err
	}
	return dlcData, nil
}

// FindDLCDataPublishedAt will try to retrieve asset dlcData at specific publish date
// from database
func FindDLCDataPublishedAt(db *gorm.DB, assetID string, publishDate time.Time) (*EventData, error) {
	dlcData := &EventData{}
	filterCondition := &EventData{
		AssetID:       assetID,
		PublishedDate: publishDate,
	}
	err := db.Where(filterCondition).First(dlcData).Error
	if err != nil {
		return nil, err
	}
	return dlcData, nil
}

// UpdateDLCDataSignatureAndValue will try to update signature and value of the DLCData if it exists
// and if the DLCdata is not already signed
func UpdateDLCDataSignatureAndValue(db *gorm.DB, assetID string, publishDate time.Time, sigs []string, values []string) (*EventData, error) {
	filterCondition := &EventData{
		AssetID:       assetID,
		PublishedDate: publishDate,
	}
	db.Logger.LogMode(logger.Info)
	tx := db.Begin()
	tx = tx.Where(filterCondition)

	var old EventData
	tx.First(&old)

	if old.Signatures != nil || old.Values != nil {
		tx.Rollback()
		return nil, errors.New("Already signed or assigned values")
	}

	tx = tx.Updates(EventData{Signatures: sigs, Values: values})

	if tx.RowsAffected == 0 {
		tx.Rollback()
	} else {
		err := tx.Commit().Error
		if err != nil {
			return nil, err
		}
	}

	return FindDLCDataPublishedAt(db, assetID, publishDate)
}
