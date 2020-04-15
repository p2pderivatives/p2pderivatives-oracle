package conf

import (
	"encoding/hex"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

// Configuration contains the application configuration parameters.
type Configuration struct {
	AppName         string
	EnvironmentName string
	paths           []string
	viper           *viper.Viper
	initialized     bool
}

// Encoding represent the encoding used.
type Encoding int

const (
	// UTF8 for UTF8 encoding
	UTF8 Encoding = iota
	// HEX for hexadecimal encoding
	HEX
)

const (
	configTagName  = "configkey"
	defaultTagName = "default"
	utf8TagValue   = "utf8"
	hexTagValue    = "hex"
)

var (
	stringSliceType = reflect.TypeOf(([]string)(nil))
	byteArrayType   = reflect.TypeOf(([]byte)(nil))
)

// String returns the name of the encoding
func (e Encoding) String() string {
	switch e {
	case UTF8:
		return "utf-8"
	case HEX:
		return "hex"
	}
	return "unknown"
}

// NewConfiguration returns a new configuration based on the given parameters.
func NewConfiguration(appName, envname string, searchPaths []string) *Configuration {
	return &Configuration{
		AppName:         appName,
		EnvironmentName: envname,

		paths:       searchPaths,
		viper:       viper.New(),
		initialized: false,
	}
}

// NewConfigurationFromReader creates a configuration using the content read
// from the given stream using the provided format.
// Supported formats are "json", "toml", "yaml", "yml", "properties", "props",
// "prop", "hcl".
func NewConfigurationFromReader(format string, in io.Reader) (*Configuration, error) {
	v := viper.New()
	v.SetConfigType(format)
	if err := v.ReadConfig(in); err != nil {
		return nil, errors.Wrapf(err, "failed to init Configuration with %s", in)
	}
	c := &Configuration{viper: v}
	c.initialized = true
	return c, nil
}

// Initialize initializes the configuration.
func (c *Configuration) Initialize() error {
	if strings.TrimSpace(c.EnvironmentName) == "" {
		return errors.Errorf("environment name is missing")
	}

	c.viper.SetConfigName(c.EnvironmentName)
	for _, path := range c.paths {
		c.viper.AddConfigPath(path)
	}

	err := c.viper.ReadInConfig()
	if err != nil {
		return errors.Wrapf(err, "failed to read config file: [%s] (suffix ommitted)", c.EnvironmentName)
	}

	c.viper.SetEnvPrefix(c.AppName)
	c.viper.AutomaticEnv()
	c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	c.initialized = true
	return nil
}

// IsInitialized returns whether the configuration is initialized.
func (c *Configuration) IsInitialized() bool {
	c.ensureInitialized()
	return c.initialized
}

// GetInt returns the values associated with the given key as an integer.
func (c *Configuration) GetInt(key string) int {
	c.ensureInitialized()
	return c.viper.GetInt(key)
}

// GetString returns the values associated with the given key as a string.
func (c *Configuration) GetString(key string) string {
	c.ensureInitialized()
	return c.viper.GetString(key)
}

// GetStringSlice returns the values associated with the given key as a string
// slice.
func (c *Configuration) GetStringSlice(key string) []string {
	c.ensureInitialized()
	return c.viper.GetStringSlice(key)
}

// GetBool returns the values associated with the given key as a boolean.
func (c *Configuration) GetBool(key string) bool {
	c.ensureInitialized()
	return c.viper.GetBool(key)
}

// GetByte returns the values associated with the given key as a byte array
// using the given encoding.
func (c *Configuration) GetByte(key string, enc Encoding) (b []byte, err error) {
	c.ensureInitialized()
	v := c.viper.GetString(key)

	switch enc {
	case UTF8:
		return []byte(v), nil
	case HEX:
		b, err = hex.DecodeString(v)
	default:
		return nil, errors.Errorf("unsupported encoding [%s]", enc)
	}

	return b, err
}

// GetDuration returns the values associated with the given key as a Duration
// If the format of the value does not match Duration 0 is returned.
func (c *Configuration) GetDuration(key string) time.Duration {
	c.ensureInitialized()
	return c.viper.GetDuration(key)
}

// SetFormat specifies the format of the configuration file.
// viper automatically detects the format so this method is mainly for testing.
func (c *Configuration) SetFormat(formatType string) {
	c.ensureInitialized()
	c.viper.SetConfigType(formatType)
}

// GetFloat64 returns the values associated with the given key as a float64.
func (c *Configuration) GetFloat64(key string) float64 {
	c.ensureInitialized()
	return c.viper.GetFloat64(key)
}

// GetInt64 returns the values associated with the given key as a int64.
func (c *Configuration) GetInt64(key string) int64 {
	c.ensureInitialized()
	return c.viper.GetInt64(key)
}

// GetUInt32 returns the values associated with the given key as a uint32.
func (c *Configuration) GetUInt32(key string) uint32 {
	c.ensureInitialized()
	return c.viper.GetUint32(key)
}

// GetUInt8 returns the values associated with the given key as a uint8.
func (c *Configuration) GetUInt8(key string) uint8 {
	c.ensureInitialized()
	return uint8(c.viper.GetUint(key))
}

// GetUInt64 returns the values associated with the given key as a uint64.
func (c *Configuration) GetUInt64(key string) uint64 {
	c.ensureInitialized()
	return c.viper.GetUint64(key)
}

// GetFloat32 returns the values associated with the given key as a float32.
func (c *Configuration) GetFloat32(key string) float32 {
	c.ensureInitialized()
	return float32(c.viper.GetFloat64(key))
}

// InitializeComponentConfig initializes a component configuration object.
// This method uses tags to read the key path where to read the value in the
// configuration file for a given field as well as to determine default values,
// and validation tags for validation (see
// https://godoc.org/gopkg.in/go-playground/validator.v9).
// Example:
// type TestConfig struct {
// 	   I        int           `configkey:"unittest.i" validate:"min=10" default:"10"`
// 	   S        string        `configkey:"unittest.s" validate:"required" default:"hoge"`
// 	   Ss       []string      `configkey:"unittest.ss" validate:"dive,required" default:"hoge,fuga"`
// 	   B        bool          `configkey:"unittest.b" default:"true"`
// 	   Utf8byte []byte        `configkey:"unittest.utf8byte,utf8" validate:"required" default:"abcde"`
// 	   Hexbytes []byte        `configkey:"unittest.hexbyte,hex" validate:"required" default:"abcd0e"`
// 	   Dr       time.Duration `configkey:"unittest.dr,duration" validate:"required" default:"1h10m10s"`
// 	   I64      int64         `configkey:"unittest.i64" validate:"min=11" default:"132904"`
// }
func (c *Configuration) InitializeComponentConfig(compConf interface{}) error {
	c.ensureInitialized()
	v := reflect.ValueOf(compConf)
	t := reflect.TypeOf(compConf)

	for v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
		t = v.Type()
	}

	if !v.CanSet() {
		panic("Configuration cannot be set, try passing it as a reference.")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tField := t.Field(i)
		tagValues := strings.Split(tField.Tag.Get(configTagName), ",")
		tag := tagValues[0]
		if tag == "" {
			continue
		}

		defaultValue := tField.Tag.Get(defaultTagName)
		if defaultValue != "" {
			if values := strings.Split(defaultValue, ","); len(values) > 1 {
				c.viper.SetDefault(tag, values)
			} else {
				c.viper.SetDefault(tag, defaultValue)
			}
		}

		switch field.Kind() {
		case reflect.Int:
			value := c.GetInt(tag)
			field.SetInt(int64(value))
		case reflect.Int64:
			if len(tagValues) > 1 && tagValues[1] == "duration" {
				value := c.GetDuration(tag)
				field.Set(reflect.ValueOf(value))
			} else {
				value := c.GetInt64(tag)
				field.SetInt(int64(value))
			}
		case reflect.String:
			value := c.GetString(tag)
			field.SetString(value)
		case reflect.Bool:
			value := c.GetBool(tag)
			field.SetBool(value)
		case reflect.Uint8:
			value := c.GetUInt8(tag)
			field.SetUint(uint64(value))
		case reflect.Uint32:
			value := c.GetUInt32(tag)
			field.SetUint(uint64(value))
		case reflect.Uint64:
			value := c.GetUInt64(tag)
			field.SetUint(value)
		case reflect.Float32:
			value := c.GetFloat32(tag)
			field.SetFloat(float64(value))
		case reflect.Float64:
			value := c.GetFloat64(tag)
			field.SetFloat(value)
		default:
			fieldType := field.Type()
			switch fieldType {
			case stringSliceType:
				value := c.GetStringSlice(tag)
				field.Set(reflect.ValueOf(value))
			case byteArrayType:
				var encoding Encoding
				var encodingTag string
				if len(tagValues) > 1 {
					encodingTag = tagValues[1]
				} else {
					encodingTag = ""
				}
				if encodingTag == "" || encodingTag == utf8TagValue {
					encoding = UTF8
				} else if encodingTag == hexTagValue {
					encoding = HEX
				} else {
					return errors.Errorf("Unknown encoding %v", encodingTag)
				}
				value, err := c.GetByte(tag, encoding)
				if err != nil {
					return errors.Errorf("Could not parse byte %v.", tag)
				}
				field.Set(reflect.ValueOf(value))
			default:
				return errors.Errorf("Unknown field type %v.", fieldType.Name())
			}
		}
	}

	validate := validator.New()
	return validate.Struct(compConf)
}

func (c *Configuration) ensureInitialized() {
	if !c.initialized {
		panic("Configuration used without being initialized.")
	}
}
