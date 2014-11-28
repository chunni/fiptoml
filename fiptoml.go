package fiptoml

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"unicode/utf8"
)

var (
	errInvalidToml = errors.New("invalid TOML doc")

	errInvalidTableKey   = errors.New("invalid table key name")
	errInvalidKeyName    = errors.New("invalid key name")
	errUtf8              = errors.New("not valid UTF-8 content")
	errEmptyKey          = errors.New("key name is empty")
	errBool              = errors.New("bool should be either true or false")
	errNumber            = errors.New("number like value can only be int, float or datetime(RFC3399)")
	errMultiString       = errors.New("invalid multi-line string")
	errStringSyntaxError = errors.New("string syntax error")
	errArray             = errors.New("Date types in an array should NOT be mixed")

	multiLineSkipR = regexp.MustCompile(`\\[\n\r\t\f ]+`)
	quoteLineR     = regexp.MustCompile(`\n`)
	intR           = regexp.MustCompile(`^[+-]?(?:0|[1-9][0-9]*)$`)
	floatR         = regexp.MustCompile(`-?(?:0|[1-9][0-9]*)\.[0-9]+$`)
	datatimeR      = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(.\d+)?(Z|[+-]\d{2}:\d{2})$`)
)

func errDuplicatedKey(key string) error {
	return errors.New(fmt.Sprint("Duplicated key: ", key))
}

func errUnsupportedValue(key string) error {
	return errors.New(fmt.Sprint("unsupported value type for key: ", key))
}

func Parse(input []byte) (doc *toml, err error) {
	doc = newToml()
	idx, delta := 0, 0

	for idx < len(input) {
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

func Load(path string) (doc *toml, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	doc, err = Parse(bytes)

	return
}

func ParseString(input string) (doc *toml, err error) {
	return Parse([]byte(input))
}
