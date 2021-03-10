package internal

import (
	"fmt"
	"os"
)

func Warnf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "warning: "+format+"\n", args...)
}
