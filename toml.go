package fiptoml

import (
	"errors"
	"strings"
	"time"
	"fmt"
	"bufio"
	"strconv"
	"reflect"
)

var (
	errValueNotFound = errors.New("Value not found")
	errTypeMismatch  = errors.New("Type mismatch")
	errNoKey         = errors.New("No key name")
)

type Toml struct {
	dict map[string]interface{}
}

func NewToml() *Toml {
	return &Toml{make(map[string]interface{})}
}

func (t *Toml) GetStringEx(key string) (val string, err error) {
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

func (t *Toml) GetString(key string, dflt string) string {
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

func (t *Toml) GetBoolEx(key string) (val bool, err error) {
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

func (t *Toml) GetBool(key string, dflt bool) bool {
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

func (t *Toml) GetIntEx(key string) (val int, err error) {
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

func (t *Toml) GetInt(key string, dflt int) int {
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

func (t *Toml) GetFloatEx(key string) (val float64, err error) {
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

func (t *Toml) GetFloat(key string, dflt float64) float64 {
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

func (t *Toml) GetDatetimeEx(key string) (val time.Time, err error) {
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

func (t *Toml) GetDatetime(key string, dflt time.Time) time.Time {
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

func (t *Toml) GetArrayEx(key string) (array interface{}, err error) {
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

func (t *Toml) GetStringArray(key string) []string {
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

func (t *Toml) GetBoolArray(key string) []bool {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	if err != nil {
		return nil
	}

	switch arr := doc.dict[fKey].(type) {
	case []bool:
		return arr
	default:
		return nil
	}
}

func (t *Toml) GetIntArray(key string) []int {
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

func (t *Toml) GetFloatArray(key string) []float64 {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	if err != nil {
		return nil
	}

	switch arr := doc.dict[fKey].(type) {
	case []float64:
		return arr
	default:
		return nil
	}
}

func (t *Toml) GetDatetimeArray(key string) []time.Time {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	if err != nil {
		fmt.Println("key,err:",key,err)
		return nil
	}

	switch arr := doc.dict[fKey].(type) {
	case []time.Time:
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

func (t *Toml) GetTableToml(key string) (table *Toml, err error) {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	switch v := doc.dict[fKey].(type) {
	case *Toml:
		table = v
	case nil:
		//err = errValueNotFound
		table = nil
	default:
		err = errTypeMismatch
	}
	return
}

func (t *Toml) GetTableArray(key string) (array []*Toml, err error) {
	fKey, doc, err := getFinalKeyAndTable(key, t)
	switch v := doc.dict[fKey].(type) {
	case []*Toml:
		array = v
	case nil:
		err = errValueNotFound
		//table = nil
	default:
		err = errTypeMismatch
	}
	return
}

func getFinalKeyAndTable(key string, doc *Toml) (finalKey string, finalToml *Toml, err error) {
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
		case *Toml:
			finalKey, finalToml, err = getFinalKeyAndTable(keys[1], v)
		default:
			err = errValueNotFound
		}
	}

	return
}

/*
func (t toml) getStruct(key string, st interface {}) ( err error) {
	switch doc := t.dict[key].(type){
	case *toml:
		v := reflect.ValueOf(st).Elem()
		tp := v.Type()
		for i:=0; i< v.NumField(); i++ {
			f := v.Field(i)
			name := tp.Field(i).Name
			switch f.Type() {
			case reflect.String:
				s, e := doc.GetString(name)
				if e != nil {
					f.SetString(s)
				}
				case reflect.Int


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
}
*/

func (t *Toml) SetValue(key string, v interface {}) {
	t.dict[key] = v
}
/*
func (t *Toml) SetTable(key string, v *Toml) {

}*/

func (t *Toml) WriteTo(writer *bufio.Writer) {
	for key := range t.dict {
		switch val := t.dict[key].(type) {
		case []*Toml:
			fmt.Fprint(writer,"[[",key,"]]\n")
			for _,st := range val {
				st.WriteTo(writer)
			}
		case *Toml:
			fmt.Fprint(writer, "[",key,"]\n")
			val.WriteTo(writer)
		default:
			fmt.Fprintln(writer,key,"=",wrapVal(val))
		}
	}
}

func wrapVal(val interface {}) string {
	switch v := val.(type) {
	case string:
		return fmt.Sprint("\"",v,"\"")
	case time.Time:
		return v.Format(time.RFC3339)
	case []time.Time:
		s := "["
		l := len(v)
		for i:=0;i < l;i++ {
			if i > 0 {
				s += ","
			}
			s += v[i].Format(time.RFC3339)
		}
		s += "]"
		return s
	case []int:
		s := "["
		l := len(v)
		for i:=0;i < l;i++ {
			if i > 0 {
				s += ","
			}
			s += strconv.Itoa(v[i])
		}
		s += "]"
		return s
	case []float64:
		s := "["
		l := len(v)
		for i:=0;i < l;i++ {
			if i > 0 {
				s += ","
			}
			s += fmt.Sprint(v[i])
		}
		s += "]"
		return s
	case []bool:
		s := "["
		l := len(v)
		for i:=0;i < l;i++ {
			if i > 0 {
				s += ","
			}
			s += strconv.FormatBool(v[i])
		}
		s += "]"
		return s
	case []string:
		s := "["
		l := len(v)
		for i:=0;i < l;i++ {
			if i > 0 {
				s += ","
			}
			s += fmt.Sprint("\"",v[i],"\"")
			fmt.Println("wrap:",i,l,s)
		}
		s += "]"
		return s
	default:
		return fmt.Sprint(v)
	}
}

//for test
func IterateTomlDoc(doc *Toml) {
	dict := doc.dict
	count := len(dict)
	fmt.Println("length of toml:", count)
	for key := range dict {
		switch val := dict[key].(type) {
		case nil:
			fmt.Println("empty value of: ", key)
		case string:
			fmt.Println(key, "=", val,"(string)")
		case *Toml:
			fmt.Println("Table:", key)
			IterateTomlDoc(val)
		case []*Toml:
			fmt.Println("Array of tables:", key)
		for i, v := range val {
			fmt.Println("Table", i)
			IterateTomlDoc(v)
		}
		default:
			fmt.Println(key, "=", val,"(", reflect.TypeOf(val),")")
		}
	}
}

