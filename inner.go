package catcher

import (
	"log"
	"strconv"
)

func atoi(s string) int {
	c, err := strconv.Atoi(s)
	if err != nil {
		c = CodeDeserializeError
		if reportBad {
			log.Print(err)
		}
	}
	return c
}
