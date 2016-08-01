// Copyright 2016 José Santos <henrique_1609@me.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package jet

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/url"
	"reflect"
	"strings"
	"text/template"
)

var defaultExtensions = []string{
	".html.jet",
	".jet.html",
	".jet",
}

var defaultVariables map[string]reflect.Value

func init() {
	defaultVariables = map[string]reflect.Value{
		"lower":     reflect.ValueOf(strings.ToLower),
		"upper":     reflect.ValueOf(strings.ToUpper),
		"hasPrefix": reflect.ValueOf(strings.HasPrefix),
		"hasSuffix": reflect.ValueOf(strings.HasSuffix),
		"repeat":    reflect.ValueOf(strings.Repeat),
		"replace":   reflect.ValueOf(strings.Replace),
		"split":     reflect.ValueOf(strings.Split),
		"trimSpace": reflect.ValueOf(strings.TrimSpace),
		"map":       reflect.ValueOf(newMap),
		"html":      reflect.ValueOf(html.EscapeString),
		"url":       reflect.ValueOf(url.QueryEscape),
		"safeHtml":  reflect.ValueOf(SafeWriter(template.HTMLEscape)),
		"safeJs":    reflect.ValueOf(SafeWriter(template.JSEscape)),
		"raw":       reflect.ValueOf(SafeWriter(unsafePrinter)),
		"unsafe":    reflect.ValueOf(SafeWriter(unsafePrinter)),
		"writeJson": reflect.ValueOf(jsonRenderer),
		"json":      reflect.ValueOf(json.Marshal),
		"isset": reflect.ValueOf(Func(func(a Arguments) reflect.Value {
			a.RequireNumOfArguments("isset", 1, 99999999999)
			for i := 0; i < len(a.argExpr); i++ {
				if !a.runtime.isSet(a.argExpr[i]) {
					return valueBoolFALSE
				}
			}
			return valueBoolTRUE
		})),
		"len": reflect.ValueOf(Func(func(a Arguments) reflect.Value {
			a.RequireNumOfArguments("len", 1, 1)

			expression := a.Get(0)
			if expression.Kind() == reflect.Ptr {
				expression = expression.Elem()
			}

			switch expression.Kind() {
			case reflect.Array, reflect.Chan, reflect.Slice, reflect.Map, reflect.String:
				return reflect.ValueOf(expression.Len())
			case reflect.Struct:
				return reflect.ValueOf(expression.NumField())
			}

			a.Panicf("inválid value type %s in len builtin", expression.Type())
			return reflect.Value{}
		})),
	}

}

func jsonRenderer(v interface{}) RendererFunc {
	return func(r *Runtime) {
		err := json.NewEncoder(r.Writer).Encode(v)
		if err != nil {
			panic(err)
		}
	}
}

func unsafePrinter(w io.Writer, b []byte) {
	w.Write(b)
}

// SafeWriter escapee func, func implementing this type will write direct in the writer skipping the escape
// faze, use this type to create special types of escapee funcs
type SafeWriter func(io.Writer, []byte)

func newMap(values ...interface{}) (nmap map[string]interface{}) {
	if len(values)%2 > 0 {
		panic("new map: invalid number of arguments on call to map")
	}
	nmap = make(map[string]interface{})

	for i := 0; i < len(values); i += 2 {
		nmap[fmt.Sprint(values[i])] = values[i+1]
	}
	return
}
