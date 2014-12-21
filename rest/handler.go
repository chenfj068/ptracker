package rest

import ()

type HandleFunc func(*RequestContext, interface{}) (interface{}, error)

//context wraped request and responseWriter
//params  wrap parameters as field
type RequestHandler interface {
	Handler(context *RequestContext, params interface{}) interface{}
}
