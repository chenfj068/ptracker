package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
	"fmt"
)

type RestApp struct {
	port       int
	listenAddr string
	rootDir    string
	*Router
}

func NewApp(port int, listenAddr, rootdir string) *RestApp {
	return &RestApp{port, listenAddr, rootdir, NewRouter()}
}

func (app *RestApp) Start() error {
	add:=app.listenAddr+":"+fmt.Sprintf("%d",app.port)
	fmt.Printf("addr :%s\n",add)
	err := http.ListenAndServe(add, app)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return nil
}

func (disp *RestApp) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Printf("%s\n","servehttp")
	uri := req.RequestURI
	if strings.Contains(uri, "?") {
		uri = strings.Split(uri, "?")[0]
	}
	route, ok := disp.getRoute(uri)
	if !ok {
		http.NotFound(writer, req)
		return
	}
	ctx := NewRequestContext(req, writer)
	var model interface{}
	if route.requestWrapper != nil {
		model = route.requestWrapper(req, route.requestType, route.url)
	} else {
		model = defaultWrapper(req, route.requestType, route.url)
	}
	if route.handlerFunc != nil {
		rv, _ := route.handlerFunc(ctx, model)
		b, _ := json.Marshal(rv)
		writer.Write(b)
	} else {
		m, ok := reflect.TypeOf(route.handler).MethodByName(route.handlerMethod)
		if ok {
			rv := m.Func.Call([]reflect.Value{reflect.ValueOf(route.handler), reflect.ValueOf(ctx), reflect.ValueOf(model)})
			vv := rv[0].Interface()
			b, _ := json.Marshal(vv)
			writer.Write(b)
		} else {
			http.NotFound(writer, req)
		}
	}
}
