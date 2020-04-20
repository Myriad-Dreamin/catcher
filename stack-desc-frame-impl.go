package catcher

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type frameImpl struct {
	pos  Caller
	code int
	err  string
}

/*
 stackErrorFromString(frameImpl.String()) == frameImpl
*/
//BenchmarkWrapChainWithOutCollectInfo-4   	 1000000	      1148 ns/op
//BenchmarkWrapChainWithOutCollectInfo-4   	  670528	      1580 ns/op
//BenchmarkWrapChainWithOutCollectInfo-4   	  460164	      2577 ns/op

func (g frameImpl) BytesSource(b *bytes.Buffer) {
	b.WriteString(magic)
	b.WriteString("pos:")
	g.pos.BytesSource(b)
	b.WriteByte(',')
	b.WriteString(magic)
	b.WriteString("code:")
	b.WriteString(strconv.Itoa(g.code))
	b.WriteByte(',')
	b.WriteString(magic)
	b.WriteString("err:")
	b.WriteString(g.err)
}

func (g frameImpl) Bytes() []byte {
	var b = bytes.NewBuffer(make([]byte, 0, 100))
	g.BytesSource(b)
	return b.Bytes()
}

func (g frameImpl) String() string {
	return string(g.Bytes())
}

func stackErrorFromBytes(s []byte) (*frameImpl, bool) {
	if !bytes.HasPrefix(s, magicBytes) {
		return &frameImpl{}, false
	}
	x := bytes.SplitN(s, magicBytes, 4)[1:4]

	return &frameImpl{
		pos:  callerFromString(string(x[0][4 : len(x[0])-1])),
		code: atoi(string(x[1][5 : len(x[1])-1])),
		err:  string(x[2][4:]),
	}, true
}

func stackErrorFromString(s string) (*frameImpl, bool) {
	if !strings.HasPrefix(s, magic) {
		return &frameImpl{}, false
	}
	x := strings.SplitN(s, magic, 4)[1:4]

	return &frameImpl{
		pos:  callerFromString(x[0][4 : len(x[0])-1]),
		code: atoi(x[1][5 : len(x[1])-1]),
		err:  x[2][4:],
	}, true
}

func wrapStringToStackError(skip int, code int, err string) Frame {
	if errorFlag == Debug {
		return &frameImpl{
			pos:  callerFromRuntimeResult(runtime.Caller(skip)),
			code: code,
			err:  err,
		}
	}

	return &frameImpl{
		pos:  Caller{},
		code: code,
		err:  err,
	}
}

func wrapCodeToStackError(skip int, code int) Frame {
	return wrapStringToStackError(skip, code, "")
}

func wrapToStackError(skip int, code int, err error) Frame {
	return wrapStringToStackError(skip, code, err.Error())
}

func (g frameImpl) GetPos() Caller {
	return g.pos
}

func (g frameImpl) GetCode() int {
	return g.code
}

func (g frameImpl) GetErr() string {
	return g.err
}

func (g frameImpl) Dump() string {
	return fmt.Sprintf(
		"<pos:%v,code:%s,err:%s>", g.pos, _codeDescriptor(g.code), g.err)
}

func (g frameImpl) RelDump(pack, rel string) (string, error) {
	p, err := g.pos.Rel(pack, rel)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"<pos:%v,code:%s,err:%s>", p, _codeDescriptor(g.code), g.err), nil
}

func (g frameImpl) Error() string {
	return g.String()
}

func (g *frameImpl) ReleaseError() {
	g.err = ""
}
