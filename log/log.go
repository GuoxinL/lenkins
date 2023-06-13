package log

import "fmt"

func Infof(format string, a ...any) {
	fmt.Println(fmt.Sprintf(format, a...))
}
func Errorf(format string, a ...any) {
	fmt.Println(fmt.Sprintf(format, a...))
}
