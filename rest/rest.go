package rest

import (
	
	"net/http"
)

type ResponseType uint8

const (
	Json      ResponseType = ResponseType(0)
	Xml       ResponseType = ResponseType(1)
	Html      ResponseType = ResponseType(2)
	PlainText ResponseType = ResponseType(3)
	//pb ?
)




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

