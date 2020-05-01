package out

import "fmt"

func Log(msg string) {
	fmt.Printf("> %s\n", msg)
}

func Logf(format string, args ...interface{}) {
	Log(fmt.Sprintf(format, args...))
}

