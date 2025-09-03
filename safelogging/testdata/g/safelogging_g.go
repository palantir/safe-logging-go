package g

import (
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type TestStruct struct {
}

// PasswordStructReceiverFn is a function.
//
// note when updating tests: need to duplicate "want" output with "t" replacing "PasswordStructReceiverFn" because the
// fact is recorded for both the function and the receiver type.
//
// safelogging:@DoNotLog
func (t TestStruct) PasswordStructReceiverFn() string { // want PasswordStructReceiverFn:`safe-logging-go/safeloggingtests/g\.PasswordStructReceiverFn: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/g\.PasswordStructReceiverFn" at .+/safelogging/testdata/g/safelogging_g\.go:16:1\)` t:`safe-logging-go/safeloggingtests/g\.t: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/g\.t" at .+/safelogging/testdata/g/safelogging_g\.go:16:1\)`
	return ""
}

// safelogging:@DoNotLog
func (t *TestStruct) PasswordStructPointerReceiverFn() string { // want PasswordStructPointerReceiverFn:`safe-logging-go/safeloggingtests/g\.PasswordStructPointerReceiverFn: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/g\.PasswordStructPointerReceiverFn" at .+/safelogging/testdata/g/safelogging_g\.go:21:1\)` t:`safe-logging-go/safeloggingtests/g\.t: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/g\.t" at .+/safelogging/testdata/g/safelogging_g\.go:21:1\)`
	return ""
}

type TestInterface interface {
	// safelogging:@DoNotLog
	PasswordInterfaceFn() string // want PasswordInterfaceFn:`safe-logging-go/safeloggingtests/g\.PasswordInterfaceFn: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/g\.PasswordInterfaceFn" at .+/safelogging/testdata/g/safelogging_g\.go:27:2\)`
}

// safelogging:@DoNotLog
func PasswordStandaloneFn() string { // want PasswordStandaloneFn:`safe-logging-go/safeloggingtests/g\.PasswordStandaloneFn: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/g\.PasswordStandaloneFn" at .+/safelogging/testdata/g/safelogging_g\.go:31:1\)`
	return ""
}

// safelogging:@DoNotLog
func DualReturnFn() (map[string]interface{}, map[string]interface{}) { // want DualReturnFn:`safe-logging-go/safeloggingtests/g\.DualReturnFn: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/g\.DualReturnFn" at .+/safelogging/testdata/g/safelogging_g\.go:36:1\)`
	return nil, nil
}

func ParamTests() {
	svc1log.SafeParam("testParam", PasswordStandaloneFn()) // want `^SafeParam called with unsafe argument: argument references "safe-logging-go/safeloggingtests/g\.PasswordStandaloneFn", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/g\.PasswordStandaloneFn" at .+safelogging_g\.go:31:1\)$`

	var testStructVar TestStruct
	svc1log.SafeParam("testParam", testStructVar.PasswordStructReceiverFn())        // want `^SafeParam called with unsafe argument: argument references "safe-logging-go/safeloggingtests/g\.PasswordStructReceiverFn", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/g\.PasswordStructReceiverFn" at .+safelogging_g\.go:16:1\)$`
	svc1log.SafeParam("testParam", testStructVar.PasswordStructPointerReceiverFn()) // want `^SafeParam called with unsafe argument: argument references "safe-logging-go/safeloggingtests/g\.PasswordStructPointerReceiverFn", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/g\.PasswordStructPointerReceiverFn" at .+safelogging_g\.go:21:1\)$`
	svc1log.SafeParam("testParam", testStructVar)                                   // no warning: ok to log struct that has functions annotated with @DoNotLog

	var multiLevel map[string][]TestStruct
	svc1log.SafeParam("testParam", multiLevel["test"][0].PasswordStructReceiverFn()) // want `^SafeParam called with unsafe argument: argument references "safe-logging-go/safeloggingtests/g\.PasswordStructReceiverFn", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/g\.PasswordStructReceiverFn" at .+safelogging_g\.go:16:1\)$`

	var testInterfaceVar TestInterface
	svc1log.SafeParam("testParam", testInterfaceVar.PasswordInterfaceFn()) // want `^SafeParam called with unsafe argument: argument references "safe-logging-go/safeloggingtests/g\.PasswordInterfaceFn", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/g\.PasswordInterfaceFn" at .+safelogging_g\.go:27:2\)$`

	werror.SafeAndUnsafeParams(map[string]interface{}{"key": PasswordStandaloneFn()}, nil) // want `^SafeAndUnsafeParams called with unsafe argument: argument references "safe-logging-go/safeloggingtests/g\.PasswordStandaloneFn", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/g\.PasswordStandaloneFn" at .+safelogging_g\.go:31:1\)$`

	// note when updating tests: need to duplicate "want" output twice, since it is flagged once for each param.
	// That is, for expected message "ExampleContent", the comment after "want" should be: `ExampleContent` `ExampleContent`
	werror.SafeAndUnsafeParams(DualReturnFn()) // want `^SafeAndUnsafeParams called with unsafe argument: argument references "safe-logging-go/safeloggingtests/g\.DualReturnFn", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/g\.DualReturnFn" at .+safelogging_g\.go:36:1\)$` `^SafeAndUnsafeParams called with unsafe argument: argument references "safe-logging-go/safeloggingtests/g\.DualReturnFn", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/g\.DualReturnFn" at .+safelogging_g\.go:36:1\)$`

	// function call checking is based on types: won't catch if function is stored in a different variable
	fnVar := testStructVar.PasswordStructReceiverFn
	svc1log.SafeParam("testParam", fnVar())

	// only direct function calls are checked: won't catch if value is stored in a variable
	password := PasswordStandaloneFn()
	svc1log.SafeParam("testParam", password)
}
