package catcher

import "strconv"

const CodeDeserializeError = -1

const (
	Prod = iota
	Debug
)

var errorFlag = Debug

// SetErrorFlag set error flag. if set, wrapToStackError will collect runtime
// information of Wrap's caller
func SetErrorFlag(f int) {
	errorFlag = f
}
func GetErrorFlag() int {
	return errorFlag
}

var _codeDescriptor = strconv.Itoa

// SetCodeDescriptor set code descriptor which will be called in frame-impl.Dump
func SetCodeDescriptor(descriptor func(int) string) {
	_codeDescriptor = descriptor
}
func GetCodeDescriptor() func(int) string {
	return _codeDescriptor
}

var reportBad = true

// SetReportBad set boolean var that describes whether inner.atoi will report bad
func SetReportBad(x bool) {
	reportBad = x
}
func GetReportBad() bool {
	return reportBad
}

var (
	magic      = "<84f4446f>"
	magicBytes = []byte(magic)
)

func SetMagic(x string) {
	magic = x
	magicBytes = []byte(magic)
}
func GetMagic() string {
	return magic
}
