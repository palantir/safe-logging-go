package config_d

import (
	"fmt"
	"time"
)

func CustomLoggingFunction(priority int, msg string, args ...any) {}

func OtherCustomLoggingFunction(priority int, msg string, args ...any) {}

type LoggerStruct struct{}

func (l LoggerStruct) StructReceiverLoggingFunction(priority int, msg string, args ...any) {}

func (l *LoggerStruct) PointerReceiverLoggingFunction(priority int, msg string, args ...any) {}

func LoggingFunctionWithGenerics[T any](priority int, msg string, args ...T) {}

func ParamTests() {
	CustomLoggingFunction(1, "String constant", "safe")

	msg := fmt.Sprintf("%v", time.Now())
	CustomLoggingFunction(1, msg, "unsafe")      // want `CustomLoggingFunction called with unsafe argument: message must be a compile-time constant`
	OtherCustomLoggingFunction(1, msg, "unsafe") // want `OtherCustomLoggingFunction called with unsafe argument: message must be a compile-time constant`

	var logger LoggerStruct
	logger.StructReceiverLoggingFunction(1, msg, "unsafe")  // want `StructReceiverLoggingFunction called with unsafe argument: message must be a compile-time constant`
	logger.PointerReceiverLoggingFunction(1, msg, "unsafe") // want `PointerReceiverLoggingFunction called with unsafe argument: message must be a compile-time constant`

	LoggingFunctionWithGenerics(1, msg, "unsafe")      // want `LoggingFunctionWithGenerics called with unsafe argument: message must be a compile-time constant`
	LoggingFunctionWithGenerics[any](1, msg, 1, "one") // want `LoggingFunctionWithGenerics called with unsafe argument: message must be a compile-time constant`
}
