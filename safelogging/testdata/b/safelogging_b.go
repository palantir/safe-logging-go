package b // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevel: @DoNotLog, safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevelFromEmbedded: @DoNotLog, safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevelFromMap: @DoNotLog, safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevelFromPtr: @DoNotLog, safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevelFromSlice: @DoNotLog\]\]`

import (
	"safe-logging-go/safeloggingtests/a"

	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type TestStructWithInheritedLevel struct {
	TestStructField a.TestStruct // want TestStructField:`safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevel\.TestStructField: @DoNotLog \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)`
}

type TestStructWithInheritedLevelFromPtr struct {
	TestStructField *a.TestStruct // want TestStructField:`safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevelFromPtr\.TestStructField: @DoNotLog \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)`
}

type TestStructWithInheritedLevelFromMap struct {
	TestStructField map[string]a.TestStruct // want TestStructField:`safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevelFromMap\.TestStructField: @DoNotLog \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)`
}

type TestStructWithInheritedLevelFromSlice struct {
	TestStructField []a.TestStruct // want TestStructField:`safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevelFromSlice\.TestStructField: @DoNotLog \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)`
}

type TestStructWithInheritedLevelFromEmbedded struct {
	a.TestStruct // want TestStruct:`safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevelFromEmbedded\.TestStruct: @DoNotLog \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)`
}

func ParamTests() {
	testStruct := TestStructWithInheritedLevel{}
	svc1log.SafeParam("testParam", testStruct)                 // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevel", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
	svc1log.SafeParam("testParam", testStruct.TestStructField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	testStruct2 := TestStructWithInheritedLevelFromPtr{}
	svc1log.SafeParam("testParam", testStruct2)                 // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevelFromPtr", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
	svc1log.SafeParam("testParam", testStruct2.TestStructField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	testStruct3 := TestStructWithInheritedLevelFromMap{}
	svc1log.SafeParam("testParam", testStruct3)                 // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevelFromMap", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
	svc1log.SafeParam("testParam", testStruct3.TestStructField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	testStruct4 := TestStructWithInheritedLevelFromSlice{}
	svc1log.SafeParam("testParam", testStruct4)                 // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevelFromSlice", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
	svc1log.SafeParam("testParam", testStruct4.TestStructField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	testStruct5 := TestStructWithInheritedLevelFromEmbedded{}
	svc1log.SafeParam("testParam", testStruct5)            // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/b\.TestStructWithInheritedLevelFromEmbedded", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
	svc1log.SafeParam("testParam", testStruct5.TestStruct) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
}
