package plugins

import (
	"fmt"
	"strings"
)

const (
	pattern = "${%s}"
)

func Replace(format, key, value string) string {
	old := fmt.Sprintf(pattern, key)
	return strings.Replace(format, old, value, -1)
}
