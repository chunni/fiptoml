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

func TestExtractTableArray(t *testing.T) {
	toml := newToml()
	input := `products]]
	name = "Hammer"
	sku = 738594937

	[[products]]

	[[products]]
	name = "Nail"
	sku = 284758393
	color = "gray"`
	_, err := extractTableArray([]byte(input), toml)
	if err != nil {
		t.Log("ExtractTableArray should work. err", err)
		t.Fail()
	}

	iterateTomlDoc(toml)
}

func TestExtractTable(t *testing.T) {
	toml := newToml()
	input := `owner]
	name = "Lance Uppercut"
	dob = 1979-05-27T07:32:00-08:00 # First class dates? Why not?
	tags = [
			"tag 1",
			"tag 2"
		]`
	idx, err := extractTable([]byte(input), toml)
	if err != nil {
		t.Log("ExtractTable should work. err", err)
		t.Fail()
	}
	if idx != len(input) {
		t.Log("ExtractTable idx:", idx, "len", len(input))
		t.Fail()
	}
	if len(toml.dict) != 1 {
		t.Log("ExtractTable, should have one item in toml")
		t.Fail()
	}

	table, _ := toml.GetTableToml("owner")
	if table == nil {
		t.Log("ExtractTable, should have one owner item")
		t.Fail()
	}

	name := table.GetString("name", "")
	dob := table.GetDatetime("dob", time.Now())
	tags := table.GetStringArray("tags")

	if name != "Lance Uppercut" ||
		dob.Year() != 1979 ||
		tags[1] != "tag 2" {
		t.Log("ExtractTable, should get name:", name, "dob:", dob, "tags:", tags)
	}

	//iterateTomlDoc(toml,t)
}

func TestExtractKeyValue(t *testing.T) {
	toml := newToml()
	input := `title = "TOML Example" #this is an TOML doc`
	idx, err := extractKeyValue([]byte(input), toml)
	if err != nil {
		t.Log("Extract key value should works. err", err)
		t.Fail()
	}
	if idx != len(input) {
		t.Log("ExtractKeyValue: should skip the whole string, idx:", idx)
		t.Fail()
	}
	if toml.GetString("title", "") != "TOML Example" {
		t.Log("ExtractKeyValue, title should be correct, but it is:", toml.GetString("title", ""))
		t.Fail()
	}

	//iterateTomlDoc(toml,t)
}

func TestExtractKeyValueSection(t *testing.T) {
	toml := newToml()
	input := `title = "TOML Example" #this is an TOML doc
		age = 3`
	idx, err := extractKeyValueSection([]byte(input), toml)
	if err != nil {
		t.Log("Extract key value should works. err", err)
		t.Fail()
	}
	if idx != len(input) {
		t.Log("ExtractKeyValue: should skip the whole string, idx:", idx)
		t.Fail()
	}
	if toml.GetString("title", "") != "TOML Example" || toml.GetInt("age", 0) != 3 {
		t.Log("ExtractKeyValue, title should be correct, but it is:", toml.GetString("title", ""))
		t.Fail()
	}
	//iterateTomlDoc(toml,t)
}

func TestExtractArray(t *testing.T) {
	input := `[1, 2, 3]`
	val, idx, err := extractArray([]byte(input))
	if err != nil {
		t.Log("Extract array should work. err:", err)
		t.Fail()
	}

	if reflect.TypeOf(val).Kind() != reflect.Slice || idx != len(input) {
		t.Log("Extract int array should be ok, val:", val, "idx:", idx)
		t.Fail()
	}

	input = `["ab", "c", "d"]`
	val, idx, err = extractArray([]byte(input))
	if err != nil {
		t.Log("Extract array should work. err:", err)
		t.Fail()
	}
	if reflect.TypeOf(val).Kind() != reflect.Slice || idx != len(input) {
		t.Log("Extract String array should be ok, val:", val, "idx:", idx)
		t.Fail()
	}

	input = `["ab", 1, "d"]`
	_, _, err = extractArray([]byte(input))
	if err == nil {
		t.Log("Extract array should NOT work because of mixed data type")
		t.Fail()
	}
}

func TestExtractStringUntil(t *testing.T) {
	input := `abcde"`
	val, idx := extractStringUntil([]byte(input), []byte{'"'})
	if val != `abcde` || idx != 6 {
		t.Log("Extract string:", val, "idx:", idx)
		t.Fail()
	}

	input = `ab"cde"""`
	val, idx = extractStringUntil([]byte(input), []byte{'"', '"', '"'})
	if val != `ab"cde` || idx != 9 {
		t.Log("Extract string:", val, "idx:", idx)
		t.Fail()
	}
}

func TestExtractString(t *testing.T) {
	input := `"I am here.\n"`
	val, idx, err := extractString([]byte(input))
	if val != "I am here.\n" || idx != len(input) {
		t.Log("Extract single string:", val, "idx:", idx, "err:", err)
		t.Fail()
	}

	input = `"""a: "hi".\n"""`
	val, idx, err = extractString([]byte(input))
	if val != "a: \"hi\".\n" || idx != len(input) {
		t.Log("Extract multi string:", val, "idx:", idx, "err:", err)
		t.Fail()
	}

	input = `"""
One
Two"""`
	val, idx, err = extractString([]byte(input))
	if val != "One\nTwo" || idx != len(input) {
		t.Log("Extract multi string:", val, "idx:", idx, "err:", err)
		t.Fail()
	}

	input = `'this is "2" chars \n'`
	val, idx, err = extractString([]byte(input))
	if val != "this is \"2\" chars \\n" || idx != len(input) {
		t.Log("Extract single literal:", val, "idx:", idx, "err:", err)
		t.Fail()
	}

	input = `'''A:I'm here.
		B:good.\t.H'''`
	lt := `A:I'm here.
		B:good.\t.H`
	val, idx, err = extractString([]byte(input))
	if val != lt || idx != len(input) {
		t.Log("Extract single literal:", val, "idx:", idx, "err:", err)
		t.Fail()
	}
}

func TestExtractBool(t *testing.T) {
	input := `true`
	val, idx, err := extractBool([]byte(input))
	if val != true || idx != len(input) {
		t.Log("Extract bool:", val, "idx:", idx, "err:", err)
		t.Fail()
	}

	input = "false"
	val, idx, err = extractBool([]byte(input))
	if val != false || idx != len(input) {
		t.Log("Extract bool:", val, "idx:", idx, "err:", err)
		t.Fail()
	}

	input = "fals"
	val, idx, err = extractBool([]byte(input))
	if err == nil {
		t.Log("Extract bool:should be error here, err:", err)
		t.Fail()
	}
}

func TestExtractNumber(t *testing.T) {
	input := "-123"
	val, idx, err := extractNumber([]byte(input))
	if val != -123 || idx != len(input) {
		t.Log("Extract int: val", val, "idx:", idx, "err:", err)
		t.Fail()
	}

	input = "1.23"
	val, idx, err = extractNumber([]byte(input))
	if val != 1.23 || idx != len(input) {
		t.Log("Extract float: val", val, "idx:", idx, "", err)
		t.Fail()
	}

	input = "2014-12-04T11:09:30+02:00"
	val, idx, err = extractNumber([]byte(input))
	if err != nil || idx != len(input) {
		t.Log(len(input), time.Date(2014, time.December, 4, 11, 9, 30, 0, time.UTC))
		t.Log("Extract datetime: val", val, "idx:", idx, "err:", err)
		t.Fail()
	}

	input = "2014-12-04 11:09:30.889Z"
	val, idx, err = extractNumber([]byte(input))
	if err == nil {
		t.Log("Extract datetime: invalid", val, "idx:", idx)
		t.Fail()
	}
}

func TestSkipLeft(t *testing.T) {
	input := "  a"
	idx := skipLeft([]byte(input))
	if idx != 2 {
		t.Log("idx =", idx)
		t.Log("SkipLeft: should skip 2")
		t.Fail()
	}

	input = `  #comment
		  a`
	idx = skipLeft([]byte(input))
	if idx != len(input)-1 {
		t.Log("idx =", idx)
		t.Log("SkipLeft: should skip all but one char")
		t.Fail()
	}
}

func TestSkipRight(t *testing.T) {
	input := `  a
		` //contains one space
	idx := skipRight([]byte(input[3:]))
	if idx != 1 {
		t.Log("idx =", idx)
		t.Log("SkipRight: skip 1")
		t.Fail()
	}

	input = ` a = b #c
		  a`
	idx = skipRight([]byte(input[6:]))
	if idx != 4 {
		t.Log("idx =", idx)
		t.Log("SkipLeft: should skip 4")
		t.Fail()
	}
}

func TestSkipUntil(t *testing.T) {
	f := func(r rune) bool {
		if r == '#' {
			return true
		} else {
			return false
		}
	}
	idx := skipUntil([]byte(`abc#di`), f, true)
	if idx != 4 {
		t.Log("idx =", idx)
		t.Log("should skip 4 ")
		t.Fail()
	}

	idx = skipUntil([]byte(`abcdi`), f, false)
	if idx != 5 {
		t.Log("idx =", idx)
		t.Log("should skip 5 char")
		t.Fail()
	}
}

func TestSkipIf(t *testing.T) {
	f := func(r rune) bool {
		if r == 'a' {
			return true
		} else {
			return false
		}
	}

	idx := skipIf([]byte(`aac#di`), f)
	if idx != 2 {
		t.Log("idx =", idx)
		t.Log("should skip 2 char")
		t.Fail()
	}
}

func TestSkipComment(t *testing.T) {
	input := `#this is a comment`
	idx := skipComments([]byte(input))
	if idx != len(input) {
		t.Log("idx =", idx)
		t.Log("should skip the whole comment")
		t.Fail()
	}
}

func TestIsSectionEnd(t *testing.T) {
	input := `
		`
	isEnd, idx := isSectionEnd([]byte(input))
	if !isEnd || idx != 1 {
		t.Log(len(input))
		t.Log("isEnd = ", isEnd, "idx = ", idx)
		t.Log("should consider as section end, and skip all bytes")
		t.Fail()
	}
	input = ` #comments here
`
	isEnd, idx = isSectionEnd([]byte(input))
	if !isEnd || idx != len(input) {
		t.Log("isEnd = ", isEnd, "idx = ", idx, "len:", len(input))
		t.Log("should consider as section end when there is comment, and skip all bytes")
		t.Fail()
	}

	input = `  abc`
	isEnd, idx = isSectionEnd([]byte(input))
	if isEnd || idx != 2 {
		t.Log("isEnd = ", isEnd, "idx = ", idx)
		t.Log("should consider as section end, and skip all bytes")
		t.Fail()
	}
}

func iterateTomlDoc(doc *toml) {
	dict := doc.dict
	count := len(dict)
	fmt.Println("length of toml:", count)
	for key := range dict {
		switch val := dict[key].(type) {
		case nil:
			fmt.Println("empty value of: ", key)
		case string:
			fmt.Println(key, "=", val,"(string)")
		case *toml:
			fmt.Println("Table:", key)
			iterateTomlDoc(val)
		case []*toml:
			fmt.Println("Array of tables:", key)
			for i, v := range val {
				fmt.Println("Table", i)
				iterateTomlDoc(v)
			}
		default:
			fmt.Println(key, "=", val,"(", reflect.TypeOf(val),")")
		}
	}
}
