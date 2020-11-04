package main

import (
	"context"
	"flag"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"p2pderivatives-oracle/internal/api"
	"p2pderivatives-oracle/internal/cryptocompare"
	"p2pderivatives-oracle/internal/database/entity"
	"p2pderivatives-oracle/internal/datafeed"
	"p2pderivatives-oracle/internal/dlccrypto"
	"p2pderivatives-oracle/internal/oracle"
	"syscall"
	"time"

	conf "github.com/cryptogarageinc/server-common-go/pkg/configuration"
	"github.com/cryptogarageinc/server-common-go/pkg/database/orm"
	"github.com/cryptogarageinc/server-common-go/pkg/log"
	"github.com/cryptogarageinc/server-common-go/pkg/rest/router"
)

var (
	configPath = flag.String("config", "", "Path to the configuration file to use.")
	appName    = flag.String("appname", "", "The name of the application. Will be use as a prefix for environment variables.")
	envname    = flag.String("e", "", "environment (ex., \"development\"). Should match with the name of the configuration file.")
	migrate    = flag.Bool("migrate", false, "If set performs a db migration before starting.")
)

// Config contains the configuration parameters for the server.
type Config struct {
	Address  string `configkey:"server.address" validate:"required"`
	TLS      bool   `configkey:"server.tls"`
	CertFile string `configkey:"server.certfile" validate:"required_with=TLS"`
	KeyFile  string `configkey:"server.keyfile" validate:"required_with=TLS"`
}

func init() {
	flag.Parse()

	if *configPath == "" {
		stdlog.Fatal("No configuration path specified")
	}

	if *appName == "" {
		stdlog.Fatal("No configuration name specified")
	}

	if *envname != "" {
		os.Setenv("P2PD_ENV", *envname)
	}
}

func main() {

	config := conf.NewConfiguration(*appName, *envname, []string{*configPath})
	err := config.Initialize()

	if err != nil {
		stdlog.Fatalf("Could not read configuration %v.", err)
	}

	logInstance := newInitializedLog(config)
	log := logInstance.Logger

	// Initialize Router
	routerInstance := newInitializedRouter(logInstance, config)

	serverConfig := &Config{}
	config.InitializeComponentConfig(serverConfig)

	srv := &http.Server{
		Addr:    serverConfig.Address,
		Handler: routerInstance.GetEngine(),
	}

	listenAndServe := func() error {
		return srv.ListenAndServe()
	}
	if serverConfig.TLS {
		certFile := serverConfig.CertFile
		keyFile := serverConfig.KeyFile
		if certFile == "" {
			log.Fatal("Need to provide the path to the certificate file")
		}
		if keyFile == "" {
			log.Fatal("Need to provide the path to the key file")
		}

		listenAndServe = func() error {
			return srv.ListenAndServeTLS(certFile, keyFile)
		}
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := listenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failing to listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	routerInstance.Finalize()
	log.Println("Server exiting")
	logInstance.Finalize()
}

func newInitializedLog(config *conf.Configuration) *log.Log {
	logConfig := &log.Config{}
	config.InitializeComponentConfig(logConfig)
	logger := log.NewLog(logConfig)
	logger.Initialize()
	return logger
}

func newInitializedOrm(config *conf.Configuration, log *log.Log) *orm.ORM {
	ormConfig := &orm.Config{}
	if err := config.InitializeComponentConfig(ormConfig); err != nil {
		panic(err)
	}
	ormInstance := orm.NewORM(ormConfig, log)
	err := ormInstance.Initialize()

	if err != nil {
		panic("Could not initialize database.")
	}

	if *migrate {
		if err := doMigration(ormInstance); err != nil {
			log.Logger.Fatalf("Could not apply migrations")
			panic(err)
		}
	}

	return ormInstance
}

func newInitializedRouter(log *log.Log, config *conf.Configuration) *router.Router {
	api := NewDefaultOracleAPI(log, config)
	routerInstance := router.NewRouter(log, api)
	err := routerInstance.Initialize()

	if err != nil {
		panic("Could not initialize router.")
	}

	return routerInstance
}

// NewDefaultOracleAPI returns a router.API with default crypto, database and datafeed services
func NewDefaultOracleAPI(l *log.Log, config *conf.Configuration) router.API {
	// Setup crypto service
	cryptoInstance := dlccrypto.NewCfdgoCryptoService()

	// Setup Oracle
	oracleConfig := &oracle.Config{}
	config.InitializeComponentConfig(oracleConfig)
	oracleInstance, err := oracle.FromConfig(oracleConfig)
	if err != nil {
		l.Logger.Fatalf("Could not create a oracle instance")
		panic(err)
	}

	// Setup orm service
	ormInstance := newInitializedOrm(config, l)

	// Setup DataFeed service
	var feedInstance datafeed.DataFeed
	datafeedConfig := config.Sub("datafeed")
	dummyFeedConfig := &datafeed.DummyConfig{}
	err = datafeedConfig.InitializeComponentConfig(dummyFeedConfig)
	if err != nil {
		ccFeedConfig := &cryptocompare.Config{}
		datafeedConfig.InitializeComponentConfig(ccFeedConfig)
		cryptoCompareClient := cryptocompare.NewClient(ccFeedConfig)
		cryptoCompareClient.Initialize()
		feedInstance = cryptoCompareClient
	} else {
		feedInstance = datafeed.NewDummyDataFeed(dummyFeedConfig)
	}

	apiConfig := &api.Config{}
	err = config.InitializeComponentConfig(apiConfig)
	if err != nil {
		panic(err)
	}
	return api.NewOracleAPI(apiConfig, l, oracleInstance, ormInstance, cryptoInstance, feedInstance)
}

func doMigration(o *orm.ORM) error {
	db := o.GetDB()
	err := db.AutoMigrate(&entity.Asset{}, &entity.DLCData{})
	err = db.Create(&entity.Asset{AssetID: "btcusd", Description: "BTC USD"}).Error
	err = db.Create(&entity.Asset{AssetID: "btcjpy", Description: "BTC JPY"}).Error
	return err
}
