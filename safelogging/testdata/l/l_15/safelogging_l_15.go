package l_15

type TestAmbiguous_15 struct {
}                        // safelogging:@Safe // 15.1 {CommentMap: *ast.GenDecl.Specs[0].Name("TestAmbiguous_15") (1/1)}
var GlobalField15 string // want GlobalField15:`safe-logging-go/safeloggingtests/l/l_15\.GlobalField15: @Safe \(comment on "safe-logging-go/safeloggingtests/l/l_15\.GlobalField15" at .+/safelogging/testdata/l/l_15/safelogging_l_15\.go:4:26\)`
