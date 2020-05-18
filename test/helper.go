package test

import (
	"io/ioutil"
	conf "p2pderivatives-oracle/internal/configuration"
	"p2pderivatives-oracle/internal/database/orm"
	"p2pderivatives-oracle/internal/log"
	"path/filepath"
	"runtime"
)

var (
	// DirectoryPath directory path containing test helpers and vectors
	DirectoryPath string
	// VectorsDirectoryPath directory path containing test vectors
	VectorsDirectoryPath string
	// ConfigDirectoryPath directory path containing configuration files
	ConfigDirectoryPath string

	configuration *conf.Configuration
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information.")
	}
	absTestDirPath, _ := filepath.Abs(filepath.Dir(filename))
	DirectoryPath = absTestDirPath
	VectorsDirectoryPath = filepath.Join(DirectoryPath, "vectors")
	ConfigDirectoryPath = filepath.Join(DirectoryPath, "config")
	envName := "unittest"
	configuration = conf.NewConfiguration("unittest", envName, []string{ConfigDirectoryPath})
	err := configuration.Initialize()
	if err != nil {
		panic("Failed to initialize configuration")
	}
}

// InitializeSubConfig initializes a sub component configuration using initialized configuration.
func InitializeSubConfig(prefixKey string, componentConfig interface{}) error {
	return configuration.Sub(prefixKey).InitializeComponentConfig(componentConfig)
}

// InitializeConfig initializes a component configuration using initialized configuration.
func InitializeConfig(componentConfig interface{}) error {
	return configuration.InitializeComponentConfig(componentConfig)
}

// NewLogger returns a test logger initialized (the logs will be discarded)
func NewLogger() *log.Log {
	logConfig := &log.Config{}
	InitializeConfig(logConfig)
	logger := log.NewLog(logConfig)
	err := logger.Initialize()
	logger.Logger.SetOutput(ioutil.Discard)
	if err != nil {
		panic("Could not initialize log.")
	}
	return logger
}

// NewOrm returns a test orm initialized migrated models.
func NewOrm(models ...interface{}) *orm.ORM {
	logger := NewLogger()
	ormConfig := &orm.Config{}
	InitializeConfig(ormConfig)
	ormInstance := orm.NewORM(ormConfig, logger)
	ormInstance.Initialize()
	orm.NewMigrator(ormInstance, models...).Initialize()
	return ormInstance
}
