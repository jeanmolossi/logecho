package logecho

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// defaultConfig set up the production default config
var defaultConfig ZapConfig = ZapConfig{
	IsDevelopment: false,
	Encoding:      JSON,
}

// DevConfig set up a development builtin config
var DevConfig ZapConfig = ZapConfig{
	IsDevelopment: true,
	Level:         zap.DebugLevel,
	Encoding:      Text,
}

// ZapConfig defines the base config to our Logger
//
// You can set Keys to preset
//
//	Message
//	Caller
//	Time
//
// Key overrides only apply to key value log formats
type ZapConfig struct {
	// IsDevelopment as the property describes, is to
	// set development configs
	IsDevelopment bool

	// MessageKey will change the prefix "message"
	// on the log.
	//
	// Only applied to key value log formats
	MessageKey string
	// CallerKey will change the prefix "caller"
	// on the log. It can be empty to remove caller
	// from log message
	//
	// Only applied to key value log formats
	CallerKey string
	// TimeKey will change the prefix "timestamp"
	// on the log.
	//
	// Only applied to key value log formats
	TimeKey string

	// Level sets the level of logs who should appear
	//
	// Default is `0` who is the same as INFO
	Level zapcore.Level

	// EncodeTime is a option to change how the timestamp
	// will be formatted.
	//
	// Default is RFC3339-formatted string
	EncodeTime zapcore.TimeEncoder

	// EncodeLevel defines how de format will be printed
	//
	// Default is Level serializer to a lowercase string
	EncodeLevel zapcore.LevelEncoder

	// Encoding changes the log format.
	//
	// Can be any of:
	//
	// 	JSON
	// 	Text
	//
	// Defaults
	//
	// Development environments Text is default
	//
	// Production (same as IsDevelopment = false) JSON is default
	Encoding Encoding
}

func (z ZapConfig) msgKey() string {
	if z.MessageKey == "" {
		return "message"
	}

	return z.MessageKey
}

func (z ZapConfig) getTimeKey() string {
	if z.TimeKey == "" {
		return "timestamp"
	}

	return z.TimeKey
}

func (z ZapConfig) getEncodeTime() zapcore.TimeEncoder {
	if z.EncodeTime == nil {
		return zapcore.RFC3339TimeEncoder
	}

	return z.EncodeTime
}

func (z ZapConfig) getEncodeLevel() zapcore.LevelEncoder {
	if z.EncodeLevel == nil {
		if z.IsDevelopment {
			return zapcore.LowercaseColorLevelEncoder
		}

		return zapcore.LowercaseLevelEncoder
	}

	return z.EncodeLevel
}

func (z ZapConfig) getEncoding() Encoding {
	if z.Encoding == "" {
		return JSON
	}

	return z.Encoding
}
