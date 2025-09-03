package innerh

import (
	"safe-logging-go/safeloggingtests/h"

	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

func ParamTests() {
	svc1log.SafeParam("testParam", h.PasswordVarExported) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/h\.PasswordVarExported" at .+safelogging_h\.go:7:1\)$`

	svc1log.SafeParam("testParam", h.PasswordConstExported) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/h\.PasswordConstExported" at .+safelogging_h\.go:10:1\)$`
}
