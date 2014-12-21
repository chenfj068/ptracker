package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

type ResponseType uint8

const (
	Json      ResponseType = ResponseType(0)
	Xml       ResponseType = ResponseType(1)
	Html      ResponseType = ResponseType(2)
	PlainText ResponseType = ResponseType(3)
	//pb ?
)

type HttpDispatcher interface {
	Dispatch(http.ResponseWriter, *http.Request)
	//map a handler to specified url pattern
	//handler  concrete handler object,not nil
	//paramType param wrap type,if nil then map
	//url     url pattern
	//method  http methods
	RegisterHandler(handler RequestHandler, paramType reflect.Type, rspType ResponseType, url string, method ...string)
	RegisterHandlerWithParser(handler RequestHandler, paramType reflect.Type, rspType ResponseType, parser ParamsParseFunc, url string, method ...string)
}

type SimpleDispatcher struct {
	handlerMap map[string]handlerSpec
}

func NewDispatcher() HttpDispatcher {
	return &SimpleDispatcher{make(map[string]handlerSpec)}
}
func (disp *SimpleDispatcher) Dispatch(respWriter http.ResponseWriter, req *http.Request) {
	//find the target handler
	spec, ok := disp.getHandler(req)
	if ok {
		http.NotFound(respWriter, req)
		return //return 404
	}
	v := spec.reqParseFunc(req, spec.reqType, spec.urlPattern) //parse request struct
	ctx := NewRequestContext(req, respWriter)
	result := spec.handler.Handler(ctx, v)
	switch spec.rspType {
	case Json:
		b, _ := json.Marshal(result)
		respWriter.Header().Set("content-type", "application/json")
		respWriter.Write(b)
	case Xml:

	case Html:

	case PlainText:
		ss := fmt.Sprintf("%v", result)
		respWriter.Header().Set("content-type", "text/plain")
		respWriter.Write([]byte(ss))
	}

}

//find the mapped handler to the request
func (disp *SimpleDispatcher) getHandler(req *http.Request) (handlerSpec, bool) {
	uri := req.RequestURI
	for pattern, spec := range disp.handlerMap {
		if match(pattern, uri) {
			return spec, true
		}
	}
	return handlerSpec{}, false
}

type handlerSpec struct {
	handler      RequestHandler  //handler
	reqType      reflect.Type    //wrapped request type
	rspType      ResponseType    //reponse type
	urlPattern   string          //url pattern
	httpMethods  []string        //http methods
	reqParseFunc ParamsParseFunc //request parameter parse function
}

//body parser
//form parser
func defaultParamParseFunc(request *http.Request, paramType reflect.Type, urlPattern string) interface{} {
	request.ParseForm()
	if paramType == nil {
		m := make(map[string][]string)
		values := request.Form
		for name, value := range values {
			m[name] = value
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
		pv := request.FormValue(name)
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

	return jp
}
func (disp *SimpleDispatcher) RegisterHandler(handler RequestHandler, paramType reflect.Type, rspType ResponseType, url string, method ...string) {
	hdSpec := handlerSpec{handler: handler,
		reqType:      paramType,
		rspType:      rspType,
		urlPattern:   url,
		reqParseFunc: defaultParamParseFunc,
		httpMethods:  method,
	}
	disp.handlerMap[url] = hdSpec

}
func (disp *SimpleDispatcher) RegisterHandlerWithParser(handler RequestHandler, paramType reflect.Type, rspType ResponseType, parser ParamsParseFunc, url string, method ...string) {
	hdSpec := handlerSpec{handler: handler,
		reqType:      paramType,
		rspType:      rspType,
		urlPattern:   url,
		reqParseFunc: parser,
		httpMethods:  method,
	}
	disp.handlerMap[url] = hdSpec
}

//httprequest
//parse return struct type
//url pattern
type ParamsParseFunc func(*http.Request, reflect.Type, string) interface{}

func MyParseFunc(request *http.Request, paramType reflect.Type, url string) interface{} {

	return nil
}

type RequestContext struct {
	Request  *http.Request
	Response http.ResponseWriter
	valueMap map[string]interface{}
}

func NewRequestContext(request *http.Request, respWriter http.ResponseWriter) *RequestContext {
	return &RequestContext{request, respWriter, make(map[string]interface{})}
}

func (context *RequestContext) SetObject(key string, value interface{}) {
	context.valueMap[key] = value
}
func (context *RequestContext) getObject(key string) interface{} {
	return context.valueMap[key]
}

type HelloHandler struct{}

func (h *HelloHandler) Handler(context *RequestContext, params interface{}) interface{} {

	return ""
}

type MyModel struct {
	Age  string `path:"age"`
	Name string `name:"name"`
	//uri pattern
	// /users/{age}?name=name1
	// url match handler how to?
}

func tt() {
	var disp HttpDispatcher
	disp = NewDispatcher()
	h := &HelloHandler{}
	t := reflect.TypeOf(MyModel{})
	disp.RegisterHandlerWithParser(h, t, Json, MyParseFunc, "/users/id", "get", "put")
	disp.RegisterHandler(h, t, PlainText, "", "get", "put")
}
