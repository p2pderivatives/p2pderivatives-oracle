// +build integration

package helper

import (
	"flag"
	"fmt"
	"log"
	"os"
	"p2pderivatives-oracle/internal/api"
	conf "p2pderivatives-oracle/internal/configuration"
	"p2pderivatives-oracle/internal/oracle"
	"path/filepath"
	"runtime"

	"github.com/go-resty/resty/v2"
)

// HttpConfig contains the configuration parameters for the server.
type HttpConfig struct {
	TLS bool `configkey:"server.tls" default:"false"`
}

var (
	// IntegrationDir path to the integration test directory
	IntegrationDir string
	// OracleBaseURL The base url of the oracle server (without protocol)
	OracleBaseURL = flag.String("oracle-base-url", "localhost:8080", "The base url of the oracle server (without protocol)")
	// Config The oracle configuration used by the oracle and that will be used for integration testing
	Config conf.Configuration
	// ServConfig Configuration of the http server
	ServConfig *HttpConfig
	// APIConfig Configuration of the rest api
	APIConfig *api.Config

	// ExpectedOracle represents the oracle to test against
	ExpectedOracle *oracle.Oracle
)

var (
	configDir  = flag.String("config-dir", "../config", "Path to the directory of the config file.")
	appName    = flag.String("appname", "p2pdoracle", "The name of the application. Will be use as a prefix for environment variables")
	configName = flag.String("config-file-name", "integration", "environment (ex., \"development\"). Should match with the name of the configuration file.")
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information.")
	}
	absIntegrationDirPath, _ := filepath.Abs(filepath.Dir(filename))
	IntegrationDir = absIntegrationDirPath
}

// InitHelper initializes integration test and configuration
func InitHelper() {
	flag.Parse()

	if *configDir == "" {
		log.Fatal("No configuration path specified")
	}

	if *appName == "" {
		log.Fatal("No configuration name specified")
	}

	if *configName != "" {
		os.Setenv("P2PD_ENV", *configName)
	}
	path := filepath.Join(IntegrationDir, *configDir)
	Config = *conf.NewConfiguration(*appName, *configName, []string{path})
	err := Config.Initialize()

	if err != nil {
		log.Fatalf("Could not read configuration with %v", path)
	}

	ServConfig = &HttpConfig{}
	err = Config.InitializeComponentConfig(ServConfig)
	if err != nil {
		log.Fatal("Could not read Server configuration.")
	}

	APIConfig = &api.Config{}
	err = Config.InitializeComponentConfig(APIConfig)
	if err != nil {
		log.Fatal("Could not read API configuration.")
	}

	oracleConfig := &oracle.Config{}
	err = Config.InitializeComponentConfig(oracleConfig)
	if err != nil {
		log.Fatal("Could not read Oracle configuration.")
	}

	// fix config by retrieving absolute filepath
	oracleConfig.KeyFile = filepath.Join(IntegrationDir, oracleConfig.KeyFile)
	if oracleConfig.KeyPassFile != "" {
		oracleConfig.KeyPassFile = filepath.Join(IntegrationDir, oracleConfig.KeyPassFile)
	}

	ExpectedOracle, err = oracle.FromConfig(oracleConfig)
	if err != nil {
		log.Fatalf("Could not instantiate Oracle from config. %v", err)
	}
}

// CreateDefaultClient returns a default Rusty rest client
// (http/https scheme will be based on configuration file if TLS enabled)
func CreateDefaultClient() *resty.Client {
	client := resty.New()
	scheme := "http"
	if ServConfig.TLS {
		scheme = "https"
	}
	client.SetHostURL(fmt.Sprintf("%s://%s", scheme, *OracleBaseURL))
	client.SetScheme(scheme)
	client.SetHeader("Accept", "application/json")
	return client
}
