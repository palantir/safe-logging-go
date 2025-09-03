package config_a // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/config/config_a\.TestStructWithFieldSafetySetViaConfig: @DoNotLog\]\]`

import (
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type TestStructWithFieldSafetySetViaConfig struct {
	FieldToMarkUnsafe   string // want FieldToMarkUnsafe:`safe-logging-go/safeloggingtests/config/config_a\.TestStructWithFieldSafetySetViaConfig\.FieldToMarkUnsafe: @Unsafe \(configuration specified log safety value for struct field "safe-logging-go/safeloggingtests/config/config_a\.TestStructWithFieldSafetySetViaConfig\.FieldToMarkUnsafe"\)`
	FieldToMarkDoNotLog string // want FieldToMarkDoNotLog:`safe-logging-go/safeloggingtests/config/config_a\.TestStructWithFieldSafetySetViaConfig\.FieldToMarkDoNotLog: @DoNotLog \(configuration specified log safety value for struct field "safe-logging-go/safeloggingtests/config/config_a\.TestStructWithFieldSafetySetViaConfig\.FieldToMarkDoNotLog"\)`
}

func ParamTests() {
	testStruct := TestStructWithFieldSafetySetViaConfig{}
	svc1log.SafeParam("testParam", testStruct.FieldToMarkUnsafe)   // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: configuration specified log safety value for struct field "safe-logging-go/safeloggingtests/config/config_a\.TestStructWithFieldSafetySetViaConfig\.FieldToMarkUnsafe"\)$`
	svc1log.SafeParam("testParam", testStruct.FieldToMarkDoNotLog) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: configuration specified log safety value for struct field "safe-logging-go/safeloggingtests/config/config_a\.TestStructWithFieldSafetySetViaConfig\.FieldToMarkDoNotLog"\)$`
	svc1log.SafeParam("testParam", testStruct)                     // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/config/config_a\.TestStructWithFieldSafetySetViaConfig", which is @DoNotLog \(reason: configuration specified log safety value for struct field "safe-logging-go/safeloggingtests/config/config_a\.TestStructWithFieldSafetySetViaConfig\.FieldToMarkDoNotLog"\)$`
}
