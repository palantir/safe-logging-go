package l_2 // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/l/l_2\.TestStruct_2: @Unsafe\]\]`

type TestStruct_2 struct { // safelogging:@Unsafe // 2.0
	Field_2_1 string // safelogging:@Unsafe // 2.1 {CommentMap: *ast.Field.Names["Field_2_1"] (1/1)} {*ast.Field.Names["Field_2_1"], *ast.Field.Comment.List[1/1]} // want Field_2_1:`safe-logging-go/safeloggingtests/l/l_2\.Field_2_1: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_2\.Field_2_1" at .+/safelogging/testdata/l/l_2/safelogging_l_2\.go:4:19\)`

	Field_2_2, // safelogging:@Unsafe // 2.2 {CommentMap: *ast.Ident.Name("Field_2_2") (1/1)} // want Field_2_2:`safe-logging-go/safeloggingtests/l/l_2\.Field_2_2: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_2\.Field_2_2" at .+/safelogging/testdata/l/l_2/safelogging_l_2\.go:6:13\)`
	Field_2_3, // safelogging:@Unsafe // 2.3 {CommentMap: *ast.Ident.Name("Field_2_3") (1/1)} // want Field_2_3:`safe-logging-go/safeloggingtests/l/l_2\.Field_2_3: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_2\.Field_2_3" at .+/safelogging/testdata/l/l_2/safelogging_l_2\.go:7:13\)`
	Field_2_4 string // safelogging:@Unsafe // 2.4 {CommentMap: *ast.Field.Names["Field_2_2", "Field_2_3", "Field_2_4"] (1/1)} {*ast.Field.Names["Field_2_2", "Field_2_3", "Field_2_4"], *ast.Field.Comment.List[1/1]} // want Field_2_4:`safe-logging-go/safeloggingtests/l/l_2\.Field_2_4: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_2\.Field_2_4" at .+/safelogging/testdata/l/l_2/safelogging_l_2\.go:8:19\)`

	Field_2_5, Field_2_6 string // safelogging:@Unsafe // 2.5 {CommentMap: *ast.Field.Names["Field_2_5", "Field_2_6"] (1/1)} {*ast.Field.Names["Field_2_5", "Field_2_6"], *ast.Field.Comment.List[1/1]} // want Field_2_5:`safe-logging-go/safeloggingtests/l/l_2\.Field_2_5: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_2\.Field_2_5" at .+/safelogging/testdata/l/l_2/safelogging_l_2\.go:10:30\)` Field_2_6:`safe-logging-go/safeloggingtests/l/l_2\.Field_2_6: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_2\.Field_2_6" at .+/safelogging/testdata/l/l_2/safelogging_l_2\.go:10:30\)`

	Field_2_7 struct { // safelogging:@Unsafe // 2.6 {CommentMap: *ast.Field.Names["Field_2_7"] (1/1)} {*ast.Field.Names["Field_2_7"], *ast.Field.Doc.List[1/1]} // want Field_2_7:`safe-logging-go/safeloggingtests/l/l_2\.TestStruct_2\.Field_2_7: @Unsafe \(comment on "safe-logging-go/safeloggingtests/l/l_2\.TestStruct_2\.Field_2_7" at .+/safelogging/testdata/l/l_2/safelogging_l_2\.go:12:21\)`
		Inner string
	}
}
