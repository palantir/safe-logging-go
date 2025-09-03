package j

import (
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

// safelogging:@DoNotLog
type Password string

func stringPasswordPairFn() (string, Password) {
	return "", ""
}

func stringParamsFn() (string, svc1log.Param, svc1log.Param) {
	return "", nil, nil
}

func mapPairFn() (map[string]interface{}, map[string]interface{}) {
	return nil, nil
}

func ParamTests() {
	svc1log.SafeParam(stringPasswordPairFn()) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/j\.Password", which is @DoNotLog \(reason: comment on "safe-logging-go/safeloggingtests/j\.Password" at .+safelogging_j\.go:8:1\)$`

	svc1log.FromContext(nil).Info(stringParamsFn()) // want `^Info called with unsafe argument: message must be a compile-time constant$`

	werror.SafeAndUnsafeParams(mapPairFn()) // safe: types do not contain any illegal references
}
