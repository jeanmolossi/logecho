package logecho

import (
	"sync"

	"go.uber.org/zap"
)

// Encoding is the string to represents the accepted
// encoding formats
type Encoding string

const (
	// JSON will set log format to a json message
	JSON Encoding = "json"
	// Text will set log format to a mixed message
	//
	// The message, timestamp and level will be printed
	// as read friendly message.
	//
	// Additional fields are
	// been printed as JSON in all cases
	Text Encoding = "console"
)

// Logecho wrapps the zap.Logger in to a logger who has
// mutex.
//
// That struct will manage read-write locks on the instance.
// The original logger is from package go.uber.org/zap
//
// More about zap logger:
//
// https://pkg.go.dev/go.uber.org/zap
type Logecho struct {
	zl *zap.Logger
	m  *sync.RWMutex
}

// NewZapWithConfig enables custom configuration to instantiate
// a new ZapLog
func NewZapWithConfig(config Config) *Logecho {
	initConfig := zap.NewProductionConfig()
	if config.IsDevelopment {
		initConfig = zap.NewDevelopmentConfig()
	}

	if config.CallerKey == "" {
		initConfig.EncoderConfig.EncodeCaller = nil
	}

	initConfig.Sampling = nil // disable sampling

	initConfig.EncoderConfig.MessageKey = config.msgKey()
	initConfig.EncoderConfig.CallerKey = config.CallerKey
	initConfig.EncoderConfig.TimeKey = config.getTimeKey()
	initConfig.EncoderConfig.EncodeTime = config.getEncodeTime()
	initConfig.EncoderConfig.EncodeLevel = config.getEncodeLevel()
	initConfig.Level = zap.NewAtomicLevelAt(config.Level)
	initConfig.Encoding = string(config.getEncoding())

	zapLog := &Logecho{
		zl: zap.Must(initConfig.Build()),
		m:  &sync.RWMutex{},
	}

	return zapLog
}

// NewZap instantiate a ZapLog with default configs
func NewZap() *Logecho {
	return NewZapWithConfig(defaultConfig)
}

// Logger is a singleton to a ZapLog instance
var Logger = NewZap()
