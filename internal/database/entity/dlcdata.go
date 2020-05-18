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

// CreateDLCData create a DLCData with a new Rvalue corresponding to an asset and publishDate
// can be serialize to avoid database concurrency problem
func CreateDLCData(db *gorm.DB, assetID string, publishDate time.Time, signingk string, rvalue string, serialize bool) (*DLCData, error) {
	var tx *gorm.DB

	// watch for concurrency issue
	if serialize {
		db.Exec("set transaction isolation level serializable") // this statement ensures synchronicity at the database level
	}

	tx = db.Begin()

	newDLCData := &DLCData{}
	err := tx.Where("asset_id = ?", assetID).Where("published_date = ?", publishDate).First(newDLCData).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		tx.Rollback()
		return nil, err
	}

	newDLCData = &DLCData{
		PublishedDate: publishDate,
		AssetID:       assetID,
		Kvalue:        signingk,
		Rvalue:        rvalue,
	}

	if err = tx.Create(newDLCData).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return newDLCData, nil
}

// FindDLCDataPublishedNear will try to retrieve the oldest dlcData which has been published between nearTime and rangeD
// limit
func FindDLCDataPublishedNear(db *gorm.DB, assetID string, nearTime time.Time, rangeD time.Duration) (*DLCData, error) {
	dlcData := &DLCData{}
	req := db.Where("asset_id = ?", assetID)
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
	req := db.Where("asset_id = ?", assetID)
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
	err := db.Where("asset_id = ?", assetID).Where("published_date = ?", publishDate).Find(dlcData).Error
	if err != nil {
		return nil, err
	}
	return dlcData, nil
}

// PutDLCData will try to update signature and value of the DLCData if it exists
func UpdateDLCDataSignatureAndValue(db *gorm.DB, assetID string, publishDate time.Time, sig string, value string) error {
	err := db.Model(&DLCData{AssetID: assetID, PublishedDate: publishDate}).Updates(DLCData{Signature: sig, Value: value}).Error
	return err
}
