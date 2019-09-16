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

func (this *JsonArray) addElement() (element *JsonElement) {
	return this.jsonBase.addElement(nil)
}

func (this *JsonArray) print(depth int, bfmt bool) string {
	return `[` + this.jsonBase.print(depth, bfmt) + `]`
}

func (this *JsonArray) Add(v interface{}) int {
	if v == nil {
		return this.AddNull()
	}
	rt := reflect.TypeOf(v)
	rtk := rt.Kind()
	switch rtk {
	case reflect.Bool:
		return this.AddBool(v.(bool))
	case reflect.Float64:
		return this.AddDouble(v.(float64))
	case reflect.Float32:
		return this.AddFloat(v.(float32))
	case reflect.Int:
		return this.AddInt(v.(int))
	case reflect.Int32:
		return this.AddInt32(v.(int32))
	case reflect.Int64:
		return this.AddInt64(v.(int64))
	case reflect.String:
		return this.AddString(v.(string))
	case reflect.Ptr:
		vname := rt.Elem().Name()
		if vname == reflect.TypeOf(JsonObject{}).Name() {
			return this.AddObject(v.(*JsonObject))
		} else if vname == reflect.TypeOf(JsonArray{}).Name() {
			return this.AddArray(v.(*JsonArray))
		} else {
			panic(fmt.Errorf(`unknown %s pointer`, vname))
		}
	default:
		panic(fmt.Errorf(`invalid type %s`, rt.Name()))
	}
}

func (this *JsonArray) AddArray(v *JsonArray) int {
	if this == v {
		panic(`current array equals append value`)
	}
	e0 := this.addElement()
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
					if v0.vArray == this {
						panic(`current json array is include in v`)
					}
					check(v0.vArray)
				}
			}
		}
		check(v)
	}
	return this.currIndex()
}

func (this *JsonArray) AddBool(v bool) int {
	e0 := this.addElement()
	if v {
		e0.value.vType = val_Type_True
	} else {
		e0.value.vType = val_Type_False
	}
	return this.currIndex()
}

func (this *JsonArray) AddDouble(v float64) int {
	e0 := this.addElement()
	e0.value.vType = val_Type_Number
	e0.value.vNumber.d(v)
	return this.currIndex()
}

func (this *JsonArray) AddFloat(v float32) int {
	return this.AddDouble(float64(v))
}

func (this *JsonArray) AddInt(v int) int {
	return this.AddInt64(int64(v))
}

func (this *JsonArray) AddInt32(v int32) int {
	return this.AddInt64(int64(v))
}

func (this *JsonArray) AddInt64(v int64) int {
	return this.AddDouble(float64(v))
}

func (this *JsonArray) AddNull() int {
	e0 := this.addElement()
	e0.value.vType = val_Type_NULL
	return this.currIndex()
}

func (this *JsonArray) AddObject(v *JsonObject) int {
	e0 := this.addElement()
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
					if v0.vArray == this {
						panic(`current json array is include in v`)
					}
				} else if v0.vType == val_Type_Object {
					check(v0.vObject)
				}
			}
		}
		check(v)
	}
	return this.currIndex()
}

func (this *JsonArray) AddString(v string) int {
	e0 := this.addElement()
	e0.value.vType = val_Type_String
	e0.value.vString = JSONString(v)
	return this.currIndex()
}

func (this *JsonArray) currIndex() int {
	return this.childCount - 1
}

func (this *JsonArray) Values() []*JsonValue {
	return this.values
}

func (this *JsonArray) Val(idx int) (val *JsonValue) {
	if e0 := this.getElement(idx); e0 == nil {
		return nil
	} else {
		return e0.value
	}
}

func (this *JsonArray) GetArray(idx int) *JsonArray {
	if val := this.Val(idx); val == nil {
		return nil
	} else {
		if val.vType == val_Type_Array {
			return val.vArray
		} else {
			return nil
		}
	}
}

func (this *JsonArray) GetBool(idx int) bool {
	if val := this.Val(idx); val == nil {
		return false
	} else {
		return val.AsBoolDef(false)
	}
}

func (this *JsonArray) GetDouble(idx int) float64 {
	if val := this.Val(idx); val == nil {
		return 0
	} else {
		return val.AsDoubleDef(0)
	}
}

func (this *JsonArray) GetFloat(idx int) float32 {
	return float32(this.GetDouble(idx))
}

func (this *JsonArray) GetInt(idx int) int {
	return int(this.GetInt64(idx))
}

func (this *JsonArray) GetInt32(idx int) int32 {
	return int32(this.GetInt64(idx))
}

func (this *JsonArray) GetInt64(idx int) int64 {
	if val := this.Val(idx); val == nil {
		return 0
	} else {
		return val.AsInt64Def(0)
	}
}

func (this *JsonArray) GetObject(idx int) *JsonObject {
	if val := this.Val(idx); val == nil {
		return nil
	} else {
		if val.vType == val_Type_Object {
			return val.vObject
		} else {
			return nil
		}
	}
}

func (this *JsonArray) GetString(idx int) string {
	if val := this.Val(idx); val == nil {
		return ``
	} else {
		return val.AsStringDef(``)
	}
}

func (this *JsonArray) insertValue(idx int) (val *JsonValue) {
	if val = this.Val(idx); val == nil {
		k := idx - this.childCount + 1
		for i := 0; i < k; i++ {
			val = this.values[this.AddNull()]
		}
	}
	return val
}

func (this *JsonArray) SetArray(idx int, arr *JsonArray) {
	val := this.insertValue(idx)
	val.vType = val_Type_Array
	val.vArray = arr
}

func (this *JsonArray) SetBool(idx int, b bool) {
	val := this.insertValue(idx)
	if b {
		val.vType = val_Type_True
	} else {
		val.vType = val_Type_False
	}
}

func (this *JsonArray) SetDouble(idx int, f64 float64) {
	val := this.insertValue(idx)
	val.vType = val_Type_Number
	val.vNumber.d(f64)
}

func (this *JsonArray) SetFloat(idx int, f32 float32) {
	this.SetDouble(idx, float64(f32))
}

func (this *JsonArray) SetInt(idx int, i int) {
	this.SetInt64(idx, int64(i))
}

func (this *JsonArray) SetInt32(idx int, i32 int32) {
	this.SetInt64(idx, int64(i32))
}

func (this *JsonArray) SetInt64(idx int, i64 int64) {
	this.SetDouble(idx, float64(i64))
}

func (this *JsonArray) SetNull(idx int) {
	val := this.insertValue(idx)
	val.vType = val_Type_NULL
}

func (this *JsonArray) SetObject(idx int, obj *JsonObject) {
	val := this.insertValue(idx)
	val.vType = val_Type_Object
	val.vObject = obj
}

func (this *JsonArray) SetString(idx int, s string) {
	val := this.insertValue(idx)
	val.vType = val_Type_String
	val.vString = JSONString(s)
}

func (this *JsonArray) A(idx int) *JsonArray {
	return this.GetArray(idx)
}

func (this *JsonArray) B(idx int) interface{} {
	return this.GetBool(idx)
}

func (this *JsonArray) D(idx int) float64 {
	return this.GetDouble(idx)
}

func (this *JsonArray) F(idx int) float32 {
	return this.GetFloat(idx)
}

func (this *JsonArray) I(idx int) int {
	return this.GetInt(idx)
}

func (this *JsonArray) I32(idx int) int32 {
	return this.GetInt32(idx)
}

func (this *JsonArray) I64(idx int) int64 {
	return this.GetInt64(idx)
}

func (this *JsonArray) O(idx int) *JsonObject {
	return this.GetObject(idx)
}

func (this *JsonArray) S(idx int) string {
	return this.GetString(idx)
}

func (this *JsonArray) Remove(idx int) bool {
	return this.delElement(idx)
}

func (this *JsonArray) AsString() (str string) {
	return this.print(0, false)
}

func (this *JsonArray) AsJson() (str string) {
	return this.print(0, true)
}
