package entity_test

import (
	"p2pderivatives-oracle/internal/database/entity"
	"p2pderivatives-oracle/test"
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

func Test_FindAsset_Present_ReturnsCorrectAsset(t *testing.T) {
	db := test.NewOrm(&entity.Asset{}).GetDB()
	expected := &entity.Asset{
		AssetID:     "test",
		Description: "Some test Asset",
	}
	db.Create(expected)
	actual, err := entity.FindAsset(db, expected.AssetID)
	assert.NoError(t, err)
	assert.Equal(t, expected.AssetID, actual.AssetID)
	assert.Equal(t, expected.Description, actual.Description)
}

func Test_FindAsset_NotPresent_ReturnsRecordNotFoundError(t *testing.T) {
	db := test.NewOrm(&entity.Asset{}).GetDB()
	_, err := entity.FindAsset(db, "invalid asset")
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}
