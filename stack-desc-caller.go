package catcher

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
)

type Caller struct {
	Func Func
	File string
	Line int
	Ok   bool
}

/*
 callerFromString(Caller.String()) == Caller
*/

var emptyCaller = "<!>"
var bytesEmptyCaller = []byte(emptyCaller)

func (c Caller) BytesSource(b *bytes.Buffer) {
	//fmt.Sprintf("<%v,%s:%d>", c.Func, c.File, c.Line)
	if c.Ok {
		b.WriteByte('<')
		c.Func.BytesSource(b)
		b.WriteByte(',')
		b.WriteString(c.File)
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(c.Line))
		b.WriteByte('>')
	} else {
		b.Write(bytesEmptyCaller)
	}
}

func (c Caller) Bytes() []byte {
	var b = bytes.NewBuffer(make([]byte, 0, 40))
	c.BytesSource(b)
	return b.Bytes()
}

func (c Caller) String() string {
	if c.Ok {
		return string(c.Bytes())
	} else {
		return emptyCaller
	}
}

func callerFromString(s string) Caller {
	s = s[1 : len(s)-1]
	if s[0] != '!' {
		for i := range s {
			if s[i] == '>' {
				for j := len(s) - 1; j >= 0; j-- {
					if s[j] == ':' {
						return Caller{
							Func: funcFromString(s[:i+1]),
							File: s[i+2 : j],
							Line: atoi(s[j+1:]),
							Ok:   true,
						}
					}
				}
			}
		}
		return Caller{}
	}
	return Caller{}
}

func (c Caller) Rel(pack string, rel string) (string, error) {
	if c.Ok {
		f, err := filepath.Rel(rel, c.File)
		if err != nil {
			return "", err
		}
		r, err := c.Func.Rel(pack, rel)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"<%s,%s:%d>", r, f, c.Line), nil
	} else {
		return "<!>", nil
	}
}

func callerFromRuntimeResult(pc uintptr, file string, line int, ok bool) Caller {
	f := runtime.FuncForPC(pc)
	fi, li := f.FileLine(f.Entry())
	return Caller{
		Func: Func{Name: f.Name(), File: fi, Line: li},
		File: file,
		Line: line,
		Ok:   ok,
	}
}
