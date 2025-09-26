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
	"testing"

	"github.com/go-viper/mapstructure/v2"
	"github.com/palantir/safe-logging-go/safelogging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Verifies that plugin can be constructed from mapstructure-decoded config.
// Effectively ensures that struct tags for "mapstructure" and "json" are consistent.
func TestPluginConfigSerDe(t *testing.T) {
	cfg := safelogging.Config{
		TypeLogSafety: &map[string]safelogging.LogSafetyType{
			"custom/type": safelogging.LogSafetyTypeSafe,
		},
		TypeLogSafetyOmitDefaults: true,
		StructFieldLogSafety: &map[string]safelogging.LogSafetyType{
			"custom/type.Field": safelogging.LogSafetyTypeDoNotLog,
		},
		StructFieldLogSafetyOmitDefaults: true,
		ConstMessageLoggingFunctions: []safelogging.ConstMessageLoggingFunction{
			{
				Function:          "fmt.Printf",
				MessageParamIndex: 0,
			},
		},
	}

	// translate from struct to mapstructure-compatible map (this is what golangci-lint does)
	var output map[string]any
	err := mapstructure.Decode(cfg, &output)
	require.NoError(t, err)

	// creating linter from mapstructure-decoded output should work
	_, err = New(output)
	assert.NoError(t, err)
}
