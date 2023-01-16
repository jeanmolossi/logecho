package logecho

import (
	"net/url"
	"os"
	"text/template"

	"github.com/labstack/echo/v4"
)

// ContextFields is the accepted field to generate
// a Fields template. It have two classes, Request fields and
// Response fields.
//
// Request fields can not be present in case that is a optional field.
//
// Response fields can not be present in case that is a optional field or
// can be present with a default request value.
//
// A Response field who appears before a request response is Status.
// That one will be 200 by default incoming from request defaults
type ContextFields struct {
	// request fields

	RequestURI      string
	RequestID       string
	TransactionID   string
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

	Status   int   // Status is Response class field. It can be wrong before request's response
	BytesOut int64 // BytesOut is Response class field. It can be wrong before request's response
}

func getTemplateFields(c echo.Context) ContextFields {
	return ContextFields{
		// request fields
		RequestURI:      c.Request().RequestURI,
		RequestID:       getXRequestID(c),
		TransactionID:   getTransactionID(c),
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
