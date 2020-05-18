package entity_test

import (
	"p2pderivatives-oracle/internal/database/entity"
	"p2pderivatives-oracle/test"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func GetInitializedDB() *gorm.DB {
	db := test.NewOrm(&entity.Asset{}, &entity.DLCData{}).GetDB()
	db.Create(&entity.Asset{AssetID: "test"})
	return db
}

func Test_CreateDLCData_NotPresent_ReturnsNoError(t *testing.T) {
	db := GetInitializedDB()
	_, err := entity.CreateDLCData(db, "test", time.Now(), "kvalue", "rvalue", false)
	assert.NoError(t, err)
}

func Test_CreateDLCData_Present_ReturnsError(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now()
	db.Create(&entity.DLCData{AssetID: "test", PublishedDate: now, Kvalue: "kvalue1", Rvalue: "kvalue2"})
	_, err := entity.CreateDLCData(db, "test", now, "kvalue2", "rvalue2", false)
	assert.Error(t, err)
}

func Test_FindDLCDataPublishedNear_NotPresent_ReturnsRecordNotFoundError(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now()
	expected := &entity.DLCData{AssetID: "test", PublishedDate: now.Add(time.Hour), Kvalue: "", Rvalue: ""}
	db.Create(expected)
	_, err := entity.FindDLCDataPublishedNear(db, expected.AssetID, now, 30*time.Minute)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func Test_FindDLCDataPublishedNear_Present_ReturnsCorrectValue(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now()
	expected := &entity.DLCData{AssetID: "test", PublishedDate: now.Add(time.Hour), Kvalue: "", Rvalue: ""}
	db.Create(expected)
	actual, err := entity.FindDLCDataPublishedNear(db, expected.AssetID, now, 2*time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, expected.AssetID, actual.AssetID)
	assert.True(t, expected.PublishedDate.Equal(actual.PublishedDate))
}

func Test_FindDLCDataPublishedBefore_NotPresent_ReturnsRecordNotFoundError(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now()
	expected := &entity.DLCData{AssetID: "test", PublishedDate: now.Add(time.Hour), Kvalue: "", Rvalue: ""}
	db.Create(expected)
	_, err := entity.FindDLCDataPublishedBefore(db, expected.AssetID, now)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func Test_FindDLCDataPublishedBefore_Present_ReturnsCorrectValue(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now()
	expected := &entity.DLCData{AssetID: "test", PublishedDate: now.Add(-1 * time.Hour), Kvalue: "", Rvalue: ""}
	db.Create(expected)
	actual, err := entity.FindDLCDataPublishedBefore(db, expected.AssetID, now)
	assert.NoError(t, err)
	assert.Equal(t, expected.AssetID, actual.AssetID)
	assert.True(t, expected.PublishedDate.Equal(actual.PublishedDate))
}

func Test_FindDLCDataPublishedAt_NotPresent_ReturnsRecordNotFoundError(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now()
	_, err := entity.FindDLCDataPublishedAt(db, "test", now)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func Test_FindDLCDataPublishedAt_Present_ReturnsCorrectValue(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now()
	expected := &entity.DLCData{AssetID: "test", PublishedDate: now, Kvalue: "", Rvalue: ""}
	db.Create(expected)
	actual, err := entity.FindDLCDataPublishedAt(db, expected.AssetID, now)
	assert.NoError(t, err)
	assert.Equal(t, expected.AssetID, actual.AssetID)
	assert.True(t, expected.PublishedDate.Equal(actual.PublishedDate))
}

func Test_UpdateDLCDataSignatureAndValue_Present_ReturnsNoError(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now()
	expected := &entity.DLCData{AssetID: "test", PublishedDate: now, Kvalue: "", Rvalue: ""}
	db.Create(expected)
	err := entity.UpdateDLCDataSignatureAndValue(db, expected.AssetID, now, "new Signature", "new value")
	assert.NoError(t, err)
}
