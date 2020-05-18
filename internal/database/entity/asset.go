package entity

import (
	"github.com/jinzhu/gorm"
)

// Asset represents an asset currency
type Asset struct {
	Base
	AssetID     string `gorm:"primarykey"`
	Description string
	DLCData     []DLCData `gorm:"foreignkey:AssetID"`
}

// FindAsset will try to find in the db the asset corresponding to the id
func FindAsset(db *gorm.DB, assetID string) (*Asset, error) {
	existingAsset := &Asset{AssetID: assetID}
	err := db.First(existingAsset).Error
	return existingAsset, err
}
