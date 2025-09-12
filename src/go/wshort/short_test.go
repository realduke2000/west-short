package wshort

import (
	"fmt"
	"testing"
)

func TestGenID(t *testing.T) {
	/*
		[83 247 60 0 203 176]
		[15 247 60 0 203 176]
		1007031513
		1007077458
	*/
	for i := 0; i < 10; i++ {
		s := generateId()
		fmt.Printf("id=%s\n", s)
	}
}
