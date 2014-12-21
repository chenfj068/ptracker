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


