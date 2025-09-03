package c // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/c\.TestStructWithCycle: @Unsafe, safe-logging-go/safeloggingtests/c\.TestStructWithUnsafeFieldAndCycle: @Unsafe\]\]`

import (
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type TestStructWithUnsafeFieldAndCycle struct {
	Cycle       *TestStructWithCycle
	UnsafeField string `safelogging:"@Unsafe"` // want UnsafeField:`safe-logging-go/safeloggingtests/c\.TestStructWithUnsafeFieldAndCycle\.UnsafeField: @Unsafe \(tag on struct field "safe-logging-go/safeloggingtests/c\.TestStructWithUnsafeFieldAndCycle\.UnsafeField"\)`
}

// TestStructWithCycle should be unsafe because TestStructWithUnsafeFieldAndCycle has an unsafe field
type TestStructWithCycle struct {
	Cycle *TestStructWithUnsafeFieldAndCycle // want Cycle:`safe-logging-go/safeloggingtests/c\.TestStructWithCycle\.Cycle: @Unsafe \(tag on struct field "safe-logging-go/safeloggingtests/c\.TestStructWithUnsafeFieldAndCycle\.UnsafeField"\)`
}

func ParamTests() {
	testStruct := TestStructWithUnsafeFieldAndCycle{}
	svc1log.SafeParam("testParam", testStruct) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/c\.TestStructWithUnsafeFieldAndCycle", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/c\.TestStructWithUnsafeFieldAndCycle\.UnsafeField"\)$`

	testStruct2 := TestStructWithCycle{}
	svc1log.SafeParam("testParam", testStruct2) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/c\.TestStructWithCycle", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/c\.TestStructWithUnsafeFieldAndCycle\.UnsafeField"\)$`
}
