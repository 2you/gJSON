package gJSON

const (
	valTypeInvalid = 1
	valTypeNull    = 100
	valTypeTrue    = 101
	valTypeFalse   = 102
	valTypeString  = 103
	valTypeArray   = 104
	valTypeObject  = 105
	valTypeNumber  = 106
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
	numberSetMark = -128
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
