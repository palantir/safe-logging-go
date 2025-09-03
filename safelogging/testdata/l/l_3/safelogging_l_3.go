package l_3 // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/l/l_3\.TestStruct_3: @Unsafe\]\]`

// safelogging:@Unsafe // 3.1 {CommentMap: *ast.GenDecl.Specs[0](*ast.TypeSpec).Name("TestStruct_3") (1/2)} {*ast.GenDecl.Doc.List[1/1]}
type TestStruct_3 struct {
	// safelogging:@Unsafe // 3.2 {CommentMap: *ast.Field.Names["Field_3_1"] (1/2)}
	Field_3_1 string // safelogging:@Unsafe // 3.3 {CommentMap: *ast.Field.Names["Field_3_1"] (2/2)} // want Field_3_1:`safe-logging-go/safeloggingtests/l/l_3\.Field_3_1: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_3\.Field_3_1" at .+/safelogging/testdata/l/l_3/safelogging_l_3\.go:6:19\)`

	// safelogging:@Unsafe // 3.4 {CommentMap: *ast.Field.Names["Field_3_2", "Field_3_3", "Field_3_4"] (1/2)}
	Field_3_2, // safelogging:@Unsafe // 3.5 {CommentMap: *ast.Ident.Name("Field_3_2") (1/1)} // want Field_3_2:`safe-logging-go/safeloggingtests/l/l_3\.Field_3_2: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_3\.Field_3_2" at .+/safelogging/testdata/l/l_3/safelogging_l_3\.go:9:13\)`
	// safelogging:@Unsafe // 3.6 {CommentMap: *ast.Ident.Name("Field_3_3") (1/2)}
	Field_3_3, // safelogging:@Unsafe // 3.7 {CommentMap: *ast.Ident.Name("Field_3_3") (2/2)} // want Field_3_3:`safe-logging-go/safeloggingtests/l/l_3\.Field_3_3: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_3\.Field_3_3" at .+/safelogging/testdata/l/l_3/safelogging_l_3\.go:11:13\)`
	// safelogging:@Unsafe // 3.8 {CommentMap: *ast.Ident.Name("Field_3_4") (1/1)}
	Field_3_4 string // safelogging:@Unsafe // 3.9 {CommentMap: *ast.Field.Names["Field_3_2", "Field_3_3", "Field_3_4"] (2/2)} // want Field_3_4:`safe-logging-go/safeloggingtests/l/l_3\.Field_3_4: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_3\.Field_3_4" at .+/safelogging/testdata/l/l_3/safelogging_l_3\.go:13:19\)`

	// safelogging:@Unsafe // 3.10 {CommentMap: *ast.Field.Names["Field_3_5", "Field_3_6"] (1/2)}
	Field_3_5, Field_3_6 string // safelogging:@Unsafe // 3.11 {CommentMap: *ast.Field.Names["Field_3_5", "Field_3_6"] (2/2)} // want Field_3_5:`safe-logging-go/safeloggingtests/l/l_3\.Field_3_5: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_3\.Field_3_5" at .+/safelogging/testdata/l/l_3/safelogging_l_3\.go:15:2\)` Field_3_6:`safe-logging-go/safeloggingtests/l/l_3\.Field_3_6: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_3\.Field_3_6" at .+/safelogging/testdata/l/l_3/safelogging_l_3\.go:15:2\)`

	// safelogging:@Unsafe // 3.12 {CommentMap: *ast.Field.Names["Field_3_7"] (1/2)}
	Field_3_7 struct { // want Field_3_7:`safe-logging-go/safeloggingtests/l/l_3\.TestStruct_3\.Field_3_7: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_3\.TestStruct_3\.Field_3_7" at .+/safelogging/testdata/l/l_3/safelogging_l_3\.go:18:2\)`
		Inner string
	} // safelogging:@Unsafe // 3.13 {CommentMap: *ast.Field.Names["Field_3_7"] (2/2)} {*ast.Field.Names["Field_3_7"], *ast.Field.Comment.List[1/1]}
} // safelogging:@Unsafe // 3.14 {CommentMap: *ast.GenDecl.Specs[0](*ast.TypeSpec).Name("TestStruct_3") (2/2)} {*ast.GenDecl.Specs[0](*ast.TypeSpec).Name("TestStruct_3"), *ast.GenDecl.Specs[0](*ast.TypeSpec).Comment.List[1/1]}
