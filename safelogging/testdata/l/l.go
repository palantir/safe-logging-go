package l

// These tests exist to test how comments are parsed/handled.
// Each comment is annotated with how it is categorized.
//
// Comment-based marking can be done by adding a comment of the form "// safelogging:{Level}" in one of 2 locations:
//   1. On the same line as the identifier
//   2. On the line immediately preceding the identifier
//
// The following are all examples of comment-based marking applying to the "Foo" struct based on rule (1):
//
// ----
//   type Foo struct { // safelogging:@Safe
//      Field string
//   }
//
//   type
//   Foo struct { // safelogging:@Safe
//      Field string
//   }
//
//   type
//   Foo struct // safelogging:@Safe
//   {
//      Field string
//   }
// ----
//
// Although it is not recommended, if there are multiple declarations on the same line and there is a comment that
// matches rule (1) on that line, it matches all declarations on that line. For example, in the following cases, both
// "Foo" and "Bar" are marked as "Safe" by the comment on the same line as the identifier:
//
// ----
// type Foo struct{}; var Bar string // safelogging:@Safe
//
// var Foo, Bar string // safelogging:@Safe
// ----
//
// If a particular comment matches a declaration based on rule (1), it is assigned to that declaration and cannot be
// used to mark another field. For example:
//
// ----
//  type Foo struct {
// 	  Field_1 string // safelogging:@Safe
//	  Field_2 string
//  }
// ----
//
// In this example, the "safelogging" comment satisfies rule (1) for Field_1 and rule (2) for Field_2. However, because
// rule (1) takes precedence, the comment is assigned to Field_1, after which it is not considered for Field_2.
//
// Similarly, for:
//
// ----
//  type Foo struct{} // safelogging:@Safe
//  var Var string
// ----
//
// The "safelogging" comment satisfies rule (1) for Foo and rule (2) for Var. However, because rule (1) takes
// precedence, the comment is assigned to Foo, after which it is not considered for Field.
//
// A comment can satisfy rule (2) even if there is text that occurs on the same line as the comment.
// For example:
//
//  type Foo struct {
//  } // safelogging:@Safe
//  var Var string
//
// In this example, the "safelogging" comment does not satisfy rule (1) (because it does not occur on the same line as
// an identifier), but it does satisfy rule (2) because it occurs on the line immediately preceding an identifier.
// Thus, in this case, "Var" would be marked as "Safe". Although this is technically correct, it is ambiguous, and in
// practice it would be better to use one of the following forms instead:
//
//  type Foo struct {
//  }
//  // safelogging:@Safe
//  var Var string
//
//  type Foo struct {
//  }
//  var Var string // safelogging:@Safe
