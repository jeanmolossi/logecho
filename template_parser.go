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

func readContext(c echo.Context) []byte {
	poolMutex.Lock()
	buf := new(bytes.Buffer)
	tpl.Execute(buf, getTemplateFields(c))
	poolMutex.Unlock()

	return buf.Bytes()
}

func getFields(fbytes []byte) []zapcore.Field {
	poolMutex.Lock()
	defer poolMutex.Unlock()

	var f map[string]interface{}
	json.Unmarshal(fbytes, &f)

	fields := make([]zapcore.Field, 0, len(f))
	for field, value := range f {
		switch v := value.(type) {
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

func configTemplate(ctxFields Fields, c echo.Context) {
	poolMutex.Lock()
	parseFields(ctxFields, buf)

	tpl = template.Must(
		template.New("fields").
			Funcs(getTemplateFuncMap(c)).
			Parse(buf.String()),
	)
	poolMutex.Unlock()
}

func parseFields(f Fields, w *bytes.Buffer) {
	w.Reset()

	lastIndex := len(f)
	indexCounter := 0

	w.WriteRune('{')
	for key, value := range f {
		indexCounter++

		w.WriteRune('"')
		w.WriteString(key)
		w.WriteString(`":`)
		writeField(w, value)

		if indexCounter != lastIndex {
			w.WriteRune(',')
		}
	}
	w.WriteRune('}')
}

func writeField(w *bytes.Buffer, f Field) {
	switch f.Type() {
	case reflect.Int, reflect.Int64:
		w.Write([]byte(f.tpl))
	default:
		w.WriteRune('"')
		w.WriteString(f.tpl)
		w.WriteRune('"')
	}
}
