package main

import (
	"fmt"
	"github.com/cjf068/ptracker/rest"
	"reflect"
)



func main(){
	
	app:=rest.NewApp(8080,"0.0.0.0","/")
	
	route:=rest.NewRoute().Func(getUser)
	route.Url("/users/byid/{id}").Produce("application/json").Methods([]string{"get"})
	app.AddRoute(route)
	
	croute:=rest.NewRoute().Url("/users/add").Methods([]string{"get","post"})
	croute.Produce("text/plain").RequestType(reflect.TypeOf(User{})).Handler(UHandler{},"CreateUser")
	
	app.AddRoute(croute)
	app.Start()
	}


func getUser(ctx *rest.RequestContext,v  interface{})(interface{},error){
	m:=v.(rest.Params)
	id,_:=m.GetInt("id")
	return User{Name:"chen",Job:"programer",Age:10,Id:int(id)},nil
}


type UHandler struct{
	
	
}

func (h UHandler)CreateUser(ctx *rest.RequestContext,v interface{})(interface{},error){
	fmt.Printf("%s\n",v)
	return "ok",nil
}