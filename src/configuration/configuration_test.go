package conf

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncodingString_WithUTF8_ReturnsUTF8(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	enc := UTF8

	// Act
	s := enc.String()

	// Assert
	assert.Equal("utf-8", s)
}

func TestEncodingString_WithHEX_ReturnsHex(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	enc := HEX

	// Act
	s := enc.String()

	// Assert
	assert.Equal("hex", s)
}

func createConfiguration(t *testing.T) *Configuration {
	path := filepath.Join("..", "..", "test", "config")
	c := NewConfiguration("core", "unittest", []string{path})
	err := c.Initialize()
	if err != nil {
		t.Fatalf("reading configuration failed: %s", err)
		return nil
	}
	if !c.IsInitialized() {
		t.Fatal("must be initialized")
	}

	return c
}

func TestConfigurationGetInt_WithIntValue_ReturnsCorrectValue(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	expected := 10

	// Act
	actual := config.GetInt("unittest.i")

	// Assert
	assert.Equal(expected, actual)
}

func TestConfigurationGetString_WithStringValue_ReturnsCorrectValue(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	expected := "hoge"

	// Act
	actual := config.GetString("unittest.s")

	// Assert
	assert.Equal(expected, actual)
}

func TestConfigurationGetStringSlice_WithStringSliceValue_ReturnsCorrectValue(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	expected := []string{"hoge", "fuga"}

	// Act
	actual := config.GetStringSlice("unittest.ss")

	// Assert
	assert.Equal(expected, actual)
}

func TestConfigurationGetBool_WithBoolValue_ReturnsCorrectValue(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	expected := true

	// Act
	actual := config.GetBool("unittest.b")

	// Assert
	assert.Equal(expected, actual)
}

func TestConfigurationGetByte_WithByteValue_ReturnsCorrectValue(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	expected := []byte{0x61, 0x62, 0x63, 0x64, 0x65} // abcde

	// Act
	actual, err := config.GetByte("unittest.utf8byte", UTF8)

	// Assert
	assert.NoError(err)
	assert.Equal(expected, actual)
}

func TestConfigurationGetByte_WithByteValueUnknownEncoding_ReturnsError(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)

	// Act
	actual, err := config.GetByte("unittest.utf8byte", 3)

	// Assert
	assert.Error(err)
	assert.Nil(actual)
}

func TestConfigurationGetDuration_WithDurationValue_ReturnsCorrectValue(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	expected, _ := time.ParseDuration("1h10m10s")

	// Act
	actual := config.GetDuration("unittest.dr")

	// Assert
	assert.Equal(expected, actual)
}

func TestConfigurationGetDuration_WithUInt32_ReturnsCorrectValue(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	expected := uint32(32)

	// Act
	actual := config.GetUInt32("unittest.ui32")

	// Assert
	assert.Equal(expected, actual)
}

func TestConfigurationGetDuration_WithUInt8_ReturnsCorrectValue(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	expected := uint8(8)

	// Act
	actual := config.GetUInt8("unittest.ui8")

	// Assert
	assert.Equal(expected, actual)
}

func TestConfiguration_WithEnvironmentVariable_ReturnsEnvironmentVariableValue(t *testing.T) {
	// Arrange
	expected := 100
	os.Setenv("CORE_PORT", strconv.Itoa(expected))

	config := createConfiguration(t)
	assert := assert.New(t)
	// Act
	actual := config.GetInt("port")
	// Assert
	assert.Equal(expected, actual)
}

func TestNewConfiguration_FromReader_ReturnsCorrectValue(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	expected := "fuga"

	// Act
	config, err := NewConfigurationFromReader("properties", strings.NewReader("hoge=fuga"))
	actual := config.GetString("hoge")

	// Assert
	assert.NoError(err)
	assert.Equal(expected, actual)
}

type UnitTestConfig struct {
	I         int           `configkey:"unittest.i" validate:"min=10" default:"10"`
	S         string        `configkey:"unittest.s" validate:"required" default:"hoge"`
	Ss        []string      `configkey:"unittest.ss" validate:"dive,required" default:"hoge,fuga"`
	B         bool          `configkey:"unittest.b" default:"true"`
	Utf8byte  []byte        `configkey:"unittest.utf8byte,utf8" validate:"required" default:"abcde"`
	Utf8NoEnc []byte        `configkey:"unittest.utf8byte" validate:"required" default:"abcde"`
	Hexbytes  []byte        `configkey:"unittest.hexbyte,hex" validate:"required" default:"abcd0e"`
	Dr        time.Duration `configkey:"unittest.dr,duration" validate:"required" default:"1h10m10s"`
	I64       int64         `configkey:"unittest.i64" validate:"min=11" default:"132904"`
	UI8       uint8         `configkey:"unittest.ui8" validate:"min=8" default:"8"`
	UI32      uint32        `configkey:"unittest.ui32" validate:"min=32" default:"32"`
	UI64      uint64        `configkey:"unittest.ui64" validate:"min=64" default:"64"`
	F32       float32       `configkey:"unittest.f32" validate:"min=3.2" default:"3.2"`
	F64       float64       `configkey:"unittest.f64" validate:"min=6.4" default:"6.4"`
	Ignored   bool
}

type InvalidTestConfig struct {
	I int `configkey:"unittest.i" validate:"min=11"`
}

type InvalidEncodingTestConfig struct {
	Byte []byte `configkey:"unittest.utf8byte,wrong"`
}

type InvalidHexTestConfig struct {
	Byte []byte `configkey:"unittest.invalidhex,hex"`
}

type UnknownTypeConfig struct {
	UnknownType []int `configkey:"unittest.unkowntype"`
}

func TestConfiguration_InitializeComponentConfig_CorrectlyInitializesConfig(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	unitTestConfig := &UnitTestConfig{}
	expectedDuration, _ := time.ParseDuration("1h10m10s")

	// Act
	config.InitializeComponentConfig(unitTestConfig)

	// Assert
	assert.Equal(10, unitTestConfig.I)
	assert.Equal("hoge", unitTestConfig.S)
	assert.Equal([]string{"hoge", "fuga"}, unitTestConfig.Ss)
	assert.Equal(true, unitTestConfig.B)
	assert.Equal([]byte{0x61, 0x62, 0x63, 0x64, 0x65}, unitTestConfig.Utf8byte)
	assert.Equal([]byte{0x61, 0x62, 0x63, 0x64, 0x65}, unitTestConfig.Utf8NoEnc)
	assert.Equal([]byte{0xab, 0xcd, 0x0e}, unitTestConfig.Hexbytes)
	assert.Equal(expectedDuration, unitTestConfig.Dr)
	assert.Equal(int64(132904), unitTestConfig.I64)
	assert.Equal(uint8(8), unitTestConfig.UI8)
	assert.Equal(uint32(32), unitTestConfig.UI32)
	assert.Equal(uint64(64), unitTestConfig.UI64)
	assert.Equal(float32(3.2), unitTestConfig.F32)
	assert.Equal(float64(6.4), unitTestConfig.F64)
}

func TestConfiguration_InitializeValidComponentConfig_NoError(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	unitTestConfig := &UnitTestConfig{}

	// Act
	err := config.InitializeComponentConfig(unitTestConfig)

	// Assert
	assert.NoError(err)
}

func TestConfiguration_InitializeInvalidComponentConfig_Error(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	invalidTestConfig := &InvalidTestConfig{}

	// Act
	err := config.InitializeComponentConfig(invalidTestConfig)

	// Assert
	assert.Error(err)
}

func TestConfiguration_InitializeInvalidEncodingComponentConfig_Error(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	invalidEncodingTestConfig := &InvalidEncodingTestConfig{}

	// Act
	err := config.InitializeComponentConfig(invalidEncodingTestConfig)

	// Assert
	assert.Error(err)
}

func TestConfiguration_InitializeInvalidByteComponentConfig_Error(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	invalidHexTestConfig := &InvalidHexTestConfig{}

	// Act
	err := config.InitializeComponentConfig(invalidHexTestConfig)

	// Assert
	assert.Error(err)
}

func TestConfiguration_InitializeUnknownTypeComponentConfig_Error(t *testing.T) {
	// Arrange
	config := createConfiguration(t)
	assert := assert.New(t)
	unknownTypeConfig := &UnknownTypeConfig{}

	// Act
	err := config.InitializeComponentConfig(unknownTypeConfig)

	// Assert
	assert.Error(err)
}

func TestConfiguration_InitializeComponentConfigNotInitialized_Panics(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	var config Configuration
	unitTestConfig := &UnitTestConfig{}

	// Act
	act := func() { config.InitializeComponentConfig(unitTestConfig) }

	// Assert
	assert.Panics(act)
}

func TestConfiguration_InitializeComponentWithEmptyConfig_HasDefaultValues(t *testing.T) {
	// Arrange
	config, _ := NewConfigurationFromReader("yaml", strings.NewReader(""))
	config.Initialize()
	assert := assert.New(t)
	unitTestConfig := &UnitTestConfig{}
	expectedDuration, _ := time.ParseDuration("1h10m10s")

	// Act
	config.InitializeComponentConfig(unitTestConfig)

	// Assert
	assert.Equal(10, unitTestConfig.I)
	assert.Equal("hoge", unitTestConfig.S)
	assert.Equal([]string{"hoge", "fuga"}, unitTestConfig.Ss)
	assert.Equal(true, unitTestConfig.B)
	assert.Equal([]byte{0x61, 0x62, 0x63, 0x64, 0x65}, unitTestConfig.Utf8byte)
	assert.Equal([]byte{0xab, 0xcd, 0x0e}, unitTestConfig.Hexbytes)
	assert.Equal(expectedDuration, unitTestConfig.Dr)
	assert.Equal(int64(132904), unitTestConfig.I64)
}

func TestConfiguration_InitializeComponentConfigPassedByValue_Panics(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	config := createConfiguration(t)

	// Act
	act := func() { config.InitializeComponentConfig(UnitTestConfig{}) }

	// Assert
	assert.Panics(act)
}
