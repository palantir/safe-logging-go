package l_1 // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/l/l_1\.TestStruct_1: @Unsafe\]\]`

// safelogging:@Unsafe 1.1 // {CommentMap: *ast.GenDecl.Specs[0](*ast.TypeSpec).Name("TestStruct_1") (3/3)}
type TestStruct_1 struct {

	// safelogging:@Unsafe // 1.2 {CommentMap: *ast.Field.Names["Field_1_1"] (1/1)}
	Field_1_1 string // want Field_1_1:`safe-logging-go/safeloggingtests/l/l_1\.Field_1_1: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_1\.Field_1_1" at .+/safelogging/testdata/l/l_1/safelogging_l_1\.go:6:2\)`

	// safelogging:@Unsafe // 1.3 {CommentMap: *ast.Field.Names["Field_1_2"] (1/1)} {*ast.Field.Names["Field_1_2"], *ast.Field.Doc.List[1/1]}
	Field_1_2 string // want Field_1_2:`safe-logging-go/safeloggingtests/l/l_1\.Field_1_2: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_1\.Field_1_2" at .+/safelogging/testdata/l/l_1/safelogging_l_1\.go:9:2\)`

	// safelogging:@Safe // 1.4 {CommentMap: *ast.Field.Names["Field_1_3", "Field_1_4", "Field_1_5"] (1/1)}
	Field_1_3, // want Field_1_3:`safe-logging-go/safeloggingtests/l/l_1\.Field_1_3: @Safe \(comment on "safe-logging-go/safeloggingtests/l/l_1\.Field_1_3" at .+/safelogging/testdata/l/l_1/safelogging_l_1\.go:12:2\)`
	// safelogging:@Unsafe // 1.5 {CommentMap: *ast.Ident.Name("Field_1_4") (1/1)}
	Field_1_4, // want Field_1_4:`safe-logging-go/safeloggingtests/l/l_1\.Field_1_4: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_1\.Field_1_4" at .+/safelogging/testdata/l/l_1/safelogging_l_1\.go:14:2\)`
	// safelogging:@DoNotLog // 1.6 {CommentMap: *ast.Ident.Name("Field_1_5") (1/1)}
	Field_1_5 string // want Field_1_5:`safe-logging-go/safeloggingtests/l/l_1\.Field_1_5: @DoNotLog \(comment on "safe-logging-go/safeloggingtests/l/l_1\.Field_1_5" at .+/safelogging/testdata/l/l_1/safelogging_l_1\.go:16:2\)`

	// safelogging:@Unsafe // 1.7 {CommentMap: *ast.Field.Names["Field_1_6", "Field_1_7"] (1/1)}
	Field_1_6, Field_1_7 string // want Field_1_6:`safe-logging-go/safeloggingtests/l/l_1\.Field_1_6: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_1\.Field_1_6" at .+/safelogging/testdata/l/l_1/safelogging_l_1\.go:19:2\)` Field_1_7:`safe-logging-go/safeloggingtests/l/l_1\.Field_1_7: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_1\.Field_1_7" at .+/safelogging/testdata/l/l_1/safelogging_l_1\.go:19:2\)`

	// safelogging:@Unsafe // 1.8 {CommentMap: *ast.Field.Names["Field_1_8"] (1/1)}
	Field_1_8 struct { // want Field_1_8:`safe-logging-go/safeloggingtests/l/l_1\.TestStruct_1\.Field_1_8: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_1\.TestStruct_1\.Field_1_8" at .+/safelogging/testdata/l/l_1/safelogging_l_1\.go:22:2\)`
		Inner string
	}
}
