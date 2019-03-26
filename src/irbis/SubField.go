package irbis

type SubField struct {
	Code  rune
	Value string
}

func NewSubField(code rune, value string) *SubField {
	return &SubField{Code: code, Value: value}
}

func (subfield *SubField) Decode(text string) {
	runes := []rune(text)
	subfield.Code = runes[0]
	subfield.Value = text[1:]
}

func (subfield *SubField) Encode() string {
	return "^" + string(subfield.Code) + subfield.Value
}
