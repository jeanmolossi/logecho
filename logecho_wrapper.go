package logecho

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap/zapcore"
)

func (z *ZapLog) acquireContext(c echo.Context, call func(f ...zapcore.Field)) {
	z.m.Lock()
	call(getFields(readContext(c))...)
	z.m.Unlock()
}

func (z *ZapLog) Print(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Debug(s, f...) })
}

func (z *ZapLog) Printf(format string, i ...interface{}) {
	z.zl.Sugar().Debugf(format, i...)
}

func (z *ZapLog) Debug(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Debug(s, f...) })
}

func (z *ZapLog) Info(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Info(s, f...) })
}

func (z *ZapLog) Warn(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Warn(s, f...) })
}

func (z *ZapLog) Error(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Error(s, f...) })
}

func (z *ZapLog) Panic(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Panic(s, f...) })
}

func (z *ZapLog) Fatal(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Fatal(s, f...) })
}
