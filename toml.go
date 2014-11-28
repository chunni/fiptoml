package fiptoml

import (
	"errors"
	"strings"
	"time"
	//	"fmt"
	"fmt"
	//	"reflect"
)

var (
	errValueNotFound = errors.New("Value not found")
	errTypeMismatch  = errors.New("Type mismatch")
	errNoKey         = errors.New("No key name")
)

type toml struct {
	dict map[string]interface{}
}

func newToml() *toml {
	return &toml{make(map[string]interface{})}
}

func (t *toml) GetStringEx(key string) (val string, err error) {
	fKey, doc, err := getFinalKeyAndTable(key, t)

	switch v := doc.dict[fKey].(type) {
	case string:
		val = v
	case nil:
		err = errValueNotFound
	default:
		err = errTypeMismatch
	}
	return
}

func (t *toml) GetString(key string, dflt string) string {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	if err != nil {
		fmt.Println("err:", err)
		return dflt
	}

	switch v := doc.dict[fKey].(type) {
	case string:
		return v
	default:
		return dflt
	}
}

func (t *toml) GetBoolEx(key string) (val bool, err error) {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	switch v := doc.dict[fKey].(type) {
	case bool:
		val = v
	case nil:
		err = errValueNotFound
	default:
		err = errTypeMismatch
	}
	return
}

func (t *toml) GetBool(key string, dflt bool) bool {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	if err != nil {
		return dflt
	}

	switch v := doc.dict[fKey].(type) {
	case bool:
		return v
	default:
		return dflt
	}
}

func (t *toml) GetIntEx(key string) (val int, err error) {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	switch v := doc.dict[fKey].(type) {
	case int:
		val = v
	case nil:
		err = errValueNotFound
	default:
		err = errTypeMismatch
	}
	return
}

func (t *toml) GetInt(key string, dflt int) int {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	if err != nil {
		return dflt
	}
	switch v := doc.dict[fKey].(type) {
	case int:
		return v
	default:
		return dflt
	}
}

func (t *toml) GetFloatEx(key string) (val float64, err error) {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	switch v := doc.dict[fKey].(type) {
	case float64:
		val = v
	case nil:
		err = errValueNotFound
	default:
		err = errTypeMismatch
	}
	return
}

func (t *toml) GetFloat(key string, dflt float64) float64 {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	if err != nil {
		return dflt
	}
	switch v := doc.dict[fKey].(type) {
	case float64:
		return v
	default:
		return dflt
	}
}

func (t *toml) GetDatetimeEx(key string) (val time.Time, err error) {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	switch v := doc.dict[fKey].(type) {
	case time.Time:
		val = v
	case nil:
		err = errValueNotFound
	default:
		err = errTypeMismatch
	}
	return
}

func (t *toml) GetDatetime(key string, dflt time.Time) time.Time {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	if err != nil {
		return dflt
	}
	switch v := doc.dict[fKey].(type) {
	case time.Time:
		return v
	default:
		return dflt
	}
}

func (t *toml) GetArrayEx(key string) (array interface{}, err error) {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	switch arr := doc.dict[fKey].(type) {
	case []string:
		array = arr
	case []bool:
		array = arr
	case []int:
		array = arr
	case []float64:
		array = arr
	case []time.Time:
		array = arr
	case nil:
		err = errValueNotFound
	default:
		err = errTypeMismatch
	}
	return
}

func (t *toml) GetStringArray(key string) []string {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	if err != nil {
		return nil
	}

	switch arr := doc.dict[fKey].(type) {
	case []string:
		return arr
	default:
		return nil
	}
}

func (t *toml) GetIntArray(key string) []int {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	if err != nil {
		return nil
	}

	switch arr := doc.dict[fKey].(type) {
	case []int:
		return arr
	default:
		return nil
	}
}

/*func (t *toml) GetArray(key string,dflt interface{}) interface {} {
	fKey, doc, err := getFinalKeyAndTable(key,t)
	if(err != nil) {
		return dflt
	}

	switch arr := doc.dict[fKey].(type) {
	case []string:
		return arr
	case []bool:
		return arr
	case []int:
		return arr
	case []float64:
		return arr
	case []time.Time:
		return arr
	default:
		return dflt
	}

}*/

func (t *toml) GetTableToml(key string) (table *toml, err error) {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	switch v := doc.dict[fKey].(type) {
	case *toml:
		table = v
	case nil:
		//err = errValueNotFound
		table = nil
	default:
		err = errTypeMismatch
	}
	return
}

func (t *toml) GetTableArray(key string) (array []*toml, err error) {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	switch v := doc.dict[fKey].(type) {
	case []*toml:
		array = v
	case nil:
		err = errValueNotFound
		//table = nil
	default:
		err = errTypeMismatch
	}
	return
}

func getFinalKeyAndTable(key string, doc *toml) (finalKey string, finalToml *toml, err error) {
	if len(key) == 0 {
		err = errNoKey
		return
	}

	keys := strings.SplitN(key, ".", 2)
	//keys[0] is always there
	if len(keys[0]) == 0 {
		err = errNoKey
		return
	}
	if len(keys) == 1 {
		finalKey, finalToml = key, doc
	} else {
		if len(keys[1]) == 0 {
			err = errNoKey
			return
		}
		switch v := doc.dict[keys[0]].(type) {
		case *toml:
			finalKey, finalToml, err = getFinalKeyAndTable(keys[1], v)
		default:
			err = errValueNotFound
		}
	}

	return
}

/*func (t toml) getStruct(key string, tp reflect.Type) (val interface {}, err error) {
	switch doc := t.dict[key].(type){
	case *toml:
		v := reflect.New(tp)
		for i:=0; i< v.NumField(); i++ {
			f := v.Field(i)
			switch t := f.Type() {
			case reflect.String:
				f.SetString(doc[f.Name])

			}
		}

		for key := range doc.dict {
			v.FieldByName(key) = doc[key] //by type recursive
		}
	case nil:
		err = errValueNotFound
	default:
		err = errTypeMismatch
	}
	return
}*/
