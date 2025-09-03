package l_16 // want package:`TypeRepToTypeToLogSafety: \[\*types\.Named: \[safe-logging-go/safeloggingtests/l/l_16\.TestSameLine_16: @Safe\]\]`

// NOTE: editors will often force-reformat this file.
// The content of this file is supposed to be a single line that contains 2 statements where both a struct and variable are declared on the same line:
//
// type TestSameLine_16 struct{}; var SameLineVar_16 string // safelogging:@Safe // 16.1 {CommentMap: *ast.File.Decls[*ast.GenDecl.Specs[0].Name("TestSameLine_16"), *ast.GenDecl.Specs[0].Names["SameLineVar_16"]] (1/1)}

type TestSameLine_16 struct{}; var SameLineVar_16 string // safelogging:@Safe // 16.1 {CommentMap: *ast.File.Decls[*ast.GenDecl.Specs[0].Name("TestSameLine_16"), *ast.GenDecl.Specs[0].Names["SameLineVar_16"]] (1/1)} // want SameLineVar_16:`safe-logging-go/safeloggingtests/l/l_16\.SameLineVar_16: @Safe \(comment on "safe-logging-go/safeloggingtests/l/l_16\.SameLineVar_16" at .+/safelogging/testdata/l/l_16/safelogging_l_16\.go:8:58\)`
