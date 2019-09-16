package gJSON

import "testing"

func TestJsonObjectParse(t *testing.T) {
	obj := ParseObject(`{"k_str":"string","k_true":true,"k_false":false,"k_null":null,"k_number_0":123,"k_number_1":1.1314,"k_object":{},"k_array":[]}`)
	if obj == nil {
		t.Fatal(`json object parse error`)
	}
	t.Log(obj.AsString())
	t.Log(obj.GetString(`k_str`))
	t.Log(obj.GetString(`k_true`))
	t.Log(obj.GetString(`k_false`))
	t.Log(obj.GetString(`k_null`))
	t.Log(obj.GetString(`k_number_0`))
	t.Log(obj.GetString(`k_number_1`))
	t.Log(obj.GetString(`k_object`))
	t.Log(obj.GetString(`k_array`))
}

func TestJsonObjectCreate(t *testing.T) {
	obj := NewObject()
	obj.Set(`k0_str`, `string`)
	obj.Set(`k0_true`, true)
	obj.Set(`k0_false`, false)
	obj.Set(`k0_null`, nil)
	obj.Set(`k0_number_0`, 123)
	obj.Set(`k0_number_1`, 1.1314)
	o1 := NewObject()
	obj.Set(`k0_object_0`, o1)
	a1 := NewArray()
	obj.Set(`k0_array_0`, a1)
	t.Log(obj.AsString())
}

func TestJsonArrayParse(t *testing.T) {
	arr := ParseArray(`["string",true,false,null,123,1.1314,{},[]]`)
	if arr == nil {
		t.Fatal(`json array parse error`)
	}
	t.Log(arr.AsString())
	for _, v := range arr.Values() {
		t.Log(v.AsStringDef(`invalid`))
	}
}

func TestJsonArrayCreate(t *testing.T) {
	arr := NewArray()
	arr.Add(`str`)
	arr.Add(true)
	arr.Add(false)
	arr.Add(nil)
	arr.Add(999)
	arr.Add(3.1415926)
	o1 := NewObject()
	arr.Add(o1)
	a1 := NewArray()
	arr.Add(a1)
	t.Log(arr.AsString())
}
