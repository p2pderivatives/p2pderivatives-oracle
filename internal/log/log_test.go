package log

import (
	"io/ioutil"
	conf "p2pderivatives-oracle/internal/configuration"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLog_WithInfoLevelConfig_HasInfoLevel(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	configProperties :=
		`log.format=json
		 log.output_stdout=true
		 log.level=info`

	assertCallback := func(
		logger *logrus.Logger, initError, readerError, logrusInitError error) {
		assert.NoError(initError)
		assert.NoError(readerError)
		assert.NoError(logrusInitError)
		assert.Equal(logrus.InfoLevel, logger.GetLevel())
	}

	// Act/Assert
	testLogHelper(configProperties, assertCallback, assert)
}

func TestLog_WithDebugLevelConfig_HasDebugLevel(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	configProperties :=
		`log.format=json
		 log.output_stdout=true
		 log.level=debug`

	assertCallback := func(
		logger *logrus.Logger, initError, readerError, logrusInitError error) {
		assert.NoError(initError)
		assert.NoError(readerError)
		assert.NoError(logrusInitError)
		assert.Equal(logrus.DebugLevel, logger.GetLevel())
	}

	// Act/Assert
	testLogHelper(configProperties, assertCallback, assert)
}

func TestLog_WithJsonFormatConfig_HasJsonFormat(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	configProperties :=
		`log.format=json
		 log.output_stdout=true
		 log.level=debug`

	assertCallback := func(
		logger *logrus.Logger, initError, readerError, logrusInitError error) {
		assert.NoError(initError)
		assert.NoError(readerError)
		assert.NoError(logrusInitError)
		_, ok := logger.Formatter.(*logrus.JSONFormatter)
		assert.True(ok)
	}

	// Act/Assert
	testLogHelper(configProperties, assertCallback, assert)
}

func TestLog_WithTomlFormatConfig_HasInitError(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	configProperties :=
		`log.format=toml
		 log.output_stdout=true
		 log.level=debug`

	assertCallback := func(
		logger *logrus.Logger, initError, readerError, logrusInitError error) {
		assert.Error(initError)
		assert.NoError(readerError)
		assert.Error(logrusInitError)
	}

	// Act/Assert
	testLogHelper(configProperties, assertCallback, assert)
}

func TestLog_WithInvalidLogLevelConfig_HasInitError(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	configProperties :=
		`log.format=json
		 log.output_stdout=true
		 log.level=hoge`

	assertCallback := func(
		logger *logrus.Logger, initError, readerError, logrusInitError error) {
		assert.NoError(initError)
		assert.NoError(readerError)
		assert.Error(logrusInitError)
	}

	// Act/Assert
	testLogHelper(configProperties, assertCallback, assert)
}

func testLogHelper(
	configProperties string,
	assertCallback func(*logrus.Logger, error, error, error),
	assert *assert.Assertions) {

	// Act
	config, readerError := conf.NewConfigurationFromReader(
		"properties", strings.NewReader(configProperties))
	logConfig := Config{}
	initError := config.InitializeComponentConfig(&logConfig)
	log := NewLog(&logConfig)
	logger, logrusInitError := log.initializeLogrus(ioutil.Discard)
	defer log.Finalize()

	// Assert
	assertCallback(logger, initError, readerError, logrusInitError)
}

func TestLog_IsInitialized(t *testing.T) {

	// Arrange
	assert := assert.New(t)
	configProperties := `log.format=text
	log.output_stdout=true
	log.level=debug
	`
	// Act
	config, err1 := conf.NewConfigurationFromReader(
		"properties", strings.NewReader(configProperties))
	logConfig := Config{}
	err2 := config.InitializeComponentConfig(&logConfig)
	l := NewLog(&logConfig)
	err3 := l.Initialize()
	defer l.Finalize()

	// Assert
	assert.NoError(err1)
	assert.NoError(err2)
	assert.NoError(err3)
	assert.True(l.IsInitialized())
}

func TestLog_InitializeFinalizeWithRotateLog_NoError(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	configProperties := `log.format=text
	log.output_stdout=false
	log.level=debug
	log.RotateCount=7
	log.RotateInterval=24h
	log.BaseName=unittest.log.%Y-%m-%d
	log.Dir=_log
	`
	config, _ := conf.NewConfigurationFromReader(
		"properties", strings.NewReader(configProperties))
	logConfig := Config{}
	config.InitializeComponentConfig(&logConfig)

	// Act
	l := NewLog(&logConfig)
	err1 := l.Initialize()
	err2 := l.Finalize()

	// Assert
	assert.NoError(err1)
	assert.NoError(err2)
}
