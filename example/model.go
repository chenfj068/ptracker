package main

import (
	"encoding/json"
)

type User struct {
	Name string `name:"name"`
	Job  string `name:"job"`
	Age  int    `name:"age"`
	Id   int    `name:"id"`
	//PP  string `path:pp`
}

func (u User) String() string {
	b, _ := json.Marshal(u)
	return string(b)
}
