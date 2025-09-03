package k // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/k\.TestStruct: @Unsafe\]\]`

import (
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type TestStruct struct {
	UnsafeField *string `safelogging:"@Unsafe"` // want UnsafeField:`safe-logging-go/safeloggingtests/k\.TestStruct\.UnsafeField: @Unsafe \(tag on struct field "safe-logging-go/safeloggingtests/k\.TestStruct\.UnsafeField"\)`
}

func ParamTests() {
	svc1log.SafeParam("testParam", TestStruct{}) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/k\.TestStruct", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/k\.TestStruct\.UnsafeField"\)$`

	// safelogging:@Allow: allow parameter usage
	svc1log.SafeParam("testParam", TestStruct{})
}
