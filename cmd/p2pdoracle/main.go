package main

import (
	"context"
	"crypto/tls"
	"flag"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"p2pderivatives-oracle/internal/api"
	conf "p2pderivatives-oracle/internal/configuration"
	"p2pderivatives-oracle/internal/database/orm"
	"p2pderivatives-oracle/internal/log"
	"p2pderivatives-oracle/internal/router"
	"syscall"
	"time"
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
	ormInstance := newInitializedOrm(config, logInstance)
	if *migrate {
		err := doMigration(ormInstance)

		if err != nil {
			log.Fatalf("Failed to apply migration %v", err)
		}
	}
	routerInstance := newInitializedRouter(config, ormInstance, logInstance)

	serverConfig := &Config{}
	config.InitializeComponentConfig(serverConfig)

	srv := &http.Server{
		Addr:    serverConfig.Address,
		Handler: routerInstance.GetEngine(),
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
		cer, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("Unable to load certificate %v", err)
			return
		}
		srv.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cer}}
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
	config.InitializeComponentConfig(ormConfig)
	ormInstance := orm.NewORM(ormConfig, log)
	err := ormInstance.Initialize()

	if err != nil {
		panic("Could not initialize database.")
	}

	return ormInstance
}

func newInitializedRouter(config *conf.Configuration, orm *orm.ORM, log *log.Log) *router.Router {
	apiConfig := &api.Config{}
	config.InitializeComponentConfig(apiConfig)
	routerInstance := router.NewRouter(orm, log)
	err := routerInstance.Initialize()

	if err != nil {
		panic("Could not initialize router.")
	}

	return routerInstance
}

func doMigration(o *orm.ORM) error {
	// TODO add future models
	return nil
}
