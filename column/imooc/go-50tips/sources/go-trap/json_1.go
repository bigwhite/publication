package main

import (
	"encoding/json"
	"fmt"
)

type person struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
	id     string `json:"id"`
}

func main() {
	p := person{
		Name:   "tony",
		Age:    20,
		Gender: "male",
		id:     "xxx-xxx-xxx-xxx",
	}

	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println("json marshal error:", err)
		return
	}

	fmt.Printf("%s\n", string(b))

	s := `{"name":"tony","age":20,"gender":"male", "id":"xxx-xxx-xxx-xxx"}`
	var p1 person
	err = json.Unmarshal([]byte(s), &p1)
	if err != nil {
		fmt.Println("json unmarshal error:", err)
		return
	}
	fmt.Printf("%#v\n", p1) //main.person{Name:"tony", Age:20, Gender:"male", id:""}
}
