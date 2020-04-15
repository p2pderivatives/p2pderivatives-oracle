package orm

import (
	"fmt"
	"time"

	"reflect"
	"unicode"

	"p2pderivatives-oracle/internal/log"

	"github.com/jinzhu/gorm"

	// Drivers for postgres and sql
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// NoLimit is used when searching without an upper limit on the number of
// returned records.
const NoLimit = -1

// ORM represent an Object Relational Mapper instance.
type ORM struct {
	config        *Config
	log           *log.Log
	connectionStr string
	enableLog     bool
	logger        *logrus.Logger
	initialized   bool
	db            *gorm.DB
}

// NewORM creates a new ORM structure with the given parameters.
func NewORM(config *Config, l *log.Log) *ORM {
	return &ORM{
		config:      config,
		log:         l,
		initialized: false,
	}
}

// Initialize initializes the ORM structure.
func (o *ORM) Initialize() error {

	if o.initialized {
		return nil
	}

	o.log.Logger.Info("ORM initialization starts")
	defer o.log.Logger.Info("ORM initialization end")

	enableLog := o.config.EnableLogging

	o.enableLog = enableLog
	o.logger = o.log.Logger

	var databaseType string

	if o.config.InMemory {
		databaseType = "sqlite3"
		o.connectionStr = ":memory:"
		o.log.Logger.Info("InMemory flag detected : Using Sqlite Inmemory DB")
	} else {
		databaseType = "postgres"
		o.connectionStr = fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s %s",
			o.config.Host,
			o.config.Port,
			o.config.DbName,
			o.config.DbUser,
			o.config.DbPassword,
			o.config.ConnectionParams)
	}

	opened, err := gorm.Open(databaseType, o.connectionStr)
	if err != nil {
		o.log.Logger.Error(err, "Could not open database.")
		return errors.Wrap(err, "failed to open database")
	}

	opened.SetLogger(o.logger)
	opened.LogMode(o.enableLog)

	o.db = opened
	o.db.DB().SetConnMaxLifetime(time.Hour)
	o.initialized = true

	return nil
}

// IsInitialized returns whether the orm is initialized.
func (o *ORM) IsInitialized() bool {
	return o.initialized
}

// Finalize releases the resources held by the orm.
func (o *ORM) Finalize() error {
	err := o.db.Close()
	if err != nil {
		return errors.Errorf("failed to close database connection")
	}
	return nil
}

// GetDB returns the DB instance associated with the orm object. Panics if the
// object is not initialized.
func (o *ORM) GetDB() *gorm.DB {
	if !o.IsInitialized() {
		panic("Trying to access uninitialized ORM object.")
	}

	return o.db
}

// GetColumnNames returns the name of the columns for the given model.
func GetColumnNames(model interface{}) []string {
	result := make([]string, 0)
	t := reflect.TypeOf(model)
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		first := rune(name[0])
		if unicode.IsUpper(first) {
			result = append(result, gorm.TheNamingStrategy.ColumnName(name))
		}
	}
	return result
}

// GetTableName returns the name of the table for the given model.
// Assumes that the globalDB is initialized, returns the default table name
// otherwise.
func (o *ORM) GetTableName(model interface{}) string {
	if o.initialized {
		return o.db.NewScope(model).GetModelStruct().TableName(o.db)
	}

	v := reflect.ValueOf(model)
	t := reflect.TypeOf(model)

	for v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
		t = v.Type()
	}
	return gorm.ToTableName(t.Name())
}

// IsRecordNotFoundError returns whether the given error is due to a requested
// record not present in the DB.
func IsRecordNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}

// NewRecordNotFoundError returns a ErrRecordNotFoundError.
func NewRecordNotFoundError() error {
	return gorm.ErrRecordNotFound
}
