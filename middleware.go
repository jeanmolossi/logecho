package logecho

import "github.com/labstack/echo/v4"

type Config struct {
	EnableLatency      bool
	EnableRequestCount bool
	EnableRequestID    bool

	Fields Fields
}

var (
	defaultTpl = Fields{
		"host":           Host,
		"request.method": Method,
		"request.path":   Path,
		"real-ip":        RealIP,
		"user-agent":     UserAgent,
	}

	defaultMiddlewareConfig = Config{
		EnableRequestCount: true,
		EnableLatency:      true,
		EnableRequestID:    true,

		Fields: defaultTpl,
	}
)

func Middleware() echo.MiddlewareFunc {
	return MiddlewareWithConfig(defaultMiddlewareConfig)
}

func MiddlewareWithTemplate(tpl Fields) echo.MiddlewareFunc {
	cfg := defaultMiddlewareConfig
	cfg.Fields = tpl

	return MiddlewareWithConfig(cfg)
}

func MiddlewareWithConfig(cfg Config) echo.MiddlewareFunc {
	if len(cfg.Fields) == 0 {
		cfg.Fields = defaultTpl
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer Logger.zl.Sync()

			configTemplate(cfg.Fields, c)

			if cfg.EnableLatency {
				initLatencyCalc(c)
			}

			if cfg.EnableRequestID {
				installXRequestID(c)
			}

			if cfg.EnableRequestCount {
				incrementRequestCounter()
			}

			var err error
			if err = next(c); err != nil {
				c.Error(err)
			}

			Logger.Info(c, "request done")

			if cfg.EnableRequestCount {
				decrementRequestCounter()
			}

			return err
		}
	}
}
