package logecho

import (
	"net/url"
	"os"
	"text/template"

	"github.com/labstack/echo/v4"
)

type ContextFields struct {
	// request fields

	RequestID       string
	RealIP          string
	Host            string
	Method          string
	Referer         string
	UserAgent       string
	Query           url.Values
	Path            string
	UrlEncodedQuery string
	BytesIn         string
	Body            string
	Params          string

	// response fields

	Status   int
	BytesOut int64
}

func getTemplateFields(c echo.Context) ContextFields {
	return ContextFields{
		// request fields
		RequestID:       getXRequestID(c),
		RealIP:          c.RealIP(),
		Host:            c.Request().Host,
		Method:          c.Request().Method,
		Referer:         c.Request().Referer(),
		UserAgent:       c.Request().UserAgent(),
		Query:           c.Request().URL.Query(),
		Path:            c.Request().URL.Path,
		UrlEncodedQuery: c.Request().URL.RawQuery,
		BytesIn:         c.Request().Header.Get(echo.HeaderContentLength),
		Body:            readBody(c),
		Params:          string(buildParams(c)),
		// response fields
		Status:   c.Response().Status,
		BytesOut: c.Response().Size,
	}
}

func getTemplateFuncMap(c echo.Context) template.FuncMap {
	return template.FuncMap{
		string(ParamField.name):      c.Param,
		string(HeaderField.name):     getHeader(c),
		string(CookieField.name):     extractCookie(c),
		string(LatencyField.name):    calcLatency(c),
		string(GetenvField.name):     func(key string) string { return os.Getenv(key) },
		string(RunningReqField.name): CurrentCount,
		string(ConcurrentField.name): TransactionCounter,
	}
}
