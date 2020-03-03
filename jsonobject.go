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

func (this *JsonObject) print(depth int, bfmt bool) string {
	return `{` + this.jsonBase.print(depth, bfmt) + `}`
}

func (this *JsonObject) Contains(k string) bool {
	return this.getElement(JSONString(k)) != nil
}

func (this *JsonObject) Has(k string) bool {
	return this.Contains(k)
}

func (this *JsonObject) GetNames() (arr *JsonArray) {
	arr = NewArray()
	for _, v := range this.values {
		arr.AddString(v.e.key.toString())
	}
	return arr
}

func (this *JsonObject) GetValues() (arr *JsonArray) {
	arr = NewArray()
	for _, v := range this.values {
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

func (this *JsonObject) KeyAsArray(k string) *JsonArray {
	return this.GetArray(k)
}

func (this *JsonObject) KeyAsDouble(k string) float64 {
	return this.GetDouble(k)
}

func (this *JsonObject) KeyAsFloat(k string) float32 {
	return this.GetFloat(k)
}

func (this *JsonObject) KeyAsInt(k string) int {
	return this.GetInt(k)
}

func (this *JsonObject) KeyAsInt32(k string) int32 {
	return this.GetInt32(k)
}

func (this *JsonObject) KeyAsInt64(k string) int64 {
	return this.GetInt64(k)
}

func (this *JsonObject) KeyAsObject(k string) *JsonObject {
	return this.GetObject(k)
}

func (this *JsonObject) KeyAsString(k string) string {
	return this.GetString(k)
}

////////////////////////////////////////////////////////////////////////////////

func (this *JsonObject) Val(k string) (val *JsonValue) {
	val, _ = this.ValErr(k)
	return val
}

func (this *JsonObject) ValErr(k string) (val *JsonValue, err error) {
	if e0 := this.getElement(JSONString(k)); e0 == nil {
		return nil, fmt.Errorf("there is no node whose key is %s", k)
	} else {
		return e0.value, nil
	}
}

func (this *JsonObject) GetArray(k string) *JsonArray {
	if val := this.Val(k); val != nil {
		if arr, err := val.AsArray(); err == nil {
			return arr
		}
		return nil
	}
	return nil
}

func (this *JsonObject) GetBool(k string) bool {
	if val := this.Val(k); val != nil {
		return val.AsBoolDef(false)
	}
	return false
}

func (this *JsonObject) GetDouble(k string) float64 {
	if val := this.Val(k); val != nil {
		return val.AsFloat64Def(0)
	}
	return 0
}

func (this *JsonObject) GetFloat(k string) float32 {
	return float32(this.GetDouble(k))
}

func (this *JsonObject) GetInt(k string) int {
	return int(this.GetInt64(k))
}

func (this *JsonObject) GetInt32(k string) int32 {
	return int32(this.GetInt64(k))
}

func (this *JsonObject) GetInt64(k string) int64 {
	if val := this.Val(k); val != nil {
		return val.AsInt64Def(0)
	}
	return 0
}

func (this *JsonObject) GetObject(k string) *JsonObject {
	if val := this.Val(k); val != nil {
		if obj, err := val.AsObject(); err == nil {
			return obj
		}
		return nil
	}
	return nil
}

func (this *JsonObject) GetString(k string) string {
	if val := this.Val(k); val != nil {
		return val.AsStringDef(``)
	}
	return ``
}

func (this *JsonObject) Remove(k string) bool {
	return this.delElement(JSONString(k))
}

func (this *JsonObject) Set(k string, v interface{}) {
	if v == nil {
		this.SetNull(k)
		return
	}
	rt := reflect.TypeOf(v)
	rtk := rt.Kind()
	switch rtk {
	case reflect.Bool:
		this.SetBool(k, v.(bool))
	case reflect.Float64:
		this.SetDouble(k, v.(float64))
	case reflect.Float32:
		this.SetFloat(k, v.(float32))
	case reflect.Int:
		this.SetInt(k, v.(int))
	case reflect.Int32:
		this.SetInt32(k, v.(int32))
	case reflect.Int64:
		this.SetInt64(k, v.(int64))
	case reflect.String:
		this.SetString(k, v.(string))
	case reflect.Ptr:
		vname := rt.Elem().Name()
		if vname == reflect.TypeOf(JsonObject{}).Name() {
			this.SetObject(k, v.(*JsonObject))
		} else if vname == reflect.TypeOf(JsonArray{}).Name() {
			this.SetArray(k, v.(*JsonArray))
		} else {
			panic(fmt.Errorf(`unknown %s pointer`, vname))
		}
	default:
		panic(fmt.Errorf(`invalid type %s`, rt.Name()))
	}
}

func (this *JsonObject) SetBool(k string, v bool) {
	e0 := this.addElement(JSONString(k))
	if v {
		e0.value.vType = val_Type_True
	} else {
		e0.value.vType = val_Type_False
	}
}

func (this *JsonObject) SetDouble(k string, v float64) {
	e0 := this.addElement(JSONString(k))
	e0.value.vType = val_Type_Number
	e0.value.vNumber.setVal(v)
}

func (this *JsonObject) SetFloat(k string, v float32) {
	this.SetDouble(k, float64(v))
}

func (this *JsonObject) SetInt(k string, v int) {
	this.SetInt64(k, int64(v))
}

func (this *JsonObject) SetInt32(k string, v int32) {
	this.SetInt64(k, int64(v))
}

func (this *JsonObject) SetInt64(k string, v int64) {
	e0 := this.addElement(JSONString(k))
	e0.value.vType = val_Type_Number
	e0.value.vNumber.setVal(v)
}

func (this *JsonObject) SetNull(k string) {
	e0 := this.addElement(JSONString(k))
	e0.value.vType = val_Type_NULL
}

func (this *JsonObject) SetString(k, v string) {
	e0 := this.addElement(JSONString(k))
	e0.value.vType = val_Type_String
	e0.value.vString = JSONString(v)
}

func (this *JsonObject) SetObject(k string, v *JsonObject) {
	e0 := this.addElement(JSONString(k))
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
					if v0.vObject == this {
						panic(`current json object is include in v`)
					}
					check(v0.vObject)
				}
			}
		}
		check(v)
	}
}

func (this *JsonObject) SetArray(k string, v *JsonArray) {
	e0 := this.addElement(JSONString(k))
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
					if v0.vObject == this {
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

func (this *JsonObject) AsString() (str string) {
	return this.print(0, false)
}

func (this *JsonObject) AsJson() (str string) {
	return this.print(0, true)
}
