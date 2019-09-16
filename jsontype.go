package gJSON

const (
	val_Type_Invalid = 1
	val_Type_NULL    = 100
	val_Type_True    = 101
	val_Type_False   = 102
	val_Type_String  = 103
	val_Type_Array   = 104
	val_Type_Object  = 105
	val_Type_Number  = 106
)

type JSONString []rune
type ValType int

type JSONNumber struct {
	dbl     float64 //数值
	iCount  int     //整数位数
	integer uint64  //整数部分
	dCount  int     //小数位数
	decimal uint64  //小数部分
	eMark   string  //指数标记
	pCount  int     //幂位数
	power   uint64  //幂
	iSign   int8    //整数或者指数部分正负符号
	pMark   string  //冥符号标记
}

const (
	number_set_mark = -128
)

type JsonElement struct {
	key   JSONString
	value *JsonValue
}

type JsonValue struct {
	e       *JsonElement
	vType   ValType
	vObject *JsonObject
	vArray  *JsonArray
	vString JSONString
	vNumber JSONNumber
}

type jsonBase struct {
	childCount int
	values     []*JsonValue
}

type JsonObject struct {
	jsonBase
}

type JsonArray struct {
	jsonBase
}

type JsonNode struct {
	vObject *JsonObject
	vArray  *JsonArray
}

type JsonTokenizer struct {
	data   JSONString
	size   int
	offset int
	err    bool
}
