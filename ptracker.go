package main

import (
	"encoding/json"
	"github.com/cjf068/ptracker/rest"
	"reflect"
	"fmt"
)

func main() {
	app := rest.NewApp(8800, "0.0.0.0", "/")
	rout := rest.NewRoute().Consume("application/json").Methods([]string{"get", "post"}).Produce("application/json")
	rout.Url("/p/lxy/{id}").RequestType(reflect.TypeOf(Model{})).Func(Echo)
	rout2:=rest.NewRoute().Consume("application/json").Handler(MyHandler{},"Handle").Produce("application/json").Url("/p2/lxy/{id}").RequestType(reflect.TypeOf(Model2{}))
	app.AddRoute(rout)
	app.AddRoute(rout2)
	app.Start()

}

type MyHandler struct{
	
}

type Model2 struct{
	Name  string `name:"name"`
	Value int64  `name:"value"`
	Id    int64  `path:"id"`
}

func (h MyHandler)Handle(ctx *rest.RequestContext,v interface{})(interface{},error){
	
	fmt.Printf("%s\n",v)
	m:=make(map[string]string)
	m["lvy"]="kk"
	return m,nil
	
}
func Echo(ctx *rest.RequestContext,v interface{})(interface{},error){
	
	fmt.Printf("%s\n",v)
	m:=make(map[string]string)
	m["lvy"]="kk"
	return m,nil
	
	}

type Model struct {
	Name  string `name:"name"`
	Value int64  `name:"value"`
	Id    int64  `path:"id"`
}

func (m Model) String() string {

	b, _ := json.Marshal(m)
	return string(b)
}

func (m Model2) String() string {

	b, _ := json.Marshal(m)
	return string(b)
}


