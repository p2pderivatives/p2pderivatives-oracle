package log

import (
	"io"
	"os"
	"path/filepath"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Log is used by the application to log information.
type Log struct {
	config      *Config
	rotateLog   *rotatelogs.RotateLogs
	Logger      *logrus.Logger
	initialized bool
}

// NewLog creates a new log structure.
func NewLog(config *Config) *Log {
	return &Log{
		config:      config,
		initialized: false,
	}
}

// Initialize sets up a log instance.
func (l *Log) Initialize() error {
	var w io.Writer

	stdoutFlag := l.config.OutputStdout
	if stdoutFlag {
		w = os.Stdout
	} else {
		rotatelog, err := l.initializeRotateLog()
		if err != nil {
			return errors.Wrap(err, "failed to initialize rotatelogs")
		}
		w = rotatelog
		l.rotateLog = rotatelog
	}

	logger, err := l.initializeLogrus(w)
	if err != nil {
		return errors.Wrap(err, "failed to initialize logrus")
	}
	l.Logger = logger

	l.initialized = true
	return nil
}

// initializeRotateLog initialize the rotating log for the log instance.
func (l *Log) initializeRotateLog() (*rotatelogs.RotateLogs, error) {
	count := l.config.RotationCount
	interval := l.config.RotationInterval

	logOption := []rotatelogs.Option{
		rotatelogs.WithRotationCount(count),
		rotatelogs.WithRotationTime(interval),
	}
	dir := l.config.LogDir
	basename := l.config.LogFileBaseName
	return rotatelogs.New(filepath.Join(dir, basename), logOption...)
}

// initializeLogrus initializes the Logrus for the log instance.
func (l *Log) initializeLogrus(writer io.Writer) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.SetOutput(writer)
	var formatter logrus.Formatter
	f := l.config.LogFormat
	switch f {
	case "json":
		formatter = &logrus.JSONFormatter{}
	case "text":
		formatter = &logrus.TextFormatter{FullTimestamp: true, QuoteEmptyFields: true}
	default:
		return nil, errors.Errorf("illegal log format [%s], specify \"text\" or \"json\" with \"log.format\" key", f)
	}
	logger.SetFormatter(formatter)
	v := l.config.LogLevel
	level, err := logrus.ParseLevel(v)
	if err != nil {
		return nil, errors.Wrapf(err, "illegal log level [%s]", v)
	}
	logger.SetLevel(level)
	return logger, nil
}

// IsInitialized returns whether the log instance is initialized.
func (l *Log) IsInitialized() bool {
	return l.initialized
}

// Finalize cleans up the resources of the log instance.
func (l *Log) Finalize() error {
	if l.rotateLog != nil {
		if err := l.rotateLog.Close(); err != nil {
			return errors.Wrap(err, "failed to close rotatelog")
		}
	}
	return nil
}

// NewEntry creates a new entry.
func (l *Log) NewEntry() *logrus.Entry {
	return logrus.NewEntry(l.Logger)
}
