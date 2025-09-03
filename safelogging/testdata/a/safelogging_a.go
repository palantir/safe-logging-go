package a // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/a\.TestStruct: @DoNotLog, safe-logging-go/safeloggingtests/a\.TestStructEmbedded: @DoNotLog, safe-logging-go/safeloggingtests/a\.TestStructEmbeddedPointer: @DoNotLog, safe-logging-go/safeloggingtests/a\.TestStructInlineDef: @Unsafe, safe-logging-go/safeloggingtests/a\.TestStructInlineNestedDef: @Unsafe, safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField: @DoNotLog, safe-logging-go/safeloggingtests/a\.TestStructWithInheritedLevel: @DoNotLog, safe-logging-go/safeloggingtests/a\.TestStructWithSafeField: @Safe, safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField: @Unsafe\]\]`

import (
	"safe-logging-go/safeloggingtests/a/innera"

	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type TestStruct struct {
	SafeField     *string `safelogging:"@Safe"`     // want SafeField:`safe-logging-go/safeloggingtests/a\.TestStruct\.SafeField: @Safe \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.SafeField"\)`
	UnsafeField   *string `safelogging:"@Unsafe"`   // want UnsafeField:`safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField: @Unsafe \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)`
	DoNotLogField *string `safelogging:"@DoNotLog"` // want DoNotLogField:`safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField: @DoNotLog \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)`
}

type TestStructWithSafeField struct {
	SafeField *string `safelogging:"@Safe"` // want SafeField:`safe-logging-go/safeloggingtests/a\.TestStructWithSafeField\.SafeField: @Safe \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithSafeField\.SafeField"\)`
}

type TestStructWithUnsafeField struct {
	UnsafeField *string `safelogging:"@Unsafe"` // want UnsafeField:`safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField: @Unsafe \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)`
}

type TestStructWithDoNotLogField struct {
	DoNotLogField *string `safelogging:"@DoNotLog"` // want DoNotLogField:`safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField: @DoNotLog \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)`
}

type TestStructWithInheritedLevel struct {
	TestStructField TestStruct // want TestStructField:`safe-logging-go/safeloggingtests/a\.TestStructWithInheritedLevel\.TestStructField: @DoNotLog \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)`
}

type TestStructEmbedded struct {
	TestStruct // want TestStruct:`safe-logging-go/safeloggingtests/a\.TestStructEmbedded\.TestStruct: @DoNotLog \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)`
}

type TestStructEmbeddedPointer struct {
	*TestStruct // want TestStruct:`safe-logging-go/safeloggingtests/a\.TestStructEmbeddedPointer\.TestStruct: @DoNotLog \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)`
}

type TestStructInlineDef struct {
	Auth struct { // want Auth:`safe-logging-go/safeloggingtests/a\.TestStructInlineDef\.Auth: @Unsafe \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructInlineDef\.Auth\.Password"\)`
		Password string `safelogging:"@Unsafe"` // want Password:`safe-logging-go/safeloggingtests/a\.TestStructInlineDef\.Auth\.Password: @Unsafe \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructInlineDef\.Auth\.Password"\)`
	}
}

type TestStructInlineNestedDef struct {
	Inner struct { // want Inner:`safe-logging-go/safeloggingtests/a\.TestStructInlineNestedDef\.Inner: @Unsafe \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructInlineNestedDef\.Inner\.Auth\.Password"\)`
		Auth struct { // want Auth:`safe-logging-go/safeloggingtests/a\.TestStructInlineNestedDef\.Inner\.Auth: @Unsafe \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructInlineNestedDef\.Inner\.Auth\.Password"\)`
			Password string `safelogging:"@Unsafe"` // want Password:`safe-logging-go/safeloggingtests/a\.TestStructInlineNestedDef\.Inner\.Auth\.Password: @Unsafe \(tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructInlineNestedDef\.Inner\.Auth\.Password"\)`
		}
	}
}

func VerifyStructFieldReference() {
	testStruct := TestStruct{}

	svc1log.SafeParam("testParam", testStruct.SafeField)     // no warning: ok to log safe field using svc1log.SafeParam
	svc1log.SafeParam("testParam", testStruct.UnsafeField)   // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
	svc1log.SafeParam("testParam", testStruct.DoNotLogField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	svc1log.SafeParams(map[string]interface{}{ // no warning: ok to log safe field using svc1log.SafeParams
		"testParam": testStruct.SafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
		"testParam": testStruct.UnsafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
		"testParam": testStruct.DoNotLogField,
	})

	werror.SafeParam("testParam", testStruct.SafeField)     // no warning: ok to log safe field using werror.SafeParam
	werror.SafeParam("testParam", testStruct.UnsafeField)   // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
	werror.SafeParam("testParam", testStruct.DoNotLogField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	werror.SafeParams(map[string]interface{}{ // no warning: ok to log safe field using werror.SafeParams
		"testParam": testStruct.SafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
		"testParam": testStruct.UnsafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
		"testParam": testStruct.DoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // no warning: ok to log safe field using safe argument of werror.SafeAndUnsafeParams
			"testParam": testStruct.SafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
			"testParam": testStruct.UnsafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
			"testParam": testStruct.DoNotLogField,
		},
		nil,
	)

	svc1log.UnsafeParam("testParam", testStruct.SafeField)     // no warning: ok to log safe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", testStruct.UnsafeField)   // no warning: ok to log unsafe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", testStruct.DoNotLogField) // want `^UnsafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log safe field using svc1log.UnsafeParams
		"testParam": testStruct.SafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log unsafe field using svc1log.UnsafeParams
		"testParam": testStruct.UnsafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // want `^UnsafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
		"testParam": testStruct.DoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log safe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": testStruct.SafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log unsafe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": testStruct.UnsafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
			"testParam": testStruct.DoNotLogField,
		},
	)
}

func VerifyStructFieldAddressReference() {
	testStruct := TestStruct{}

	svc1log.SafeParam("testParam", &testStruct.SafeField)     // no warning: ok to log safe field using svc1log.SafeParam
	svc1log.SafeParam("testParam", &testStruct.UnsafeField)   // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
	svc1log.SafeParam("testParam", &testStruct.DoNotLogField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	svc1log.SafeParams(map[string]interface{}{ // no warning: ok to log safe field using svc1log.SafeParams
		"testParam": &testStruct.SafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
		"testParam": &testStruct.UnsafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
		"testParam": &testStruct.DoNotLogField,
	})

	werror.SafeParam("testParam", &testStruct.SafeField)     // no warning: ok to log safe field using werror.SafeParam
	werror.SafeParam("testParam", &testStruct.UnsafeField)   // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
	werror.SafeParam("testParam", &testStruct.DoNotLogField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	werror.SafeParams(map[string]interface{}{ // no warning: ok to log safe field using werror.SafeParams
		"testParam": &testStruct.SafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
		"testParam": &testStruct.UnsafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
		"testParam": &testStruct.DoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // no warning: ok to log safe field using safe argument of werror.SafeAndUnsafeParams
			"testParam": &testStruct.SafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
			"testParam": &testStruct.UnsafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
			"testParam": &testStruct.DoNotLogField,
		},
		nil,
	)

	svc1log.UnsafeParam("testParam", &testStruct.SafeField)     // no warning: ok to log safe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", &testStruct.UnsafeField)   // no warning: ok to log unsafe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", &testStruct.DoNotLogField) // want `^UnsafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log safe field using svc1log.UnsafeParams
		"testParam": &testStruct.SafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log unsafe field using svc1log.UnsafeParams
		"testParam": &testStruct.UnsafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // want `^UnsafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
		"testParam": &testStruct.DoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log safe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": &testStruct.SafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log unsafe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": &testStruct.UnsafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
			"testParam": &testStruct.DoNotLogField,
		},
	)
}

func VerifyStructFieldAddressDereference() {
	testStruct := TestStruct{}

	svc1log.SafeParam("testParam", *testStruct.SafeField)     // no warning: ok to log safe field using svc1log.SafeParam
	svc1log.SafeParam("testParam", *testStruct.UnsafeField)   // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
	svc1log.SafeParam("testParam", *testStruct.DoNotLogField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	svc1log.SafeParams(map[string]interface{}{ // no warning: ok to log safe field using svc1log.SafeParams
		"testParam": *testStruct.SafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
		"testParam": *testStruct.UnsafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
		"testParam": *testStruct.DoNotLogField,
	})

	werror.SafeParam("testParam", *testStruct.SafeField)     // no warning: ok to log safe field using werror.SafeParam
	werror.SafeParam("testParam", *testStruct.UnsafeField)   // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
	werror.SafeParam("testParam", *testStruct.DoNotLogField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	werror.SafeParams(map[string]interface{}{ // no warning: ok to log safe field using werror.SafeParams
		"testParam": *testStruct.SafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
		"testParam": *testStruct.UnsafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
		"testParam": *testStruct.DoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // no warning: ok to log safe field using safe argument of werror.SafeAndUnsafeParams
			"testParam": *testStruct.SafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.UnsafeField"\)$`
			"testParam": *testStruct.UnsafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
			"testParam": *testStruct.DoNotLogField,
		},
		nil,
	)

	svc1log.UnsafeParam("testParam", *testStruct.SafeField)     // no warning: ok to log safe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", *testStruct.UnsafeField)   // no warning: ok to log unsafe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", *testStruct.DoNotLogField) // want `^UnsafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log safe field using svc1log.UnsafeParams
		"testParam": *testStruct.SafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log unsafe field using svc1log.UnsafeParams
		"testParam": *testStruct.UnsafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // want `^UnsafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
		"testParam": *testStruct.DoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log safe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": *testStruct.SafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log unsafe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": *testStruct.UnsafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`
			"testParam": *testStruct.DoNotLogField,
		},
	)
}

func VerifyStructReference() {
	testStructWithSafeField := TestStructWithSafeField{}
	testStructWithUnsafeField := TestStructWithUnsafeField{}
	testStructWithDoNotLogField := TestStructWithDoNotLogField{}

	svc1log.SafeParam("testParam", testStructWithSafeField)     // no warning: ok to log safe struct using svc1log.SafeParam
	svc1log.SafeParam("testParam", testStructWithUnsafeField)   // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
	svc1log.SafeParam("testParam", testStructWithDoNotLogField) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`

	svc1log.SafeParams(map[string]interface{}{ // no warning: ok to log safe struct using svc1log.SafeParams
		"testParam": testStructWithSafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
		"testParam": testStructWithUnsafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": testStructWithDoNotLogField,
	})

	werror.SafeParam("testParam", testStructWithSafeField)     // no warning: ok to log safe struct using werror.SafeParam
	werror.SafeParam("testParam", testStructWithUnsafeField)   // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
	werror.SafeParam("testParam", testStructWithDoNotLogField) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`

	werror.SafeParams(map[string]interface{}{ // no warning: ok to log safe field using werror.SafeParams
		"testParam": testStructWithSafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
		"testParam": testStructWithUnsafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": testStructWithDoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // no warning: ok to log safe field using safe argument of werror.SafeAndUnsafeParams
			"testParam": testStructWithSafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
			"testParam": testStructWithUnsafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
			"testParam": testStructWithDoNotLogField,
		},
		nil,
	)

	svc1log.UnsafeParam("testParam", testStructWithSafeField)     // no warning: ok to log safe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", testStructWithUnsafeField)   // no warning: ok to log unsafe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", testStructWithDoNotLogField) // want `^UnsafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`

	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log safe field using svc1log.UnsafeParams
		"testParam": testStructWithSafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log unsafe field using svc1log.UnsafeParams
		"testParam": testStructWithUnsafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // want `^UnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": testStructWithDoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log safe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": testStructWithSafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log unsafe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": testStructWithUnsafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
			"testParam": testStructWithDoNotLogField,
		},
	)
}

func VerifyImportedStructReference() {
	testImportedStructWithSafeField := innera.InnerTestStructWithSafeField{}
	testImportedStructWithUnsafeField := innera.InnerTestStructWithUnsafeField{}
	testImportedStructWithDoNotLogField := innera.InnerTestStructWithDoNotLogField{}

	svc1log.SafeParam("testParam", testImportedStructWithSafeField)     // no warning: ok to log safe struct using svc1log.SafeParam
	svc1log.SafeParam("testParam", testImportedStructWithUnsafeField)   // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField\.UnsafeField"\)$`
	svc1log.SafeParam("testParam", testImportedStructWithDoNotLogField) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`

	svc1log.SafeParams(map[string]interface{}{ // no warning: ok to log safe struct using svc1log.SafeParams
		"testParam": testImportedStructWithSafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField\.UnsafeField"\)$`
		"testParam": testImportedStructWithUnsafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": testImportedStructWithDoNotLogField,
	})

	werror.SafeParam("testParam", testImportedStructWithSafeField)     // no warning: ok to log safe struct using werror.SafeParam
	werror.SafeParam("testParam", testImportedStructWithUnsafeField)   // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField\.UnsafeField"\)$`
	werror.SafeParam("testParam", testImportedStructWithDoNotLogField) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`

	werror.SafeParams(map[string]interface{}{ // no warning: ok to log safe field using werror.SafeParams
		"testParam": testImportedStructWithSafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField\.UnsafeField"\)$`
		"testParam": testImportedStructWithUnsafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": testImportedStructWithDoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // no warning: ok to log safe field using safe argument of werror.SafeAndUnsafeParams
			"testParam": testImportedStructWithSafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField\.UnsafeField"\)$`
			"testParam": testImportedStructWithUnsafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`
			"testParam": testImportedStructWithDoNotLogField,
		},
		nil,
	)

	svc1log.UnsafeParam("testParam", testImportedStructWithSafeField)     // no warning: ok to log safe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", testImportedStructWithUnsafeField)   // no warning: ok to log unsafe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", testImportedStructWithDoNotLogField) // want `^UnsafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`

	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log safe field using svc1log.UnsafeParams
		"testParam": testImportedStructWithSafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log unsafe field using svc1log.UnsafeParams
		"testParam": testImportedStructWithUnsafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // want `^UnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": testImportedStructWithDoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log safe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": testImportedStructWithSafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log unsafe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": testImportedStructWithUnsafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`
			"testParam": testImportedStructWithDoNotLogField,
		},
	)
}

func VerifyStructLiteral() {
	svc1log.SafeParam("testParam", TestStructWithSafeField{})     // no warning: ok to log safe struct using svc1log.SafeParam
	svc1log.SafeParam("testParam", TestStructWithUnsafeField{})   // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
	svc1log.SafeParam("testParam", TestStructWithDoNotLogField{}) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`

	svc1log.SafeParams(map[string]interface{}{ // no warning: ok to log safe struct using svc1log.SafeParams
		"testParam": TestStructWithSafeField{},
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
		"testParam": TestStructWithUnsafeField{},
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": TestStructWithDoNotLogField{},
	})

	werror.SafeParam("testParam", TestStructWithSafeField{})     // no warning: ok to log safe struct using werror.SafeParam
	werror.SafeParam("testParam", TestStructWithUnsafeField{})   // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
	werror.SafeParam("testParam", TestStructWithDoNotLogField{}) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`

	werror.SafeParams(map[string]interface{}{ // no warning: ok to log safe field using werror.SafeParams
		"testParam": TestStructWithSafeField{},
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
		"testParam": TestStructWithUnsafeField{},
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": TestStructWithDoNotLogField{},
	})

	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // no warning: ok to log safe field using safe argument of werror.SafeAndUnsafeParams
			"testParam": TestStructWithSafeField{},
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
			"testParam": TestStructWithUnsafeField{},
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
			"testParam": TestStructWithDoNotLogField{},
		},
		nil,
	)

	svc1log.UnsafeParam("testParam", TestStructWithSafeField{})     // no warning: ok to log safe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", TestStructWithUnsafeField{})   // no warning: ok to log unsafe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", TestStructWithDoNotLogField{}) // want `^UnsafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`

	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log safe field using svc1log.UnsafeParams
		"testParam": TestStructWithSafeField{},
	})
	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log unsafe field using svc1log.UnsafeParams
		"testParam": TestStructWithUnsafeField{},
	})
	svc1log.UnsafeParams(map[string]interface{}{ // want `^UnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": TestStructWithDoNotLogField{},
	})

	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log safe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": TestStructWithSafeField{},
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log unsafe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": TestStructWithUnsafeField{},
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
			"testParam": TestStructWithDoNotLogField{},
		},
	)
}

func VerifyImportedStructLiteral() {
	svc1log.SafeParam("testParam", innera.InnerTestStructWithSafeField{})     // no warning: ok to log safe struct using svc1log.SafeParam
	svc1log.SafeParam("testParam", innera.InnerTestStructWithUnsafeField{})   // want `^SafeParam called with unsafe argument: argument is of type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField\.UnsafeField"\)$`
	svc1log.SafeParam("testParam", innera.InnerTestStructWithDoNotLogField{}) // want `^SafeParam called with unsafe argument: argument is of type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`

	svc1log.SafeParams(map[string]interface{}{ // no warning: ok to log safe struct using svc1log.SafeParams
		"testParam": innera.InnerTestStructWithSafeField{},
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is of type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField\.UnsafeField"\)$`
		"testParam": innera.InnerTestStructWithUnsafeField{},
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is of type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": innera.InnerTestStructWithDoNotLogField{},
	})

	werror.SafeParam("testParam", innera.InnerTestStructWithSafeField{})     // no warning: ok to log safe struct using werror.SafeParam
	werror.SafeParam("testParam", innera.InnerTestStructWithUnsafeField{})   // want `^SafeParam called with unsafe argument: argument is of type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField\.UnsafeField"\)$`
	werror.SafeParam("testParam", innera.InnerTestStructWithDoNotLogField{}) // want `^SafeParam called with unsafe argument: argument is of type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`

	werror.SafeParams(map[string]interface{}{ // no warning: ok to log safe field using werror.SafeParams
		"testParam": innera.InnerTestStructWithSafeField{},
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is of type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField\.UnsafeField"\)$`
		"testParam": innera.InnerTestStructWithUnsafeField{},
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument is of type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": innera.InnerTestStructWithDoNotLogField{},
	})

	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // no warning: ok to log safe field using safe argument of werror.SafeAndUnsafeParams
			"testParam": innera.InnerTestStructWithSafeField{},
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument is of type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithUnsafeField\.UnsafeField"\)$`
			"testParam": innera.InnerTestStructWithUnsafeField{},
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument is of type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`
			"testParam": innera.InnerTestStructWithDoNotLogField{},
		},
		nil,
	)

	svc1log.UnsafeParam("testParam", innera.InnerTestStructWithSafeField{})     // no warning: ok to log safe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", innera.InnerTestStructWithUnsafeField{})   // no warning: ok to log unsafe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", innera.InnerTestStructWithDoNotLogField{}) // want `^UnsafeParam called with unsafe argument: argument is of type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`

	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log safe field using svc1log.UnsafeParams
		"testParam": innera.InnerTestStructWithSafeField{},
	})
	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log unsafe field using svc1log.UnsafeParams
		"testParam": innera.InnerTestStructWithUnsafeField{},
	})
	svc1log.UnsafeParams(map[string]interface{}{ // want `^UnsafeParams called with unsafe argument: argument is of type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": innera.InnerTestStructWithDoNotLogField{},
	})

	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log safe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": innera.InnerTestStructWithSafeField{},
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log unsafe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": innera.InnerTestStructWithUnsafeField{},
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument is of type "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a/innera\.InnerTestStructWithDoNotLogField\.DoNotLogField"\)$`
			"testParam": innera.InnerTestStructWithDoNotLogField{},
		},
	)
}

func VerifyStructAddressReference() {
	testStructWithSafeField := TestStructWithSafeField{}
	testStructWithUnsafeField := TestStructWithUnsafeField{}
	testStructWithDoNotLogField := TestStructWithDoNotLogField{}

	svc1log.SafeParam("testParam", &testStructWithSafeField)     // no warning: ok to log safe struct using svc1log.SafeParam
	svc1log.SafeParam("testParam", &testStructWithUnsafeField)   // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
	svc1log.SafeParam("testParam", &testStructWithDoNotLogField) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`

	svc1log.SafeParams(map[string]interface{}{ // no warning: ok to log safe struct using svc1log.SafeParams
		"testParam": &testStructWithSafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
		"testParam": &testStructWithUnsafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": &testStructWithDoNotLogField,
	})

	werror.SafeParam("testParam", &testStructWithSafeField)     // no warning: ok to log safe struct using werror.SafeParam
	werror.SafeParam("testParam", &testStructWithUnsafeField)   // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
	werror.SafeParam("testParam", &testStructWithDoNotLogField) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`

	werror.SafeParams(map[string]interface{}{ // no warning: ok to log safe field using werror.SafeParams
		"testParam": &testStructWithSafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
		"testParam": &testStructWithUnsafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": &testStructWithDoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // no warning: ok to log safe field using safe argument of werror.SafeAndUnsafeParams
			"testParam": &testStructWithSafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
			"testParam": &testStructWithUnsafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
			"testParam": &testStructWithDoNotLogField,
		},
		nil,
	)

	svc1log.UnsafeParam("testParam", &testStructWithSafeField)     // no warning: ok to log safe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", &testStructWithUnsafeField)   // no warning: ok to log unsafe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", &testStructWithDoNotLogField) // want `^UnsafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`

	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log safe field using svc1log.UnsafeParams
		"testParam": &testStructWithSafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log unsafe field using svc1log.UnsafeParams
		"testParam": &testStructWithUnsafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // want `^UnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": &testStructWithDoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log safe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": &testStructWithSafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log unsafe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": &testStructWithUnsafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
			"testParam": &testStructWithDoNotLogField,
		},
	)
}

func VerifyStructAddressDereference() {
	testStructWithSafeField := &TestStructWithSafeField{}
	testStructWithUnsafeField := &TestStructWithUnsafeField{}
	testStructWithDoNotLogField := &TestStructWithDoNotLogField{}

	svc1log.SafeParam("testParam", *testStructWithSafeField)     // no warning: ok to log safe struct using svc1log.SafeParam
	svc1log.SafeParam("testParam", *testStructWithUnsafeField)   // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
	svc1log.SafeParam("testParam", *testStructWithDoNotLogField) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`

	svc1log.SafeParams(map[string]interface{}{ // no warning: ok to log safe struct using svc1log.SafeParams
		"testParam": *testStructWithSafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
		"testParam": *testStructWithUnsafeField,
	})
	svc1log.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": *testStructWithDoNotLogField,
	})

	werror.SafeParam("testParam", *testStructWithSafeField)     // no warning: ok to log safe struct using werror.SafeParam
	werror.SafeParam("testParam", *testStructWithUnsafeField)   // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
	werror.SafeParam("testParam", *testStructWithDoNotLogField) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`

	werror.SafeParams(map[string]interface{}{ // no warning: ok to log safe field using werror.SafeParams
		"testParam": *testStructWithSafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
		"testParam": *testStructWithUnsafeField,
	})
	werror.SafeParams(map[string]interface{}{ // want `^SafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": *testStructWithDoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // no warning: ok to log safe field using safe argument of werror.SafeAndUnsafeParams
			"testParam": *testStructWithSafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithUnsafeField\.UnsafeField"\)$`
			"testParam": *testStructWithUnsafeField,
		},
		nil,
	)
	werror.SafeAndUnsafeParams(
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
			"testParam": *testStructWithDoNotLogField,
		},
		nil,
	)

	svc1log.UnsafeParam("testParam", *testStructWithSafeField)     // no warning: ok to log safe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", *testStructWithUnsafeField)   // no warning: ok to log unsafe field using svc1log.UnsafeParam
	svc1log.UnsafeParam("testParam", *testStructWithDoNotLogField) // want `^UnsafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`

	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log safe field using svc1log.UnsafeParams
		"testParam": *testStructWithSafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // no warning: ok to log unsafe field using svc1log.UnsafeParams
		"testParam": *testStructWithUnsafeField,
	})
	svc1log.UnsafeParams(map[string]interface{}{ // want `^UnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
		"testParam": *testStructWithDoNotLogField,
	})

	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log safe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": *testStructWithSafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // no warning: ok to log unsafe field using unsafe argument of werror.SafeAndUnsafeParams
			"testParam": *testStructWithUnsafeField,
		},
	)
	werror.SafeAndUnsafeParams(
		nil,
		map[string]interface{}{ // want `^SafeAndUnsafeParams called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructWithDoNotLogField\.DoNotLogField"\)$`
			"testParam": *testStructWithDoNotLogField,
		},
	)
}

func ParamTests() {
	testStruct := TestStruct{}

	// slice that contains unsafe struct should be flagged
	svc1log.SafeParam("testParam", []TestStruct{testStruct}) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStruct", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	// array that contains unsafe struct should be flagged
	svc1log.SafeParam("testParam", [1]TestStruct{testStruct}) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStruct", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	// map that contains unsafe struct should be flagged
	svc1log.SafeParam("testParam", map[string]TestStruct{"test": testStruct}) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStruct", which is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	testStructWithInheritedLevel := TestStructWithInheritedLevel{}
	svc1log.SafeParam("testParam", testStructWithInheritedLevel.TestStructField) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStruct\.DoNotLogField"\)$`

	var testStructInlineDefVar TestStructInlineDef
	svc1log.SafeParam("testParam", testStructInlineDefVar.Auth.Password) // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructInlineDef\.Auth\.Password"\)$`
	svc1log.SafeParam("testParam", testStructInlineDefVar.Auth)          // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructInlineDef\.Auth\.Password"\)$`
	svc1log.SafeParam("testParam", testStructInlineDefVar)               // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructInlineDef", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructInlineDef\.Auth\.Password"\)$`

	var testStructInlineNestedDefVar TestStructInlineNestedDef
	svc1log.SafeParam("testParam", testStructInlineNestedDefVar.Inner.Auth.Password) // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructInlineNestedDef\.Inner\.Auth\.Password"\)$`
	svc1log.SafeParam("testParam", testStructInlineNestedDefVar.Inner.Auth)          // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructInlineNestedDef\.Inner\.Auth\.Password"\)$`
	svc1log.SafeParam("testParam", testStructInlineNestedDefVar.Inner)               // want `^SafeParam called with unsafe argument: argument is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructInlineNestedDef\.Inner\.Auth\.Password"\)$`
	svc1log.SafeParam("testParam", testStructInlineNestedDefVar)                     // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.TestStructInlineNestedDef", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.TestStructInlineNestedDef\.Inner\.Auth\.Password"\)$`

	type innerStruct struct {
		Password string `safelogging:"@Unsafe"` // want Password:`safe-logging-go/safeloggingtests/a\.innerStruct\.Password: @Unsafe \(tag on struct field "safe-logging-go/safeloggingtests/a\.innerStruct\.Password"\)`
	}
	var innerStructVar innerStruct
	svc1log.SafeParam("testParam", innerStructVar) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.innerStruct", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.innerStruct\.Password"\)$`

	anonymousStructVar := struct {
		Password string `safelogging:"@Unsafe"`
	}{}
	svc1log.SafeParam("testParam", anonymousStructVar) // want `^SafeParam called with unsafe argument: argument references type "safe-logging-go/safeloggingtests/a\.\[Anonymous Struct \(defined inline\)\]", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.\[Anonymous Struct \(defined inline\)\]\.Password"\)$`

	svc1log.SafeParam("testParam", struct { // want `^SafeParam called with unsafe argument: argument references "safe-logging-go/safeloggingtests/a\.\[Anonymous Struct \(defined inline\)\]", which is @Unsafe \(reason: tag on struct field "safe-logging-go/safeloggingtests/a\.\[Anonymous Struct \(defined inline\)\]\.Password"\)$`
		Password string `safelogging:"@Unsafe"`
	}{})
}
