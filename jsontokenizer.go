package gJSON

func newJsonTokenizer(s JSONString) (tok *JsonTokenizer) {
	tok = new(JsonTokenizer)
	tok.data = s
	tok.size = len(s)
	tok.offset = 0
	tok.err = false
	return tok
}

func (tok *JsonTokenizer) eof() bool {
	if tok == nil {
		return true
	}
	return tok.offset == tok.size
}

func (tok *JsonTokenizer) seek(v int) {
	tok.offset += v
}

func (tok *JsonTokenizer) isStrTag() bool {
	if tok.size-tok.offset < 2 {
		return false
	}
	return tok.data[tok.offset] == '"'
}

func (tok *JsonTokenizer) currCharIs(c rune) (b bool) {
	if tok.eof() {
		return false
	}
	return tok.data[tok.offset] == c
}

func (tok *JsonTokenizer) currStrIs(s string) (b bool) {
	if tok.eof() {
		return false
	}
	size := len(JSONString(s))
	if tok.size-tok.offset < size {
		return false
	}
	header := tok.data[tok.offset : tok.offset+size]
	return string(header) == s
}

func (tok *JsonTokenizer) getCurrChar() rune {
	return tok.at(tok.offset)
}

func (tok *JsonTokenizer) at(i int) rune {
	return tok.data[i]
}

func (tok *JsonTokenizer) currIsArrayBegin() bool {
	return tok.currCharIs('[')
}

func (tok *JsonTokenizer) currIsArrayEnd() bool {
	return tok.currCharIs(']')
}

func (tok *JsonTokenizer) currIsObjectBegin() bool {
	return tok.currCharIs('{')
}

func (tok *JsonTokenizer) currIsObjectEnd() bool {
	return tok.currCharIs('}')
}

func (tok *JsonTokenizer) currStrSize() int {
	if !tok.isStrTag() {
		return 0
	}
	var idx int
	for idx = tok.offset + 1; idx < tok.size; idx++ {
		if tok.data[idx] == '"' && tok.data[idx-1] != '\\' {
			break
		}
	}
	return idx - tok.offset
}

func (tok *JsonTokenizer) indexTrim(v int) {
	if tok.eof() {
		return
	}
	tok.seek(v)
	var idx int
	var c rune
	for idx = tok.offset; idx < tok.size; idx++ {
		c = tok.data[idx]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			break
		}
	}
	tok.offset = idx
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func (tok *JsonTokenizer) parseObject() (obj *JsonObject) {
	if !tok.currIsObjectBegin() {
		return nil
	}
	obj = NewObject()
	tok.indexTrim(1)
	if tok.currIsObjectEnd() {
		tok.seek(1)
		return obj //空的JSON对象
	}
	tok.indexTrim(0)
	k_str := tok.parseString()
	tok.indexTrim(0)
	e0 := obj.addElement(k_str)
	if !tok.currCharIs(':') {
		return nil
	}
	tok.indexTrim(1)
	if e0.value.parse(tok); tok.err {
		return nil
	}
	tok.indexTrim(0)
	for tok.currCharIs(',') {
		tok.indexTrim(1)
		k_str = tok.parseString()
		tok.indexTrim(0)
		e0 = obj.addElement(k_str)
		if !tok.currCharIs(':') {
			return nil
		}
		tok.indexTrim(1)
		if e0.value.parse(tok); tok.err {
			return nil
		}
		tok.indexTrim(0)
	}

	if tok.currIsObjectEnd() {
		tok.indexTrim(1)
		return obj
	}
	tok.err = true
	return nil
}

func (tok *JsonTokenizer) parseArray() (arr *JsonArray) {
	if !tok.currIsArrayBegin() {
		return nil
	}
	arr = NewArray()
	tok.indexTrim(1)
	if tok.currIsArrayEnd() {
		tok.seek(1)
		return arr
	}
	tok.indexTrim(0)
	e0 := arr.addElement()
	tok.indexTrim(0)
	if e0.value.parse(tok); tok.err {
		return nil
	}
	tok.indexTrim(0)
	for tok.currCharIs(',') {
		tok.indexTrim(1)
		e0 = arr.addElement()
		tok.indexTrim(0)
		if e0.value.parse(tok); tok.err {
			return nil
		}
		tok.indexTrim(0)
	}

	if tok.currIsArrayEnd() {
		tok.indexTrim(1)
		return arr
	}
	tok.err = true
	return nil
}

func (tok *JsonTokenizer) parseNumber() (num JSONNumber) {
	num.init()
	if tok.currCharIs('-') {
		num.setISign(-1)
		tok.seek(1)
	}

	for tok.getCurrChar() >= '0' && tok.getCurrChar() <= '9' {
		num.incInteger(tok.getCurrChar())
		tok.seek(1)
	}

	if num.iCount < 1 {
		tok.err = true
		return
	}

	if tok.getCurrChar() == '.' {
		tok.seek(1)
		for tok.getCurrChar() >= '0' && tok.getCurrChar() <= '9' {
			num.incDecimal(tok.getCurrChar())
			tok.seek(1)
		}
	}
	c := tok.getCurrChar()
	if c == 'e' || c == 'E' {
		num.eMark = string(c)
		tok.seek(1)
		if tok.currCharIs('+') {
			num.pMark = `+`
			tok.seek(1)
		} else if tok.currCharIs('-') {
			num.pMark = `-`
			tok.seek(1)
		}

		for tok.getCurrChar() >= '0' && tok.getCurrChar() <= '9' {
			num.incPower(tok.getCurrChar())
			tok.seek(1)
		}
	}
	return
}

func (tok *JsonTokenizer) parseString() (str JSONString) {
	var idx1, idx2 int
	var uc uint
	if tok.eof() || !tok.isStrTag() {
		return nil
	} //非字符串起始标志
	out := make(JSONString, tok.currStrSize())
	idx2 = 0
	for idx1 = tok.offset + 1; idx1 < tok.size && tok.data[idx1] != '"'; idx1++ {
		if tok.data[idx1] != '\\' {
			out[idx2] = tok.data[idx1]
		} else {
			idx1++
			switch tok.data[idx1] {
			case 'b':
				out[idx2] = '\b'
			case 'f':
				out[idx2] = '\f'
			case 'n':
				out[idx2] = '\n'
			case 'r':
				out[idx2] = '\r'
			case 't':
				out[idx2] = '\t'
			case 'u': //transcode utf16 to utf8
				idx1++
				uc = parse_hex4(tok.data[idx1:])
				out[idx2] = rune(uc)
				idx1 += 3
			default:
				out[idx2] = tok.data[idx1]
			}
		}
		idx2++
	}

	if idx1 < tok.size {
		if tok.data[idx1] == '"' {
			idx1++
		}
		tok.offset = idx1
	} else {
		tok.err = true
	}
	return out[:idx2]
}
