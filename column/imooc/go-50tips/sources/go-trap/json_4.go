package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type foo struct {
	Name string
	Age  int
}

func fixedJsonUnmarshal() {
	s := `{"age": 23, "name": "tony"}`
	var f foo
	err := json.Unmarshal([]byte(s), &f)
	if err != nil {
		fmt.Println("json unmarshal error:", err)
		return
	}

	fmt.Printf("%#v\n", f)
}

func flexibleJsonUnmarshal() {
	s := `{"age": 23, "name": "tony"}`
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		fmt.Println("json unmarshal error:", err)
		return
	}

	age := m["age"].(int)
	fmt.Println("age =", age)
}

func flexibleJsonUnmarshalImproved() {

	s := `{"age": 23, "name": "tony"}`
	m := map[string]interface{}{}

	d := json.NewDecoder(strings.NewReader(s))
	d.UseNumber()

	err := d.Decode(&m)
	if err != nil {
		fmt.Println("json decode error:", err)
		return
	}

	age, err := m["age"].(json.Number).Int64()
	if err != nil {
		fmt.Println("json.Number.Int64 error:", err)
		return
	}
	fmt.Println("age =", age)
}

func main() {
	fixedJsonUnmarshal()
	flexibleJsonUnmarshalImproved()
}
