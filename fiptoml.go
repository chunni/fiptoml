/*
Package fiptoml provides fast, reliable and easy-to-use parser for TOML doc.
 */
package fiptoml

import (
	"io/ioutil"
	"unicode/utf8"
	"bufio"
	"os"
)

func Parse(input []byte) (doc *Toml, err error) {
	doc = NewToml()
	idx, delta := 0, 0
	l := len(input)
	for idx < l {
		idx += skipLeft(input[idx:])
		r, w := utf8.DecodeRune(input[idx:])
		switch r {
		case utf8.RuneError:
			goto Utf8Err
		case '[':
			idx += w

			r1, w1 := utf8.DecodeRune(input[idx:])
			switch r1 {
			case utf8.RuneError:
				goto Utf8Err
			case '[':
				idx += w1
				delta, err = extractTableArray(input[idx:], doc)
			case '\n', '\r', '\f':
				goto KeyErr
			default:
				delta, err = extractTable(input[idx:], doc)
			}
		default:
			delta, err = extractKeyValueSection(input[idx:], doc)

		}
		idx += delta
	}

	return

Utf8Err:
	err = errUtf8
	return
KeyErr:
	err = errInvalidKeyName
	return
}

func Load(path string) (doc *Toml, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	doc, err = Parse(bytes)

	return
}

func ParseString(input string) (doc *Toml, err error) {
	return Parse([]byte(input))
}

func Write(doc *Toml, path string) (err error) {
	if doc == nil {
		return
	}

	file, err := os.Create(path)
	defer file.Close()

	if err != nil {
		return
	}

	writer := bufio.NewWriter(file)
	doc.WriteTo(writer)
	err = writer.Flush()
	return
}

