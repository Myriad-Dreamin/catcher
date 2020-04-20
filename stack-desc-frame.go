package catcher

import "fmt"

type Frames []Frame
type Frame interface {
	Dump() string
	RelDump(pack, rel string) (string, error)
	String() string
	Bytes() []byte
	Error() string
	ReleaseError()
	GetPos() Caller
	GetCode() int
	GetErr() string
}

func (fs Frames) String() (res string) {
	for i := range fs {
		res += fmt.Sprintf("%d <- %s\n", len(fs)-i, fs[i].Dump())
	}
	return res
}

func (fs Frames) Rel(pack, rel string) (res string, err error) {
	var c string
	for i := range fs {
		c, err = fs[i].RelDump(pack, rel)
		if err != nil {
			return "", err
		}
		res += fmt.Sprintf("%d <- %s\n", len(fs)-i, c)
	}
	return res, nil
}
