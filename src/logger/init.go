package logger

import (
	"errors"
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	EnvLoggerEnvironment = "ZAP_ENVIRONMENT"
	DefLoggerEnvironment = "development"
	EnvLoggerLevel       = "ZAP_LEVEL"
	DefLoggerLevel       = "debug"
)

func init() {
	// Initializing zap environment
	zapEnvironment, ok := os.LookupEnv(EnvLoggerEnvironment)
	if !ok {
		zapEnvironment = DefLoggerEnvironment
	}

	// Initializing zap level
	zapLevel, ok := os.LookupEnv(EnvLoggerLevel)
	if !ok {
		zapLevel = DefLoggerLevel
	}

	_, err := GetLogger(
		&LoggerConfig{
			Environment: ToLoggerEnvironment(zapEnvironment),
			Level:       ToLoggerLevel(zapLevel),
		},
	)
	if err != nil {
		log.Fatal("unexpected error while getting logger", zap.Error(err))
	}
}

type LoggerEnvironment int64

const (
	Development LoggerEnvironment = iota + 1
	Production
	DevelopmentJson
)

// String casts LoggerEnvironment type to string value
func (s LoggerEnvironment) String() string {
	switch s {
	case Development:
		return "development"
	case Production:
		return "production"
	case DevelopmentJson:
		return "development-json"
	default:
		return "unknown"
	}
}

// CheckOnDefine checks if zap environment is defined and returns nil though
func (s LoggerEnvironment) CheckOnDefine() error {
	switch s {
	case Development:
		return nil
	case Production:
		return nil
	case DevelopmentJson:
		return nil
	}
	return fmt.Errorf(
		"logger environment not defined, use one of predefined logger environment: %s, %s, %s",
		Development.String(),
		Production.String(),
		DevelopmentJson.String(),
	)
}

// CheckZapLevel checks if zap level is defined and returns nil though
func CheckZapLevel(level zapcore.Level) error {
	if level == zapcore.DebugLevel ||
		level == zapcore.InfoLevel ||
		level == zapcore.WarnLevel ||
		level == zapcore.ErrorLevel ||
		level == zapcore.DPanicLevel ||
		level == zapcore.PanicLevel ||
		level == zapcore.FatalLevel {
		return nil
	}

	return errors.New(
		"zap level not defined, use predefined zap level:" +
			fmt.Sprintf("\n"+
				"\t\t%s\n"+ // Debug
				"\t\t%s\n"+ // Info
				"\t\t%s\n"+ // Warn
				"\t\t%s\n"+ // Error
				"\t\t%s\n"+ // DPanic
				"\t\t%s\n"+ // Panic
				"\t\t%s\n", // Fatal
				zapcore.DebugLevel.String(),
				zapcore.InfoLevel.String(),
				zapcore.WarnLevel.String(),
				zapcore.ErrorLevel.String(),
				zapcore.DPanicLevel.String(),
				zapcore.PanicLevel.String(),
				zapcore.FatalLevel.String(),
			),
	)
}

// ToLoggerEnvironment converts string to LoggerEnvironment
func ToLoggerEnvironment(arg string) LoggerEnvironment {
	switch arg {
	case "development":
		return Development
	case "production":
		return Production
	case "development-json":
		return DevelopmentJson
	default:
		return Development
	}
}

// ToLoggerLevel converts string to LoggerLevel
func ToLoggerLevel(arg string) zapcore.Level {
	switch arg {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zap.DebugLevel
	}
}
