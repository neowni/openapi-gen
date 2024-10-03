package common

import (
	"fmt"
	"os"
)

func Log(format string, a ...any) {
	_, err := os.Stdout.Write([]byte(fmt.Sprintf(format+"\n", a...)))
	if err != nil {
		panic(err)
	}
}
