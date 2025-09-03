package l_6 // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/l/l_6\.TestStructCommentEmbeddedInMultilineComment_6: @Unsafe\]\]`

type TestStructCommentEmbeddedInMultilineComment_6 struct { // safelogging:@Unsafe // 6.1 {CommentMap: *ast.Field.Names["UnsafeField_6"] (1/2)}
	/*


		// safelogging:@Unsafe // 6.2 {CommentMap: *ast.Field.Names["UnsafeField_6"] (2/2)}
	*/

	UnsafeField_6 string
}
