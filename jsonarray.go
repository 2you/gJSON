package gJSON

import (
	"fmt"
	"reflect"
)

//JsonArray
////////////////////////////////////////////////////////////////////////////////////////////////////

func NewArray() (arr *JsonArray) {
	arr = new(JsonArray)
	arr.init()
	return arr
}

func ParseArray(s string) (arr *JsonArray) {
	tok := newJsonTokenizer(JSONString(s))
	arr = tok.parseArray()
	if tok.err || !tok.eof() {
		return nil
	}
	return arr
}

func (a *JsonArray) addElement() (element *JsonElement) {
	return a.jsonBase.addElement(nil)
}

func (a *JsonArray) print(depth int, bfmt bool) string {
	return `[` + a.jsonBase.print(depth, bfmt) + `]`
}

func (a *JsonArray) Add(v interface{}) int {
	if v == nil {
		return a.AddNull()
	}
	rt := reflect.TypeOf(v)
	rtk := rt.Kind()
	switch rtk {
	case reflect.Bool:
		return a.AddBool(v.(bool))
	case reflect.Float64:
		return a.AddDouble(v.(float64))
	case reflect.Float32:
		return a.AddFloat(v.(float32))
	case reflect.Int:
		return a.AddInt(v.(int))
	case reflect.Int32:
		return a.AddInt32(v.(int32))
	case reflect.Int64:
		return a.AddInt64(v.(int64))
	case reflect.String:
		return a.AddString(v.(string))
	case reflect.Ptr:
		vname := rt.Elem().Name()
		if vname == reflect.TypeOf(JsonObject{}).Name() {
			return a.AddObject(v.(*JsonObject))
		} else if vname == reflect.TypeOf(JsonArray{}).Name() {
			return a.AddArray(v.(*JsonArray))
		} else {
			panic(fmt.Errorf(`unknown %s pointer`, vname))
		}
	default:
		panic(fmt.Errorf(`invalid type %s`, rt.Name()))
	}
}

func (a *JsonArray) AddArray(v *JsonArray) int {
	if a == v {
		panic(`current array equals append value`)
	}
	e0 := a.addElement()
	if v == nil {
		e0.value.vType = val_Type_NULL
	} else {
		e0.value.vType = val_Type_Array
		e0.value.vArray = v
		var check func(arr *JsonArray)
		check = func(arr *JsonArray) {
			for i := 0; i < arr.childCount; i++ {
				v0 := arr.values[i]
				if v0.vType == val_Type_Array {
					if v0.vArray == a {
						panic(`current json array is include in v`)
					}
					check(v0.vArray)
				}
			}
		}
		check(v)
	}
	return a.currIndex()
}

func (a *JsonArray) AddBool(v bool) int {
	e0 := a.addElement()
	if v {
		e0.value.vType = val_Type_True
	} else {
		e0.value.vType = val_Type_False
	}
	return a.currIndex()
}

func (a *JsonArray) AddDouble(v float64) int {
	e0 := a.addElement()
	e0.value.vType = val_Type_Number
	e0.value.vNumber.setVal(v)
	return a.currIndex()
}

func (a *JsonArray) AddFloat(v float32) int {
	return a.AddDouble(float64(v))
}

func (a *JsonArray) AddInt(v int) int {
	return a.AddInt64(int64(v))
}

func (a *JsonArray) AddInt32(v int32) int {
	return a.AddInt64(int64(v))
}

func (a *JsonArray) AddInt64(v int64) int {
	e0 := a.addElement()
	e0.value.vType = val_Type_Number
	e0.value.vNumber.setVal(v)
	return a.currIndex()
}

func (a *JsonArray) AddNull() int {
	e0 := a.addElement()
	e0.value.vType = val_Type_NULL
	return a.currIndex()
}

func (a *JsonArray) AddObject(v *JsonObject) int {
	e0 := a.addElement()
	if v == nil {
		e0.value.vType = val_Type_NULL
	} else {
		e0.value.vType = val_Type_Object
		e0.value.vObject = v
		var check func(obj *JsonObject)
		check = func(obj *JsonObject) {
			for i := 0; i < obj.childCount; i++ {
				v0 := obj.values[i]
				if v0.vType == val_Type_Array {
					if v0.vArray == a {
						panic(`current json array is include in v`)
					}
				} else if v0.vType == val_Type_Object {
					check(v0.vObject)
				}
			}
		}
		check(v)
	}
	return a.currIndex()
}

func (a *JsonArray) AddString(v string) int {
	e0 := a.addElement()
	e0.value.vType = val_Type_String
	e0.value.vString = JSONString(v)
	return a.currIndex()
}

func (a *JsonArray) currIndex() int {
	return a.childCount - 1
}

func (a *JsonArray) Values() []*JsonValue {
	return a.values
}

func (a *JsonArray) Val(idx int) (val *JsonValue) {
	if e0 := a.getElement(idx); e0 == nil {
		return nil
	} else {
		return e0.value
	}
}

func (a *JsonArray) GetArray(idx int) *JsonArray {
	if val := a.Val(idx); val == nil {
		return nil
	} else {
		if val.vType == val_Type_Array {
			return val.vArray
		} else {
			return nil
		}
	}
}

func (a *JsonArray) GetBool(idx int) bool {
	if val := a.Val(idx); val == nil {
		return false
	} else {
		return val.AsBoolDef(false)
	}
}

func (a *JsonArray) GetDouble(idx int) float64 {
	if val := a.Val(idx); val == nil {
		return 0
	} else {
		return val.AsDoubleDef(0)
	}
}

func (a *JsonArray) GetFloat(idx int) float32 {
	return float32(a.GetDouble(idx))
}

func (a *JsonArray) GetInt(idx int) int {
	return int(a.GetInt64(idx))
}

func (a *JsonArray) GetInt32(idx int) int32 {
	return int32(a.GetInt64(idx))
}

func (a *JsonArray) GetInt64(idx int) int64 {
	if val := a.Val(idx); val == nil {
		return 0
	} else {
		return val.AsInt64Def(0)
	}
}

func (a *JsonArray) GetObject(idx int) *JsonObject {
	if val := a.Val(idx); val == nil {
		return nil
	} else {
		if val.vType == val_Type_Object {
			return val.vObject
		} else {
			return nil
		}
	}
}

func (a *JsonArray) GetString(idx int) string {
	if val := a.Val(idx); val == nil {
		return ``
	} else {
		return val.AsStringDef(``)
	}
}

func (a *JsonArray) insertValue(idx int) (val *JsonValue) {
	if val = a.Val(idx); val == nil {
		k := idx - a.childCount + 1
		for i := 0; i < k; i++ {
			val = a.values[a.AddNull()]
		}
	}
	return val
}

func (a *JsonArray) SetArray(idx int, arr *JsonArray) {
	val := a.insertValue(idx)
	val.vType = val_Type_Array
	val.vArray = arr
}

func (a *JsonArray) SetBool(idx int, b bool) {
	val := a.insertValue(idx)
	if b {
		val.vType = val_Type_True
	} else {
		val.vType = val_Type_False
	}
}

func (a *JsonArray) SetDouble(idx int, f64 float64) {
	val := a.insertValue(idx)
	val.vType = val_Type_Number
	val.vNumber.setVal(f64)
}

func (a *JsonArray) SetFloat(idx int, f32 float32) {
	a.SetDouble(idx, float64(f32))
}

func (a *JsonArray) SetInt(idx int, i int) {
	a.SetInt64(idx, int64(i))
}

func (a *JsonArray) SetInt32(idx int, i32 int32) {
	a.SetInt64(idx, int64(i32))
}

func (a *JsonArray) SetInt64(idx int, i64 int64) {
	a.SetDouble(idx, float64(i64))
}

func (a *JsonArray) SetNull(idx int) {
	val := a.insertValue(idx)
	val.vType = val_Type_NULL
}

func (a *JsonArray) SetObject(idx int, obj *JsonObject) {
	val := a.insertValue(idx)
	val.vType = val_Type_Object
	val.vObject = obj
}

func (a *JsonArray) SetString(idx int, s string) {
	val := a.insertValue(idx)
	val.vType = val_Type_String
	val.vString = JSONString(s)
}

func (a *JsonArray) A(idx int) *JsonArray {
	return a.GetArray(idx)
}

func (a *JsonArray) B(idx int) interface{} {
	return a.GetBool(idx)
}

func (a *JsonArray) D(idx int) float64 {
	return a.GetDouble(idx)
}

func (a *JsonArray) F(idx int) float32 {
	return a.GetFloat(idx)
}

func (a *JsonArray) I(idx int) int {
	return a.GetInt(idx)
}

func (a *JsonArray) I32(idx int) int32 {
	return a.GetInt32(idx)
}

func (a *JsonArray) I64(idx int) int64 {
	return a.GetInt64(idx)
}

func (a *JsonArray) O(idx int) *JsonObject {
	return a.GetObject(idx)
}

func (a *JsonArray) S(idx int) string {
	return a.GetString(idx)
}

func (a *JsonArray) Remove(idx int) bool {
	return a.delElement(idx)
}

func (a *JsonArray) AsString() (str string) {
	return a.print(0, false)
}

func (a *JsonArray) AsJson() (str string) {
	return a.print(0, true)
}
