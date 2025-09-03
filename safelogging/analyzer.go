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
	"encoding/json"
	"fmt"
	"go/types"
	"maps"
	"sync"
	"unicode"

	"github.com/palantir/safe-logging-go/safelogging/internal/filecomments"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"
)

const (
	AnalyzerName       = "safelogging"
	JSONConfigFlagName = "json-config"
)

// Analyzer is defined as its own struct to make it possible for an instance of analysis.Analyzer to track state across
// analysis runs but without using global variables/state. The "Run" function of analysis.Analyzer delegates to a method
// on Analyzer, which allows it to use state stored in the struct. Flag registration and config loading is also
// performed using member variables of the struct. This makes it easier to do things like running tests for analyzers
// using different flag configurations in the same run.
//
// Clients that just want an analyzer can call "safelogging.NewAnalyzer().Analyzer()" to get an instance of an
// *analysis.Analyzer.
type Analyzer struct {
	// flag for analyzer
	logSafetyConfigJSONBytesFlagVal string

	// variables related to loading parameter
	once     sync.Once
	param    Param
	paramErr error

	analyzer *analysis.Analyzer
}

func (a *Analyzer) Analyzer() *analysis.Analyzer {
	return a.analyzer
}

func (a *Analyzer) getParam() (Param, error) {
	a.once.Do(func() {
		var cfg Config
		if a.logSafetyConfigJSONBytesFlagVal != "" {
			if err := json.Unmarshal([]byte(a.logSafetyConfigJSONBytesFlagVal), &cfg); err != nil {
				a.paramErr = errors.Wrapf(err, "failed to unmarshal JSON configuration bytes from %q", a.logSafetyConfigJSONBytesFlagVal)
				return
			}
		}
		a.param = cfg.ToParam()
	})
	return a.param, a.paramErr
}

func NewAnalyzer() *Analyzer {
	analyzer := &Analyzer{}
	analyzer.analyzer = &analysis.Analyzer{
		Name: AnalyzerName,
		Doc:  "verifies that witchcraft logging and error calls log safely",
		FactTypes: []analysis.Fact{
			new(PackageTypeLogSafetyInfo),
			new(LogSafetyInfo),
		},
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return analyzer.doRun(pass)
		},
	}
	analyzer.analyzer.Flags.StringVar(&analyzer.logSafetyConfigJSONBytesFlagVal, JSONConfigFlagName, "", "JSON config for safe logging")

	return analyzer
}

func (a *Analyzer) doRun(pass *analysis.Pass) (interface{}, error) {
	param, err := a.getParam()
	if err != nil {
		return nil, err
	}

	allLogSafetyInfos := &logSafetyInfoBundle{
		typeLogSafety:   make(map[types.Type]LogSafetyInfo),
		primaryTypes:    make(map[types.Type]struct{}),
		objectLogSafety: make(map[types.Object]LogSafetyInfo),
	}

	if allPkgFacts := pass.AllPackageFacts(); len(allPkgFacts) > 0 {
		// record all type representations of types.Type (*types.Named, etc.) that have log safety info
		allTypeReps := make(map[string]struct{})
		for _, pkgFact := range allPkgFacts {
			pkgLogSafetyFact := pkgFact.Fact.(*PackageTypeLogSafetyInfo)

			for k := range pkgLogSafetyFact.TypeRepToTypeToLogSafety {
				allTypeReps[k] = struct{}{}
			}
		}

		// collect type string -> log safety info from all facts
		typeLogSafetyInfoFromFacts := make(map[string]LogSafetyInfo)
		for _, pkgFact := range allPkgFacts {
			pkgLogSafetyFact := pkgFact.Fact.(*PackageTypeLogSafetyInfo)
			for _, typeStrToTypeToLogSafetyInfoMap := range pkgLogSafetyFact.TypeRepToTypeToLogSafety {
				maps.Copy(typeLogSafetyInfoFromFacts, typeStrToTypeToLogSafetyInfoMap)
			}
		}

		// add all log safety infos for types in the pass to allLogSafetyInfos
		addAllLogSafetyInfos(typeLogSafetyInfoFromFacts, pass.TypesInfo, allTypeReps, allLogSafetyInfos)
	}

	for _, objFact := range pass.AllObjectFacts() {
		objLogSafetyInfo := objFact.Fact.(*LogSafetyInfo)
		allLogSafetyInfos.objectLogSafety[objFact.Object] = *objLogSafetyInfo
	}

	fileCommentRetriever := filecomments.NewRetriever(safeLoggingCommentPrefix, pass.Fset, pass.Files)
	commentLogSafetyTracker := NewCommentBasedLogSafetyTracker(pass.TypesInfo.Defs, pass.Fset, pass.Files)

	logSafetyInfoForPkg, allTypesMapWithUnderlyingTypes, err := computeLogSafetyInfo(pass, allLogSafetyInfos, commentLogSafetyTracker, param)
	if err != nil {
		return nil, err
	}
	// remove entries that have unmarked log level
	for k, v := range logSafetyInfoForPkg.typeLogSafety {
		if v.LogSafety == LogSafetyTypeUnmarked {
			delete(logSafetyInfoForPkg.typeLogSafety, k)
			delete(logSafetyInfoForPkg.primaryTypes, k)
		}
	}
	for k, v := range logSafetyInfoForPkg.objectLogSafety {
		if v.LogSafety == LogSafetyTypeUnmarked {
			delete(logSafetyInfoForPkg.objectLogSafety, k)
		}
	}

	pkgLogSafetyInfoFact := createPackageTypeLogSafetyInfo(logSafetyInfoForPkg.typeLogSafety, logSafetyInfoForPkg.primaryTypes)
	if len(pkgLogSafetyInfoFact.TypeRepToTypeToLogSafety) > 0 {
		pass.ExportPackageFact(pkgLogSafetyInfoFact)
	}
	for k, v := range logSafetyInfoForPkg.objectLogSafety {
		pass.ExportObjectFact(k, &v)
	}

	allLogSafetyInfos.addAllFromBundle(logSafetyInfoForPkg)

	findLogCalls(pass, allLogSafetyInfos, allTypesMapWithUnderlyingTypes, commentLogSafetyTracker, fileCommentRetriever, param)

	return nil, nil
}

func addAllLogSafetyInfos(
	typeLogSafetyInfoFromFacts map[string]LogSafetyInfo,
	typesInfo *types.Info,
	typeReps map[string]struct{},
	allLogSafetyInfos *logSafetyInfoBundle,
) {
	seenTypes := make(map[types.Type]struct{})
	for _, typeAndValue := range typesInfo.Types {
		if currType := typeAndValue.Type; currType != nil {
			if _, ok := seenTypes[currType]; ok {
				// type has already been encountered: skip
				continue
			}
			seenTypes[currType] = struct{}{}

			currTypeRep := fmt.Sprintf("%T", currType)
			if _, ok := typeReps[currTypeRep]; !ok {
				// skip type.Type types that do not have log safety info associated with them
				continue
			}

			currTypeStr := currType.String()
			logSafetyInfo, ok := typeLogSafetyInfoFromFacts[currTypeStr]
			if !ok {
				// not a type that contains log safety info: skip
				continue
			}
			// found type: remove from map
			delete(typeLogSafetyInfoFromFacts, currTypeStr)

			// add log safety info for type
			allLogSafetyInfos.typeLogSafety[currType] = logSafetyInfo
			// add log safety info for underlying types as well
			allLogSafetyInfos.typeLogSafety[currType.Underlying()] = logSafetyInfo

			// if all types have been found, break
			if len(typeLogSafetyInfoFromFacts) == 0 {
				break
			}
		}
	}
}

func createPackageTypeLogSafetyInfo(typeSafetyMap map[types.Type]LogSafetyInfo, primaryTypes map[types.Type]struct{}) *PackageTypeLogSafetyInfo {
	info := &PackageTypeLogSafetyInfo{
		TypeRepToTypeToLogSafety: make(map[string]map[string]LogSafetyInfo),
	}
	for t := range primaryTypes {
		isExported := false
		switch typ := t.(type) {
		case *types.Named:
			isExported = isExportedIdentifier(typ.Obj().Name())
		case *types.Alias:
			isExported = isExportedIdentifier(typ.Obj().Name())
		}
		if !isExported {
			// skip non-exported types
			continue
		}
		typeRepStr := fmt.Sprintf("%T", t)
		currMap, ok := info.TypeRepToTypeToLogSafety[typeRepStr]
		if !ok {
			currMap = make(map[string]LogSafetyInfo)
			info.TypeRepToTypeToLogSafety[typeRepStr] = currMap
		}
		currMap[t.String()] = typeSafetyMap[t]
	}
	return info
}

func isExportedIdentifier(in string) bool {
	return len(in) > 0 && unicode.IsUpper(rune(in[0])) && unicode.IsLetter(rune(in[0]))
}
