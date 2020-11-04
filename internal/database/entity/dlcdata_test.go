package entity_test

import (
	"p2pderivatives-oracle/internal/database/entity"
	"p2pderivatives-oracle/test"
	"testing"
	"time"

	"gorm.io/gorm"
	"github.com/stretchr/testify/assert"
)

func GetInitializedDB() *gorm.DB {
	db := test.NewOrm(&entity.Asset{}, &entity.DLCData{}).GetDB()
	db.Create(&entity.Asset{AssetID: "test"})
	return db
}

func Test_CreateDLCData_NotPresent_ReturnsCorrectValue(t *testing.T) {
	// arrange
	db := GetInitializedDB()
	expected := &entity.DLCData{
		PublishedDate: time.Now().UTC(),
		AssetID:       "test",
		Rvalue:        "rvalue",
		Kvalue:        "kvalue",
	}

	// act
	actual, err := entity.CreateDLCData(
		db,
		expected.AssetID,
		expected.PublishedDate,
		expected.Kvalue,
		expected.Rvalue)

	// assert
	assertSub := assert.New(t)
	assertSub.NoError(err)
	assertDLCDataEqual(assertSub, expected, actual)
}

func Test_CreateDLCData_Present_ReturnsError(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now().UTC()
	inDB := &entity.DLCData{AssetID: "test", PublishedDate: now, Kvalue: "kvalue1", Rvalue: "rvalue2"}
	db.Create(inDB)
	_, err := entity.CreateDLCData(db, inDB.AssetID, inDB.PublishedDate, inDB.Kvalue, inDB.Rvalue)
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

func Test_UpdateDLCDataSignatureAndValue_Present_ReturnsUpdated(t *testing.T) {
	// arrange
	db := GetInitializedDB()
	now := time.Now().UTC()
	expected := &entity.DLCData{
		AssetID:       "test",
		PublishedDate: now,
		Kvalue:        "",
		Rvalue:        "",
		Signature:     "new Signature",
		Value:         "new value"}
	db.Create(expected)

	// act
	actual, err := entity.UpdateDLCDataSignatureAndValue(
		db,
		expected.AssetID,
		expected.PublishedDate,
		expected.Signature, expected.Value)

	// assert
	assertSub := assert.New(t)
	assertSub.NoError(err)
	assertDLCDataEqual(assertSub, expected, actual)
}

func assertDLCDataEqual(assertSub *assert.Assertions, expected *entity.DLCData, actual *entity.DLCData) {
	assertSub.Equal(expected.AssetID, actual.AssetID)
	assertSub.Equal(expected.PublishedDate, actual.PublishedDate)
	assertSub.Equal(expected.Kvalue, actual.Kvalue)
	assertSub.Equal(expected.Rvalue, actual.Rvalue)
	assertSub.Equal(expected.Signature, actual.Signature)
	assertSub.Equal(expected.Value, actual.Value)
}
