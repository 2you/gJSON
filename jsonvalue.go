package gJSON

import (
	"fmt"
	"strconv"
)

func (v *JsonValue) parse(tok *JsonTokenizer) {
	if tok.eof() {
		tok.err = true
		return
	}
	idx := tok.offset
	if tok.currCharIs('"') {
		v.vType = valTypeString
		v.vString = tok.parseString()
	} else if tok.currCharIs('[') {
		v.vType = valTypeArray
		v.vArray = tok.parseArray()
	} else if tok.currCharIs('{') {
		v.vType = valTypeObject
		v.vObject = tok.parseObject()
	} else if (tok.currCharIs('-')) || ((tok.data[idx] >= '0') && (tok.data[idx] <= '9')) {
		v.vType = valTypeNumber
		v.vNumber = tok.parseNumber()
	} else if tok.currStrIs("null") {
		v.vType = valTypeNull
		tok.seek(4)
	} else if tok.currStrIs("false") {
		v.vType = valTypeFalse
		tok.seek(5)
	} else if tok.currStrIs("true") {
		v.vType = valTypeTrue
		tok.seek(4)
	} else {
		tok.err = true
	}
}

func (v *JsonValue) AsArray() (arr *JsonArray, err error) {
	if v.IsArray() {
		return v.vArray, nil
	}
	return nil, fmt.Errorf("v node is not json array")
}

func (v *JsonValue) AsBool() (bool, error) {
	if v.vType == valTypeTrue {
		return true, nil
	} else if v.vType == valTypeFalse {
		return false, nil
	}
	return false, fmt.Errorf(`json node is not true or false`)
}

func (v *JsonValue) AsBoolDef(def bool) bool {
	if b, e := v.AsBool(); e == nil {
		return b
	} else {
		switch v.vType {
		case valTypeNumber:
			if f, _ := v.AsFloat64(); f == 0 {
				return false
			}
			return true
		case valTypeString:
			if f, e := strconv.ParseFloat(string(v.vString), 64); e == nil {
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

func (v *JsonValue) AsDouble() (float64, error) {
	if f, e := v.AsFloat64(); e == nil {
		return f, nil
	} else {
		return 0, e
	}
}

func (v *JsonValue) AsDoubleDef(def float64) (val float64) {
	return v.AsFloat64Def(def)
}

func (v *JsonValue) AsFloat32() (val float32, err error) {
	if f, e := v.AsFloat64(); e == nil {
		return float32(f), nil
	} else {
		return 0, e
	}
}

func (v *JsonValue) AsFloat32Def(def float32) (val float32) {
	return float32(v.AsFloat64Def(float64(def)))
}

func (v *JsonValue) AsFloat64() (val float64, err error) {
	if v.IsNumber() {
		return v.vNumber.toFloat64(), nil
	}
	return 0, fmt.Errorf(`json node is not number`)
}

func (v *JsonValue) AsFloat64Def(def float64) (val float64) {
	switch v.vType {
	case valTypeNumber:
		return v.vNumber.toFloat64()
	case valTypeString:
		if f64, err := strconv.ParseFloat(string(v.vString), 64); err == nil {
			return f64
		} else {
			return def
		}
	case valTypeTrue:
		return 1
	case valTypeFalse:
		return 0
	default:
		return def
	}
}

func (v *JsonValue) AsInt() (int, error) {
	if i, e := v.AsFloat64(); e == nil {
		return int(i), nil
	} else {
		return 0, e
	}
}

func (v *JsonValue) AsIntDef(def int) (val int) {
	return int(v.AsFloat64Def(float64(def)))
}

func (v *JsonValue) AsInt32() (int32, error) {
	if i, e := v.AsFloat64(); e == nil {
		return int32(i), nil
	} else {
		return 0, e
	}
}

func (v *JsonValue) AsInt32Def(def int32) (val int32) {
	return int32(v.AsFloat64Def(float64(def)))
}

func (v *JsonValue) AsInt64() (val int64, err error) {
	if f, e := v.AsFloat64(); e == nil {
		return int64(f), nil
	} else {
		return 0, fmt.Errorf(`json node is not number`)
	}
}

func (v *JsonValue) AsInt64Def(def int64) (val int64) {
	return int64(v.AsFloat64Def(float64(def)))
}

func (v *JsonValue) AsObject() (obj *JsonObject, err error) {
	if v.IsObject() {
		return v.vObject, nil
	}
	return nil, fmt.Errorf("v node is not json object")
}

func (v *JsonValue) AsString() (str string, err error) {
	if v.IsString() {
		return v.vString.toString(), nil
	}
	return ``, fmt.Errorf(`json node is not string`)
}

func (v *JsonValue) AsStringDef(def string) (val string) {
	switch v.vType {
	case valTypeObject:
		return v.vObject.AsString()
	case valTypeArray:
		return v.vArray.AsString()
	case valTypeString:
		return v.vString.toString()
	case valTypeNumber:
		return v.vNumber.toString()
	case valTypeTrue:
		return `true`
	case valTypeFalse:
		return `false`
	case valTypeNull:
		return `null`
	default:
		return def
	}
}

func (v *JsonValue) IsArray() bool {
	return v.vType == valTypeArray
}

func (v *JsonValue) IsBool() bool {
	return v.IsTrue() || v.IsFalse()
}

func (v *JsonValue) IsTrue() bool {
	return v.vType == valTypeTrue
}

func (v *JsonValue) IsFalse() bool {
	return v.vType == valTypeFalse
}

func (v *JsonValue) IsNull() bool {
	return v.vType == valTypeNull
}

func (v *JsonValue) IsNumber() bool {
	return v.vType == valTypeNumber
}

func (v *JsonValue) IsObject() bool {
	return v.vType == valTypeObject
}

func (v *JsonValue) IsString() bool {
	return v.vType == valTypeString
}

func (v *JsonValue) print_value(depth int, bfmt bool) string {
	out := ``
	switch v.vType {
	case valTypeNull:
		out = `null`
	case valTypeFalse:
		out = `false`
	case valTypeTrue:
		out = `true`
	case valTypeNumber:
		out = v.print_number()
	case valTypeString:
		out = v.print_string()
	case valTypeArray:
		out = v.print_array(depth, bfmt)
	case valTypeObject:
		out = v.print_object(depth, bfmt)
	}
	return out
}

func (v *JsonValue) print_object(depth int, bfmt bool) string {
	return v.vObject.print(depth, bfmt)
}

func (v *JsonValue) print_array(depth int, bfmt bool) string {
	return v.vArray.print(depth, bfmt)
}

func (v *JsonValue) print_string() string {
	return v.print_string_ptr(v.vString)
}

func (v *JsonValue) print_string_ptr(s JSONString) string {
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

func (v *JsonValue) print_number() string {
	return v.vNumber.toString()
}
