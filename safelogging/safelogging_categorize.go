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
	"go/token"
	"go/types"
	"maps"
	"reflect"
	"regexp"

	"github.com/palantir/safe-logging-go/safelogging/internal/filecomments"
	"golang.org/x/tools/go/analysis"
)

func computeLogSafetyInfo(
	pass *analysis.Pass,
	dependentPkgLogSafetyInfos *logSafetyInfoBundle,
	commentLogSafetyTracker *CommentBasedLogSafetyTracker,
	param Param,
) (*logSafetyInfoBundle, map[types.Type]*ast.Ident, error) {
	// collect all the types to analyze from the type definitions
	typeMaps := collectTypes(pass.TypesInfo.Defs, pass.TypesInfo.Types)

	// log safety information for all struct fields defined in package
	logSafetyInfoForPkg := &logSafetyInfoBundle{
		typeLogSafety:   make(map[types.Type]LogSafetyInfo),
		primaryTypes:    make(map[types.Type]struct{}),
		objectLogSafety: make(map[types.Object]LogSafetyInfo),
	}

	// a copy of typeInfoMaps.structTypesMap where every key is reinserted as its underlying type with the same value
	allTypesMapWithUnderlyingTypes := make(map[types.Type]*ast.Ident)
	for k, v := range typeMaps.structTypesMap {
		allTypesMapWithUnderlyingTypes[k] = v
		allTypesMapWithUnderlyingTypes[k.Underlying()] = v
	}
	maps.Copy(allTypesMapWithUnderlyingTypes, typeMaps.namedTypesMap)

	cachedReturnValues := make(map[types.Type]*computeLogSafetyInfoForTypeReturnType)

	// iterate over all struct type definitions and compute log safety for type and fields
	for structType, structASTIdent := range typeMaps.structTypesMap {
		typeLogSafetyInfo, typeFieldsLogSafetyInfo, cyclicOnly := computeLogSafetyInfoForType(
			structType,
			structASTIdent,
			pass.TypesInfo.Defs[structASTIdent],
			pass,
			logSafetyInfoForPkg,
			dependentPkgLogSafetyInfos,
			allTypesMapWithUnderlyingTypes,
			nil,
			cachedReturnValues,
			commentLogSafetyTracker,
			param,
		)
		if typeLogSafetyInfo == nil {
			// no type safety information because it is a struct that has only cyclical fields
			if cyclicOnly {
				continue
			}
			panic(fmt.Sprintf("failed to compute log safety info for type %v\n", structType))
		}
		logSafetyInfoForPkg.typeLogSafety[structType] = *typeLogSafetyInfo
		// record that this type is a primary type
		logSafetyInfoForPkg.primaryTypes[structType] = struct{}{}
		// record type safety information for underlying type as well as named type
		logSafetyInfoForPkg.typeLogSafety[structType.Underlying()] = *typeLogSafetyInfo
		maps.Copy(logSafetyInfoForPkg.objectLogSafety, typeFieldsLogSafetyInfo)
	}

	// iterate over all anonymous struct definitions and compute log safety for types
	for anonStructType, anonStructAstStructType := range typeMaps.anonStructTypesMap {
		anonStructLogSafetyValue, _, cyclicOnly := computeStructLogSafetyInfo(
			"[Anonymous Struct (defined inline)]",
			anonStructType,
			anonStructAstStructType,
			pass,
			logSafetyInfoForPkg,
			dependentPkgLogSafetyInfos,
			allTypesMapWithUnderlyingTypes,
			nil,
			cachedReturnValues,
			commentLogSafetyTracker,
			param,
		)
		if anonStructLogSafetyValue == nil {
			// no type safety information because it is a struct that has only cyclical fields
			if cyclicOnly {
				continue
			}
			panic(fmt.Sprintf("failed to compute log safety info for type %v\n", anonStructType))
		}

		underlyingType := anonStructType.Underlying()
		// if underlying type already has a log safety populated for it, do not update, as the information from the
		// non-anonymous struct is likely more useful
		if _, ok := logSafetyInfoForPkg.typeLogSafety[underlyingType]; !ok {
			logSafetyInfoForPkg.typeLogSafety[underlyingType] = *anonStructLogSafetyValue
		}
	}

	// iterate over all named type definitions and compute log safety for type
	for namedType, namedTypeASTIdent := range typeMaps.namedTypesMap {
		typeLogSafetyInfo, _, cyclicOnly := computeLogSafetyInfoForType(
			namedType,
			namedTypeASTIdent,
			pass.TypesInfo.Defs[namedTypeASTIdent],
			pass,
			logSafetyInfoForPkg,
			dependentPkgLogSafetyInfos,
			allTypesMapWithUnderlyingTypes,
			nil,
			cachedReturnValues,
			commentLogSafetyTracker,
			param,
		)
		if typeLogSafetyInfo == nil {
			// no type safety information because it is a struct that has only cyclical fields
			if cyclicOnly {
				continue
			}
			panic(fmt.Sprintf("failed to compute log safety info for type %v", namedType))
		}
		logSafetyInfoForPkg.typeLogSafety[namedType] = *typeLogSafetyInfo
	}

	// iterate over all function definitions and compute log safety for functions
	for fnObject, fnASTIdent := range typeMaps.functionsMap {
		typeLogSafetyInfo := computeLogSafetyInfoASTIdentFromComment(fnASTIdent, pass, commentLogSafetyTracker)
		logSafetyInfoForPkg.objectLogSafety[fnObject] = *typeLogSafetyInfo
	}

	// iterate over all variable and constant definitions and compute log safety for them
	for varConstObject, varConstASTIdent := range typeMaps.varConstsMap {
		typeLogSafetyInfo := computeLogSafetyInfoASTIdentFromComment(varConstASTIdent, pass, commentLogSafetyTracker)

		// set log safety for variable or constant object if it has not already been set or if the comment specifies a more restrictive value
		if currVal, ok := logSafetyInfoForPkg.objectLogSafety[varConstObject]; !ok || currVal.LogSafety < typeLogSafetyInfo.LogSafety {
			logSafetyInfoForPkg.objectLogSafety[varConstObject] = *typeLogSafetyInfo
		}
	}

	return logSafetyInfoForPkg, allTypesMapWithUnderlyingTypes, nil
}

type typeInfoMaps struct {
	// named type definitions (defined types/type aliases)
	namedTypesMap map[types.Type]*ast.Ident

	// struct type definitions
	structTypesMap map[types.Type]*ast.Ident

	// anonymous struct type definitions
	anonStructTypesMap map[*types.Struct]*ast.StructType

	// function definitions
	functionsMap map[types.Object]*ast.Ident

	// variable and constant declarations
	varConstsMap map[types.Object]*ast.Ident
}

func collectTypes(typeDefs map[*ast.Ident]types.Object, typesMap map[ast.Expr]types.TypeAndValue) *typeInfoMaps {
	typeMapsVar := &typeInfoMaps{
		namedTypesMap:      make(map[types.Type]*ast.Ident),
		structTypesMap:     make(map[types.Type]*ast.Ident),
		anonStructTypesMap: make(map[*types.Struct]*ast.StructType),
		functionsMap:       make(map[types.Object]*ast.Ident),
		varConstsMap:       make(map[types.Object]*ast.Ident),
	}

	for typeDefAstIdent, typeDefInPkg := range typeDefs {
		// add struct type information if this is a struct type
		addNamedType(typeDefAstIdent, typeDefInPkg, typeMapsVar)

		// add function information if this is a function
		addFunction(typeDefAstIdent, typeDefInPkg, typeMapsVar)

		// add var/const information
		addVarConst(typeDefAstIdent, typeDefInPkg, typeMapsVar)
	}

	// typesMap contains type definitions. There are some scenarios in which a struct is in this map but not in the
	// "typeDefs" map: for example, some anonymous structs fall into this category
	for astExpr, typeAndValue := range typesMap {
		if astStructType, ok := astExpr.(*ast.StructType); ok {
			if typeStruct, ok := typeAndValue.Type.(*types.Struct); ok {
				typeMapsVar.anonStructTypesMap[typeStruct] = astStructType
			}
		}
	}

	return typeMapsVar
}

func addNamedType(typeDefAstIdent *ast.Ident, typeDefInPkg types.Object, typeMapsVar *typeInfoMaps) {
	// skip nodes that do not define type
	if typeDefAstIdent.Obj == nil || typeDefAstIdent.Obj.Kind != ast.Typ {
		return
	}
	typeSpecDecl, ok := typeDefAstIdent.Obj.Decl.(*ast.TypeSpec)
	if !ok {
		return
	}

	// if type of object is not *types.Named or *types.Alias, skip: this can happen for embedded struct fields that
	// are themselves structs, where the type will be *types.Var
	switch typeDefInPkg.Type().(type) {
	case *types.Named, *types.Alias:
		// named or alias type: process

		switch typeSpecDecl.Type.(type) {
		case *ast.StructType:
			// struct type definition: add to structTypesMap
			typeMapsVar.structTypesMap[typeDefInPkg.Type()] = typeDefAstIdent
		default:
			// type def that is not a struct: named/alias type, interface, etc. Record as named type.
			typeMapsVar.namedTypesMap[typeDefInPkg.Type()] = typeDefAstIdent
		}
	}
}

func addFunction(typeDefAstIdent *ast.Ident, typeDefInPkg types.Object, typeMapsVar *typeInfoMaps) {
	fnType, ok := typeDefInPkg.(*types.Func)
	if !ok {
		return
	}
	typeMapsVar.functionsMap[fnType] = typeDefAstIdent
}

func addVarConst(typeDefAstIdent *ast.Ident, typeDefInPkg types.Object, typeMapsVar *typeInfoMaps) {
	switch typeDefType := typeDefInPkg.(type) {
	case *types.Var, *types.Const:
		typeMapsVar.varConstsMap[typeDefType] = typeDefAstIdent
	}
}

type computeLogSafetyInfoForTypeReturnType struct {
	typeSafetyValue       *LogSafetyInfo
	typeFieldSafetyValues map[types.Object]LogSafetyInfo
}

// computeLogSafetyInfoForType computes and records the log safety information for a specific type.
// Currently, supported types are struct declarations, type definitions/aliases that resolve to structs, and named types
// that are annotated with comments.
func computeLogSafetyInfoForType(
	targetType types.Type,
	typeDefASTIdent *ast.Ident,
	typeDefObject types.Object,
	pass *analysis.Pass,
	currPgLogSafetyInfos, // read-only, and may not be fully populated
	dependentPkgLogSafetyInfos *logSafetyInfoBundle,
	typesToASTIdentsForCurrentPkg map[types.Type]*ast.Ident,
	typePathToCurrentType []types.Type,
	cachedReturnValues map[types.Type]*computeLogSafetyInfoForTypeReturnType,
	commentLogSafetyTracker *CommentBasedLogSafetyTracker,
	param Param,
) (typeSafetyValue *LogSafetyInfo, typeFieldSafetyValues map[types.Object]LogSafetyInfo, cycle bool) {

	// if type was already processed, return
	for _, typeInPath := range typePathToCurrentType {
		if typeOrUnderlyingTypeMatches(targetType, typeInPath) {
			// indicates that there was a cycle
			return nil, nil, true
		}
	}

	// if value was already computed in this pass, return it
	if val, ok := cachedReturnValues[targetType]; ok {
		return val.typeSafetyValue, val.typeFieldSafetyValues, false
	}

	// record return value if not part of a cycle
	defer func() {
		if !cycle {
			cachedReturnValues[targetType] = &computeLogSafetyInfoForTypeReturnType{
				typeSafetyValue:       typeSafetyValue,
				typeFieldSafetyValues: typeFieldSafetyValues,
			}
		}
	}()

	// add current Type to path
	typePathToCurrentType = append(typePathToCurrentType, targetType)

	// if AST for type is not a type declaration, nothing to do
	if typeDefASTIdent.Obj == nil || typeDefASTIdent.Obj.Kind != ast.Typ {
		return nil, nil, false
	}

	typeSpecDecl, ok := typeDefASTIdent.Obj.Decl.(*ast.TypeSpec)
	if !ok {
		return nil, nil, false
	}

	switch astType := typeSpecDecl.Type.(type) {
	// type is a struct declaration
	case *ast.StructType:
		// types.Object is easier to work with for reading tags, so get *types.Struct.
		// Need to resolve to "Underlying()" here because if the type def object is a "Named.Type", it is hard to extract
		// the *types.Struct object otherwise. However, the key for the type should still be "typeDefObject.Type()".
		underlyingStructType, isStruct := typeDefObject.Type().Underlying().(*types.Struct)
		if !isStruct {
			return nil, nil, false
		}

		// compute and return safety value for struct
		return computeStructLogSafetyInfo(
			typeSpecDecl.Name.Name,
			underlyingStructType,
			astType,
			pass,
			currPgLogSafetyInfos,
			dependentPkgLogSafetyInfos,
			typesToASTIdentsForCurrentPkg,
			typePathToCurrentType,
			cachedReturnValues,
			commentLogSafetyTracker,
			param,
		)
	default:
		// get log safety value specified in comment (if no comment, will be LogSafetyTypeUnmarked)
		logSafetyFromComment, commentPos := commentLogSafetyTracker.LogSafetyForIdentAtPos(astType.Pos())

		// type is *not* a struct declaration: it is most likely a name or alias type.
		// Compute the log safety based on the type (the computeTypeLogSafetyInfo resolves the typeDefObject.Type() into
		// its base types).
		logSafetyInfoFromType, cyclicFromType := computeTypeLogSafetyInfo(
			typeDefObject.Type(),
			pass,
			currPgLogSafetyInfos,
			dependentPkgLogSafetyInfos,
			typesToASTIdentsForCurrentPkg,
			typePathToCurrentType,
			cachedReturnValues,
			commentLogSafetyTracker,
			param,
		)

		// if logSafetyInfoFromType is nil, indicates a cycle: if there is no comment, return nil
		if logSafetyFromComment == LogSafetyTypeUnmarked && cyclicFromType {
			return nil, nil, true
		}

		var overallLogSafety LogSafetyType
		var reason string
		if logSafetyInfoFromType.LogSafety > logSafetyFromComment {
			overallLogSafety = logSafetyInfoFromType.LogSafety
			reason = logSafetyInfoFromType.Reason
		} else {
			overallLogSafety = logSafetyFromComment
			reason = commentOnElementMessage(typeName(pass.Pkg, typeDefASTIdent.Name), commentPos)
		}
		return &LogSafetyInfo{
			Identifier: typeName(pass.Pkg, typeSpecDecl.Name.Name),
			Reason:     reason,
			LogSafety:  overallLogSafety,
		}, nil, false
	}
}

func typeOrUnderlyingTypeMatches(t1, t2 types.Type) bool {
	t1Underlying := t1.Underlying()
	t2Underlying := t2.Underlying()
	return t1 == t2 || t1Underlying == t2 || t1 == t2Underlying || t1Underlying == t2Underlying
}

func computeStructLogSafetyInfo(
	structNameVal string,
	underlyingStructType *types.Struct,
	astStructType *ast.StructType,
	pass *analysis.Pass,
	currPkgLogSafetyInfos, // read-only, and may not be fully populated
	dependentPkgLogSafetyInfos *logSafetyInfoBundle,
	typesToASTIdentsForCurrentPkg map[types.Type]*ast.Ident,
	typePathToCurrentType []types.Type,
	cachedReturnTypes map[types.Type]*computeLogSafetyInfoForTypeReturnType,
	commentLogSafetyTracker *CommentBasedLogSafetyTracker,
	param Param,
) (typeSafetyValue *LogSafetyInfo, objectSafetyValues map[types.Object]LogSafetyInfo, cycleFieldsOnly bool) {

	// create a single array of *ast.Ident for the field names.
	// Required because "number of fields" for AST type vs. object type are computed differently.
	//
	// Consider struct:
	//    type MemProfileRecord struct {
	//	    AllocBytes, FreeBytes     int64
	//	    AllocObjects, FreeObjects int64
	//	    Stack                     []uintptr
	//    }
	//
	// For AST, this struct is considered to have 3 "fields", where field 1 and 2 each have 2 "names".
	// For object, this struct is considered to have 5 "fields".
	// Because the logic iterates over object fields, expand out the AST fields to create a single list.
	var astNameFields []*ast.Ident
	for _, field := range astStructType.Fields.List {
		// if field does not have any names, then it is likely an embedded type: the "Type" field should be an
		// *ast.Ident and should be added as such
		if len(field.Names) == 0 {
			fieldTypeExpr := field.Type

			// if field.Type is a wrapped expression, unwrap it until it cannot be unwrapped further
			for {
				unwrapped := unwrapASTExpr(fieldTypeExpr)
				// no more wrapping: continue
				if unwrapped == fieldTypeExpr {
					break
				}
				fieldTypeExpr = unwrapped
			}

			if fieldTypeAstIdent, ok := fieldTypeExpr.(*ast.Ident); ok {
				astNameFields = append(astNameFields, fieldTypeAstIdent)
			} else {
				panic(fmt.Sprintf("unexpected AST identifier in struct field type: %T", fieldTypeExpr))
			}
			continue
		}

		// if field has names, add them directly as identifiers
		for _, name := range field.Names {
			astNameFields = append(astNameFields, name)
		}
	}

	// start struct log safety as value from comment (if no comment, will start as LogSafetyTypeUnmarked).
	// Even if comment sets a value, the fields may make the struct "less safe", so still consider them.
	overallStructLogSafety, commentPos := commentLogSafetyTracker.LogSafetyForIdentAtPos(astStructType.Pos())
	overallStructSafetyReason := commentOnElementMessage(typeName(pass.Pkg, structNameVal), commentPos)

	objectSafetyValues = make(map[types.Object]LogSafetyInfo)
	numCycleFields := 0
	// Iterate over fields of struct
	for i := 0; i < underlyingStructType.NumFields(); i++ {
		// Determine log safety of field based on struct tag if one exists
		fieldLogSafetyFromStructTag := LogSafetyTypeUnmarked
		if value, ok := reflect.StructTag(underlyingStructType.Tag(i)).Lookup("safelogging"); ok {
			switch value {
			case "@Safe":
				fieldLogSafetyFromStructTag = LogSafetyTypeSafe
			case "@Unsafe":
				fieldLogSafetyFromStructTag = LogSafetyTypeUnsafe
			case "@DoNotLog":
				fieldLogSafetyFromStructTag = LogSafetyTypeDoNotLog
			}
		}

		// Set the log safety for the types.Object of the struct field
		fieldNameASTIdent := astNameFields[i]
		structFieldObj := pass.TypesInfo.Defs[fieldNameASTIdent]

		fieldLogSafetyFromList := LogSafetyTypeUnmarked
		fullyQualifiedName := typeName(pass.Pkg, structNameVal+"."+fieldNameASTIdent.Name)
		if value, ok := param.structFieldSafetyMap[fullyQualifiedName]; ok {
			fieldLogSafetyFromList = value
		}

		// compute field log safety based on type
		fieldLogSafetyInfoFromType, fieldCycleFromType := computeTypeLogSafetyInfo(
			structFieldObj.Type(),
			pass,
			currPkgLogSafetyInfos,
			dependentPkgLogSafetyInfos,
			typesToASTIdentsForCurrentPkg,
			typePathToCurrentType,
			cachedReturnTypes,
			commentLogSafetyTracker,
			param,
		)

		// if field is an inline struct definition, compute safety based on the struct definition
		var fieldLogSafetyInfoFromInlineStructDef LogSafetyInfo
		if typesStruct, ok := structFieldObj.Type().(*types.Struct); ok {
			// if type of field is *types.Struct, this means that the field is a literal struct definition
			// (anonymous struct defined as field): compute safety based on the struct directly (since the definition is
			// inline, it is not in the top-level types map)
			if astField, ok := fieldNameASTIdent.Obj.Decl.(*ast.Field); ok {
				// compute and return safety value for struct.
				// Include fields: since the struct definition is inline, this is the only chance to compute the field
				// safety values, and need to return the object safety of the fields so that check can flag if these objects
				// are referenced directly.
				computedLogSafetyInfo, computedObjectLogSafetyValues, computedCycleFieldsOnly := computeStructLogSafetyInfo(
					structNameVal+"."+fieldNameASTIdent.Name,
					typesStruct,
					astField.Type.(*ast.StructType),
					pass,
					currPkgLogSafetyInfos,
					dependentPkgLogSafetyInfos,
					typesToASTIdentsForCurrentPkg,
					typePathToCurrentType,
					cachedReturnTypes,
					commentLogSafetyTracker,
					param,
				)
				if computedLogSafetyInfo != nil && !computedCycleFieldsOnly {
					maps.Copy(objectSafetyValues, computedObjectLogSafetyValues)
					fieldLogSafetyInfoFromInlineStructDef = *computedLogSafetyInfo
				}
			}
		}

		// set overallFieldLogSafety to be most restrictive value between tag and type
		var overallFieldLogSafety LogSafetyType
		var reason string
		if fieldLogSafetyFromStructTag > overallFieldLogSafety {
			overallFieldLogSafety = fieldLogSafetyFromStructTag
			reason = fmt.Sprintf("tag on struct field %q", typeName(pass.Pkg, structNameVal+"."+fieldNameASTIdent.Name))
		}
		if fieldLogSafetyFromList > overallFieldLogSafety {
			overallFieldLogSafety = fieldLogSafetyFromList
			reason = fmt.Sprintf("configuration specified log safety value for struct field %q", typeName(pass.Pkg, structNameVal+"."+fieldNameASTIdent.Name))
		}
		if fieldLogSafetyInfoFromType.LogSafety > overallFieldLogSafety {
			overallFieldLogSafety = fieldLogSafetyInfoFromType.LogSafety
			reason = fieldLogSafetyInfoFromType.Reason
		}
		if fieldLogSafetyInfoFromInlineStructDef.LogSafety > overallFieldLogSafety {
			overallFieldLogSafety = fieldLogSafetyInfoFromInlineStructDef.LogSafety
			reason = fieldLogSafetyInfoFromInlineStructDef.Reason
		}

		// field was cyclic and did not contribute to log safety
		if overallFieldLogSafety == LogSafetyTypeUnmarked && fieldCycleFromType {
			numCycleFields++
		}

		if overallFieldLogSafety != LogSafetyTypeUnmarked {
			objectSafetyValues[structFieldObj] = LogSafetyInfo{
				Identifier: typeName(pass.Pkg, structNameVal+"."+fieldNameASTIdent.Name),
				Reason:     reason,
				LogSafety:  overallFieldLogSafety,
			}
		}

		// ensure that overallStructLogSafety is the "least safe" of its fields
		if overallFieldLogSafety > overallStructLogSafety {
			overallStructLogSafety = overallFieldLogSafety
			overallStructSafetyReason = reason
		}
	}

	// if there is at least one field and all are cyclic, return values that indicate that struct's status should not be
	// persisted (but may be revisited later)
	if numCycleFields != 0 && numCycleFields == underlyingStructType.NumFields() {
		return nil, objectSafetyValues, true
	}

	return &LogSafetyInfo{
		Identifier: typeName(pass.Pkg, structNameVal),
		Reason:     overallStructSafetyReason,
		LogSafety:  overallStructLogSafety,
	}, objectSafetyValues, false
}

// returns the computed LogSafetyInfo. If the return value is nil, indicates that there is a cycle, and the information
// from this pass should not be persisted/stored.
func computeTypeLogSafetyInfo(
	targetType types.Type,
	pass *analysis.Pass,
	currPgLogSafetyInfos, // read-only, and may not be fully populated
	dependentPkgLogSafetyInfos *logSafetyInfoBundle,
	typesToASTIdentsForCurrentPkg map[types.Type]*ast.Ident,
	typePathToCurrentType []types.Type,
	cachedReturnTypes map[types.Type]*computeLogSafetyInfoForTypeReturnType,
	commentLogSafetyTracker *CommentBasedLogSafetyTracker,
	param Param,
) (info LogSafetyInfo, cycle bool) {
	// Get all the base types for the type. In most cases, this is the declared type of the field. However, if the type
	// of the field is a container like a slice, array, or map, the component type(s) will be different from the
	// declared type.
	targetTypeAllBaseTypes := getAllBaseTypes(targetType)

	var logSafetyFromType LogSafetyInfo
	numCycleBaseTypes := 0
	for _, currTargetTypeBaseType := range targetTypeAllBaseTypes {
		var currTypeSafety LogSafetyInfo
		if gotLogSafetyValueForType, ok := dependentPkgLogSafetyInfos.typeLogSafety[currTargetTypeBaseType]; ok {
			// log safety value for type was recorded from dependent packages: use value
			currTypeSafety = gotLogSafetyValueForType
		} else if gotLogSafetyValueForType, ok := currPgLogSafetyInfos.typeLogSafety[currTargetTypeBaseType]; ok {
			// log safety value for type was recorded in current package: use value
			currTypeSafety = gotLogSafetyValueForType
		} else if astIdentInPkgForType, ok := typesToASTIdentsForCurrentPkg[currTargetTypeBaseType]; ok {
			// type does not have log safety recorded, but is defined in this package.
			// Compute the log safety for the type recursively and use result.
			// The intermediate computation results are not recorded, because they may depend on the completion of
			// type safety computation for the current type.
			computedSafetyInfo, _, computedCycle := computeLogSafetyInfoForType(
				currTargetTypeBaseType,
				astIdentInPkgForType,
				pass.TypesInfo.Defs[astIdentInPkgForType],
				pass,
				currPgLogSafetyInfos,
				dependentPkgLogSafetyInfos,
				typesToASTIdentsForCurrentPkg,
				typePathToCurrentType,
				cachedReturnTypes,
				commentLogSafetyTracker,
				param,
			)

			if computedSafetyInfo == nil {
				currTypeSafety = LogSafetyInfo{}
			} else {
				currTypeSafety = *computedSafetyInfo
			}

			if computedCycle {
				numCycleBaseTypes++
			}
		}

		if currTypeSafety.LogSafety > logSafetyFromType.LogSafety {
			logSafetyFromType = currTypeSafety
		}
	}

	// consider type a cycle only if all base types are cycles
	cycle = numCycleBaseTypes != 0 && numCycleBaseTypes == len(targetTypeAllBaseTypes)

	return logSafetyFromType, cycle
}

func commentOnElementMessage(elemName string, pos token.Position) string {
	return fmt.Sprintf("comment on %q at %v", elemName, pos)
}

func computeLogSafetyInfoASTIdentFromComment(
	astIdent *ast.Ident,
	pass *analysis.Pass,
	commentLogSafetyTracker *CommentBasedLogSafetyTracker,
) (typeSafetyValue *LogSafetyInfo) {
	logSafety, pos := commentLogSafetyTracker.LogSafetyForIdentAtPos(astIdent.Pos())
	astIdentTypeName := typeName(pass.Pkg, astIdent.Name)
	return &LogSafetyInfo{
		Identifier: astIdentTypeName,
		Reason:     commentOnElementMessage(astIdentTypeName, pos),
		LogSafety:  logSafety,
	}
}

const suppressWarningsRegexpString = "^" + safeLoggingCommentPrefix + "@Allow: .+$"

var suppressWarningsRegexp = regexp.MustCompile(suppressWarningsRegexpString)

func nodeHasSuppressWarningsComment(astNodePos token.Pos, fileCommentRetriever filecomments.Retriever) bool {
	targetComment, _ := fileCommentRetriever.CommentOnLineRelativeToPos(astNodePos, -1)
	return suppressWarningsRegexp.MatchString(targetComment)
}

func typeName(pkg *types.Package, identifierName string) string {
	return pkg.Path() + "." + identifierName
}

func unwrapASTExpr(expr ast.Expr) ast.Expr {
	// unwrap star
	// Example: "type Foo struct { *pkg.Val }"
	if starExpr, ok := expr.(*ast.StarExpr); ok {
		return starExpr.X
	}

	// unwrap selector
	// Example: "type Foo struct { _ pkg.Val }"
	if selectorExpr, ok := expr.(*ast.SelectorExpr); ok {
		return selectorExpr.Sel
	}

	// unwrap index list
	// Example: type Foo struct { children [nChildren]int }
	if indexListExpr, ok := expr.(*ast.IndexListExpr); ok {
		return indexListExpr.X
	}

	// unwrap index expression (generics)
	// Example: type LastSuccess[T any] struct { Atomic[T] }
	if indexExpr, ok := expr.(*ast.IndexExpr); ok {
		return indexExpr.X
	}

	// not wrapped: return expression itself
	return expr
}

func getAllBaseTypes(t types.Type) []types.Type {
	allBaseTypes := getBaseTypes(t, true)
	allBaseTypes = append(allBaseTypes, getBaseTypes(t, false)...)

	seen := make(map[types.Type]bool)
	var returnBaseTypes []types.Type
	for _, baseType := range allBaseTypes {
		if seen[baseType] {
			continue
		}
		seen[baseType] = true
		returnBaseTypes = append(returnBaseTypes, baseType)
	}
	return returnBaseTypes
}

// getBaseTypes returns all component types of the given types.Type. If the provided type has an underlying type,
// returns results for both the provided type and its underlying type.
//
// For example, for "type NamedKey string; type NamedValue int" and a type "map[NamedKey]NamedValue", calling this
// function with "convertToUnderlying=true" would return "(string, int)", while calling it with a value of "false" would
// return "(MyNamedType, MyNamedValue)".
//
// Conceptually, returns the concrete types of the values that could be logged by providing an instance of the given
// type. As such, for types like function signatures or channels, will return nil.
func getBaseTypes(t types.Type, getUnderlying bool) []types.Type {
	// convert type to underlying type
	if getUnderlying {
		t = t.Underlying()
	}

	switch typ := t.(type) {
	case *types.Named, *types.Basic, *types.Struct:
		return []types.Type{typ}
	case *types.Pointer:
		return getBaseTypes(typ.Elem(), getUnderlying)
	case *types.Array:
		return getBaseTypes(typ.Elem(), getUnderlying)
	case *types.Slice:
		return getBaseTypes(typ.Elem(), getUnderlying)
	case *types.Map:
		return append(getBaseTypes(typ.Key(), getUnderlying), getBaseTypes(typ.Elem(), getUnderlying)...)
	case *types.Signature, *types.Chan:
		// return nil because fields of these types do not contain any actual data
		return nil
	}
	return nil
}
