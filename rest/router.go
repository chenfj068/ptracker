package rest

//request , handler manager
import (
	"fmt"
	"reflect"
)


type Router struct {
	routeMap map[string]*Route
}

type Route struct {
	url            string
	methods        []string
	consume        string
	produce        string
	requestType    reflect.Type
	requestWrapper RequestWrapper
	handlerFunc    HandleFunc
	handler        interface{}
	handlerMethod  string
}

func NewRoute() *Route {

	return &Route{}
}

//target url pattern
//a/b/c  a/b/*/c/*/d
func (r *Route) Url(url string) *Route {
	r.url = url
	return r
}

//target method
//accept all methods if not set
func (r *Route) Methods(methods []string) *Route {
	r.methods = methods
	return r

}

//target mime type
//eg application/json
func (r *Route) Consume(mime string) *Route {
	r.consume = mime
	return r
}

//reponse mime type
//eg application/json
func (r *Route) Produce(mime string) *Route {
	r.produce = mime
	return r
}

//request param wraped type
func (r *Route) RequestType(t reflect.Type) *Route {
	r.requestType = t
	return r
}

//customized param wrapper
func (r *Route) RequestWrapper(wrapper RequestWrapper) *Route {
	r.requestWrapper = wrapper
	return r
}

//handler must implement
//Handler(context *RequestContext, params interface{}) interface{},error
func (r *Route) Handler(h interface{}, method string) *Route {
	r.handler = h
	r.handlerMethod = method
	return r
}

func (r *Route) Func(h HandleFunc) *Route {
	r.handlerFunc = h
	return r
}

func NewRouter() *Router {
	return &Router{make(map[string]*Route)}
}

func (r *Router) getRoute(url string) (*Route, bool) {

	for p, route := range r.routeMap {
		fmt.Printf("%s,%s\n", p, url)
		if match(p, url) {
			return route, true
		}
	}
	return nil, false
}
func (r *Router) AddRoute(route *Route) {
	r.routeMap[route.url] = route
}

