package f // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/f\.MyStruct: @DoNotLog\]\]`

import (
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

// safelogging:@DoNotLog
type Password string

type MyStruct struct {
	Username string
	Password Password // want Password:`safe-logging-go/safeloggingtests/f\.MyStruct\.Password: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/f\.Password" at .+/safelogging/testdata/f/safelogging_f\.go:7:1\)`
}

func GetPassword() Password { return "" }

func GetPasswordSlice() []Password { return nil }

func ParamTests() {
	// should flag if unsafe type is passed in via variable
	var passwordVar Password
	svc1log.SafeParam("testParam", passwordVar) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`

	// should flag if unsafe type is passed in via construction
	svc1log.SafeParam("testParam", Password("secret")) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`

	// should flag if unsafe type is passed in via container types
	svc1log.SafeParam("testParam", map[string]Password{}) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`
	svc1log.SafeParam("testParam", []Password{})          // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`
	svc1log.SafeParam("testParam", [1]Password{})         // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`

	// should flag if unsafe type is passed in via container type variable
	var passwordMap map[string]Password
	svc1log.SafeParam("testParam", passwordMap) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`
	var passwordSlice []Password
	svc1log.SafeParam("testParam", passwordSlice) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`
	var passwordArray [1]Password
	svc1log.SafeParam("testParam", passwordArray) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`

	// should flag if unsafe type is passed in via construction in a type
	svc1log.SafeParam("testParam", map[string]any{ // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`
		"testParam": Password("secret"),
	})

	// should flag if unsafe type is passed in via construction in a type
	svc1log.SafeParam("testParam", map[string]Password{ // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`
		"testParam": "secret",
	})

	// should flag if unsafe type is passed in via pointer type
	svc1log.SafeParam("testParam", (*Password)(nil))  // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`
	svc1log.SafeParam("testParam", (**Password)(nil)) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`

	// should flag if unsafe type is passed in via pointer type variable
	var passwordPointerVar *Password
	svc1log.SafeParam("testParam", passwordPointerVar) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`

	// should flag if function returns an unsafe type
	svc1log.SafeParam("testParam", GetPassword()) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`

	// should flag if function returns a container type tha references an unsafe type
	svc1log.SafeParam("testParam", GetPasswordSlice()) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`

	var myStructVar MyStruct
	svc1log.SafeParam("testParam", myStructVar)          // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/f\.MyStruct", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`
	svc1log.SafeParam("testParam", myStructVar.Password) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/f\.Password" at .+safelogging_f\.go:7:1\)$`

	// type checking is based on safety: won't catch if type comes in more permissive type
	var anyVar any
	anyVar = Password("password")
	svc1log.SafeParam("testParam", anyVar)
}
