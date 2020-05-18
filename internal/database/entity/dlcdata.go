package entity

import (
	"time"

	"github.com/jinzhu/gorm"
)

// DLCData represents the db model of the oracle data rvalue/signature relative an asset
type DLCData struct {
	Base
	PublishedDate time.Time `gorm:"primary_key"`
	AssetID       string    `gorm:"primary_key"`
	Rvalue        string    `gorm:"unique;not null"`
	Signature     string
	Value         string
	Asset         Asset `gorm:"association_foreignkey:AssetID" json:"-"`

	// TODO should be stored somewhere secure
	Kvalue string `gorm:"unique;not null" json:"-"`
}

// IsSigned returns true if the Signature is set
func (m *DLCData) IsSigned() bool {
	return m.Signature != ""
}

// CreateDLCData will try to create a DLCData with a new Rvalue corresponding to an asset and publishDate
// if already in db, it will return the value found with no error
func CreateDLCData(db *gorm.DB, assetID string, publishDate time.Time, signingk string, rvalue string) (*DLCData, error) {
	tx := db.Begin()

	newDLCData := &DLCData{
		PublishedDate: publishDate,
		AssetID:       assetID,
		Kvalue:        signingk,
		Rvalue:        rvalue,
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
	err := db.Where(filterCondition).Find(dlcData).Error
	if err != nil {
		return nil, err
	}
	return dlcData, nil
}

// UpdateDLCDataSignatureAndValue will try to update signature and value of the DLCData if it exists
// and if the DLCdata is not already signed
func UpdateDLCDataSignatureAndValue(db *gorm.DB, assetID string, publishDate time.Time, sig string, value string) (*DLCData, error) {
	tx := db.Begin()
	filterCondition := &DLCData{
		AssetID:       assetID,
		PublishedDate: publishDate,
	}
	tx = tx.Model(filterCondition)
	// ensure that the signature and value are empty, doesn't work in using filterCondition (ignored)
	tx = tx.Where("signature = ?", "").Where("value = ?", "")

	tx = tx.Updates(DLCData{Signature: sig, Value: value})
	if err := tx.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

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
