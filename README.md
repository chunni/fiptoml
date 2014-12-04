#fiptoml - Golang TOML parser(deserializer)/writer(serializer)

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

#### Parse(deserialize) TOML

**Load a config file:**

```
toml, err := fiptoml.Load("./config/config.toml")
```
**Parse a byte array:**

```
toml, err := fiptoml.Parse([]byte(example))
```

**Parse a string:**

```
toml, err := fiptoml.ParseString(example)
```

**Get the values:**

You may get value quickly by set a default value in case something goes wrong.
```
title := toml.GetString("title","")
owner := toml.GetString("owner.name","")
files := toml.GetStringArray("files")
```
Or, you may check the error yourself to ensure your config file is valid.
```
title, err := toml.GetStringEx("title")
if err != nil {
    //handle the error
}
```

#### Write/serialize TOML
**Form a TOML document:**

```
toml := NewToml()
toml.SetValue("title", "A Perfect Trip")
toml.SetValue("days", 21)
toml.SetValue("start", time.Now())
toml.SetValue("enabled",true)
toml.SetValue("guys",[]string{"Tony","Tim","Abby"})
```
**Serialize it to a writer:**
```
toml.WriteTo(writer)
```
**Directly write it to a file:**
```
fiptoml.Write(toml,"./config/out.toml")
```

Please refer to the test file [fiptoml_test.go](https://github.com/chunni/fiptoml/blob/master/fiptoml_test.go) for working examples.

### API list
- `func Load(path string) (doc *toml, err error)`
- `func Parse(input []byte) (doc *toml, err error)`
- `func ParseString(input string) (doc *toml, err error)`
- `func Write(doc *Toml, path string) (err error)`

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
- `func (t *Toml) WriteTo(writer *bufio.Writer)`
