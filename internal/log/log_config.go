package log

import "time"

// Config contains the configuration parameters for the log.
type Config struct {
	OutputStdout     bool          `configkey:"log.output_stdout"`
	RotationCount    int           `configkey:"log.rotation_counts" validate:"required_without=OutputStdout"`
	RotationInterval time.Duration `configkey:"log.rotation_interval,duration" validate:"required_without=OutputStdout"`
	LogDir           string        `configkey:"log.dir" validate:"required_without=OutputStdout"`
	LogFileBaseName  string        `configkey:"log.basename" validate:"required_without=OutputStdout"`
	LogFormat        string        `configkey:"log.format" validate:"eq=json|eq=text"`
	LogLevel         string        `configkey:"log.level" validate:"required"`
}
