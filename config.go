package logecho

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultConfig ZapConfig = ZapConfig{
	IsDevelopment: false,
	Encoding:      JSON,
}

var DevConfig ZapConfig = ZapConfig{
	IsDevelopment: true,
	Level:         zap.DebugLevel,
	Encoding:      Text,
}

type ZapConfig struct {
	IsDevelopment bool

	MessageKey string
	CallerKey  string
	TimeKey    string

	Level       zapcore.Level
	EncodeTime  zapcore.TimeEncoder
	EncodeLevel zapcore.LevelEncoder
	Encoding    Encoding
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
