package wordmaker

import (
	"fmt"
)

var DEBUG bool

func Debug(val string) {
	if DEBUG {
		fmt.Print(val)
		fmt.Print("\n")
	}
}

func Debugf(val string, params ...interface{}) {
	if DEBUG {
		fmt.Printf(val, params...)
		fmt.Print("\n")
	}
}
