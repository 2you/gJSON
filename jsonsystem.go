package gJSON

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func locaFromFile(fname string) (data JSONString, err error) {
	f, e := os.Open(fname)
	if e != nil {
		return nil, e
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	info, e := f.Stat()
	if e != nil {
		return nil, e
	}
	buf := make([]byte, info.Size())
	if _, e = f.Read(buf); e != nil {
		return nil, e
	}
	return JSONString(string(buf)), nil
}

func parse_hex4(txt JSONString) uint {
	var h rune = 0
	if txt[0] >= '0' && txt[0] <= '9' {
		h += txt[0] - '0'
	} else if txt[0] >= 'A' && txt[0] <= 'F' {
		h += 10 + txt[0] - 'A'
	} else if txt[0] >= 'a' && txt[0] <= 'f' {
		h += 10 + txt[0] - 'a'
	} else {
		return 0
	}
	h = h << 4
	txt = txt[1:]
	if txt[0] >= '0' && txt[0] <= '9' {
		h += txt[0] - '0'
	} else if txt[0] >= 'A' && txt[0] <= 'F' {
		h += 10 + txt[0] - 'A'
	} else if txt[0] >= 'a' && txt[0] <= 'f' {
		h += 10 + txt[0] - 'a'
	} else {
		return 0
	}
	h = h << 4
	txt = txt[1:]
	if txt[0] >= '0' && txt[0] <= '9' {
		h += txt[0] - '0'
	} else if txt[0] >= 'A' && txt[0] <= 'F' {
		h += 10 + txt[0] - 'A'
	} else if txt[0] >= 'a' && txt[0] <= 'f' {
		h += 10 + txt[0] - 'a'
	} else {
		return 0
	}
	h = h << 4
	txt = txt[1:]
	if txt[0] >= '0' && txt[0] <= '9' {
		h += txt[0] - '0'
	} else if txt[0] >= 'A' && txt[0] <= 'F' {
		h += 10 + txt[0] - 'A'
	} else if txt[0] >= 'a' && txt[0] <= 'f' {
		h += 10 + txt[0] - 'a'
	} else {
		return 0
	}
	return uint(h)
}

func StringTrim(str string) string {
	return strings.Trim(str, " \t\n\r")
}

//JSONNode
////////////////////////////////////////////////////////////////////////////////////////////////////
func New(ss ...string) *JsonNode {
	var data string
	n := len(ss)
	if n == 0 {
		data = `{}`
	} else if n == 1 {
		data = StringTrim(ss[0])
	} else {
		return nil
	}

	if len(data) < 2 {
		return nil
	}
	node := new(JsonNode)
	if data[0] == '{' {
		if node.vObject = ParseObject(data); node.vObject == nil {
			return nil
		}
	} else if data[0] == '[' {
		if node.vArray = ParseArray(data); node.vArray == nil {
			return nil
		}
	} else {
		return nil
	}
	return node
}

func LoadFromFile(fname string) (node *JsonNode, err error) {
	data, err := locaFromFile(fname)
	if err != nil {
		return nil, err
	}
	node = New(data.toString())
	if node == nil {
		return nil, fmt.Errorf(`parse error`)
	}
	return node, nil
}

func (this *JsonNode) AsString() string {
	if this.vObject != nil {
		return this.vObject.AsString()
	}

	if this.vArray != nil {
		return this.vArray.AsString()
	}
	return ``
}

func (this *JsonNode) AsObject() *JsonObject {
	return this.vObject
}

func (this *JsonNode) AsArray() *JsonArray {
	return this.vArray
}

func (this *JsonNode) IsArray() bool {
	return this.vArray != nil
}

func (this *JsonNode) IsObject() bool {
	return this.vObject != nil
}

func (this *JsonNode) Val(k interface{}) *JsonValue {
	rt := reflect.TypeOf(k)
	rtk := rt.Kind()
	if this.IsObject() {
		if rtk == reflect.String {
			return this.vObject.Val(k.(string))
		} else {
			return nil
		}
	}

	if !this.IsArray() {
		return nil
	}
	var idx int
	switch rtk {
	case reflect.Int:
		idx = k.(int)
	case reflect.Uint:
		idx = int(k.(uint))
	case reflect.Int8:
		idx = int(k.(int8))
	case reflect.Uint8:
		idx = int(k.(uint8))
	case reflect.Int16:
		idx = int(k.(int16))
	case reflect.Uint16:
		idx = int(k.(uint16))
	case reflect.Int32:
		idx = int(k.(int32))
	case reflect.Uint32:
		idx = int(k.(uint32))
	case reflect.Int64:
		idx = int(k.(int64))
	case reflect.Uint64:
		idx = int(k.(uint64))
	default:
		return nil
	}
	return this.vArray.Val(idx)
}

//JSONString
////////////////////////////////////////////////////////////////////////////////////////////////////
func (s JSONString) toString() string {
	return string(s)
}

func (s JSONString) equals(v string) bool {
	return s.toString() == v
}

func (s JSONString) equalsA(v JSONString) bool {
	return s.toString() == v.toString()
}

//JSONNumber
////////////////////////////////////////////////////////////////////////////////////////////////////
func (n *JSONNumber) init() {
	n.dbl = 0
	n.iCount = 0
	n.integer = 0
	n.dCount = 0
	n.decimal = 0
	n.eMark = ``
	n.pCount = 0
	n.power = 0
	n.iSign = 1
	n.pMark = ``
}

func (n *JSONNumber) isExp() bool {
	return n.eMark == `e` || n.eMark == `E`
}

func (n *JSONNumber) incInteger(v rune) {
	n.iCount++
	n.integer = (n.integer * 10) + uint64(v-'0')
}

func (n *JSONNumber) incDecimal(v rune) {
	n.dCount++
	n.decimal = (n.decimal * 10) + uint64(v-'0')
}

func (n *JSONNumber) incPower(v rune) {
	n.pCount++
	n.power = (n.power * 10) + uint64(v-'0')
}

func (n *JSONNumber) setISign(i int8) {
	n.iSign = i
}

//func (n *JSONNumber) d(v float64) {
//	n.iSign = number_set_mark
//	n.dbl = v
//}

func (n *JSONNumber) setVal(v interface{}) {
	n.init()
	switch v.(type) {
	case float64:
		n.dbl = v.(float64)
	case int64:
		i := v.(int64)
		if i < 0 {
			n.setISign(-1)
			n.integer = uint64(0 - i)
		} else {
			n.setISign(1)
			n.integer = uint64(i)
		}
		s := fmt.Sprintf(`%d`, n.integer)
		n.iCount = len(s)
	default:
		panic(fmt.Errorf(`unsupported type`))
	}
}

func (n *JSONNumber) toString() string {
	if n.iSign == number_set_mark {
		return fmt.Sprintf("%g", n.dbl)
	}
	s0 := fmt.Sprintf("%d", n.integer)
	if n.decimal > 0 {
		s1 := fmt.Sprintf("%d", n.decimal)
		z0 := n.dCount
		for i := len(s1) - 1; i >= 0; i-- {
			if s1[i] == '0' {
				z0--
			} else {
				break
			}
		}
		s0 += `.` + s1[:z0]
	}

	if n.integer > 0 || n.decimal > 0 {
		if n.iSign < 0 {
			s0 = `-` + s0
		}
	} else {
		return `0`
	}

	if n.isExp() {
		if n.power > 0 {
			s0 += n.eMark + n.pMark + fmt.Sprintf("%d", n.power)
		}
	}
	return s0
}

func (n *JSONNumber) toFloat64() float64 {
	if n.iSign == number_set_mark {
		return n.dbl
	}
	s := n.toString()
	d, _ := strconv.ParseFloat(s, 64)
	return d
}

//JsonElement
////////////////////////////////////////////////////////////////////////////////////////////////////
func newJsonElement() (element *JsonElement) {
	element = new(JsonElement)
	element.key = nil
	element.value = element.newJsonValue()
	return element
}

func (this *JsonElement) newJsonValue() (obj *JsonValue) {
	obj = new(JsonValue)
	obj.e = this
	obj.vType = val_Type_Invalid
	obj.vObject = nil
	obj.vArray = nil
	obj.vString = nil
	obj.vNumber.setVal(float64(0))
	return obj
}

//jsonBase
////////////////////////////////////////////////////////////////////////////////////////////////////
func (this *jsonBase) init() {
	this.childCount = 0
	this.values = make([]*JsonValue, 0, 16)
}

func (this *jsonBase) addElement(k interface{}) (element *JsonElement) {
	s := JSONString{}
	if k != nil {
		s = k.(JSONString)
		element = this.getElement(s)
	}

	if element == nil {
		element = newJsonElement()
		this.values = append(this.values, element.value)
		if k != nil {
			element.key = s
		}
		this.childCount++
	}
	return element
}

func (this *jsonBase) delElement(k interface{}) bool {
	var idx int
	e0 := this.getElement(k)
	if e0 == nil {
		return false
	}

	if reflect.TypeOf(k).Kind() == reflect.Int {
		idx = k.(int)
	} else {
		idx = this.indexOf(e0)
	}
	this.values = append(this.values[0:idx], this.values[idx+1:]...)
	this.childCount--
	return true
}

func (this *jsonBase) getElement(k interface{}) (element *JsonElement) {
	if reflect.TypeOf(k).Kind() == reflect.Int {
		idx := k.(int)
		if idx >= 0 && idx < this.childCount {
			element = this.values[idx].e
		}
	} else {
		for _, v := range this.values {
			if v.e.key.toString() == k.(JSONString).toString() {
				element = v.e
				break
			}
		}
	}
	return element
}

func (this *jsonBase) indexOf(element *JsonElement) int {
	if element == nil {
		return -1
	}

	for i := 0; i < this.childCount; i++ {
		if this.values[i].e == element {
			return i
		}
	}
	return -1
}

func (this *jsonBase) print(depth int, bfmt bool) string {
	var i, j int
	var out string
	if this.childCount == 0 {
		return ``
	}
	depth++
	if bfmt {
		out += "\n"
	}

	for i = 0; i < this.childCount; i++ {
		v0 := this.values[i]
		if bfmt {
			for j = 0; j < depth; j++ {
				out += "\t"
			}
		}

		if v0.e.key != nil {
			out += v0.print_string_ptr(v0.e.key)
			out += `:`
			if bfmt {
				out += " "
			}
		}
		out += v0.print_value(depth, bfmt)
		if i != this.childCount-1 {
			out += `,`
		}

		if bfmt {
			out += "\n"
		}

	}

	if bfmt {
		for i = 0; i < depth-1; i++ {
			out += "\t"
		}
	}
	return out
}
