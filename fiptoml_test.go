package fiptoml

import (
	"reflect"
	"testing"
	"time"
	"fmt"
)

const (
	example = `#this is an TOML doc
	title = "TOML Example"
	files = [
		"chapter 1",
		"chapter 2"
	]

[owner]
name = "Lance Uppercut"
dob = 1979-05-27T07:32:00-08:00 # First class dates? Why not?
tags = [
		"tag 1",
		"tag 2"
	]

[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
connection_max = 5000
enabled = true

	[[products]]
	name = "Hammer"
	sku = 738594937

	[[products]]

	[[products]]
	name = "Nail"
	sku = 284758393
	color = "gray"
	`
)

func ExampleLoad() {
	toml, err := Load("./config/config.toml")
	if err != nil {
		fmt.Println("ParseFile should work")
		return
	}
	title := toml.GetString("title","")
	datebaseEnabled := toml.GetBool("database.enabled",false)
	fmt.Println(title)
	fmt.Println(datebaseEnabled)
	// Output:
	// TOML Example
	// true
}

func ExampleParseString() {
	toml, err := ParseString(example)
	if err != nil {
		fmt.Println("ParseFile should work")
		return
	}
	owner := toml.GetString("owner.name","")
	fmt.Println(owner)

	ports := toml.GetIntArray("database.ports")
	fmt.Println(ports)
	// Output:
	// Lance Uppercut
	// [8001 8001 8002]
}

func ExampleParse() {
	toml, err := Parse([]byte(example))
	if err != nil {
		fmt.Println("ParseFile should work")
		return
	}
	files := toml.GetStringArray("files")
	fmt.Println(files)

	connectionMax := toml.GetInt("database.connection_max",-1)
	fmt.Println(connectionMax)
	// Output:
	// [chapter 1 chapter 2]
	// 5000
}

func TestWrite(t *testing.T) {
	toml := NewToml()
	toml.SetValue("name", "Chunni")
	toml.SetValue("days", 21)
	toml.SetValue("dob", time.Now())
	toml.SetValue("enabled",true)

	toml.SetValue("files",[]string{"a","b","c"})

	Write(toml,"./config/out.toml")
}

