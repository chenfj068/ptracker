package rest

//wrap request param to struct
//translate response object to byte array
import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

//request
//uri pattern
type RequestWrapper func(*http.Request, reflect.Type, string) interface{}

func defaultWrapper(request *http.Request, paramType reflect.Type, uri string) interface{} {
	request.ParseForm()
	idx := strings.IndexByte(uri, '?')
	var urlValues url.Values
	if idx > 0 && (idx+1) < len(uri) {
		qs := uri[idx+1:]
		urlValues, _ = url.ParseQuery(qs)
	}
	if paramType == nil {
		m := make(map[string][]string)
		values := request.Form
		for name, value := range values {
			m[name] = value
		}

		if urlValues != nil {
			for name, value := range urlValues {
				m[name] = value
			}
		}
		return m
	}
	if request.Header.Get("content-type") == "application/json" {
		decoder := json.NewDecoder(request.Body)
		t := reflect.New(paramType)
		decoder.Decode(&t)
		return t

	}
	jp := reflect.New(paramType)
	jl := jp.Elem()
	fn := paramType.NumField()
	for i := 0; i < fn; i++ {
		v := jl.Field(i)
		name := paramType.Field(i).Tag.Get("name")
		if name == "" {
			name = paramType.Field(i).Name
		}
		fmt.Printf("name %s\n", name)
		pv := request.FormValue(name)

		if pv == "" {
			pv = urlValues.Get(name)
		}
		if pv == "" {
			continue
		}
		switch v.Type().Kind() {
		case reflect.Int:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Int8:
			i, _ := strconv.Atoi(pv)
			v.SetInt(int64(i))
		case reflect.Float32:
		case reflect.Float64:
			f, _ := strconv.ParseFloat(pv, 64)
			v.SetFloat(f)
		case reflect.String:
			v.SetString(pv)
		case reflect.Bool:
			b, _ := strconv.ParseBool(pv)
			v.SetBool(b)
		}
	}

	return jp.Interface()
}
