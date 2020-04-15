package test

import (
	"fmt"
	conf "p2pderivatives-oracle/src/configuration"
	"p2pderivatives-oracle/src/database/orm"
	"p2pderivatives-oracle/src/log"
	"path/filepath"
	"runtime"
)

// GetTestConfig returns a configuration for unit tests.
func GetTestConfig() *conf.Configuration {
	envName := "unittest"
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information.")
	}
	fmt.Printf("called by %s", filename)
	confPath, _ := filepath.Abs(filepath.Join(filepath.Dir(filename), "config"))
	config := conf.NewConfiguration("unittest", envName, []string{confPath})
	err := config.Initialize()
	if err != nil {
		panic("Failed to initialize configuration")
	}

	return config
}

// GetTestLogger returns a logger for unit tests.
func GetTestLogger(config *conf.Configuration) *log.Log {
	logConfig := log.Config{}
	config.InitializeComponentConfig(&logConfig)
	log := log.NewLog(&logConfig)
	err := log.Initialize()
	if err != nil {
		panic("Could not initialize log.")
	}

	return log
}

// InitializeORM initializes the global db for unit tests.
func InitializeORM(models ...interface{}) *orm.ORM {
	config := GetTestConfig()
	logger := GetTestLogger(config)
	ormConfig := orm.Config{}
	config.InitializeComponentConfig(&ormConfig)
	ormInstance := orm.NewORM(&ormConfig, logger)
	ormInstance.Initialize()
	orm.NewMigrator(ormInstance, models...).Initialize()
	return ormInstance
}
