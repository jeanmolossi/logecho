package logecho

import (
	"sync"

	"go.uber.org/zap"
)

type Encoding string

const (
	JSON Encoding = "json"
	Text Encoding = "console"
)

type ZapLog struct {
	zl *zap.Logger
	m  *sync.RWMutex
}

func NewZapWithConfig(zcfg ZapConfig) *ZapLog {
	mCfg := zap.NewProductionConfig()
	if zcfg.IsDevelopment {
		mCfg = zap.NewDevelopmentConfig()
	}

	if zcfg.CallerKey == "" {
		mCfg.EncoderConfig.EncodeCaller = nil
	}

	mCfg.Sampling = nil // disable sampling

	mCfg.EncoderConfig.MessageKey = zcfg.msgKey()
	mCfg.EncoderConfig.CallerKey = zcfg.CallerKey
	mCfg.EncoderConfig.TimeKey = zcfg.getTimeKey()
	mCfg.EncoderConfig.EncodeTime = zcfg.getEncodeTime()
	mCfg.EncoderConfig.EncodeLevel = zcfg.getEncodeLevel()
	mCfg.Level = zap.NewAtomicLevelAt(zcfg.Level)
	mCfg.Encoding = string(zcfg.getEncoding())

	zapLog := &ZapLog{
		zl: zap.Must(mCfg.Build()),
		m:  &sync.RWMutex{},
	}

	return zapLog
}

func NewZap() *ZapLog {
	return NewZapWithConfig(defaultConfig)
}

var Logger = NewZapWithConfig(DevConfig)
