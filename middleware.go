package logecho

import "github.com/labstack/echo/v4"

// MiddlewareConfig is a struct to set a custom
// config to echo.Middleware if you want to remove or add a config.
type MiddlewareConfig struct {
	// EnableLatency will enable latency calc.
	//
	// IMPORTANT: To see latency at log point you should
	// add the following to your Fields template:
	//
	// 	logecho.Fields{
	// 		"latency": logecho.LatencyString, // or any latency field
	// 	}
	EnableLatency bool

	// EnableRequestCount will enable request counter. It counts
	// how many request are running and was run in current transcation.
	//
	// Look to the examples:
	//
	// E.g.1.: Your app receive a single request, then you have the
	// following metrics:
	//
	//	- running requests	-> 1
	//	- transaction requests -> 1
	//
	// E.g.2.: Your app receive 100 requests, but 20 was finished. Then you
	// have the following metrics
	//
	//	- running requests	-> 80
	//	- transaction requests -> 100
	EnableRequestCount bool

	// EnableRequestID will enable transaction request ID. It will
	// set request ID on the context and the response header X-Request-ID.
	//
	// Useful to tracing all logs from a single request.
	EnableRequestID bool

	// Fields will set how aditional fields will be printed on log
	// messages.
	//
	// It set's the key and what will be logged
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

	defaultMiddlewareConfig = MiddlewareConfig{
		EnableRequestCount: true,
		EnableLatency:      true,
		EnableRequestID:    true,

		Fields: defaultTpl,
	}
)

// Middleware is a default config middleware. Look the default configs:
//
// RequestID, Latency and RequestCount are enabled by default
func Middleware() echo.MiddlewareFunc {
	return MiddlewareWithConfig(defaultMiddlewareConfig)
}

// MiddlewareWithTemplate is the middleware default config, but in this
// you can set up the template fields who you want in log
func MiddlewareWithTemplate(tpl Fields) echo.MiddlewareFunc {
	cfg := defaultMiddlewareConfig
	cfg.Fields = tpl

	return MiddlewareWithConfig(cfg)
}

// MiddlewareWithConfig is the middleware with a custom config.
//
// It will set default template if has no Fields in the config.
//
// Additional that it will log message "request done" on end call
func MiddlewareWithConfig(cfg MiddlewareConfig) echo.MiddlewareFunc {
	if len(cfg.Fields) == 0 {
		cfg.Fields = defaultTpl
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer Logger.zl.Sync()

			configTemplate(cfg.Fields, c)
			installTransactionID(c)

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

			Logger.Info(c, "handled request")

			if cfg.EnableRequestCount {
				decrementRequestCounter()
			}

			return err
		}
	}
}
