package orm_test

import (
	"p2pderivatives-oracle/internal/database/orm"
	"p2pderivatives-oracle/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	Name string
}

func TestOrm_GetColumnNames(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	type TestStruct struct {
		Hoge         string
		FugaFuga     string
		PiyoPiyoPiyo string
		unexposed    string
		//Snake_Snake  string	// Don't use underscores in Go names
	}
	expected := []string{"hoge", "fuga_fuga", "piyo_piyo_piyo"}

	// Act
	names := orm.GetColumnNames(TestStruct{unexposed: "unexposed"})

	// Assert
	assert.Equal(expected, names)
}

func TestOrmGetTableName_Initialized_ReturnsCorrectName(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	ormInstance := test.NewOrm(&TestModel{})

	// Act
	name := ormInstance.GetTableName(&TestModel{})

	// Assert
	assert.Equal("test_models", name)
}

func TestOrmGetTableName_Uninitialized_ReturnsCorrectName(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	ormInstance := orm.ORM{}

	// Act
	name := ormInstance.GetTableName(&TestModel{})

	// Assert
	assert.Equal("test_model", name)
}

func TestOrmInitializeFinalize_NoError(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	ormConfig := &orm.Config{}
	test.InitializeConfig(ormConfig)
	l := test.NewLogger()
	ormInstance := orm.NewORM(ormConfig, l)

	// Act
	err := ormInstance.Initialize()
	err2 := ormInstance.Finalize()

	// Assert
	assert.NoError(err)
	assert.NoError(err2)
}

func TestOrmInitialize_IsInitialized_True(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	ormConfig := &orm.Config{}
	test.InitializeConfig(ormConfig)
	l := test.NewLogger()
	ormInstance := orm.NewORM(ormConfig, l)
	ormInstance.Initialize()
	defer ormInstance.Finalize()

	// Assert
	assert.True(ormInstance.IsInitialized())
}

func TestOrmGetDB_Initialized_Succeeds(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	orm := test.NewOrm()

	// Act
	db := orm.GetDB()

	// Assert
	assert.NotNil(db)
}

func TestOrmGetDB_NotInitialized_Panics(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	orm := orm.ORM{}

	// Act
	act := func() { orm.GetDB() }

	// Assert
	assert.Panics(act)
}
