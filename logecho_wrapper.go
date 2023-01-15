package logecho

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap/zapcore"
)

func (z *Logecho) acquireContext(c echo.Context, call func(f ...zapcore.Field)) {
	z.m.Lock()
	call(getFields(readContext(c))...)
	z.m.Unlock()
}

func (z *Logecho) Print(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Debug(s, f...) })
}

func (z *Logecho) Printf(format string, i ...interface{}) {
	z.zl.Sugar().Debugf(format, i...)
}

func (z *Logecho) Debug(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Debug(s, f...) })
}

func (z *Logecho) Info(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Info(s, f...) })
}

func (z *Logecho) Warn(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Warn(s, f...) })
}

func (z *Logecho) Error(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Error(s, f...) })
}

func (z *Logecho) Panic(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Panic(s, f...) })
}

func (z *Logecho) Fatal(c echo.Context, s string) {
	z.acquireContext(c, func(f ...zapcore.Field) { z.zl.Fatal(s, f...) })
}
