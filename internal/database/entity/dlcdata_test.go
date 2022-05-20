package entity_test

import (
	"p2pderivatives-oracle/internal/database/entity"
	"p2pderivatives-oracle/test"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func GetInitializedDB() *gorm.DB {
	db := test.NewOrm(&entity.Asset{}, &entity.EventData{}).GetDB()
	db.Create(&entity.Asset{AssetID: "test"})
	return db
}

func Test_CreateDLCData_NotPresent_ReturnsCorrectValue(t *testing.T) {
	// arrange
	db := GetInitializedDB()
	expected := &entity.EventData{
		PublishedDate: time.Now().UTC(),
		AssetID:       "test",
		Nonces:        []string{"rvalue"},
		Kvalues:       []string{"kvalue"},
	}

	// act
	actual, err := entity.CreateEventData(
		db,
		expected.AssetID,
		expected.PublishedDate,
		expected.Kvalues,
		expected.Nonces,
		2,
		"e7d5da6e6193a8161437a860d41efe8af7c4c9073a1e75913e663ad59c092b0e0263942a600984f3352de5d089e4769b9448f63f279559408d3e3b089ddbdbc0",
	)

	// assert
	assertSub := assert.New(t)
	assertSub.NoError(err)
	assertDLCDataEqual(assertSub, expected, actual)
}

func Test_CreateDLCData_Present_ReturnsError(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now().UTC()
	inDB := &entity.EventData{AssetID: "test", PublishedDate: now, Kvalues: []string{"kvalue1"}, Nonces: []string{"rvalue2"}}
	db.Create(inDB)
	_, err := entity.CreateEventData(db, inDB.AssetID, inDB.PublishedDate, inDB.Kvalues, inDB.Nonces, 2, "e7d5da6e6193a8161437a860d41efe8af7c4c9073a1e75913e663ad59c092b0e0263942a600984f3352de5d089e4769b9448f63f279559408d3e3b089ddbdbc0")
	assert.Error(t, err)
}

func Test_FindDLCDataPublishedNear_NotPresent_ReturnsRecordNotFoundError(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now()
	expected := &entity.EventData{AssetID: "test", PublishedDate: now.Add(time.Hour), Kvalues: []string{""}, Nonces: []string{""}}
	db.Create(expected)
	_, err := entity.FindDLCDataPublishedNear(db, expected.AssetID, now, 30*time.Minute)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func Test_FindDLCDataPublishedNear_Present_ReturnsCorrectValue(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now()
	expected := &entity.EventData{AssetID: "test", PublishedDate: now.Add(time.Hour), Kvalues: []string{""}, Nonces: []string{""}}
	db.Create(expected)
	actual, err := entity.FindDLCDataPublishedNear(db, expected.AssetID, now, 2*time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, expected.AssetID, actual.AssetID)
	assert.True(t, expected.PublishedDate.Equal(actual.PublishedDate))
}

func Test_FindDLCDataPublishedBefore_NotPresent_ReturnsRecordNotFoundError(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now()
	expected := &entity.EventData{AssetID: "test", PublishedDate: now.Add(time.Hour), Kvalues: []string{""}, Nonces: []string{""}}
	db.Create(expected)
	_, err := entity.FindDLCDataPublishedBefore(db, expected.AssetID, now)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func Test_FindDLCDataPublishedBefore_Present_ReturnsCorrectValue(t *testing.T) {
	db := GetInitializedDB()
	now := time.Now()
	expected := &entity.EventData{AssetID: "test", PublishedDate: now.Add(-1 * time.Hour), Kvalues: []string{""}, Nonces: []string{""}}
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
	expected := &entity.EventData{AssetID: "test", PublishedDate: now, Kvalues: []string{""}, Nonces: []string{""}}
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
	expected := &entity.EventData{
		AssetID:       "test",
		PublishedDate: now,
		Kvalues:       []string{""},
		Nonces:        []string{""},
		Signatures:    nil,
		Values:        nil,
	}
	db.Create(expected)

	// act
	actual, err := entity.UpdateDLCDataSignatureAndValue(
		db,
		expected.AssetID,
		expected.PublishedDate,
		expected.Signatures, expected.Values)

	// assert
	assertSub := assert.New(t)
	assertSub.NoError(err)
	assertDLCDataEqual(assertSub, expected, actual)
}

func assertDLCDataEqual(assertSub *assert.Assertions, expected *entity.EventData, actual *entity.EventData) {
	assertSub.Equal(expected.AssetID, actual.AssetID)
	assertSub.Equal(expected.PublishedDate, actual.PublishedDate)
	assertSub.Equal(expected.Kvalues, actual.Kvalues)
	assertSub.Equal(expected.Nonces, actual.Nonces)
	assertSub.Equal(expected.Signatures, actual.Signatures)
	assertSub.Equal(expected.Values, actual.Values)
}
