package catcher

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type Func struct {
	Name string
	File string
	Line int
}

/*
 funcFromString(Func.String()) == Func
*/

func (f Func) BytesSource(b *bytes.Buffer) {
	//	fmt.Sprintf("<%s,%s:%d>", f.Name, f.File, f.Line)
	b.WriteByte('<')
	b.WriteString(f.Name)
	b.WriteByte(',')
	b.WriteString(f.File)
	b.WriteByte(':')
	b.WriteString(strconv.Itoa(f.Line))
	b.WriteByte('>')
}

func (f Func) Bytes() []byte {
	var b = bytes.NewBuffer(make([]byte, 0, 25))
	f.BytesSource(b)
	return b.Bytes()
}

func (f Func) String() string {
	return string(f.Bytes())
}

func funcFromString(s string) Func {
	s = s[1 : len(s)-1]
	for i := range s {
		if s[i] == ',' {
			for j := len(s) - 1; j >= 0; j-- {
				if s[j] == ':' {
					return Func{
						Name: s[:i],
						File: s[i+1 : j],
						Line: atoi(s[j+1:]),
					}
				}
			}
		}
	}
	return Func{}
}

func (f Func) Rel(pack string, rel string) (string, error) {
	r2, err := filepath.Rel(rel, f.File)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("<%s,%s:%d>",
		strings.TrimPrefix(strings.TrimPrefix(f.Name, pack), "/"), r2, f.Line), nil
}
