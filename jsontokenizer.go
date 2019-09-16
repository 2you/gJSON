package gJSON

func newJsonTokenizer(s JSONString) (tok *JsonTokenizer) {
	tok = new(JsonTokenizer)
	tok.data = s
	tok.size = len(s)
	tok.offset = 0
	tok.err = false
	return tok
}

func (this *JsonTokenizer) eof() bool {
	if this == nil {
		return true
	}
	return this.offset == this.size
}

func (this *JsonTokenizer) seek(v int) {
	this.offset += v
}

func (this *JsonTokenizer) isStrTag() bool {
	if this.size-this.offset < 2 {
		return false
	}
	return this.data[this.offset] == '"'
}

func (this *JsonTokenizer) currCharIs(c rune) (b bool) {
	if this.eof() {
		return false
	}
	return this.data[this.offset] == c
}

func (this *JsonTokenizer) currStrIs(s string) (b bool) {
	if this.eof() {
		return false
	}
	size := len(JSONString(s))
	if this.size-this.offset < size {
		return false
	}
	header := this.data[this.offset : this.offset+size]
	return string(header) == s
}

func (this *JsonTokenizer) getCurrChar() rune {
	return this.at(this.offset)
}

func (this *JsonTokenizer) at(i int) rune {
	return this.data[i]
}

func (this *JsonTokenizer) currIsArrayBegin() bool {
	return this.currCharIs('[')
}

func (this *JsonTokenizer) currIsArrayEnd() bool {
	return this.currCharIs(']')
}

func (this *JsonTokenizer) currIsObjectBegin() bool {
	return this.currCharIs('{')
}

func (this *JsonTokenizer) currIsObjectEnd() bool {
	return this.currCharIs('}')
}

func (this *JsonTokenizer) currStrSize() int {
	if !this.isStrTag() {
		return 0
	}
	var idx int
	for idx = this.offset + 1; idx < this.size; idx++ {
		if this.data[idx] == '"' && this.data[idx-1] != '\\' {
			break
		}
	}
	return idx - this.offset
}

func (this *JsonTokenizer) indexTrim(v int) {
	if this.eof() {
		return
	}
	this.seek(v)
	var idx int
	var c rune
	for idx = this.offset; idx < this.size; idx++ {
		c = this.data[idx]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			break
		}
	}
	this.offset = idx
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func (this *JsonTokenizer) parseObject() (obj *JsonObject) {
	if !this.currIsObjectBegin() {
		return nil
	}
	obj = NewObject()
	this.indexTrim(1)
	if this.currIsObjectEnd() {
		this.seek(1)
		return obj //空的JSON对象
	}
	this.indexTrim(0)
	k_str := this.parseString()
	this.indexTrim(0)
	e0 := obj.addElement(k_str)
	if !this.currCharIs(':') {
		return nil
	}
	this.indexTrim(1)
	if e0.value.parse(this); this.err {
		return nil
	}
	this.indexTrim(0)
	for this.currCharIs(',') {
		this.indexTrim(1)
		k_str = this.parseString()
		this.indexTrim(0)
		e0 = obj.addElement(k_str)
		if !this.currCharIs(':') {
			return nil
		}
		this.indexTrim(1)
		if e0.value.parse(this); this.err {
			return nil
		}
		this.indexTrim(0)
	}

	if this.currIsObjectEnd() {
		this.indexTrim(1)
		return obj
	}
	this.err = true
	return nil
}

func (this *JsonTokenizer) parseArray() (arr *JsonArray) {
	if !this.currIsArrayBegin() {
		return nil
	}
	arr = NewArray()
	this.indexTrim(1)
	if this.currIsArrayEnd() {
		this.seek(1)
		return arr
	}
	this.indexTrim(0)
	e0 := arr.addElement()
	this.indexTrim(0)
	if e0.value.parse(this); this.err {
		return nil
	}
	this.indexTrim(0)
	for this.currCharIs(',') {
		this.indexTrim(1)
		e0 = arr.addElement()
		this.indexTrim(0)
		if e0.value.parse(this); this.err {
			return nil
		}
		this.indexTrim(0)
	}

	if this.currIsArrayEnd() {
		this.indexTrim(1)
		return arr
	}
	this.err = true
	return nil
}

func (this *JsonTokenizer) parseNumber() (num JSONNumber) {
	num.init()
	if this.currCharIs('-') {
		num.setISign(-1)
		this.seek(1)
	}

	for this.getCurrChar() >= '0' && this.getCurrChar() <= '9' {
		num.incInteger(this.getCurrChar())
		this.seek(1)
	}

	if num.iCount < 1 {
		this.err = true
		return
	}

	if this.getCurrChar() == '.' {
		this.seek(1)
		for this.getCurrChar() >= '0' && this.getCurrChar() <= '9' {
			num.incDecimal(this.getCurrChar())
			this.seek(1)
		}
	}
	c := this.getCurrChar()
	if c == 'e' || c == 'E' {
		num.eMark = string(c)
		this.seek(1)
		if this.currCharIs('+') {
			num.pMark = `+`
			this.seek(1)
		} else if this.currCharIs('-') {
			num.pMark = `-`
			this.seek(1)
		}

		for this.getCurrChar() >= '0' && this.getCurrChar() <= '9' {
			num.incPower(this.getCurrChar())
			this.seek(1)
		}
	}
	return
}

func (this *JsonTokenizer) parseString() (str JSONString) {
	var idx1, idx2 int
	var uc uint
	if this.eof() || !this.isStrTag() {
		return nil
	} //非字符串起始标志
	out := make(JSONString, this.currStrSize())
	idx2 = 0
	for idx1 = this.offset + 1; idx1 < this.size && this.data[idx1] != '"'; idx1++ {
		if this.data[idx1] != '\\' {
			out[idx2] = this.data[idx1]
		} else {
			idx1++
			switch this.data[idx1] {
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
				uc = parse_hex4(this.data[idx1:])
				out[idx2] = rune(uc)
				idx1 += 3
			default:
				out[idx2] = this.data[idx1]
			}
		}
		idx2++
	}

	if idx1 < this.size {
		if this.data[idx1] == '"' {
			idx1++
		}
		this.offset = idx1
	} else {
		this.err = true
	}
	return out[:idx2]
}
