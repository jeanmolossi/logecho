package logecho

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"sync"
	"text/template"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	tpl       = new(template.Template)
	poolMutex = new(sync.RWMutex)
	buf       = new(bytes.Buffer)
)

// readContext gets the buffer and execute template writing
// result in a buffer.
//
// Returns the result of the template execution in bytes
func readContext(c echo.Context) []byte {
	poolMutex.RLock()
	buf := new(bytes.Buffer)
	tpl.Execute(buf, getTemplateFields(c))
	poolMutex.RUnlock()

	return buf.Bytes()
}

// getFields will receive a slice of bytes who contains
// bytes from template execution to build zapcore.Field slice.
func getFields(fbytes []byte) []zapcore.Field {
	// lock bytes and map work until that build
	poolMutex.Lock()
	defer poolMutex.Unlock()

	// expect a json format of bytes to unmarshal
	// into a map[string]interface{}
	var f map[string]interface{}
	json.Unmarshal(fbytes, &f)

	fields := make([]zapcore.Field, 0, len(f))
	for field, value := range f {
		switch v := value.(type) {
		// by default any number will income as float64
		case float64:
			fields = append(fields, zap.Float64(field, v))
		case int64:
			fields = append(fields, zap.Int64(field, v))
		case int:
			fields = append(fields, zap.Int(field, v))
		case url.Values:
			fields = append(fields, zap.String(field, v.Encode()))
		case string:
			fields = append(fields, zap.String(field, v))
		default:
			log.Println("[DEBUG] field", field, "was added as string but is ", reflect.TypeOf(v))
			fields = append(fields, zap.String(field, fmt.Sprintf("%v", v)))
			continue
		}
	}

	return fields
}

// configTemplate will read ctxFields were defined in MiddlewareConfig and
// build the template format to write on log
func configTemplate(ctxFields Fields, c echo.Context) {
	// lock tpl work to avoid concurrent requests
	poolMutex.Lock()
	parseFields(ctxFields, buf)

	tpl = template.Must(
		template.New("fields").
			Funcs(getTemplateFuncMap(c)).
			Parse(buf.String()),
	)
	poolMutex.Unlock()
}

// parseFields will receive Fields tpl config and write into a buffer
// with JSON format
func parseFields(f Fields, w *bytes.Buffer) {
	w.Reset()

	lastIndex := len(f)
	indexCounter := 0

	// Start string with {
	w.WriteRune('{')
	for key, value := range f {
		indexCounter++

		// at this point the buffer contains the string as:
		// {"
		w.WriteRune('"')
		// {"keyValue
		w.WriteString(key)
		// {"keyValue":
		w.WriteString(`":`)
		// {"keyValue": "value"
		writeField(w, value)

		if indexCounter != lastIndex {
			// when has more than one field adds ,
			// and continue loop
			// {"keyValue": "value",
			w.WriteRune(',')
		}
	}
	// {"keyValue": "value"}
	w.WriteRune('}')
}

// writeField will check Type from field and write correctly on
// buffer. To quoted fields write "{{ Tpl }}", to unquoted
// fields write bytes as {{ Tpl }}
func writeField(w *bytes.Buffer, f Field) {
	switch f.Type() {
	case reflect.Int, reflect.Int64, reflect.Bool:
		w.Write([]byte(f.tpl))
	default:
		w.WriteRune('"')
		w.WriteString(f.tpl)
		w.WriteRune('"')
	}
}
