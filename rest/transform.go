package rest

//wrap request param to struct
//translate response object to byte array
import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

var (
	not_found error
)

//request
//target wrapped type
//uri pattern
type RequestWrapper func(*http.Request, reflect.Type, string) interface{}

func defaultWrapper(request *http.Request, paramType reflect.Type, uri string) interface{} {
	request.ParseForm()
	reqUrl := request.URL.Path
	reqUri := request.RequestURI
	idx := strings.IndexByte(reqUri, '?')
	var urlValues url.Values
	if idx > 0 && (idx+1) < len(reqUri) {
		qs := reqUri[idx+1:]
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

		vv := strings.Split(reqUrl, "/")
		cc := strings.Split(uri, "/")
		for i, v := range cc {
			if strings.Contains(v, "{") {
				m[v[1:len(v)-1]] = []string{vv[i]}
			}
		}
		return Params(m)
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

		var pv string
		name := paramType.Field(i).Tag.Get("name")
		pathVariable := paramType.Field(i).Tag.Get("path")
		if name == "" && pathVariable == "" {
			name = paramType.Field(i).Name
		}
		if pathVariable != "" {
			ss := strings.Split(uri, "/")
			idx, sep := 0, ""
			for idx, sep = range ss {
				if sep == "{"+pathVariable+"}" {
					break
				}
			}
			ss2 := strings.Split(reqUrl, "/")
			if idx < len(ss2) {
				pv = ss2[idx]
			}
		} else {
			fmt.Printf("name %s\n", name)
			pv = request.FormValue(name)

			if pv == "" {
				pv = urlValues.Get(name)
			}
			if pv == "" {
				continue
			}
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

type Params map[string][]string

func (p Params) Contains(name string) bool {
	if _, ok := p[name]; ok {
		return true
	}

	return false
}

func (p Params) GetInt(name string) (int64, error) {
	if v, ok := p[name]; ok {
		i, er := strconv.ParseInt(v[0], 10, 64)
		return i, er
	} else {
		return -1, errors.New("param not found:" + name)
	}
}

func (p Params) GetString(name string) (string, error) {
	if v, ok := p[name]; ok {
		return v[0], nil
	} else {
		return "", notFound(name)
	}
}

func (p Params) GetStringArray(name string) ([]string, error) {
	if v, ok := p[name]; ok {
		return v, nil
	} else {
		return nil, notFound(name)
	}
}

func (p Params) GetFloat(name string) (float64, error) {
	if v, ok := p[name]; ok {
		f, er := strconv.ParseFloat(v[0], 64)
		return f, er
	} else {
		return 0.0, notFound(name)
	}
}

func (p Params) GetIntArray(name string) ([]int64, error) {
	if v, ok := p[name]; ok {
		ri := make([]int64, len(v), len(v))
		for i, vi := range v {
			r, er := strconv.ParseInt(vi, 10, 64)
			if er != nil {
				return nil, er
			}
			ri[i] = r
		}
		return ri, nil
	} else {
		return nil, notFound(name)
	}
}

func (p Params) ParamNames() []string {
	names := make([]string, 0, len(p))
	for k, _ := range p {
		names = append(names, k)
	}
	return names
}
func notFound(name string) error {

	return errors.New("param not found:" + name)
}
