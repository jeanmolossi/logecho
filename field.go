package logecho

import (
	"bytes"
	"reflect"
)

type (
	// FuncField represents a Field who can be added
	// to template fields but is a function and not a
	// plain value.
	//
	// Explain:
	//
	// 	Header is a FuncField because that should runs
	// 	as a function at log execution.
	FuncField struct {
		// name represents the field name
		name string

		// result is the result Kind to FuncField
		// func execution
		//
		// Header should returns a string. The result will
		// be reflect.String
		result reflect.Kind
	}

	// Field is plain valued field on Fields template
	Field struct {
		// tpl is the Template string format.
		//
		// Example:
		//
		// 	{{ .UserAgent }}
		tpl string

		// kind is the field type.
		//
		// UserAgent as Example has kind as reflect.String
		// because UserAgent is a string
		kind reflect.Kind
	}

	// Fields is a mapping to key->Field. It will be used
	// to set key to value in the Fields template.
	//
	// Example:
	//
	//	logecho.Fields{
	//		"user-agent": logecho.UserAgent
	// 	}
	//
	// Log will be printed like:
	//
	//	{"user-agent":"curl/7.54"}
	Fields map[string]Field
)

var (
	RequestURI      = Field{"{{ .RequestURI }}", reflect.String}      // Built-in field to RequestURI - Request class field
	RequestID       = Field{"{{ .RequestID }}", reflect.String}       // Built-in field to RequestID - Request class field
	RealIP          = Field{"{{ .RealIP }}", reflect.String}          // Built-in field to RealIP - Request class field
	Host            = Field{"{{ .Host }}", reflect.String}            // Built-in field to Host - Request class field
	Method          = Field{"{{ .Method }}", reflect.String}          // Built-in field to Method - Request class field
	Referer         = Field{"{{ .Referer }}", reflect.String}         // Built-in field to Referer - Request class field
	UserAgent       = Field{"{{ .UserAgent }}", reflect.String}       // Built-in field to UserAgent - Request class field
	Query           = Field{"{{ .Query }}", reflect.String}           // Built-in field to Query - Request class field
	Path            = Field{"{{ .Path }}", reflect.String}            // Built-in field to Path - Request class field
	UrlEncodedQuery = Field{"{{ .UrlEncodedQuery }}", reflect.String} // Built-in field to UrlEncodedQuery - Request class field
	BytesIn         = Field{"{{ .BytesIn }}", reflect.String}         // Built-in field to BytesIn - Request class field
	Body            = Field{"{{ .Body }}", reflect.String}            // Built-in field to Body - Request class field
	Params          = Field{"{{ .Params }}", reflect.String}          // Built-in field to Params - Request class field

	Status   = Field{"{{ .Status }}", reflect.Int}     // Built-in field to Status - Response class field. It can be empty or with default value
	BytesOut = Field{"{{ .BytesOut }}", reflect.Int64} // Built-in field to BytesOut - Response class field. It can be empty or with default value

	LatencyInMicroS = Field{"{{ Latency \"us\" }}", reflect.Int}        // Built-in field to Latency in microseconds
	LatencyInNs     = Field{"{{ Latency \"ns\" }}", reflect.Int}        // Built-in field to Latency in nanoseconds
	LatencyInMs     = Field{"{{ Latency \"ms\" }}", reflect.Int}        // Built-in field to Latency in milliseconds
	LatencyInSec    = Field{"{{ Latency \"s\" }}", reflect.Int}         // Built-in field to Latency in seconds
	LatencyString   = Field{"{{ Latency \"string\" }}", reflect.String} // Built-in field to Latency in string format
)

var (
	ParamField      = FuncField{"Param", reflect.String}        // Built-in FuncField to Param function
	HeaderField     = FuncField{"Header", reflect.String}       // Built-in FuncField to Header function
	CookieField     = FuncField{"Cookie", reflect.String}       // Built-in FuncField to Cookie function
	LatencyField    = FuncField{"Latency", reflect.String}      // Built-in FuncField to Latency function
	GetenvField     = FuncField{"Getenv", reflect.String}       // Built-in FuncField to Getenv function
	RunningReqField = FuncField{"RunningReqField", reflect.Int} // Built-in FuncField to RunningReq function
	ConcurrentField = FuncField{"Concurrent", reflect.Int}      // Built-in FuncField to Concurrent function
)

var (
	RunningRequests    Field = FuncFieldWithArgs(RunningReqField) // Built-in Field to get RunningRequests values
	ConcurrentRequests Field = FuncFieldWithArgs(ConcurrentField) // Built-in Field to get ConcurrentRequests values
)

var (
	// Header is a built-in function who accepts headers to define in key.
	//
	// Example:
	//
	//	logecho.Fields{
	//		// it will get header x-origin in context and set in log msg
	//		"x-origin": logecho.Header("x-origin")
	// 	}
	Header = header

	// Param is a built-in function who accepts param to define in key.
	//
	// Example:
	//
	//	logecho.Fields{
	//		// it will get path param id and set in log msg
	//		"id": logecho.Param("id")
	// 	}
	Param = param

	// Cookie is a built-in function who accepts cookie to define in key.
	//
	// Example:
	//
	//	logecho.Fields{
	//		// it will get path param id and set in log msg
	//		"session": logecho.Cookie("session")
	// 	}
	Cookie = cookie

	// Getenv is a built-in function who accepts env to define in key.
	//
	// Example:
	//
	//	logecho.Fields{
	//		// it will get path param id and set in log msg
	//		"app-name": logecho.Getenv("APP_NAME")
	// 	}
	Getenv = getenv
)

func header(headers ...string) Field {
	return FuncFieldWithArgs(HeaderField, headers...)
}

func param(paramNames ...string) Field {
	return FuncFieldWithArgs(ParamField, paramNames...)
}

func cookie(cookies ...string) Field {
	return FuncFieldWithArgs(CookieField, cookies...)
}

func getenv(env string) Field {
	return FuncFieldWithArgs(GetenvField, env)
}

// FuncFieldWithArgs is to build a Field template.
//
// FuncFieldWithArgs build templates to functions who can
// receive arguments.
//
// Example:
//
//	customField := FuncField{"Custom",reflect.String}
//	f := FuncFieldWithArgs(customField, "arg1", "arg2")
//
// Will build the following template:
//
//	"{{ Custom \"arg1\" \"arg2\" }}"
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

// Tpl will return built tpl to the Field
//
//	"{{ .Field }}"
//
// Or if a FuncField result
//
//	"{{ Field \"arg\" }}"
func (f Field) Tpl() string {
	return f.tpl
}

// Type will return field's kind
func (f Field) Type() reflect.Kind {
	return f.kind
}
