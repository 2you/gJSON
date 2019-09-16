package gJSON

import (
	"fmt"
	"strconv"
)

func (this *JsonValue) parse(tok *JsonTokenizer) {
	if tok.eof() {
		tok.err = true
		return
	}
	idx := tok.offset
	if tok.currCharIs('"') {
		this.vType = val_Type_String
		this.vString = tok.parseString()
	} else if tok.currCharIs('[') {
		this.vType = val_Type_Array
		this.vArray = tok.parseArray()
	} else if tok.currCharIs('{') {
		this.vType = val_Type_Object
		this.vObject = tok.parseObject()
	} else if (tok.currCharIs('-')) || ((tok.data[idx] >= '0') && (tok.data[idx] <= '9')) {
		this.vType = val_Type_Number
		this.vNumber = tok.parseNumber()
	} else if tok.currStrIs("null") {
		this.vType = val_Type_NULL
		tok.seek(4)
	} else if tok.currStrIs("false") {
		this.vType = val_Type_False
		tok.seek(5)
	} else if tok.currStrIs("true") {
		this.vType = val_Type_True
		tok.seek(4)
	} else {
		tok.err = true
	}
}

func (this *JsonValue) AsArray() (arr *JsonArray, err error) {
	if this.IsArray() {
		return this.vArray, nil
	}
	return nil, fmt.Errorf("this node is not json array")
}

func (this *JsonValue) AsBool() (bool, error) {
	if this.vType == val_Type_True {
		return true, nil
	} else if this.vType == val_Type_False {
		return false, nil
	}
	return false, fmt.Errorf(`json node is not true or false`)
}

func (this *JsonValue) AsBoolDef(def bool) bool {
	if b, e := this.AsBool(); e == nil {
		return b
	} else {
		switch this.vType {
		case val_Type_Number:
			if f, _ := this.AsFloat64(); f == 0 {
				return false
			}
			return true
		case val_Type_String:
			if f, e := strconv.ParseFloat(string(this.vString), 64); e == nil {
				if n := int(f); n == 0 {
					return false
				} else {
					return true
				}
			} else {
				return def
			}
		default:
			return def
		}
	}
}

func (this *JsonValue) AsDouble() (float64, error) {
	if f, e := this.AsFloat64(); e == nil {
		return f, nil
	} else {
		return 0, e
	}
}

func (this *JsonValue) AsDoubleDef(def float64) (val float64) {
	return this.AsFloat64Def(def)
}

func (this *JsonValue) AsFloat32() (val float32, err error) {
	if f, e := this.AsFloat64(); e == nil {
		return float32(f), nil
	} else {
		return 0, e
	}
}

func (this *JsonValue) AsFloat32Def(def float32) (val float32) {
	return float32(this.AsFloat64Def(float64(def)))
}

func (this *JsonValue) AsFloat64() (val float64, err error) {
	if this.IsNumber() {
		return this.vNumber.toFloat64(), nil
	}
	return 0, fmt.Errorf(`json node is not number`)
}

func (this *JsonValue) AsFloat64Def(def float64) (val float64) {
	switch this.vType {
	case val_Type_Number:
		return this.vNumber.toFloat64()
	case val_Type_String:
		if f64, err := strconv.ParseFloat(string(this.vString), 64); err == nil {
			return f64
		} else {
			return def
		}
	case val_Type_True:
		return 1
	case val_Type_False:
		return 0
	default:
		return def
	}
}

func (this *JsonValue) AsInt() (int, error) {
	if i, e := this.AsFloat64(); e == nil {
		return int(i), nil
	} else {
		return 0, e
	}
}

func (this *JsonValue) AsIntDef(def int) (val int) {
	return int(this.AsFloat64Def(float64(def)))
}

func (this *JsonValue) AsInt32() (int32, error) {
	if i, e := this.AsFloat64(); e == nil {
		return int32(i), nil
	} else {
		return 0, e
	}
}

func (this *JsonValue) AsInt32Def(def int32) (val int32) {
	return int32(this.AsFloat64Def(float64(def)))
}

func (this *JsonValue) AsInt64() (val int64, err error) {
	if f, e := this.AsFloat64(); e == nil {
		return int64(f), nil
	} else {
		return 0, fmt.Errorf(`json node is not number`)
	}
}

func (this *JsonValue) AsInt64Def(def int64) (val int64) {
	return int64(this.AsFloat64Def(float64(def)))
}

func (this *JsonValue) AsObject() (obj *JsonObject, err error) {
	if this.IsObject() {
		return this.vObject, nil
	}
	return nil, fmt.Errorf("this node is not json object")
}

func (this *JsonValue) AsString() (str string, err error) {
	if this.IsString() {
		return this.vString.toString(), nil
	}
	return ``, fmt.Errorf(`json node is not string`)
}

func (this *JsonValue) AsStringDef(def string) (val string) {
	switch this.vType {
	case val_Type_Object:
		return this.vObject.AsString()
	case val_Type_Array:
		return this.vArray.AsString()
	case val_Type_String:
		return this.vString.toString()
	case val_Type_Number:
		return this.vNumber.toString()
	case val_Type_True:
		return `true`
	case val_Type_False:
		return `false`
	case val_Type_NULL:
		return `null`
	default:
		return def
	}
}

func (this *JsonValue) IsArray() bool {
	return this.vType == val_Type_Array
}

func (this *JsonValue) IsBool() bool {
	return this.IsTrue() || this.IsFalse()
}

func (this *JsonValue) IsTrue() bool {
	return this.vType == val_Type_True
}

func (this *JsonValue) IsFalse() bool {
	return this.vType == val_Type_False
}

func (this *JsonValue) IsNull() bool {
	return this.vType == val_Type_NULL
}

func (this *JsonValue) IsNumber() bool {
	return this.vType == val_Type_Number
}

func (this *JsonValue) IsObject() bool {
	return this.vType == val_Type_Object
}

func (this *JsonValue) IsString() bool {
	return this.vType == val_Type_String
}

func (this *JsonValue) print_value(depth int, bfmt bool) string {
	out := ``
	switch this.vType {
	case val_Type_NULL:
		out = `null`
	case val_Type_False:
		out = `false`
	case val_Type_True:
		out = `true`
	case val_Type_Number:
		out = this.print_number()
	case val_Type_String:
		out = this.print_string()
	case val_Type_Array:
		out = this.print_array(depth, bfmt)
	case val_Type_Object:
		out = this.print_object(depth, bfmt)
	}
	return out
}

func (this *JsonValue) print_object(depth int, bfmt bool) string {
	return this.vObject.print(depth, bfmt)
}

func (this *JsonValue) print_array(depth int, bfmt bool) string {
	return this.vArray.print(depth, bfmt)
}

func (this *JsonValue) print_string() string {
	return this.print_string_ptr(this.vString)
}

func (this *JsonValue) print_string_ptr(s JSONString) string {
	flag := 0
	if s == nil {
		return `""`
	}

	for _, c := range s {
		if (c > 0 && c < 32) || c == '"' || c == '\\' {
			flag |= 1
		} else {
			flag |= 0
		}
	}

	if flag == 0 {
		return `"` + s.toString() + `"`
	}
	i := 0
	var out []rune
	out = append(out, '"')
	for i < len(s) {
		c := s[i]
		i++
		if uint(c) > 31 && c != '"' && c != '\\' {
			out = append(out, c)
		} else {
			out = append(out, '\\')
			switch c {
			case '\\':
				out = append(out, '\\')
			case '"':
				out = append(out, '"')
			case '\b':
				out = append(out, 'b')
			case '\f':
				out = append(out, 'f')
			case '\n':
				out = append(out, 'n')
			case '\r':
				out = append(out, 'r')
			case '\t':
				out = append(out, 't')
			default:
				out = append(out, []rune(fmt.Sprintf("u%04x", c))...)
			}
		}
	}
	out = append(out, '"')
	return string(out)
}

func (this *JsonValue) print_number() string {
	return this.vNumber.toString()
}
