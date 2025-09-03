package l_5 // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/l/l_5\.TestStructCommentAfterNewline_5: @Unsafe\]\]`

// safelogging:@Unsafe // 5.1 {CommentMap: *ast.GenDecl.Specs[0].Name("TestStructCommentAfterNewline_5") (1/1)} {*ast.GenDecl.Doc.List[1/1]}
type TestStructCommentAfterNewline_5 struct {

	// safelogging:@Unsafe // 5.2 {CommentMap: *ast.Field.Names["UnsafeField_5"] (1/1)}

	UnsafeField_5 string
}
