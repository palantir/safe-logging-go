package i

import (
	"fmt"
	"safe-logging-go/safeloggingtests/i/inneri"

	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type TestStruct struct {
	Name string
}

func ParamTests() {
	// safe: argument is a string constant
	svc1log.FromContext(nil).Info("Safe constant message")
	svc1log.FromContext(nil).Info("Safe constant " + "message")

	svc1log.FromContext(nil).Info(fmt.Sprintf("%s", "Safe constant")) // want `Info called with unsafe argument: message must be a compile-time constant`

	// safe: argument is a reference to a constant
	const safeMessageConst = ""
	svc1log.FromContext(nil).Info(safeMessageConst)
	svc1log.FromContext(nil).Info(safeMessageConst + safeMessageConst)
	svc1log.FromContext(nil).Info(safeMessageConst + "message")

	svc1log.FromContext(nil).Info(fmt.Sprintf("%s", safeMessageConst)) // want `Info called with unsafe argument: message must be a compile-time constant`

	// safe: argument is a reference to a constant from another package
	svc1log.FromContext(nil).Info(inneri.ExportedConst)

	svc1log.FromContext(nil).Info(inneri.ExportedVar) // want `Info called with unsafe argument: message must be a compile-time constant`

	var unsafeMessageVar string
	svc1log.FromContext(nil).Info(unsafeMessageVar) // want `Info called with unsafe argument: message must be a compile-time constant`

	svc1log.FromContext(nil).Info(fmt.Sprintf("%s", unsafeMessageVar)) // want `Info called with unsafe argument: message must be a compile-time constant`

	svc1log.FromContext(nil).Info(stringFn()) // want `Info called with unsafe argument: message must be a compile-time constant`

	var testStructVar TestStruct
	svc1log.FromContext(nil).Info(testStructVar.Name) // want `Info called with unsafe argument: message must be a compile-time constant`

	svc1log.FromContext(nil).Info(TestStruct{ // want `Info called with unsafe argument: message must be a compile-time constant`
		Name: stringFn(),
	}.Name)
}

func stringFn() string {
	return ""
}
