package catcher

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"
)

type DatabaseCode int

type DuplicateError struct{ Data string }

func (d DuplicateError) Error() string {
	return fmt.Sprintf(`duplicate key at field "%v"`, d.Data)
}

type UserDefinedError struct{ Data []interface{} }

func (d UserDefinedError) Error() string {
	return fmt.Sprintf("user defined: %v", d.Data)
}

func DatabaseError(code DatabaseCode, errDesc ...interface{}) error {
	switch code {
	case 1:
		return WrapN(BaseSkip+1, int(code), DuplicateError{Data: errDesc[0].(string)})
	default:
		return WrapN(BaseSkip+1, int(code), UserDefinedError{Data: errDesc})
	}
}

func TestExampleDatabaseError(t *testing.T) {
	err := DatabaseError(1, "id")
	fmt.Println(d.Describe(err))
}

var d = Describer{
	Pack: "github.com/Myriad-Dreamin",
	Rel:  handle(filepath.Abs("")),
}

func TestExampleSimpleNestWrap(t *testing.T) {
	err := outerLogic()
	fmt.Println(Describe(err))
}

func TestExampleFromString(t *testing.T) {
	frameToTransfer := Wrap(233, errors.New("QAQ"))
	if frameToTransfer == nil {
		t.Fatal("empty")
	}
	transferObject := frameToTransfer.Error()
	rawFrame, ok := FromString(transferObject)
	if !ok {
		t.Fatal("deserialize error, or not a frame")
	}
	fmt.Println(rawFrame.GetPos().Func.Name, rawFrame.GetCode(), rawFrame.GetErr())
	fmt.Println(rawFrame.GetPos().File, rawFrame.GetPos().Line)
}

func innerLogic() error {
	return Wrap(233, errors.New("QAQ"))
}

func outerLogic() error {
	err := innerLogic()
	if err != nil {
		return Wrap(666, err)
	}
	return nil
}

func TestExampleNestWrap(t *testing.T) {
	err := outerLogic()
	fmt.Println(d.Describe(err))
}

func TestExampleWrap(t *testing.T) {
	frame := Wrap(233, errors.New("QAQ"))
	fmt.Println(d.Describe(frame))
}

func handle(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}
