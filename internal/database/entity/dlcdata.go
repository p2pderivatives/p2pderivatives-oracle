package entity

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DLCData represents the db model of the oracle data rvalue/signature relative an asset
type DLCData struct {
	Base
	PublishedDate time.Time   `gorm:"primary_key"`
	AssetID       string      `gorm:"primary_key"`
	Rvalues       StringArray `gorm:"not null"`
	Signatures    StringArray
	Values        StringArray
	Asset         Asset `gorm:"foreignkey:AssetID" json:"-"`

	// TODO should be stored somewhere secure
	Kvalues StringArray `gorm:"not null" json:"-"`
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

// IsSigned returns true if the Signature is set
func (m *DLCData) IsSigned() bool {
	return len(m.Signatures) > 0
}

// CreateDLCData will try to create a DLCData with a new Rvalue corresponding to an asset and publishDate
// if already in db, it will return the value found with no error
func CreateDLCData(db *gorm.DB, assetID string, publishDate time.Time, signingks []string, rvalues []string) (*DLCData, error) {
	tx := db.Begin()

	newDLCData := &DLCData{
		PublishedDate: publishDate,
		AssetID:       assetID,
		Kvalues:       signingks,
		Rvalues:       rvalues,
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
func FindDLCDataPublishedNear(db *gorm.DB, assetID string, nearTime time.Time, rangeD time.Duration) (*DLCData, error) {
	dlcData := &DLCData{}
	filterCondition := &DLCData{
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
func FindDLCDataPublishedBefore(db *gorm.DB, assetID string, nearTime time.Time) (*DLCData, error) {
	dlcData := &DLCData{}
	filterCondition := &DLCData{
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
func FindDLCDataPublishedAt(db *gorm.DB, assetID string, publishDate time.Time) (*DLCData, error) {
	dlcData := &DLCData{}
	filterCondition := &DLCData{
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
func UpdateDLCDataSignatureAndValue(db *gorm.DB, assetID string, publishDate time.Time, sigs []string, values []string) (*DLCData, error) {
	filterCondition := &DLCData{
		AssetID:       assetID,
		PublishedDate: publishDate,
	}
	db.Logger.LogMode(logger.Info)
	tx := db.Begin()
	tx = tx.Where(filterCondition)

	var old DLCData
	tx.First(&old)

	if old.Signatures != nil || old.Values != nil {
		tx.Rollback()
		return nil, errors.New("Already signed or assigned values")
	}

	tx = tx.Updates(DLCData{Signatures: sigs, Values: values})

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
