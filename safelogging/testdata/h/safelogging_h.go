package h

import (
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

// safelogging:@DoNotLog
var PasswordVarExported string // want PasswordVarExported:`safe-logging-go/safeloggingtests/h\.PasswordVarExported: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/h\.PasswordVarExported" at .+/safelogging/testdata/h/safelogging_h\.go:7:1\)`

// safelogging:@DoNotLog
const PasswordConstExported = "password" // want PasswordConstExported:`safe-logging-go/safeloggingtests/h\.PasswordConstExported: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/h\.PasswordConstExported" at .+/safelogging/testdata/h/safelogging_h\.go:10:1\)`

// safelogging:@DoNotLog
var passwordVarPrivate string // want passwordVarPrivate:`safe-logging-go/safeloggingtests/h\.passwordVarPrivate: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/h\.passwordVarPrivate" at .+/safelogging/testdata/h/safelogging_h\.go:13:1\)`

// safelogging:@DoNotLog
const passwordConstPrivate = "password" // want passwordConstPrivate:`safe-logging-go/safeloggingtests/h\.passwordConstPrivate: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/h\.passwordConstPrivate" at .+/safelogging/testdata/h/safelogging_h\.go:16:1\)`

func ParamTests() {
	svc1log.SafeParam("testParam", PasswordVarExported) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/h\.PasswordVarExported" at .+safelogging_h\.go:7:1\)$`
	svc1log.SafeParam("testParam", passwordVarPrivate)  // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/h\.passwordVarPrivate" at .+safelogging_h\.go:13:1\)$`

	svc1log.SafeParam("testParam", PasswordConstExported) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/h\.PasswordConstExported" at .+safelogging_h\.go:10:1\)$`
	svc1log.SafeParam("testParam", passwordConstPrivate)  // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/h\.passwordConstPrivate" at .+safelogging_h\.go:16:1\)$`

	// safelogging:@DoNotLog
	localVar := "Message" // want localVar:`safe-logging-go/safeloggingtests/h\.localVar: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/h\.localVar" at .+/safelogging/testdata/h/safelogging_h\.go:26:2\)`
	localVar = "Updated message"
	svc1log.SafeParam("testParam", localVar) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/h\.localVar" at .+safelogging_h\.go:26:2\)$`

	// safelogging:@DoNotLog
	var otherLocalVar string                      // want otherLocalVar:`safe-logging-go/safeloggingtests/h\.otherLocalVar: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/h\.otherLocalVar" at .+/safelogging/testdata/h/safelogging_h\.go:31:2\)`
	svc1log.SafeParam("testParam", otherLocalVar) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/h\.otherLocalVar" at .+safelogging_h\.go:31:2\)$`

	// no warning: safety is not tracked through assignments
	okVar := PasswordVarExported
	svc1log.SafeParam("testParam", okVar)
}
