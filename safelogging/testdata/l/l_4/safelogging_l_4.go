package l_4 // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/l/l_4\.TestStructMultiLineDef_4: @Unsafe\]\]`

// safelogging:@Safe // 4.1 {CommentMap: *ast.GenDecl.Specs[0](*ast.TypeSpec).Name("TestStructMultiLineDef_4") (1/2)}
type // safelogging:@Safe // 4.2 {CommentMap: *ast.TypeSpec.Name("TestStructMultiLineDef_4") (1/2)}
// safelogging:@Safe // 4.3 {CommentMap: *ast.TypeSpec.Name("TestStructMultiLineDef_4") (2/2)}
TestStructMultiLineDef_4 struct // safelogging:@Unsafe // 4.4 {CommentMap: *ast.FieldList.List[0].Names["SafeField_4"] (1/2)}
// safelogging:@Unsafe // 4.5 {CommentMap: *ast.FieldList.List[0].Names["SafeField_4"] (2/2)}
{ // safelogging:@DoNotLog // 4.6 {CommentMap: *ast.Field.Names["SafeField_4"] (1/2)}
	SafeField_4 string // safelogging:@Safe // 4.7 {CommentMap: *ast.Field.Names["SafeField_4"] (2/2)} {*ast.FieldList.List[0].Names["SafeField_4"], *ast.FieldList.List[0].Comment.List[1/1]} // want SafeField_4:`safe-logging-go/safeloggingtests/l/l_4\.SafeField_4: @Safe \(comment on "safe-logging-go/safeloggingtests/l/l_4\.SafeField_4" at .+/safelogging/testdata/l/l_4/safelogging_l_4\.go:9:21\)`
} // safelogging:@Safe 4.8 // {CommentMap: *ast.GenDecl.Specs[0](*ast.TypeSpec).Name("TestStructMultiLineDef_4") (2/2)}
