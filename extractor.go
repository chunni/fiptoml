package fiptoml

import (
	"unicode/utf8"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func extractTableArray(input []byte, doc *toml) (idx int, err error) {
	name, idx, err := extractTableName(input, true)
	if err != nil {
		return
	}

	//array, err :=doc.GetTableArray(name)
	subDoc := newToml()
	delta := 0

	switch array := doc.dict[name].(type) {
	case nil:
		doc.dict[name] = []*toml{subDoc}
	case []*toml:
		doc.dict[name] = append(array, subDoc)
	default:
		goto DupKey
	}

	delta, err = extractKeyValueSection(input[idx:], subDoc)
	idx += delta
	return

DupKey:
	err = errDuplicatedKey(name)
	return
}

func extractTable(input []byte, doc *toml) (idx int, err error) {
	name, idx, err := extractTableName(input, false)

	delta := 0
	if err != nil {
		return
	}

	var subDoc *toml
	switch v := doc.dict[name].(type) {
	case nil:
		subDoc = newToml()
		doc.dict[name] = subDoc
	case *toml:
		subDoc = v
	default:
		goto DupKey
	}
	delta, err = extractKeyValueSection(input[idx:], subDoc)
	idx += delta
	return

DupKey:
	err = errDuplicatedKey(name)
	return
}

func extractKeyValueSection(input []byte, doc *toml) (idx int, err error) {
	for idx < len(input) {
		shouldEnd, delta := isSectionEnd(input[idx:])
		idx += delta
		if shouldEnd {
			break
		} else {
			idx += skipLeft(input[idx:])
			delta, err = extractKeyValue(input[idx:], doc)
			idx += delta
		}
	}
	return
}

func extractKeyValue(input []byte, doc *toml) (idx int, err error) {
	key, idx := extractKey(input)
	if idx == 0 {
		err = errEmptyKey
	} else if doc.dict[key] != nil {
		err = errDuplicatedKey(key)
	} else {
		idx += skipSpaceAndEquals(input[idx:])
		val, delta, err := extractValue(input[idx:])
		if err != nil {
			return idx, err
		}
		idx += delta
		idx += skipRight(input[idx:])
		if val == nil {
			err = errUnsupportedValue(key)
			return idx, err
		}
		doc.dict[key] = val
	}
	return
}

func extractKey(input []byte) (key string, idx int) {
	i := 0
L:
	for i < len(input) {
		r, w := utf8.DecodeRune(input[i:])
		switch r {
		case ' ', '\t', '=':
			idx = i
			key = string(input[:i])
			break L
		default:
			i += w
		}
	}
	return
}

func extractValue(input []byte) (val interface{}, idx int, err error) {
	switch input[0] {
	case '"', '\'':
		val, idx, err = extractString(input)
	case 't', 'f':
		val, idx, err = extractBool(input)
	case '[':
		val, idx, err = extractArray(input)
	case '+', '-', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		val, idx, err = extractNumber(input)
	default:
		val = nil
	}
	return
}

//single/multiple line string/literal
func extractString(input []byte) (val string, idx int, err error) {
	length := len(input)
	if input[0] == '"' {
		if length > 3 && input[1] == '"' && input[2] == '"' {
			val, idx, err = extractMultiString(input)
		} else {
			val, idx, err = extractSingleString(input)
		}
	} else if input[0] == '\'' {
		if length > 3 && input[1] == '\'' && input[2] == '\'' {
			val, idx = extractStringUntil(input[3:], []byte{'\'', '\'', '\''})
			idx += 3
		} else {
			val, idx = extractStringUntil(input[1:], []byte{'\''})
			idx += 1
		}
	} else {
		err = errStringSyntaxError
	}
	return

}

func unquote(in string) (out string, err error) {
	if len(in) == 0 {
		return
	}

	in = fmt.Sprintf("%c%s%c", '"', in, '"')
	out, err = strconv.Unquote(in)
	return
}

func extractMultiString(input []byte) (val string, idx int, err error) {
	val, idx = extractStringUntil(input[3:], []byte{'"', '"', '"'})
	idx += 3
	val = multiLineSkipR.ReplaceAllString(val, "")
	val = quoteLineR.ReplaceAllString(val, `\n`)
	val = strings.Replace(val, `"`, `\"`, -1)
	val, err = unquote(val)
	return
}

func extractSingleString(input []byte) (val string, idx int, err error) {
	val, idx = extractStringUntil(input[1:], []byte{'"'})
	val, err = unquote(val)
	idx += 1
	return
}

func extractStringUntil(input []byte, arr []byte) (val string, idx int) {
	from := skipLeft(input)

	delta := skipUntilArray(input[from:], arr)
	idx = from + delta //+ 2 add quotes before and after
	val = string(input[from:idx])

	idx += len(arr)
	return
}

func extractBool(input []byte) (val bool, idx int, err error) {
	switch input[0] {
	case 't':
		if len(input) > 3 && sliceEquals(input[1:4], []byte{'r', 'u', 'e'}) {
			val = true
			idx = 4
		} else {
			goto ErrBool
		}
	case 'f':
		if len(input) > 4 && sliceEquals(input[1:5], []byte{'a', 'l', 's', 'e'}) {
			val = false
			idx = 5
		} else {
			goto ErrBool
		}
	default:
		goto ErrBool
	}
	return
ErrBool:
	err = errBool
	return
}

//for int, float, datetime
func extractNumber(input []byte) (val interface{}, idx int, err error) {
	idx = skipUntilSpace(input)
	str := string(input[0:idx])
	if datatimeR.MatchString(str) {
		val, err = time.Parse(time.RFC3339, str)
	} else if intR.MatchString(str) {
		val, err = strconv.Atoi(str)
	} else if floatR.MatchString(str) {
		val, err = strconv.ParseFloat(str, 64)
	} else {
		err = errNumber
	}
	return
}

func extractArray(input []byte) (val interface{}, idx int, err error) {
	from := 1 + skipLeft(input[1:])
	idx = from + skipUntilChar(input[from:], ']')
	str := string(input[from:idx])
	idx += 1
	vals := strings.Split(str, ",")

	length := len(vals)
	if length < 1 {
		return
	}

	val0, _, err := extractValue([]byte(strings.TrimSpace(vals[0])))
	switch v0 := val0.(type) {
	case string:
		arr := make([]string, length, length)
		arr[0] = v0
		for i := 1; i < length; i++ {
			v, _, err := extractString([]byte(strings.TrimSpace(vals[i])))
			if err != nil {
				goto ErrArray
			}
			arr[i] = v
		}
		val = arr
	case int:
		arr := make([]int, length, length)
		arr[0] = v0
		for i := 1; i < length; i++ {
			v, err := strconv.Atoi(strings.TrimSpace(vals[i]))
			if err != nil {
				goto ErrArray
			}
			arr[i] = v
		}
		val = arr
	case float64:
		arr := make([]float64, length, length)
		arr[0] = v0
		for i := 1; i < length; i++ {
			v, err := strconv.ParseFloat(strings.TrimSpace(vals[i]), 64)
			if err != nil {
				goto ErrArray
			}
			arr[i] = v
		}
		val = arr
	case bool:
		arr := make([]bool, length, length)
		arr[0] = v0
		for i := 1; i < length; i++ {
			v, err := strconv.ParseBool(strings.TrimSpace(vals[i]))
			if err != nil {
				goto ErrArray
			}
			arr[i] = v
		}
		val = arr
	case time.Time:
		arr := make([]time.Time, length, length)
		arr[0] = v0
		for i := 1; i < length; i++ {
			v, err := time.Parse(time.RFC3339, strings.TrimSpace(vals[i]))
			if err != nil {
				goto ErrArray
			}
			arr[i] = v
		}
		val = arr
	}

	return
ErrArray:
	err = errArray
	return
}

func extractTableName(input []byte, isArray bool) (name string, idx int, err error) {
	//todo: build the whole thing for multiple ...
	i := 0
L:
	for i < len(input) {
		r, w := utf8.DecodeRune(input[i:])
		switch r {
		case ']':
			if i == 0 {
				goto InvalidKey
			}
			r, w = utf8.DecodeRune(input[i+w:])
			if (isArray && r == ']') ||
					(!isArray && (unicode.IsSpace(r)) || r == '#') {
				idx = i + w
				name = string(input[:i])
				break L
			} else {
				goto InvalidKey
			}
		case utf8.RuneError:
			goto InvalidUtf8
		case ' ', '\t', '\n', '\f', '\r':
			goto InvalidKey
		default:
			i += w
		}
	}
	idx += skipRight(input[idx:])

	return

InvalidUtf8:
	err = errUtf8
	return

InvalidKey:
	err = errInvalidTableKey
	return
}

func getTableValue(input []byte) (err error) {
	length := len(input)
	for i := 0; i < length; i++ {
		if input[i] == '\n' {
			//if()
		}
	}
	return
}

func isLineEnd(r rune) bool {
	switch r {
	case '\n', '\r', '\f':
		return true
	default:
		return false
	}
}

func isSpace(r rune) bool {
	switch r {
	case ' ', '\t':
		return true
	default:
		return false
	}
}

func isStringEnd(r rune) bool {
	if r == '"' {
		return true
	} else {
		return false
	}
}

func isSectionEnd(input []byte) (bool, int) {
	i := 0
	for i < len(input) {
		r, w := utf8.DecodeRune(input[i:])
		switch r {
		case '\n', '\r', '\f':
			return true, i + w
		case ' ', '\t':
			i += w
		case '#':
			i += skipComments(input[i:])
		default:
			return false, i
		}
	}
	return true, i
}

//Skip comments and space and new line before a key
func skipLeft(input []byte) (skip int) {
	i := 0
	for i < len(input) {
		r, w := utf8.DecodeRune(input[i:])
		if r == '#' {
			i += skipComments(input[i:])
		} else if unicode.IsSpace(r) {
			i += w
		} else {
			return i
		}
	}
	return i
}

//skip comments and spaces, and the line break after a value
func skipRight(input []byte) (skip int) {
	return skipUntil(input, isLineEnd, true)
}

func skipComments(input []byte) int {
	return skipUntil(input, isLineEnd, false)
}

func skipSpaceAndEquals(input []byte) int {
	return skipIf(input, func(r rune) bool {
			switch r {
			case ' ', '\t', '=':
				return true
			default:
				return false
			}
		})
}

func skipUntilSpace(input []byte) int {
	return skipUntil(input, func(r rune) bool {
			if r == '#' || unicode.IsSpace(r) {
				return true
			} else {
				return false
			}
		}, false)
}

func skipIf(input []byte, f func(rune) bool) int {
	i := 0
	for i < len(input) {
		r, w := utf8.DecodeRune(input[i:])
		if !f(r) {
			return i + w - 1
		} else {
			i += w
		}
	}
	return i
}

func skipUntil(input []byte, f func(rune) bool, include bool) int {
	i := 0
	for i < len(input) {
		r, w := utf8.DecodeRune(input[i:])
		if f(r) {
			if include {
				return i + w
			}
			return i
		} else {
			i += w
		}
	}
	return i
}

func sliceEquals(s1 []byte, s2 []byte) bool {
	length := len(s1)
	if len(s2) != length {
		return false
	}
	for i := 0; i < length; i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func skipUntilChar(input []byte, char rune) int {
	i := 0
	for i < len(input) {
		r, w := utf8.DecodeRune(input[i:])
		if r == char {
			return i
		} else {
			i += w
		}
	}
	return i
}

func skipUntilArray(input []byte, arr []byte) int {
	i := 0
	length := len(input)
	arrLength := len(arr)
	for ; i < length; i++ {
		if input[i] == arr[0] {
			if i+arrLength <= length &&
					sliceEquals(input[i+1:i+arrLength], arr[1:]) {
				return i
			}
		}
	}
	return i
}
