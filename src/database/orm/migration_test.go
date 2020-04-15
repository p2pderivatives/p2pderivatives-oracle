package orm_test

import (
	"p2pderivatives-oracle/src/database/orm"
	"p2pderivatives-oracle/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestMigrationModel struct {
	Name string
}

func TestMigrationInitialize_IsInitialized(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	config := test.GetTestConfig()
	ormConfig := orm.Config{}
	config.InitializeComponentConfig(&ormConfig)
	ormInstance := orm.NewORM(&ormConfig, test.GetTestLogger(config))
	ormInstance.Initialize()
	migrator := orm.NewMigrator(ormInstance, &TestMigrationModel{})

	// Act
	migrator.Initialize()

	// Assert
	assert.True(migrator.IsInitialized())
}

func TestMigrationInitialize_HasCorrectTable(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	config := test.GetTestConfig()
	ormConfig := orm.Config{}
	config.InitializeComponentConfig(&ormConfig)
	ormInstance := orm.NewORM(&ormConfig, test.GetTestLogger(config))
	ormInstance.Initialize()
	migrator := orm.NewMigrator(ormInstance, &TestMigrationModel{})
	migrator.Initialize()

	// Act
	var result []TestMigrationModel
	err := ormInstance.GetDB().Find(&result).Error

	// Assert
	assert.NoError(err)
}
