package config_d

import (
	"fmt"
	"time"
)

func CustomLoggingFunction(priority int, msg string, args ...any) {
	fmt.Printf("%d: %s - %v\n", priority, msg, args)
}

func ParamTests() {
	CustomLoggingFunction(1, "String constant", "safe")

	msg := fmt.Sprintf("%v", time.Now())
	CustomLoggingFunction(1, msg, "unsafe") // want `CustomLoggingFunction called with unsafe argument: message must be a compile-time constant`
}
