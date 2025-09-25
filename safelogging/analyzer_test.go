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
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	patterns := []string{
		"safe-logging-go/safeloggingtests/a",
		"safe-logging-go/safeloggingtests/b",
		"safe-logging-go/safeloggingtests/c",
		"safe-logging-go/safeloggingtests/d",
		"safe-logging-go/safeloggingtests/e",
		"safe-logging-go/safeloggingtests/f",
		"safe-logging-go/safeloggingtests/g",
		"safe-logging-go/safeloggingtests/h",
		"safe-logging-go/safeloggingtests/h/innerh",
		"safe-logging-go/safeloggingtests/i",
		"safe-logging-go/safeloggingtests/j",
		"safe-logging-go/safeloggingtests/k",
	}
	for i := 1; i <= 16; i++ {
		patterns = append(patterns, fmt.Sprintf("safe-logging-go/safeloggingtests/l/l_%d", i))
	}

	analysistest.Run(t, analysistest.TestData(), NewAnalyzer().Analyzer(), patterns...)
}

func TestAnalyzerWithConfig(t *testing.T) {
	analyzer := NewAnalyzer()
	cfgBytes, err := json.Marshal(Config{
		StructFieldLogSafety: &map[string]LogSafetyType{
			"safe-logging-go/safeloggingtests/config/config_a.TestStructWithFieldSafetySetViaConfig.FieldToMarkUnsafe":   LogSafetyTypeUnsafe,
			"safe-logging-go/safeloggingtests/config/config_a.TestStructWithFieldSafetySetViaConfig.FieldToMarkDoNotLog": LogSafetyTypeDoNotLog,
		},
		ConstMessageLoggingFunctions: []ConstMessageLoggingFunction{
			{
				Function:          "safe-logging-go/safeloggingtests/config/config_d.CustomLoggingFunction",
				MessageParamIndex: 1,
			},
			{
				Function:          "safe-logging-go/safeloggingtests/config/config_d.OtherCustomLoggingFunction",
				MessageParamIndex: 1,
			},
			{
				Function:          "(safe-logging-go/safeloggingtests/config/config_d.LoggerStruct).StructReceiverLoggingFunction",
				MessageParamIndex: 1,
			},
			{
				Function:          "(*safe-logging-go/safeloggingtests/config/config_d.LoggerStruct).PointerReceiverLoggingFunction",
				MessageParamIndex: 1,
			},
			{
				Function:          "safe-logging-go/safeloggingtests/config/config_d.LoggingFunctionWithGenerics",
				MessageParamIndex: 1,
			},
		},
	})
	require.NoError(t, err)
	err = analyzer.Analyzer().Flags.Set("json-config", string(cfgBytes))
	require.NoError(t, err)
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer(),
		"safe-logging-go/safeloggingtests/config/config_a",
		"safe-logging-go/safeloggingtests/config/config_d",
	)
}

func TestAnalyzerWithDefaultConfig(t *testing.T) {
	const pkgPatterns = "safe-logging-go/safeloggingtests/config/config_b"

	analyzer := NewAnalyzer()
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer(), pkgPatterns)

	analyzer = NewAnalyzer()
	cfg := Config{
		TypeLogSafety:        nil,
		StructFieldLogSafety: nil,
	}
	cfgBytes, err := json.Marshal(cfg)
	require.NoError(t, err)
	err = analyzer.Analyzer().Flags.Set("json-config", string(cfgBytes))
	require.NoError(t, err)
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer(), pkgPatterns)
}

func TestAnalyzerWithEmptyConfigNoDefaults(t *testing.T) {
	analyzer := NewAnalyzer()
	cfg := Config{
		TypeLogSafety:                    &map[string]LogSafetyType{},
		TypeLogSafetyOmitDefaults:        true,
		StructFieldLogSafety:             &map[string]LogSafetyType{},
		StructFieldLogSafetyOmitDefaults: true,
	}
	cfgBytes, err := json.Marshal(cfg)
	require.NoError(t, err)
	err = analyzer.Analyzer().Flags.Set("json-config", string(cfgBytes))
	require.NoError(t, err)
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer(),
		"safe-logging-go/safeloggingtests/config/config_c",
	)
}
