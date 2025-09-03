package l_12 // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/l/l_12\.TestStructCommentOnSameLineNoContent_12: @Unsafe\]\]`

type TestStructCommentOnSameLineNoContent_12 struct{} // safelogging:@Unsafe // 12.1 {CommentMap: *ast.GenDecl.Specs[0](*ast.TypeSpec).Name("TestStructCommentOnSameLineNoContent_12") (2/2)}

type TestStructCommentOnSameLineNoContent_12_1 struct { // safelogging:@Unsafe 13.1 {CommentMap: *ast.FieldList (1/1)}
} // safelogging:@Unsafe // 12.2 {CommentMap: *ast.File (1/1)}
