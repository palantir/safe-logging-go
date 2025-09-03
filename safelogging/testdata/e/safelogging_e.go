package e // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/e\.TestStructOuter: @Unsafe, safe-logging-go/safeloggingtests/e\.TestStructTaggedWithComment: @Unsafe\]\]`

import (
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

// safelogging:@Unsafe
type TestStructTaggedWithComment struct {
	Name string
}

// safelogging:@Unsafe
type TestInterfaceWithComment interface {
}

type TestNamedTypeTestStructTaggedWithComment TestStructTaggedWithComment

type TestStructOuter struct { // safelogging:@Unsafe
	inner TestStructInner
}

type TestStructInner struct {
	Name string
}

type TestNamedTypeTestStructOuter TestStructTaggedWithComment

type TestNamedTypeTestStructOuterPointer *TestStructTaggedWithComment

func ParamTests() {
	var testStructTaggedWithCommentVar TestStructTaggedWithComment
	svc1log.SafeParam("testParam", testStructTaggedWithCommentVar) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/e\.TestStructTaggedWithComment", which is @Unsafe \(reason: comment on "safe-logging-go/safeloggingtests/e\.TestStructTaggedWithComment" at .+safelogging_e\.go:7:1\)$`

	var testInterfaceWithCommentVar TestInterfaceWithComment
	svc1log.SafeParam("testParam", testInterfaceWithCommentVar) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/e\.TestInterfaceWithComment", which is @Unsafe \(reason: comment on "safe-logging-go/safeloggingtests/e\.TestInterfaceWithComment" at .+safelogging_e\.go:12:1\)$`

	var testNamedTypeTestStructTaggedWithCommentVar TestNamedTypeTestStructTaggedWithComment
	svc1log.SafeParam("testParam", testNamedTypeTestStructTaggedWithCommentVar) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/e\.TestNamedTypeTestStructTaggedWithComment", which is @Unsafe \(reason: comment on "safe-logging-go/safeloggingtests/e\.TestStructTaggedWithComment" at .+safelogging_e\.go:7:1\)$`

	var testNamedTypeTestStructOuterVar TestNamedTypeTestStructOuter
	svc1log.SafeParam("testParam", testNamedTypeTestStructOuterVar) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/e\.TestNamedTypeTestStructOuter", which is @Unsafe \(reason: comment on "safe-logging-go/safeloggingtests/e\.TestStructTaggedWithComment" at .+safelogging_e\.go:7:1\)$`

	var testNamedTypeTestStructOuterPointerVar TestNamedTypeTestStructOuterPointer
	svc1log.SafeParam("testParam", testNamedTypeTestStructOuterPointerVar) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/e\.TestNamedTypeTestStructOuterPointer", which is @Unsafe \(reason: comment on "safe-logging-go/safeloggingtests/e\.TestStructTaggedWithComment" at .+safelogging_e\.go:7:1\)$`
}
