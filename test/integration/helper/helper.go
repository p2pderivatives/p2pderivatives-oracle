package helper

import (
	"flag"
	"fmt"
	"log"
	"os"
	"p2pderivatives-oracle/internal/api"
	conf "p2pderivatives-oracle/internal/configuration"

	"github.com/go-resty/resty/v2"
)

// Config contains the configuration parameters for the server.
type Config struct {
	TLS bool `configkey:"server.tls" default:"false"`
}

var (
	// OracleBaseURL The base url of the oracle server (without protocol)
	OracleBaseURL = flag.String("oracle-base-url", "localhost:8080", "The base url of the oracle server (without protocol)")
	// OracleConfig The oracle configuration used by the oracle and that will be used for integration testing
	OracleConfig conf.Configuration
	// ServerConfig Configuration of the http server
	ServerConfig *Config
	// APIConfig Configuration of the rest api
	APIConfig *api.Config
)

var (
	configPath = flag.String("abs-config", "", "Path to the configuration file to use.")
	appName    = flag.String("appname", "", "The name of the application. Will be use as a prefix for environment variables.")
	envname    = flag.String("e", "", "environment (ex., \"development\"). Should match with the name of the configuration file.")
)

// InitHelper initializes integration test and configuration
func InitHelper() {
	flag.Parse()

	if *configPath == "" {
		log.Fatal("No configuration path specified")
	}

	if *appName == "" {
		log.Fatal("No configuration name specified")
	}

	if *envname != "" {
		os.Setenv("P2PD_ENV", *envname)
	}

	OracleConfig = *conf.NewConfiguration(*appName, *envname, []string{*configPath})
	err := OracleConfig.Initialize()

	if err != nil {
		log.Fatalf("Could not read configuration %v.", err)
	}

	ServerConfig = &Config{}
	err = OracleConfig.InitializeComponentConfig(ServerConfig)
	if err != nil {
		log.Fatalf("Could not read Server configuration %v.", err)
	}

	APIConfig = &api.Config{}
	err = OracleConfig.InitializeComponentConfig(APIConfig)
	if err != nil {
		log.Fatalf("Could not read API configuration %v.", err)
	}
}

// CreateDefaultClient returns a default Rusty rest client
// (http/https scheme will be based on configuration file if TLS enabled)
func CreateDefaultClient() *resty.Client {
	client := resty.New()
	scheme := "http"
	if ServerConfig.TLS {
		scheme = "https"
	}
	client.SetHostURL(fmt.Sprintf("%s://%s", scheme, *OracleBaseURL))
	client.SetScheme(scheme)
	client.SetHeader("Accept", "application/json")
	return client
}
