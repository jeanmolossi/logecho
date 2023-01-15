package logecho

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

type ContextOption func(r *http.Request, w *httptest.ResponseRecorder, c echo.Context)

func NewContext(options ...ContextOption) echo.Context {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/path", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	if len(options) > 0 {
		for _, option := range options {
			option(req, rec, ctx)
		}
	}

	return ctx
}

func WithParams(params map[string]string) ContextOption {
	return func(r *http.Request, w *httptest.ResponseRecorder, c echo.Context) {
		names := make([]string, 0, len(params))
		values := make([]string, 0, len(params))

		b := new(bytes.Buffer)
		b.WriteString(c.Path())
		for name, value := range params {
			if name == "" {
				continue
			}

			names = append(names, name)
			values = append(values, value)
			b.WriteRune('/')
			b.WriteRune(':')
			b.WriteString(name)
		}

		c.SetPath(b.String())
		c.SetParamNames(names...)
		c.SetParamValues(values...)
	}
}
