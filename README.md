#fiptoml - Golang TOML parser

fiptoml is a [TOML](https://github.com/toml-lang/toml) parser for Golang. It designed to be fast, reliable, and easy-to-use.

This library is compatible with TOML version [v0.3.1](https://github.com/toml-lang/toml/blob/master/versions/toml-v0.3.1.md).

## Why TOML
It is simple, you may even understand the spec at one glance, yet meets various needs of config in my real world projects.
And, OK, I admit, I'm just too lazy to get to YAML, which looks scary to me.

## How to use fiptoml
### Get it
`go get "github.com/chunni/fiptoml"`

### Use it

#### Import the package

`import "github.com/chuni/fiptoml"`

#### Parse your config

To load a config file:

```
toml, err := fiptoml.Load("./config/config.toml")
title := toml.GetString("title","")
```

To parse a string:

```
toml, err := ParseString(example)
if err != nil {
    fmt.Println("ParseFile should work")
    return
}
owner := toml.GetString("owner.name","")
fmt.Println(owner)
```

Or, to parse a byte array:

```
toml, err := Parse([]byte(example))
if err != nil {
	fmt.Println("ParseFile should work")
	return
}
files := toml.GetStringArray("files")
fmt.Println(files)
```

Please refer to the test file [fiptoml_test.go](https://github.com/chunni/fiptoml/blob/master/fiptoml_test.go) for working examples.

### API list
- `func Load(path string) (doc *toml, err error)`
- `func Parse(input []byte) (doc *toml, err error)`
- `func ParseString(input string) (doc *toml, err error)`

`type toml struct`

- `func (t *toml) GetString(key string, dflt string) string`
- `func (t *toml) GetStringEx(key string) (val string, err error)`
- `func (t *toml) GetBool(key string, dflt bool) bool`
- `func (t *toml) GetBoolEx(key string) (val bool, err error)`
- `func (t *toml) GetInt(key string, dflt int) int`
- `func (t *toml) GetIntEx(key string) (val int, err error)`
- `func (t *toml) GetFloat(key string, dflt float64) float64`
- `func (t *toml) GetFloatEx(key string) (val float64, err error)`
- `func (t *toml) GetDatetime(key string, dflt time.Time) time.Time`
- `func (t *toml) GetDatetimeEx(key string) (val time.Time, err error)`
- `func (t *toml) GetStringArray(key string) []string`
- `func (t *toml) GetBoolArray(key string) []bool`
- `func (t *toml) GetIntArray(key string) []int`
- `func (t *toml) GetFloatArray(key string) []float64`
- `func (t *toml) GetDatetimeArray(key string) []time.Time`
- `func (t *toml) GetTableToml(key string) (table *toml, err error)`
- `func (t *toml) GetTableArray(key string) (array []*toml, err error)`
