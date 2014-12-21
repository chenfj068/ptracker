package main

import (
	"encoding/json"
	"github.com/cjf068/ptracker/rest"
	"net/http"
	"reflect"
	"fmt"
)

func main() {
	router := rest.NewRouter()
	app := rest.NewApp(8800, "0.0.0.0", "/", router)
	rout := rest.NewRoute().Consume("application/json").Methods([]string{"get", "post"}).Produce("application/json")
	rout.Url("/p/lxy").RequestType(reflect.TypeOf(Model{})).Func(Echo)
	app.AddRoute(rout)
	app.Start()

}

func Echo(ctx *rest.RequestContext,v interface{})(interface{},error){
	
	fmt.Printf("%s\n",v)
	m:=make(map[string]string)
	m["lvy"]="love"
	return m,nil
	
	}

type Model struct {
	Name  string `name:"name"`
	Value int64  `name:"value"`
	Id    int64  `name:"id"`
}

func (m Model) String() string {

	b, _ := json.Marshal(m)
	return string(b)
}

type UserHandler struct {
}

func (h *UserHandler) Handler(context *rest.RequestContext, params interface{}) interface{} {

	return nil
}

type ServerHandler struct {
	disp rest.HttpDispatcher
}

func (h ServerHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	h.disp.Dispatch(writer, req)

}
func NewServerHandler() http.Handler {
	disp := rest.NewDispatcher()
	h := ServerHandler{disp}
	//func (disp *SimpleDispatcher) RegisterHandler(handler RequestHandler, paramType reflect.Type, rspType ResponseType, url string, method ...string) {
	//disp.RegisterHandler()
	return h
}
