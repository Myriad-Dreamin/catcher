package catcher

import (
	"fmt"
	"strconv"
)

func atoi(s string) int {
	c, err := strconv.Atoi(s)
	if err != nil {
		c = CodeDeserializeError
		if reportBad {
			fmt.Println(err)
		}
	}
	return c
}
