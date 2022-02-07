package gJSON

import (
	"fmt"
	"reflect"
)

//JsonObject
////////////////////////////////////////////////////////////////////////////////////////////////////
func NewObject() (obj *JsonObject) {
	obj = new(JsonObject)
	obj.init()
	return obj
}

func ParseObject(s string) (obj *JsonObject) {
	tok := newJsonTokenizer(JSONString(s))
	obj = tok.parseObject()
	if tok.err || !tok.eof() {
		return nil
	}
	return obj
}

func (O *JsonObject) print(depth int, bfmt bool) string {
	return `{` + O.jsonBase.print(depth, bfmt) + `}`
}

func (O *JsonObject) Contains(k string) bool {
	return O.getElement(JSONString(k)) != nil
}

func (O *JsonObject) Has(k string) bool {
	return O.Contains(k)
}

func (O *JsonObject) GetNames() (arr *JsonArray) {
	arr = NewArray()
	for _, v := range O.values {
		arr.AddString(v.e.key.toString())
	}
	return arr
}

func (O *JsonObject) GetValues() (arr *JsonArray) {
	arr = NewArray()
	for _, v := range O.values {
		switch v.vType {
		case val_Type_NULL:
			arr.AddNull()
		case val_Type_True:
			arr.AddBool(true)
		case val_Type_False:
			arr.AddBool(false)
		case val_Type_String:
			arr.AddString(v.vString.toString())
		case val_Type_Array:
			arr.AddArray(v.vArray)
		case val_Type_Object:
			arr.AddObject(v.vObject)
		case val_Type_Number:
			arr.AddDouble(v.vNumber.toFloat64())
		default:
			panic(fmt.Sprintf(`unknow type %d`, v.vType))
		}
	}
	return arr
}

////////////////////////////////////////////////////////////////////////////////

func (O *JsonObject) KeyAsArray(k string) *JsonArray {
	return O.GetArray(k)
}

func (O *JsonObject) KeyAsDouble(k string) float64 {
	return O.GetDouble(k)
}

func (O *JsonObject) KeyAsFloat(k string) float32 {
	return O.GetFloat(k)
}

func (O *JsonObject) KeyAsInt(k string) int {
	return O.GetInt(k)
}

func (O *JsonObject) KeyAsInt32(k string) int32 {
	return O.GetInt32(k)
}

func (O *JsonObject) KeyAsInt64(k string) int64 {
	return O.GetInt64(k)
}

func (O *JsonObject) KeyAsObject(k string) *JsonObject {
	return O.GetObject(k)
}

func (O *JsonObject) KeyAsString(k string) string {
	return O.GetString(k)
}

////////////////////////////////////////////////////////////////////////////////

func (O *JsonObject) Val(k string) (val *JsonValue) {
	val, _ = O.ValErr(k)
	return val
}

func (O *JsonObject) ValErr(k string) (val *JsonValue, err error) {
	if e0 := O.getElement(JSONString(k)); e0 == nil {
		return nil, fmt.Errorf("there is no node whose key is %s", k)
	} else {
		return e0.value, nil
	}
}

func (O *JsonObject) GetArray(k string) *JsonArray {
	if val := O.Val(k); val != nil {
		if arr, err := val.AsArray(); err == nil {
			return arr
		}
		return nil
	}
	return nil
}

func (O *JsonObject) GetBool(k string) bool {
	if val := O.Val(k); val != nil {
		return val.AsBoolDef(false)
	}
	return false
}

func (O *JsonObject) GetDouble(k string) float64 {
	if val := O.Val(k); val != nil {
		return val.AsFloat64Def(0)
	}
	return 0
}

func (O *JsonObject) GetFloat(k string) float32 {
	return float32(O.GetDouble(k))
}

func (O *JsonObject) GetInt(k string) int {
	return int(O.GetInt64(k))
}

func (O *JsonObject) GetInt32(k string) int32 {
	return int32(O.GetInt64(k))
}

func (O *JsonObject) GetInt64(k string) int64 {
	if val := O.Val(k); val != nil {
		return val.AsInt64Def(0)
	}
	return 0
}

func (O *JsonObject) GetObject(k string) *JsonObject {
	if val := O.Val(k); val != nil {
		if obj, err := val.AsObject(); err == nil {
			return obj
		}
		return nil
	}
	return nil
}

func (O *JsonObject) GetString(k string) string {
	if val := O.Val(k); val != nil {
		return val.AsStringDef(``)
	}
	return ``
}

func (O *JsonObject) Remove(k string) bool {
	return O.delElement(JSONString(k))
}

func (O *JsonObject) Set(k string, v interface{}) {
	if v == nil {
		O.SetNull(k)
		return
	}
	rt := reflect.TypeOf(v)
	rtk := rt.Kind()
	switch rtk {
	case reflect.Bool:
		O.SetBool(k, v.(bool))
	case reflect.Float64:
		O.SetDouble(k, v.(float64))
	case reflect.Float32:
		O.SetFloat(k, v.(float32))
	case reflect.Int:
		O.SetInt(k, v.(int))
	case reflect.Int32:
		O.SetInt32(k, v.(int32))
	case reflect.Int64:
		O.SetInt64(k, v.(int64))
	case reflect.String:
		O.SetString(k, v.(string))
	case reflect.Ptr:
		vname := rt.Elem().Name()
		if vname == reflect.TypeOf(JsonObject{}).Name() {
			O.SetObject(k, v.(*JsonObject))
		} else if vname == reflect.TypeOf(JsonArray{}).Name() {
			O.SetArray(k, v.(*JsonArray))
		} else {
			panic(fmt.Errorf(`unknown %s pointer`, vname))
		}
	default:
		panic(fmt.Errorf(`invalid type %s`, rt.Name()))
	}
}

func (O *JsonObject) SetBool(k string, v bool) {
	e0 := O.addElement(JSONString(k))
	if v {
		e0.value.vType = val_Type_True
	} else {
		e0.value.vType = val_Type_False
	}
}

func (O *JsonObject) SetDouble(k string, v float64) {
	e0 := O.addElement(JSONString(k))
	e0.value.vType = val_Type_Number
	e0.value.vNumber.setVal(v)
}

func (O *JsonObject) SetFloat(k string, v float32) {
	O.SetDouble(k, float64(v))
}

func (O *JsonObject) SetInt(k string, v int) {
	O.SetInt64(k, int64(v))
}

func (O *JsonObject) SetInt32(k string, v int32) {
	O.SetInt64(k, int64(v))
}

func (O *JsonObject) SetInt64(k string, v int64) {
	e0 := O.addElement(JSONString(k))
	e0.value.vType = val_Type_Number
	e0.value.vNumber.setVal(v)
}

func (O *JsonObject) SetNull(k string) {
	e0 := O.addElement(JSONString(k))
	e0.value.vType = val_Type_NULL
}

func (O *JsonObject) SetString(k, v string) {
	e0 := O.addElement(JSONString(k))
	e0.value.vType = val_Type_String
	e0.value.vString = JSONString(v)
}

func (O *JsonObject) SetObject(k string, v *JsonObject) {
	e0 := O.addElement(JSONString(k))
	if v == nil {
		e0.value.vType = val_Type_NULL
	} else {
		e0.value.vType = val_Type_Object
		e0.value.vObject = v
		var check func(obj *JsonObject)
		check = func(obj *JsonObject) {
			for i := 0; i < obj.childCount; i++ {
				v0 := obj.values[i]
				if v0.vType == val_Type_Object {
					if v0.vObject == O {
						panic(`current json object is include in v`)
					}
					check(v0.vObject)
				}
			}
		}
		check(v)
	}
}

func (O *JsonObject) SetArray(k string, v *JsonArray) {
	e0 := O.addElement(JSONString(k))
	if v == nil {
		e0.value.vType = val_Type_NULL
	} else {
		e0.value.vType = val_Type_Array
		e0.value.vArray = v
		var check func(arr *JsonArray)
		check = func(arr *JsonArray) {
			for i := 0; i < arr.childCount; i++ {
				v0 := arr.values[i]
				if v0.vType == val_Type_Object {
					if v0.vObject == O {
						panic(`current json object is include in v`)
					} else if v0.vType == val_Type_Array {
						check(v0.vArray)
					}
				}
			}
		}
		check(v)
	}
}

func (O *JsonObject) AsString() (str string) {
	return O.print(0, false)
}

func (O *JsonObject) AsJson() (str string) {
	return O.print(0, true)
}
