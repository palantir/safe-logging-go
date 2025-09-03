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

	"github.com/golangci/plugin-module-register/register"
	"github.com/palantir/safe-logging-go/safelogging"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("safelogging", New)
}

type Settings struct {
	// TypeLogSafety is a map from fully qualified type name identifier to the log safety for that type.
	// The safety value in this map can make a type less safe, but not more safe (for example, if a struct type is
	// determined to be unsafe based on its fields, marking it as safe using this configuration will not make it safe).
	//
	// If the value is nil, uses a set of preconfigured defaults. If the value is present but empty, uses none.
	TypeLogSafety *map[string]safelogging.LogSafetyType `json:"typeLogSafety,omitempty"`

	// StructFieldLogSafety is a map from fully qualified struct field identifier to the log safety for that field. The
	// type safety for a struct is the "least safe" of all of its types/fields (recursively) and any markings or safety
	// configured for the struct itself.
	//
	// If the value is nil, uses a set of preconfigured defaults. If the value is present but empty, uses none.
	StructFieldLogSafety *map[string]safelogging.LogSafetyType `json:"structFieldLogSafety,omitempty"`
}

type Plugin struct {
	configJSONBytes []byte
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[Settings](settings)
	if err != nil {
		return nil, err
	}
	configJSONBytes, err := json.Marshal(s)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal safelogging settings to JSON")
	}
	return &Plugin{
		configJSONBytes: configJSONBytes,
	}, nil
}

func (f *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	analyzer := safelogging.NewAnalyzer().Analyzer()
	if err := analyzer.Flags.Set(safelogging.JSONConfigFlagName, string(f.configJSONBytes)); err != nil {
		return nil, errors.Wrapf(err, "failed to set safelogging JSON config flag %q", safelogging.JSONConfigFlagName)
	}
	return []*analysis.Analyzer{
		analyzer,
	}, nil
}

func (f *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
