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
	"sort"
	"strings"
)

type LogSafetyType int

const (
	LogSafetyTypeUnmarked LogSafetyType = iota
	LogSafetyTypeSafe
	LogSafetyTypeUnsafe
	LogSafetyTypeDoNotLog

	safeLoggingCommentPrefix = "// safelogging:"
)

func toLogSafetyType(in string) LogSafetyType {
	for logSafetyTypeVal := LogSafetyTypeUnmarked; logSafetyTypeVal <= LogSafetyTypeDoNotLog; logSafetyTypeVal++ {
		if logSafetyTypeVal.String() == in {
			return logSafetyTypeVal
		}
	}
	return LogSafetyTypeUnmarked
}

func (t LogSafetyType) String() string {
	switch t {
	case LogSafetyTypeUnmarked:
		return "@Unmarked"
	case LogSafetyTypeSafe:
		return "@Safe"
	case LogSafetyTypeUnsafe:
		return "@Unsafe"
	case LogSafetyTypeDoNotLog:
		return "@DoNotLog"
	}
	return fmt.Sprintf("%d", t)
}

func (t LogSafetyType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *LogSafetyType) UnmarshalJSON(data []byte) error {
	var strValue string
	if err := json.Unmarshal(data, &strValue); err != nil {
		return err
	}
	*t = toLogSafetyType(strValue)
	return nil
}

type LogSafetyInfo struct {
	// fully qualified identifier for entity that has safety specified
	Identifier string

	// Reason for safety level
	Reason string

	LogSafety LogSafetyType
}

func (l *LogSafetyInfo) AFact() {}

func (l *LogSafetyInfo) String() string {
	return fmt.Sprintf("%s: %s (%s)", l.Identifier, l.LogSafety, l.Reason)
}

type PackageTypeLogSafetyInfo struct {
	// TypeRepToTypeToLogSafety is a 2-level nested map. The outer map is a map from the fmt.Sprintf("%T", types.Type)
	// to a map from the String() representation of a types.Type to its LogSafetyInfo.
	//
	// For the purposes of this check, the information that is most relevant is the log safety of a types.Type.
	// It would be most straightforward to represent this as a map from types.Type to LogSafetyInfo.
	//
	// However, facts need to be serializable, and types.Type is an interface that is not serializable. As a proxy for
	// this, the String() representation of the type is used as a key, which results in a map from string to
	// LogSafetyInfo. An analysis pass has a collection of types.Type values for the pass, so it is possible to match
	// the string representation back to a types.Type.
	//
	// This poses a different problem: computing the String() representation of a types.Type can be expensive.
	// Furthermore, the total number of types in a pass can be quite large (10k+), as it includes all types that can be
	// referenced by the package for the pass (including from all its dependencies). On the other hand, the number of
	// types that declare type safety is usually much smaller (typically <100). Because of this, converting all of the
	// types in a pass to their String() representation to match against the map keys can be quite inefficient.
	//
	// In order to reduce the search space, the types are bucketed based on the actual concrete type of the types.Type
	// (for example, *types.Struct, *types.Named, *types.Pointer, etc.). In practice, most types with type safety are
	// a *types.Named, but a pass can have many thousands of other types, so this bucketing can significantly reduce the
	// search space.
	TypeRepToTypeToLogSafety map[string]map[string]LogSafetyInfo
}

func (t *PackageTypeLogSafetyInfo) AFact() {}

func (t *PackageTypeLogSafetyInfo) String() string {
	builder := &strings.Builder{}
	writeMapVals(builder, "TypeRepToTypeToLogSafety", t.TypeRepToTypeToLogSafety, func(k string, v map[string]LogSafetyInfo) string {
		innerBuilder := &strings.Builder{}
		writeMapVals(innerBuilder, k, v, logSafetyStringer[string])
		return innerBuilder.String()
	})
	return builder.String()
}

type logSafetyInfoBundle struct {
	// map from types.Type to information about the log safety for that type.
	// Contains all primary types and their underlying types.
	typeLogSafety map[types.Type]LogSafetyInfo

	// set that contains all primary types. This set is a subset of the key set of typeLogSafety, and contains only the
	// types that are defined in the pass (and not their underlying types, unless they are also primary types).
	primaryTypes map[types.Type]struct{}

	// map from types.Object to information about its log safety
	objectLogSafety map[types.Object]LogSafetyInfo
}

func (l *logSafetyInfoBundle) addAllFromBundle(bundle *logSafetyInfoBundle) {
	maps.Copy(l.typeLogSafety, bundle.typeLogSafety)
	maps.Copy(l.primaryTypes, bundle.primaryTypes)
	maps.Copy(l.objectLogSafety, bundle.objectLogSafety)
}

func (l *logSafetyInfoBundle) String() string {
	builder := &strings.Builder{}

	writeMapVals(builder, "typeLogSafety", l.typeLogSafety, logSafetyStringer[types.Type])
	builder.WriteString(", ")

	builder.WriteString("primaryTypes: ")
	var primaryTypeEntries []string
	for t := range l.primaryTypes {
		primaryTypeEntries = append(primaryTypeEntries, t.String())
	}
	sort.Strings(primaryTypeEntries)
	builder.WriteString("[")
	builder.WriteString(strings.Join(primaryTypeEntries, ", "))
	builder.WriteString("], ")

	writeMapVals(builder, "objectLogSafety", l.objectLogSafety, logSafetyStringer[types.Object])

	return builder.String()
}

func logSafetyStringer[K any](_ K, v LogSafetyInfo) string {
	return fmt.Sprintf("%s: %s", v.Identifier, v.LogSafety)
}

func writeMapVals[K comparable, V any](builder *strings.Builder, mapName string, mapVal map[K]V, stringer func(K, V) string) {
	builder.WriteString(mapName)
	builder.WriteString(": [")
	var entries []string
	for k, v := range mapVal {
		entries = append(entries, stringer(k, v))
	}
	sort.Strings(entries)
	builder.WriteString(strings.Join(entries, ", "))
	builder.WriteString("]")
}
