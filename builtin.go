package logecho

import (
	"bytes"
	"io"
	"time"

	"github.com/labstack/echo/v4"
)

var paramBuf = &bytes.Buffer{}

func buildParams(c echo.Context) []byte {
	paramBuf.Reset()
	keys := c.ParamNames()
	lastIndex := len(keys)

	paramBuf.WriteRune('{')
	for i, key := range keys {
		paramBuf.WriteRune('"')
		paramBuf.WriteString(key)
		paramBuf.WriteString(`":"`)
		paramBuf.WriteString(c.Param(key))
		paramBuf.WriteRune('"')
		if i+1 != lastIndex {
			paramBuf.WriteRune(',')
		}
	}
	paramBuf.WriteRune('}')
	return paramBuf.Bytes()
}

func readBody(c echo.Context) string {
	bytes, _ := io.ReadAll(c.Request().Body)
	return string(bytes)
}

func extractCookie(c echo.Context) func(string) string {
	return func(name string) string {
		cookie, err := c.Cookie(name)
		if err != nil {
			return ""
		}

		return cookie.Value
	}
}

func getHeader(c echo.Context) func(headers ...string) string {
	return func(headers ...string) string {
		for _, key := range headers {
			header := c.Request().Header.Get(key)
			if header != "" {
				return header
			}
		}

		return ""
	}
}

func initLatencyCalc(c echo.Context) {
	c.Set("start", time.Now())
}

func getStartFromCtx(c echo.Context) time.Time {
	if start, ok := c.Get("start").(time.Time); ok {
		return start
	}

	return time.Now()
}

func calcLatency(c echo.Context) func(string) interface{} {
	return func(scale string) interface{} {
		start := getStartFromCtx(c)
		stop := time.Since(start)

		switch scale {
		case "s":
			return stop.Seconds()
		case "ms":
			return stop.Milliseconds()
		case "us":
			return stop.Microseconds()
		case "ns":
			return stop.Nanoseconds()
		default:
			return stop.String()
		}

	}
}
