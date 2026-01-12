package jsonconv

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

type Person struct {
	FirstName string `json:"first_Name"`
	LastName  string
	Age       int
}

func TestJson(t *testing.T) {
	data := `{
			"first_Name": "John",
			"LastName": "Doe",
			"Age": 30
		}`
	var p *Person
	err := json.Unmarshal([]byte(data), &p)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
	}
	fmt.Println("---", data)
	fmt.Println("---", AnyToJsonIndent(p))

}

func TestJsonToAnyIgnoreCase(t *testing.T) {
	data := `{
			"firstname": "John",
			"last_Name": "Doe",
			"age": 30,
			"me":[1,2,3]
		}`
	type Person struct {
		FirstName string `json:"firstName"`
		LastName  string
		Age       int
		Me        []int
	}

	var p Person
	fmt.Println("---", AnyToJsonIndent(p))

	var it interface{}
	it = &p
	if err := JsonToAnyIgnoreCase(data, it); err != nil {
		t.Errorf("JsonToAnyIgnoreCase() error = %v", err)
	}

	fmt.Println("---", data)
	fmt.Println("---", AnyToJsonIndent(p))
}

func TestJsonToAny(t *testing.T) {
	data := `{
			"first_Name": "John",
			"LastName": "Doe",
			"Age": 30
		}`
	var p Person
	p = JsonToAnyNE[Person](data)
	fmt.Println("---", data)
	fmt.Println("---", AnyToJsonIndent(p))

}

// 转换成下划线json
func TestAnyToJsonCamel(t *testing.T) {
	p := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
	}

	m := AnyToMapNE(p)
	fmt.Println("---", m)

	mm := make(map[string]any)
	for k, v := range m {
		nk := strings.ToLower(ExtractLetters(k))
		mm[nk] = v
	}
	fmt.Println("---", mm)
}
