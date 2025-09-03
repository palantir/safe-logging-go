package d // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/d\.StructWithNamedTypeField: @DoNotLog\]\]`

import (
	"safe-logging-go/safeloggingtests/a"

	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type TestStructNamedTypeImported a.TestStruct

type TestStructAliasImported = a.TestStruct

type TestStructNamedTypeImportedPointer *a.TestStruct

type TestStructAliasImportedPointer = *a.TestStruct

type StructWithNamedTypeField struct {
	NamedTypeField TestStructNamedTypeImported // want NamedTypeField:`safe-logging-go/safeloggingtests/d\.StructWithNamedTypeField\.NamedTypeField: @DoNotLog \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)`
}

type TestStructNamedTypeLocal = StructWithNamedTypeField

type TestStructAliasLocal = StructWithNamedTypeField

type TestStructNamedTypeLocalPointer *StructWithNamedTypeField

type TestStructAliasLocalPointer = *StructWithNamedTypeField

type LocalTestStruct struct{}

type AnnotatedNamedType LocalTestStruct // safelogging:@Unsafe

type AnnotatedTypeAlias = LocalTestStruct // safelogging:@Unsafe

func ParamTests() {
	var testStructNamedTypeImportedVar TestStructNamedTypeImported
	// no warning: ok to log safe field using SafeParam
	svc1log.SafeParam("testParam", testStructNamedTypeImportedVar.SafeField)
	svc1log.SafeParam("testParam", testStructNamedTypeImportedVar.UnsafeField) // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
	svc1log.SafeParam("testParam", testStructNamedTypeImportedVar)             // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/d\.TestStructNamedTypeImported", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	var testStructAliasImportedVar TestStructAliasImported
	// no warning: ok to log safe field using SafeParam
	svc1log.SafeParam("testParam", testStructAliasImportedVar.SafeField)
	svc1log.SafeParam("testParam", testStructAliasImportedVar.UnsafeField) // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
	svc1log.SafeParam("testParam", testStructAliasImportedVar)             // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStruct", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	var testStructNamedTypeImportedPointerVar TestStructNamedTypeImportedPointer
	testStructNamedTypeImportedPointerVar = &a.TestStruct{}
	// no warning: ok to log safe field using SafeParam
	svc1log.SafeParam("testParam", testStructNamedTypeImportedPointerVar.SafeField)
	svc1log.SafeParam("testParam", testStructNamedTypeImportedPointerVar.UnsafeField) // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
	svc1log.SafeParam("testParam", testStructNamedTypeImportedPointerVar)             // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/d\.TestStructNamedTypeImportedPointer", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	var testStructAliasImportedPointerVar TestStructAliasImportedPointer
	testStructAliasImportedPointerVar = &a.TestStruct{}
	// no warning: ok to log safe field using SafeParam
	svc1log.SafeParam("testParam", testStructAliasImportedPointerVar.SafeField)
	svc1log.SafeParam("testParam", testStructAliasImportedPointerVar.UnsafeField) // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
	svc1log.SafeParam("testParam", testStructAliasImportedPointerVar)             // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStruct", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	var testStructWithNamedTypeFieldVar StructWithNamedTypeField
	svc1log.SafeParam("testParam", testStructWithNamedTypeFieldVar.NamedTypeField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
	svc1log.SafeParam("testParam", testStructWithNamedTypeFieldVar)                // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/d\.StructWithNamedTypeField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	var testStructNamedTypeLocalVar TestStructNamedTypeLocal
	svc1log.SafeParam("testParam", testStructNamedTypeLocalVar.NamedTypeField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
	svc1log.SafeParam("testParam", testStructNamedTypeLocalVar)                // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/d\.StructWithNamedTypeField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	var testStructAliasLocalVar TestStructAliasLocal
	svc1log.SafeParam("testParam", testStructAliasLocalVar.NamedTypeField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
	svc1log.SafeParam("testParam", testStructAliasLocalVar)                // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/d\.StructWithNamedTypeField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	var testStructNamedTypeLocalPointerVar TestStructNamedTypeLocalPointer
	testStructNamedTypeLocalPointerVar = &StructWithNamedTypeField{}
	svc1log.SafeParam("testParam", testStructNamedTypeLocalPointerVar.NamedTypeField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
	svc1log.SafeParam("testParam", testStructNamedTypeLocalPointerVar)                // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/d\.TestStructNamedTypeLocalPointer", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	var testAnnotatedNamedType AnnotatedNamedType
	svc1log.SafeParam("testParam", testAnnotatedNamedType) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/d\.AnnotatedNamedType", which is @Unsafe \(reason: comment on "safe-logging-go/safeloggingtests/d\.AnnotatedNamedType" at .+safelogging_d\.go:31:41\)$`

	var testAnnotatedTypeAlias AnnotatedTypeAlias
	svc1log.SafeParam("testParam", testAnnotatedTypeAlias) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/d\.AnnotatedTypeAlias", which is @Unsafe \(reason: comment on "safe-logging-go/safeloggingtests/d\.AnnotatedTypeAlias" at .+safelogging_d\.go:33:43\)$`

	var testStructAliasLocalPointerVar TestStructAliasLocalPointer
	testStructAliasLocalPointerVar = &StructWithNamedTypeField{}
	// no warning: ok to log safe field using SafeParam
	svc1log.SafeParam("testParam", testStructAliasLocalPointerVar.NamedTypeField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
	svc1log.SafeParam("testParam", testStructAliasLocalPointerVar)                // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/d\.StructWithNamedTypeField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
}
