package orm

// Config contains the configuration parameter to set up the orm.
type Config struct {
	EnableLogging    bool   `configkey:"database.log"` // Whether to enable logging of the database
	InMemory         bool   `configkey:"database.inmemory" default:"false"`
	Host             string `configkey:"database.host" validate:"required"`
	Port             string `configkey:"database.port" validate:"required"`
	DbName           string `configkey:"database.dbname" default:"postgres"`
	DbUser           string `configkey:"database.dbuser" default:"postgres"`
	DbPassword       string `configkey:"database.dbpassword" validate:"required"`
	ConnectionParams string `configkey:"database.connectionParams"` // Postgres sql connection parameters separate by space
}
