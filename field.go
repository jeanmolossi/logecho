package logecho

import (
	"bytes"
	"reflect"
)

type (
	FuncField struct {
		name   string
		result reflect.Kind
	}

	Field struct {
		tpl  string
		kind reflect.Kind
	}

	Fields map[string]Field
)

var (
	RequestID       = Field{"{{ .RequestID }}", reflect.String}
	RealIP          = Field{"{{ .RealIP }}", reflect.String}
	Host            = Field{"{{ .Host }}", reflect.String}
	Method          = Field{"{{ .Method }}", reflect.String}
	Referer         = Field{"{{ .Referer }}", reflect.String}
	UserAgent       = Field{"{{ .UserAgent }}", reflect.String}
	Query           = Field{"{{ .Query }}", reflect.String}
	Path            = Field{"{{ .Path }}", reflect.String}
	UrlEncodedQuery = Field{"{{ .UrlEncodedQuery }}", reflect.String}
	BytesIn         = Field{"{{ .BytesIn }}", reflect.String}
	Body            = Field{"{{ .Body }}", reflect.String}
	Params          = Field{"{{ .Params }}", reflect.String}

	Status   = Field{"{{ .Status }}", reflect.Int}
	BytesOut = Field{"{{ .BytesOut }}", reflect.Int64}

	LatencyInMicroS = Field{"{{ Latency \"us\" }}", reflect.Int}
	LatencyInNs     = Field{"{{ Latency \"ns\" }}", reflect.Int}
	LatencyInMs     = Field{"{{ Latency \"ms\" }}", reflect.Int}
	LatencyInSec    = Field{"{{ Latency \"s\" }}", reflect.Int}
	LatencyString   = Field{"{{ Latency \"string\" }}", reflect.String}
)

var (
	ParamField      = FuncField{"Param", reflect.String}
	HeaderField     = FuncField{"Header", reflect.String}
	CookieField     = FuncField{"Cookie", reflect.String}
	LatencyField    = FuncField{"Latency", reflect.String}
	GetenvField     = FuncField{"Getenv", reflect.String}
	RunningReqField = FuncField{"RunningReqField", reflect.Int}
	ConcurrentField = FuncField{"Concurrent", reflect.Int}
)

var (
	RunningRequests    Field = FuncFieldWithArgs(RunningReqField)
	ConcurrentRequests Field = FuncFieldWithArgs(ConcurrentField)
)

var (
	Header = header
)

func header(headers ...string) Field {
	return FuncFieldWithArgs(HeaderField, headers...)
}

func FuncFieldWithArgs(field FuncField, args ...string) Field {
	b := new(bytes.Buffer)
	b.WriteString("{{")
	b.WriteString(field.name)
	for _, arg := range args {
		if arg == "" {
			continue
		}

		b.WriteRune(' ')
		b.WriteString(`"`)
		b.WriteString(arg)
		b.WriteString(`"`)
	}
	b.WriteString("}}")

	return Field{b.String(), field.result}
}

func (f Field) Tpl() string {
	return f.tpl
}

func (f Field) Type() reflect.Kind {
	return f.kind
}
