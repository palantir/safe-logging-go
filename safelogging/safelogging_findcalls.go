// Copyright 2025 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package safelogging

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/types"
	"iter"
	"maps"
	"slices"
	"strings"

	"github.com/palantir/safe-logging-go/safelogging/internal/filecomments"
	"golang.org/x/tools/go/analysis"
)

// FuncRef is a reference to a specific function. Matches the string representation of *types.Func, which is of the
// form "func (*net/http.Client).Do(req *net/http.Request) (*net/http.Response, error)".
type FuncRef string

const (
	// service.1 log param constructors: arguments must be correct safety level
	svc1logSafeParam    FuncRef = "github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.SafeParam"
	svc1logSafeParams   FuncRef = "github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.SafeParams"
	svc1logUnsafeParam  FuncRef = "github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.UnsafeParam"
	svc1logUnsafeParams FuncRef = "github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.UnsafeParams"

	// svc1log logging functions: message must be compile-time constant (params safety is covered by param construction)
	svc1logDebug FuncRef = "(github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Logger).Debug"
	svc1logInfo  FuncRef = "(github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Logger).Info"
	svc1logWarn  FuncRef = "(github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Logger).Warn"
	svc1logError FuncRef = "(github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Logger).Error"

	// werror params constructors: arguments must be correct safety level
	werrorSafeParam           FuncRef = "github.com/palantir/witchcraft-go-error.SafeParam"
	werrorSafeParams          FuncRef = "github.com/palantir/witchcraft-go-error.SafeParams"
	werrorUnsafeParam         FuncRef = "github.com/palantir/witchcraft-go-error.UnsafeParam"
	werrorUnsafeParams        FuncRef = "github.com/palantir/witchcraft-go-error.UnsafeParams"
	werrorSafeAndUnsafeParams FuncRef = "github.com/palantir/witchcraft-go-error.SafeAndUnsafeParams"
)

type FuncRefHandler func(
	funcRef FuncRef,
	id *ast.Ident,
	call *ast.CallExpr,
	pass *analysis.Pass,
	logSafetyInfos *logSafetyInfoBundle,
	allTypesMapWithUnderlyingTypes map[types.Type]*ast.Ident,
	commentLogSafetyTracker *CommentBasedLogSafetyTracker,
	fileCommentRetriever filecomments.Retriever,
	param Param,
)

func getFuncRefHandlers(compileTimeConstantArgFuncRefs map[FuncRef]int) map[FuncRef]FuncRefHandler {
	handlers := map[FuncRef]FuncRefHandler{
		svc1logSafeParam:    checkSingleParamThatIsSecondArgToFuncIsSafeToLog(LogSafetyTypeSafe),
		svc1logSafeParams:   checkSingleParamThatIsFirstArgToFuncIsSafeToLog(LogSafetyTypeSafe),
		svc1logUnsafeParam:  checkSingleParamThatIsSecondArgToFuncIsSafeToLog(LogSafetyTypeUnsafe),
		svc1logUnsafeParams: checkSingleParamThatIsFirstArgToFuncIsSafeToLog(LogSafetyTypeUnsafe),

		werrorSafeParam:    checkSingleParamThatIsSecondArgToFuncIsSafeToLog(LogSafetyTypeSafe),
		werrorSafeParams:   checkSingleParamThatIsFirstArgToFuncIsSafeToLog(LogSafetyTypeSafe),
		werrorUnsafeParam:  checkSingleParamThatIsSecondArgToFuncIsSafeToLog(LogSafetyTypeUnsafe),
		werrorUnsafeParams: checkSingleParamThatIsFirstArgToFuncIsSafeToLog(LogSafetyTypeUnsafe),
		werrorSafeAndUnsafeParams: checkParamsAtArgIdxToFuncIsSafeToLog(map[int]LogSafetyType{
			0: LogSafetyTypeSafe,
			1: LogSafetyTypeUnsafe,
		}),

		svc1logDebug: checkFirstArgIsCompileTimeConstant,
		svc1logInfo:  checkFirstArgIsCompileTimeConstant,
		svc1logWarn:  checkFirstArgIsCompileTimeConstant,
		svc1logError: checkFirstArgIsCompileTimeConstant,
	}
	for k, v := range compileTimeConstantArgFuncRefs {
		// do not allow overriding of built-in handlers
		if _, ok := handlers[k]; ok {
			continue
		}
		handlers[k] = checkArgAtIndexIsCompileTimeConstant(v)
	}
	return handlers
}

func checkSingleParamThatIsFirstArgToFuncIsSafeToLog(permittedSafetyLevel LogSafetyType) FuncRefHandler {
	return checkParamsAtArgIdxToFuncIsSafeToLog(map[int]LogSafetyType{
		0: permittedSafetyLevel,
	})
}

func checkSingleParamThatIsSecondArgToFuncIsSafeToLog(permittedSafetyLevel LogSafetyType) FuncRefHandler {
	return checkParamsAtArgIdxToFuncIsSafeToLog(map[int]LogSafetyType{
		1: permittedSafetyLevel,
	})
}

func checkParamsAtArgIdxToFuncIsSafeToLog(argIdxToSafetyLevel map[int]LogSafetyType) FuncRefHandler {
	sortedKeys := slices.Collect(maps.Keys(argIdxToSafetyLevel))
	slices.Sort(sortedKeys)
	return func(
		funcRef FuncRef,
		id *ast.Ident,
		call *ast.CallExpr,
		pass *analysis.Pass,
		logSafetyInfos *logSafetyInfoBundle,
		allTypesMapWithUnderlyingTypes map[types.Type]*ast.Ident,
		commentLogSafetyTracker *CommentBasedLogSafetyTracker,
		fileCommentRetriever filecomments.Retriever,
		param Param,
	) {
		for _, argIdx := range sortedKeys {
			permittedSafetyLevel := argIdxToSafetyLevel[argIdx]
			var paramArgExpr ast.Expr

			if argIdx >= len(call.Args) {
				// special case: function is called with fewer arguments than the argument that had safety declared.
				// Should only be possible in 2 cases:
				//  1. Function is provided with exactly 1 parameter, which is a call to another function. In this case,
				//     the function being called returns a tuple of multiple types, where the number of return types
				//     matches (or is compatible with) the number of declared arguments.
				//  2. Safety was declared on the last function parameter, that parameter is a vararg, and it was
				//     omitted on the function call
				//
				// Currently, the implementation does not handle case (2): currently, the function type safety is
				// hard-coded and it is known that there are no safety markings that match category (2). If this is
				// changed in the future, the implementation will need to be updated to handle this.

				// should only be possible if provided argument is a function call that returns multiple values
				isOK := false
				if len(call.Args) == 1 {
					if callExpr, ok := call.Args[0].(*ast.CallExpr); ok {
						if typeAndValue, ok := pass.TypesInfo.Types[callExpr]; ok {
							if tuple, ok := typeAndValue.Type.(*types.Tuple); ok {
								if tuple != nil && tuple.Len() > 1 {
									isOK = true
								}
							}
						}
					}
				}
				if !isOK {
					panic(fmt.Sprintf("failed to check safety of parameter at index %d for function %s: len(call.Args) == %d", argIdx, funcRef, len(call.Args)))
				}

				// if argument is function that returns multiple values, provide that expression for check.
				// The "checkExprSafetyLevelViolationInFunctionCall" function will use the "argIdx" parameter to check
				// the proper type.
				paramArgExpr = call.Args[0]
			} else {
				paramArgExpr = call.Args[argIdx]
			}
			issue := checkExprSafetyLevelViolationInFunctionCall(paramArgExpr, id, call, argIdx, pass, permittedSafetyLevel, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param)
			if issue == nil {
				continue
			}
			issue.Pos = paramArgExpr.Pos()
			issue.End = paramArgExpr.End()

			if nodeHasSuppressWarningsComment(id.Pos(), fileCommentRetriever) {
				continue
			}
			pass.Report(*issue)
		}
	}
}

func checkArgAtIndexIsCompileTimeConstant(idx int) FuncRefHandler {
	return func(
		_ FuncRef,
		id *ast.Ident,
		call *ast.CallExpr,
		pass *analysis.Pass,
		_ *logSafetyInfoBundle,
		_ map[types.Type]*ast.Ident,
		_ *CommentBasedLogSafetyTracker,
		fileCommentRetriever filecomments.Retriever,
		_ Param,
	) {
		// if argument index is out of range, do not report issue
		if idx >= len(call.Args) {
			return
		}

		if argIsCompileTimeConstant := checkExpressionIsCompileTimeConstant(call.Args[idx], id, call, pass); argIsCompileTimeConstant {
			return
		}
		if nodeHasSuppressWarningsComment(id.Pos(), fileCommentRetriever) {
			return
		}
		pass.Report(analysis.Diagnostic{
			Pos: call.Pos(),
			Message: fmt.Sprintf(
				"%s called with unsafe argument: message must be a compile-time constant",
				id.Name,
			),
		})
	}
}

var checkFirstArgIsCompileTimeConstant = checkArgAtIndexIsCompileTimeConstant(0)

func checkExpressionIsCompileTimeConstant(expr ast.Expr, id *ast.Ident, call *ast.CallExpr, pass *analysis.Pass) bool {
	switch exprVal := expr.(type) {
	default:
		return false
	case *ast.BasicLit:
		// basic literal value is safe
		return true
	case *ast.BinaryExpr:
		// check safety for LHS of expression
		return checkExpressionIsCompileTimeConstant(exprVal.X, id, call, pass) &&
			// check safety for RHS of expression
			checkExpressionIsCompileTimeConstant(exprVal.Y, id, call, pass)
	case *ast.SelectorExpr:
		// get the types.Selection for the selection expression and check if it is a constant
		if selection, ok := pass.TypesInfo.Selections[exprVal]; ok {
			switch selection.Obj().(type) {
			case constant.Value, *types.Const:
				return true
			}
		}

		// selector may select name from another package: verify that it is a constant
		if selectorObj, ok := pass.TypesInfo.Uses[exprVal.Sel]; ok {
			switch selectorObj.(type) {
			case constant.Value, *types.Const:
				return true
			}
		}
		return false
	case *ast.Ident:
		// get the types.Object for the identifier expression
		expressionObj := pass.TypesInfo.Uses[exprVal]
		if _, ok := expressionObj.(*types.Const); ok {
			return true
		}
		return false
	}
}

// checkExprSafetyLevelViolationInFunctionCall checks whether the provided expression (which should be a parameter to a
// function call) contains a reference that violates the specified permitted safety level.
//
// If a violation is found, calls "Report" on the provided "*analysis.Pass". The *ast.Ident and *ast.CallExpr are the
// values for the function call.
//
// This call may result in multiple "Report" calls: for example, if the provided expression is an *ast.BinaryExpression,
// the LHS and RHS may both contain expressions that have violations that must be reported.
func checkExprSafetyLevelViolationInFunctionCall(
	expr ast.Expr,
	id *ast.Ident,
	call *ast.CallExpr,
	argIdx int,
	pass *analysis.Pass,
	permittedSafetyLevel LogSafetyType,
	logSafetyInfos *logSafetyInfoBundle,
	allTypesMapWithUnderlyingTypes map[types.Type]*ast.Ident,
	commentLogSafetyTracker *CommentBasedLogSafetyTracker,
	fileCommentRetriever filecomments.Retriever,
	param Param,
) *analysis.Diagnostic {

	// "expr" is the expression that constitutes the argument to the function.
	switch exprVal := expr.(type) {
	case *ast.CallExpr:
		// determine safety based on type of call expression
		typeAndValue, ok := pass.TypesInfo.Types[expr]
		if !ok {
			break
		}

		typeAndValueType := typeAndValue.Type
		// if call expression returns multiple values, use the type of the return argument that corresponds to this one
		if tuple, ok := typeAndValue.Type.(*types.Tuple); ok && tuple != nil && tuple.Len() > 1 {
			typeAndValueType = tuple.At(argIdx).Type()
		}

		if leastSafeLogSafetyInfo := getLogSafetyForType(typeAndValueType, logSafetyInfos, param); leastSafeLogSafetyInfo.LogSafety > permittedSafetyLevel {
			return &analysis.Diagnostic{
				Pos: exprVal.Pos(),
				Message: fmt.Sprintf(
					"%s called with unsafe argument: argument references type %q, which is %s (reason: %s)",
					id.Name,
					leastSafeLogSafetyInfo.Identifier,
					leastSafeLogSafetyInfo.LogSafety,
					leastSafeLogSafetyInfo.Reason,
				),
			}
		}

		// determine if function itself is safe
		switch fnAstType := exprVal.Fun.(type) {
		case *ast.Ident:
			fnObj := pass.TypesInfo.Uses[fnAstType]
			if gotSafetyValue, ok := logSafetyInfos.objectLogSafety[fnObj]; ok {
				if leastSafeLogSafetyInfo := gotSafetyValue; leastSafeLogSafetyInfo.LogSafety > permittedSafetyLevel {
					return &analysis.Diagnostic{
						Pos: exprVal.Pos(),
						Message: fmt.Sprintf(
							"%s called with unsafe argument: argument references %q, which is %s (reason: %s)",
							id.Name,
							leastSafeLogSafetyInfo.Identifier,
							leastSafeLogSafetyInfo.LogSafety,
							leastSafeLogSafetyInfo.Reason,
						),
					}
				}
			}
		case *ast.SelectorExpr:
			fnObj := pass.TypesInfo.Uses[fnAstType.Sel]
			if gotSafetyValue, ok := logSafetyInfos.objectLogSafety[fnObj]; ok {
				if leastSafeLogSafetyInfo := gotSafetyValue; leastSafeLogSafetyInfo.LogSafety > permittedSafetyLevel {
					return &analysis.Diagnostic{
						Pos: exprVal.Pos(),
						Message: fmt.Sprintf(
							"%s called with unsafe argument: argument references %q, which is %s (reason: %s)",
							id.Name,
							leastSafeLogSafetyInfo.Identifier,
							leastSafeLogSafetyInfo.LogSafety,
							leastSafeLogSafetyInfo.Reason,
						),
					}
				}
			}
		}
	case *ast.CompositeLit:
		// check safety of declared type
		if issue := checkExprSafetyLevelViolationInFunctionCall(exprVal.Type, id, call, argIdx, pass, permittedSafetyLevel, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param); issue != nil {
			return issue
		}

		// check safety of types of each element
		for _, elt := range exprVal.Elts {
			if issue := checkExprSafetyLevelViolationInFunctionCall(elt, id, call, argIdx, pass, permittedSafetyLevel, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param); issue != nil {
				return issue
			}
		}
	case *ast.StructType:
		typeAndValue, ok := pass.TypesInfo.Types[expr]
		if !ok {
			break
		}
		// struct literal is defined inline (anonymous): compute safety of the struct
		if inlineStructLogSafetyValue, _, _ := computeStructLogSafetyInfo(
			"[Anonymous Struct (defined inline)]",
			typeAndValue.Type.(*types.Struct),
			exprVal,
			pass,
			logSafetyInfos,
			logSafetyInfos,
			allTypesMapWithUnderlyingTypes,
			nil,
			make(map[types.Type]*computeLogSafetyInfoForTypeReturnType),
			commentLogSafetyTracker,
			param,
		); inlineStructLogSafetyValue != nil && inlineStructLogSafetyValue.LogSafety > permittedSafetyLevel {
			return &analysis.Diagnostic{
				Pos: exprVal.Pos(),
				Message: fmt.Sprintf(
					"%s called with unsafe argument: argument references %q, which is %s (reason: %s)",
					id.Name,
					inlineStructLogSafetyValue.Identifier,
					inlineStructLogSafetyValue.LogSafety,
					inlineStructLogSafetyValue.Reason,
				),
			}
		}
	case *ast.KeyValueExpr:
		// check key and value
		if issue := checkExprSafetyLevelViolationInFunctionCall(exprVal.Key, id, call, argIdx, pass, permittedSafetyLevel, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param); issue != nil {
			return issue
		}
		if issue := checkExprSafetyLevelViolationInFunctionCall(exprVal.Value, id, call, argIdx, pass, permittedSafetyLevel, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param); issue != nil {
			return issue
		}
	case *ast.ArrayType:
		// check safety of array type
		return checkExprSafetyLevelViolationInFunctionCall(exprVal.Elt, id, call, argIdx, pass, permittedSafetyLevel, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param)
	case *ast.SliceExpr:
		// check safety of slice type
		return checkExprSafetyLevelViolationInFunctionCall(exprVal.X, id, call, argIdx, pass, permittedSafetyLevel, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param)
	case *ast.MapType:
		// check safety of key and value types
		if issue := checkExprSafetyLevelViolationInFunctionCall(exprVal.Key, id, call, argIdx, pass, permittedSafetyLevel, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param); issue != nil {
			return issue
		}
		if issue := checkExprSafetyLevelViolationInFunctionCall(exprVal.Value, id, call, argIdx, pass, permittedSafetyLevel, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param); issue != nil {
			return issue
		}
	case *ast.StarExpr:
		// check safety of expression being dereferenced
		return checkExprSafetyLevelViolationInFunctionCall(exprVal.X, id, call, argIdx, pass, permittedSafetyLevel, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param)
	case *ast.UnaryExpr:
		// check safety of expression being operated on
		return checkExprSafetyLevelViolationInFunctionCall(exprVal.X, id, call, argIdx, pass, permittedSafetyLevel, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param)
	case *ast.BinaryExpr:
		// check safety for LHS of expression
		if issue := checkExprSafetyLevelViolationInFunctionCall(exprVal.X, id, call, argIdx, pass, permittedSafetyLevel, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param); issue != nil {
			return issue
		}
		// check safety for RHS of expression
		if issue := checkExprSafetyLevelViolationInFunctionCall(exprVal.Y, id, call, argIdx, pass, permittedSafetyLevel, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param); issue != nil {
			return issue
		}
	case *ast.SelectorExpr:
		// get the types.Selection for the selection expression
		if selection, ok := pass.TypesInfo.Selections[exprVal]; ok {
			// selection.Obj() is the types.Object that represents the struct field.
			// These fields are populated in the logSafetyInfos map, so look up safety level for the field.
			if structArgSafetyValue := logSafetyInfos.objectLogSafety[selection.Obj()]; structArgSafetyValue.LogSafety > permittedSafetyLevel {
				return &analysis.Diagnostic{
					Pos: exprVal.Pos(),
					Message: fmt.Sprintf(
						"%s called with unsafe argument: argument is %s (reason: %s)",
						id.Name,
						structArgSafetyValue.LogSafety,
						structArgSafetyValue.Reason,
					),
				}
			}
		}

		// selector may be a type, variable, or constant from another package
		if selectorObj, ok := pass.TypesInfo.Uses[exprVal.Sel]; ok {
			// selector object is a type name (likely a struct literal declaration for a struct from another package):
			// check if type is safe
			if selectorObjTypName, ok := selectorObj.(*types.TypeName); ok {
				if leastSafeLogSafetyInfo := getLogSafetyForType(selectorObjTypName.Type(), logSafetyInfos, param); leastSafeLogSafetyInfo.LogSafety > permittedSafetyLevel {
					return &analysis.Diagnostic{
						Pos: exprVal.Pos(),
						Message: fmt.Sprintf(
							"%s called with unsafe argument: argument is of type %q, which is %s (reason: %s)",
							id.Name,
							leastSafeLogSafetyInfo.Identifier,
							leastSafeLogSafetyInfo.LogSafety,
							leastSafeLogSafetyInfo.Reason,
						),
					}
				}
			}

			// check if var or const has log safety specified
			if varConstArgSafetyValue := logSafetyInfos.objectLogSafety[selectorObj]; varConstArgSafetyValue.LogSafety > permittedSafetyLevel {
				return &analysis.Diagnostic{
					Pos: exprVal.Pos(),
					Message: fmt.Sprintf(
						"%s called with unsafe argument: argument is %s (reason: %s)",
						id.Name,
						varConstArgSafetyValue.LogSafety,
						varConstArgSafetyValue.Reason,
					),
				}
			}
		}
	case *ast.Ident:
		// get the types.Object for the Identifier expression
		expressionObj := pass.TypesInfo.Uses[exprVal]

		// determine log safety level for expression's type and report
		if leastSafeLogSafetyInfo := getLogSafetyForType(expressionObj.Type(), logSafetyInfos, param); leastSafeLogSafetyInfo.LogSafety > permittedSafetyLevel {
			return &analysis.Diagnostic{
				Pos: exprVal.Pos(),
				Message: fmt.Sprintf(
					"%s called with unsafe argument: argument references type %q, which is %s (reason: %s)",
					id.Name,
					leastSafeLogSafetyInfo.Identifier,
					leastSafeLogSafetyInfo.LogSafety,
					leastSafeLogSafetyInfo.Reason,
				),
			}
		}

		if leastSafeLogSafetyInfo, ok := logSafetyInfos.objectLogSafety[expressionObj]; ok && leastSafeLogSafetyInfo.LogSafety > permittedSafetyLevel {
			return &analysis.Diagnostic{
				Pos: exprVal.Pos(),
				Message: fmt.Sprintf(
					"%s called with unsafe argument: argument is %s (reason: %s)",
					id.Name,
					leastSafeLogSafetyInfo.LogSafety,
					leastSafeLogSafetyInfo.Reason,
				),
			}
		}
	}
	return nil
}

func getLogSafetyForType(targetType types.Type, logSafetyInfos *logSafetyInfoBundle, param Param) LogSafetyInfo {
	var leastSafeLogSafetyInfo LogSafetyInfo

	// look up safety value based on type information
	if gotSafetyValue, ok := logSafetyInfos.typeLogSafety[targetType]; ok {
		leastSafeLogSafetyInfo = gotSafetyValue
	}

	setTypeSafetyFromBuiltinMapFn := func(namedType *types.Named) {
		identifier := namedType.Obj().Name()
		if pkg := namedType.Obj().Pkg(); pkg != nil {
			identifier = pkg.Path() + "." + identifier
		}
		if safetyValue, ok := param.typeSafetyMap[identifier]; ok && safetyValue > leastSafeLogSafetyInfo.LogSafety {
			leastSafeLogSafetyInfo = LogSafetyInfo{
				Identifier: identifier,
				Reason:     fmt.Sprintf("configuration specified log safety value for type %q", identifier),
				LogSafety:  safetyValue,
			}
		}
	}

	if namedType, ok := targetType.(*types.Named); ok {
		setTypeSafetyFromBuiltinMapFn(namedType)
	}

	// also get type safety information of all base types of type
	baseTypes := getAllBaseTypes(targetType)
	// a given semantic named type may have multiple entries in the typeLogSafetyMap: one for the actual named type, and
	// one for its underlying type. In such a case, prefer the named type, since its failure message/output will be
	// more informative.
	assignedFromNamedType := false
	for _, baseType := range baseTypes {
		if safetyValue, ok := logSafetyInfos.typeLogSafety[baseType]; ok && safetyValue.LogSafety >= leastSafeLogSafetyInfo.LogSafety {
			if assignedFromNamedType {
				continue
			}
			if _, currTypeIsNamedType := baseType.(*types.Named); currTypeIsNamedType {
				assignedFromNamedType = true
			}
			leastSafeLogSafetyInfo = safetyValue
		}

		if namedType, ok := baseType.(*types.Named); ok {
			setTypeSafetyFromBuiltinMapFn(namedType)
		}
	}

	return leastSafeLogSafetyInfo
}

func findLogCalls(
	pass *analysis.Pass,
	logSafetyInfos *logSafetyInfoBundle,
	allTypesMapWithUnderlyingTypes map[types.Type]*ast.Ident,
	commentLogSafetyTracker *CommentBasedLogSafetyTracker,
	fileCommentRetriever filecomments.Retriever,
	param Param,
) {
	funcRefHandlers := getFuncRefHandlers(param.constMessageLoggingFunctions)

	// Return a mapping from *ast.Ident to FuncRef for all relevant functions.
	// The calls to these functions are what will be analyzed for the check.
	astIdentToFuncRefMap := createASTIdentToFuncRefMap(pass.TypesInfo.Uses, maps.Keys(funcRefHandlers))

	// go through AST of all files to find call sites of identified functions
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			// function call references type *ast.CallExpr
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// identifier for the function call: *ast.Ident if called directly as an identifier ("Foo()"),
			// *ast.SelectorExpr if invoked on a package or receiver ("otherpkg.Foo()").
			// May be *ast.IndexExpr if function is a generic function with type parameters ("Foo[T]()").
			var id *ast.Ident
			switch fun := call.Fun.(type) {
			case *ast.Ident:
				id = fun
			case *ast.SelectorExpr:
				id = fun.Sel
			case *ast.IndexExpr:
				xIdent, ok := fun.X.(*ast.Ident) // panic if not ident
				if ok {
					id = xIdent
				}
			}

			funcRefForCall, ok := astIdentToFuncRefMap[id]

			// the id for this CallExpr is not in the map of call references
			if !ok {
				return true
			}

			// get the handler for this function reference and perform the check operation
			funcRefHandler := funcRefHandlers[funcRefForCall]
			funcRefHandler(funcRefForCall, id, call, pass, logSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param)

			return true
		})
	}
}

// createASTIdentToFuncRefMap returns a map from *ast.Ident to FuncRef for all the function references in the specified
// package. If the "funcRefs" argument is non-empty, then only function signature that match an element in "funcRefs"
// are included; otherwise, all function references are returned.
func createASTIdentToFuncRefMap(uses map[*ast.Ident]types.Object, funcRefs iter.Seq[FuncRef]) map[*ast.Ident]FuncRef {
	// map from identifiers to the function reference
	identToFuncRefs := make(map[*ast.Ident]FuncRef)

	funcRefsMap := make(map[FuncRef]struct{})
	for funcRef := range funcRefs {
		funcRefsMap[funcRef] = struct{}{}
	}

	var keys []*ast.Ident
	for k := range uses {
		keys = append(keys, k)
	}

	for _, id := range keys {
		obj := uses[id]
		funcPtr, ok := obj.(*types.Func)
		if !ok {
			continue
		}

		currFuncRef := FuncRef(normalizeFunctionString(funcPtr.String()))
		if len(funcRefsMap) > 0 {
			if _, ok := funcRefsMap[currFuncRef]; !ok {
				// if funcRefsMap is non-empty, skip any entries that don't match the signature
				continue
			}
		}
		// record function reference
		identToFuncRefs[id] = currFuncRef
	}
	return identToFuncRefs
}

func normalizeFunctionString(in string) string {
	// remove "func " prefix
	normalized := strings.TrimPrefix(in, "func ")

	// remove generics
	if openSquareBraceIdx := strings.Index(normalized, "["); openSquareBraceIdx != -1 {
		normalized = normalized[:openSquareBraceIdx]
	}

	// remove parameter and return value if there is an open paren after index 0 (index 0 is OK, because that is how receivers are denoted)
	if openParenIdx := strings.LastIndex(normalized, "("); openParenIdx > 0 {
		normalized = normalized[:openParenIdx]
	}

	// result should be just function name
	return normalized
}
