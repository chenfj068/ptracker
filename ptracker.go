package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func main() {
	testMatch()
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServe(":8880", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func testMatch() {

	url, uri := "a/b/{d}/g", "a/b/d/g"
	fmt.Printf("%v\n", match(url, uri))
	url,uri="a/b/*/g","a/b/d/f/g"
	fmt.Printf("%v\n", match(url, uri))
	url,uri="a/b/*/c/g","a/b/d/g"
	fmt.Printf("%v\n", match(url, uri))

}
func match(url, uri string) bool {
	reqParts := strings.Split(uri, "/")
	patParts := strings.Split(url, "/")
	fmt.Printf("%v\n,%v\n",reqParts,patParts)
	i, j := 0, 0
	for i < len(reqParts) && j < len(patParts) {
		if reqParts[i] == patParts[j] {
			i++
			j++
		} else if patParts[j] == "*" {
			if(j+1<len(patParts)&&patParts[j+1]==reqParts[i]){
				j++
				j++
			}
			i++
			
		} else if strings.Contains(patParts[j],"{") {
			i++
			j++
		} else {
			break
		}
	}
	fmt.Printf("%d,%d\n",i,j)
	if i == len(reqParts) && j == len(patParts) {
		return true
	}

	return false

}

type DispatcherHandler struct {
}

type HttpDispatcher interface {
	Dispatch(http.ResponseWriter, *http.Request)
	RegisterHandler(handler *RequestHandler, url string, method ...string)
	RegisterHandler2(handler *RequestHandler, extractor ParamsExtractor, url string, method ...string)
}

type MyDispatcher struct {
}

func (disp *MyDispatcher) Dispatch(http.ResponseWriter, *http.Request) {

}
func (disp *MyDispatcher) RegisterHandler(handler *RequestHandler, url string, method ...string) {

}
func (disp *MyDispatcher) RegisterHandler2(handler *RequestHandler, extractor ParamsExtractor, url string, method ...string) {

}

type MyParamsExtractor struct {
}

func (e *MyParamsExtractor) Extract(request *http.Request) interface{} {
	return nil
}

func tt() {
	var disp HttpDispatcher
	disp = &MyDispatcher{}
	disp.RegisterHandler2(nil, &MyParamsExtractor{}, "/users/id", "get", "put")
	disp.RegisterHandler(nil, "", "get", "put")
}

type ParamsExtractor interface {
	Extract(request *http.Request) interface{}
}
type RequestContext struct {
	Request  *http.Request
	Response *http.ResponseWriter
}

//context wraped request and responseWriter
//params  wrap parameters as field
type RequestHandler interface {
	Handler(context *RequestContext, params interface{})
}
