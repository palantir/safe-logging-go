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
	TypeLogSafety *map[string]LogSafetyType `json:"type-log-safety,omitempty" mapstructure:"type-log-safety,omitempty"`

	// If true, omits the default TypeLogSafety values and uses only those specified in the TypeLogSafety map.
	TypeLogSafetyOmitDefaults bool `json:"type-log-safety-disable-defaults,omitempty" mapstructure:"type-log-safety-disable-defaults,omitempty"`

	// StructFieldLogSafety is a map from fully qualified struct field identifier to the log safety for that field. The
	// type safety for a struct is the "least safe" of all of its types/fields (recursively) and any markings or safety
	// configured for the struct itself.
	StructFieldLogSafety *map[string]LogSafetyType `json:"struct-field-log-safety,omitempty" mapstructure:"struct-field-log-safety,omitempty"`

	// If true, omits the default StructFieldLogSafety values and uses only those specified in the StructFieldLogSafety map.
	StructFieldLogSafetyOmitDefaults bool `json:"struct-field-log-safety-disable-defaults,omitempty" mapstructure:"struct-field-log-safety-disable-defaults,omitempty"`

	// ConstMessageLoggingFunctions is a list of functions are checked to ensure that the parameter at a specified index
	// is a constant string. Currently, the check only supports checking one parameter per function -- if the provided
	// slice contains the same function multiple times, the last entry will take precedence. This configuration can add
	// to the default set of functions, but cannot override them.
	ConstMessageLoggingFunctions []ConstMessageLoggingFunction `json:"const-message-logging-functions,omitempty" mapstructure:"const-message-logging-functions,omitempty"`
}

type ConstMessageLoggingFunction struct {
	Function          FuncRef `json:"function" mapstructure:"function"`
	MessageParamIndex int     `json:"message-param-index" mapstructure:"message-param-index"`
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

	constMessageLoggingFunctions := make(map[FuncRef]int)
	for _, fn := range c.ConstMessageLoggingFunctions {
		constMessageLoggingFunctions[fn.Function] = fn.MessageParamIndex
	}
	p.constMessageLoggingFunctions = constMessageLoggingFunctions

	return p
}

type Param struct {
	// map from fully qualified type name identifier to log safety for that type
	typeSafetyMap map[string]LogSafetyType

	// map from fully qualified struct field identifier to log safety for that field
	structFieldSafetyMap map[string]LogSafetyType

	// map from function reference to the index of the parameter that should be a constant message string. Any keys
	// that match builtin configuration will not have any effect (the builtin configuration can be added to, but not
	// overridden).
	constMessageLoggingFunctions map[FuncRef]int
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
