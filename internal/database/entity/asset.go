package entity

import (
	"gorm.io/gorm"
)

// Asset represents an asset currency
type Asset struct {
	Timestamp
	AssetID     string `gorm:"primarykey"`
	Description string
	DLCData     []EventData `gorm:"foreignkey:AssetID"`
}

// FindAsset will try to find in the db the asset corresponding to the id
func FindAsset(db *gorm.DB, assetID string) (*Asset, error) {
	existingAsset := &Asset{AssetID: assetID}
	err := db.First(existingAsset).Error
	return existingAsset, err
}
