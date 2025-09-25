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

type Plugin struct {
	configJSONBytes []byte
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[safelogging.Config](settings)
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
