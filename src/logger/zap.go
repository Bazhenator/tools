package logger

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerConfig struct is a zap logger configuration
type LoggerConfig struct {
	Environment LoggerEnvironment
	Level       zapcore.Level
}

func NewLoggerConfig() (*LoggerConfig, error) {
	var multiErr error

	loggerEnvironment, ok := os.LookupEnv(EnvLoggerEnvironment)
	if !ok {
		multierr.AppendInto(&multiErr, fmt.Errorf("%v is not defined", EnvLoggerEnvironment))
	}

	loggerLevel, ok := os.LookupEnv(EnvLoggerLevel)
	if !ok {
		multierr.AppendInto(&multiErr, fmt.Errorf("%v is not defined", EnvLoggerLevel))
	}

	if multiErr != nil {
		return nil, multiErr
	}

	return &LoggerConfig{
			Environment: ToLoggerEnvironment(loggerEnvironment),
			Level:       ToLoggerLevel(loggerLevel),
	}, nil
}

// checkEmptyFields checks if necessary fields is empty
func (c *LoggerConfig) checkEmptyFields() error {
	var multiErr error

	if err := c.Environment.CheckOnDefine(); err != nil {
		multierr.AppendInto(&multiErr, err)
	}

	if err := CheckZapLevel(c.Level); err != nil {
		multierr.AppendInto(&multiErr, err)
	}

	return multiErr
}

// GetLogger inits created logger instance
func GetLogger(c *LoggerConfig) (*zap.Logger, error) {
	if err := c.checkEmptyFields(); err != nil {
		log.Printf("error while checking config: %s", err.Error())
		return nil, err
	}

	var err error
	var l *zap.Logger
	switch c.Environment {
	case Development:
		lConfig := zap.Config{
			Level:             zap.NewAtomicLevelAt(c.Level),
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: true,
			Sampling:          nil,
			Encoding:          "console",
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:          "M",
				LevelKey:            "L",
				TimeKey:             "T",
				NameKey:             "",
				CallerKey:           "C",
				FunctionKey:         "F",
				StacktraceKey:       "",
				LineEnding:          zapcore.DefaultLineEnding,
				EncodeLevel:         zapcore.CapitalColorLevelEncoder,
				EncodeTime:          zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
				EncodeDuration:      nil,
				EncodeCaller:        zapcore.ShortCallerEncoder,
				EncodeName:          nil,
				ConsoleSeparator:    " | ",
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			InitialFields:    nil,
		}

		l, err = lConfig.Build()
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
	case DevelopmentJson:
		lConfig := zap.Config{
			Level:             zap.NewAtomicLevelAt(c.Level),
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: true,
			Sampling:          nil,
			Encoding:          "json",
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:          "M",
				LevelKey:            "L",
				TimeKey:             "T",
				NameKey:             "",
				CallerKey:           "C",
				FunctionKey:         "F",
				StacktraceKey:       "",
				LineEnding:          zapcore.DefaultLineEnding,
				EncodeLevel:         zapcore.CapitalColorLevelEncoder,
				EncodeTime:          zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
				EncodeDuration:      nil,
				EncodeCaller:        zapcore.ShortCallerEncoder,
				EncodeName:          nil,
				ConsoleSeparator:    " | ",
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			InitialFields:    nil,
		}

		l, err = lConfig.Build()
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
	default:
		l = zap.NewExample()
		l.Warn("app environment not set, set default")
	}

	zap.ReplaceGlobals(l)

	l.Info(fmt.Sprintf("logger instance created with config: %+v", *c))
	return l, nil
}
