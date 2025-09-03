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
	"maps"
)

type Config struct {
	// TypeLogSafety is a map from fully qualified type name identifier to the log safety for that type.
	// The safety value in this map can make a type less safe, but not more safe (for example, if a struct type is
	// determined to be unsafe based on its fields, marking it as safe using this configuration will not make it safe).
	// The values in this map are applied on top of the default.
	TypeLogSafety *map[string]LogSafetyType `json:"typeLogSafety,omitempty" yaml:"type-log-safety,omitempty"`

	// If true, omits the default TypeLogSafety values and uses only those specified in the TypeLogSafety map.
	TypeLogSafetyOmitDefaults bool `json:"typeLogSafetyDisableDefaults,omitempty" yaml:"type-log-safety-disable-defaults,omitempty"`

	// StructFieldLogSafety is a map from fully qualified struct field identifier to the log safety for that field. The
	// type safety for a struct is the "least safe" of all of its types/fields (recursively) and any markings or safety
	// configured for the struct itself.
	StructFieldLogSafety *map[string]LogSafetyType `json:"structFieldLogSafety,omitempty" yaml:"struct-field-log-safety,omitempty"`

	// If true, omits the default StructFieldLogSafety values and uses only those specified in the StructFieldLogSafety map.
	StructFieldLogSafetyOmitDefaults bool `json:"structFieldLogSafetyDisableDefaults,omitempty" yaml:"struct-field-log-safety-disable-defaults,omitempty"`
}

func (c *Config) ToParam() Param {
	var p Param

	typeSafetyMap := make(map[string]LogSafetyType)
	if !c.TypeLogSafetyOmitDefaults {
		typeSafetyMap = builtinTypeSafetyMap()
	}
	if c.TypeLogSafety != nil {
		maps.Copy(typeSafetyMap, *c.TypeLogSafety)
	}
	p.typeSafetyMap = typeSafetyMap

	structFieldSafetyMap := make(map[string]LogSafetyType)
	if !c.StructFieldLogSafetyOmitDefaults {
		structFieldSafetyMap = builtinStructFieldSafetyMap()
	}
	if c.StructFieldLogSafety != nil {
		maps.Copy(structFieldSafetyMap, *c.StructFieldLogSafety)
	}
	p.structFieldSafetyMap = structFieldSafetyMap

	return p
}

type Param struct {
	// map from fully qualified type name identifier to log safety for that type
	typeSafetyMap map[string]LogSafetyType

	// map from fully qualified struct field identifier to log safety for that field
	structFieldSafetyMap map[string]LogSafetyType
}

func builtinTypeSafetyMap() map[string]LogSafetyType {
	return map[string]LogSafetyType{
		// http.Header can often contain sensitive information like authentication header values.
		"net/http.Header": LogSafetyTypeDoNotLog,
	}
}

func builtinStructFieldSafetyMap() map[string]LogSafetyType {
	return map[string]LogSafetyType{
		"github.com/aws/aws-sdk-go-v2/aws.Credentials.AccessKeyID":     LogSafetyTypeDoNotLog,
		"github.com/aws/aws-sdk-go-v2/aws.Credentials.SecretAccessKey": LogSafetyTypeDoNotLog,
		"github.com/aws/aws-sdk-go-v2/aws.Credentials.SessionToken":    LogSafetyTypeDoNotLog,

		"k8s.io/client-go/transport.Config.Password":    LogSafetyTypeDoNotLog,
		"k8s.io/client-go/transport.Config.BearerToken": LogSafetyTypeDoNotLog,
		"k8s.io/client-go/transport.TLSConfig.KeyData":  LogSafetyTypeDoNotLog,

		"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient.BasicAuth.Password":    LogSafetyTypeDoNotLog,
		"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient.ClientConfig.APIToken": LogSafetyTypeDoNotLog,
	}
}
