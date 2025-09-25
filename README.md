<p align="right">
<a href="https://autorelease.general.dmz.palantir.tech/palantir/safe-logging-go"><img src="https://img.shields.io/badge/Perform%20an-Autorelease-success.svg" alt="Autorelease"></a>
</p>

# safe-logging-go
`safe-logging-go` defines the `safelogging` analyzer.

`safelogging` is a [Go Analyzer](https://pkg.go.dev/golang.org/x/tools/go/analysis) that provides an `analysis.Analyzer`
that flags instances where unsafe content may be logged using witchcraft loggers. Conceptually, it is analogous to
aspects of the [safe-logging](https://github.com/palantir/safe-logging) and [gradle-baseline](https://github.com/palantir/gradle-baseline)
Java projects/checks.

At a high level, `safelogging` checks for 2 different categories:
* Checks that the message provided to a logger is a *compile-time constant*
* Checks that the values provided to typed parameters (`SafeParam` and `UnsafeParam`) do not violate safety constraints:
  `SafeParam` cannot be provided with a value that is `Unsafe` or `DoNotLog`, and `UnsafeParam` cannot be provided with
  a value that is `DoNotLog`

# Usage

## Standalone
This project produces a binary that uses `golang.org/x/tools/go/analysis/singlechecker` to provide a single-analysis
checker program. The binary can be invoked directly. Packages to be checked should be provided as arguments. Run the
binary with the `-h` flag for more information.

## As an `analysis.Analyzer` in Go code
This project provides an implementation of `analysis.Analyzer` that can be used in other projects that operate on an
`analysis.Analyzer`. To obtain an instance of an analyzer, import "github.com/palantir/safe-logging-go/safelogging"
and call `safelogging.NewAnalyzer().Analyzer()`.

## As a `golangci-lint` linter
This project provides the package `github.com/palantir/safe-logging-go/golangcilint/safelogging`, which
provides a `golangci-lint` linter module.

The module can be specified as a plugin for a `golangci-lint` build per the [golangci-lint documentation](https://golangci-lint.run/docs/plugins/module-plugins/#configuration-example):

```
- module: 'github.com/palantir/safe-logging-go'
  import: 'github.com/palantir/safe-logging-go/golangcilint/safelogging'
```

## Configuration
This check can be configured to treat specific struct fields and types as having a specified log safety type. This is
useful in instances where there are types or struct fields that are known to be "Unsafe" or "DoNotLog", but the author
does not have the ability to modify the code to add the struct tags or comment-based annotations necessary to mark them
directly.

The check can also be configured to verify that the argument passed to a function at a particular index is a
compile-time constant by specifying the function and parameter index in the "ConstMessageLoggingFunctions"
configuration.

When run as a standalone program, the configuration is specified as JSON and provided using a flag.

The configuration is defined in the [safelogging/config.go] file as follows:

```
type Config struct {
	// TypeLogSafety is a map from fully qualified type name identifier to the log safety for that type.
	// The safety value in this map can make a type less safe, but not more safe (for example, if a struct type is
	// determined to be unsafe based on its fields, marking it as safe using this configuration will not make it safe).
	// The values in this map are applied on top of the default.
	TypeLogSafety *map[string]LogSafetyType `json:"typeLogSafety,omitempty" mapstructure:"type-log-safety,omitempty"`

	// If true, omits the default TypeLogSafety values and uses only those specified in the TypeLogSafety map.
	TypeLogSafetyOmitDefaults bool `json:"typeLogSafetyDisableDefaults,omitempty" mapstructure:"type-log-safety-disable-defaults,omitempty"`

	// StructFieldLogSafety is a map from fully qualified struct field identifier to the log safety for that field. The
	// type safety for a struct is the "least safe" of all of its types/fields (recursively) and any markings or safety
	// configured for the struct itself.
	StructFieldLogSafety *map[string]LogSafetyType `json:"structFieldLogSafety,omitempty" mapstructure:"struct-field-log-safety,omitempty"`

	// If true, omits the default StructFieldLogSafety values and uses only those specified in the StructFieldLogSafety map.
	StructFieldLogSafetyOmitDefaults bool `json:"structFieldLogSafetyDisableDefaults,omitempty" mapstructure:"struct-field-log-safety-disable-defaults,omitempty"`

	// ConstMessageLoggingFunctions is a list of functions are checked to ensure that the parameter at a specified index
	// is a constant string. Currently, the check only supports checking one parameter per function -- if the provided
	// slice contains the same function multiple times, the last entry will take precedence. This configuration can add
	// to the default set of functions, but cannot override them.
	ConstMessageLoggingFunctions []ConstMessageLoggingFunction `json:"constMessageLoggingFunctions,omitempty" mapstructure:"const-message-logging-functions,omitempty"`
}
```

If configuration is not provided, then the following defaults are used for type safety:

```
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
```

The first parameter (the one at index 0) for the following functions are always checked to ensure that they are
compile-time constants (configuration can be used to check parameters for additional functions, but the check for these
functions cannot be overridden):

```
func (github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Logger).Debug(msg string, params ...github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Param)
func (github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Logger).Info(msg string, params ...github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Param)
func (github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Logger).Warn(msg string, params ...github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Param)
func (github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Logger).Error(msg string, params ...github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Param)
```

## Manually suppressing check errors
There are instances in which a usage flagged by the check is deemed to be acceptable. For example, consider the
following:

```
func doWork(ctx context.Context, input string) {
    var msg string
    switch input {
    case "a":
        msg = "Message one"
    case "b":
        msg = "Message two"
    default:
        msg = "Default"
    }
    svc1log.FromContext(ctx).Info(msg)
}
```

In the code above, `msg` is *effectively* a constant, but because it is a non-compile-time constant used as a logger
message, the check will flag it.

One way to fix this is to restructure the code so that only constants are passed. However, another option is to use
comment-based suppression to signal to the check that the violation should not cause a failure.

This can be done by adding a comment of the form `// safelogging:@Allow: [reason]` to the line before the failure.
For example:

```
func doWork(ctx context.Context, input string) {
    var msg string
    switch input {
    case "a":
        msg = "Message one"
    case "b":
        msg = "Message two"
    default:
        msg = "Default"
    }
    // safelogging:@Allow: content of msg is known to be a compile-time constant
    svc1log.FromContext(ctx).Info(msg)
}
```

# Design

## Compile-time constant message check
It is considered best practice for log messages to use compile-time constants. This makes log messages easier to search
and also ensures that unsafe content is not included in a logger message. In instances where a logger message does
contain a variable or runtime value, it is almost always preferable to have a fixed message and to log the variable
portion as a parameter.

The following functions are considered logging functions:
* `github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Logger.Debug(msg string, params ...github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Param)`
* `github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Logger.Info(msg string, params ...github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Param)`
* `github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Logger.Warn(msg string, params ...github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Param)`
* `github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Logger.Error(msg string, params ...github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Param)`

The check verifies that the first parameter (the one at index 0) is a value that is known at compile-time. Generally,
this means that it is  either a string literal (`logger.Info("Message")`), a reference to a string constant
(`const msg = "Message"; logger.Info(msg)`), or a concatenation of either
(`const msg = "Message"; logger.Info(msg + " content")`).

The check can be configured to verify this property for additional functions (the user must specify the identifier of
the function and the index of the parameter that should be checked as part of the analyzer configuration).

The analyzer only checks the above -- in particular, it does not perform code analysis to determine usage patterns that
semantically result in a constant output. For example, code such as `func msg() string { return "Message" }; logger.Info(msg())`,
`logger.Info(fmt.Sprintf("Number %d", 7))`, `msg := "Message"; logger.Info(msg)` are all written such that the logger
message is a de facto constant, but all of these usages would be flagged by the check because the message isn't
guaranteed to be a constant.

## Param safety check
Witchcraft logging allows the construction of "Safe" and "Unsafe" parameters that can be provided to logging and error
functions. These parameters take a name and an object, where the object can be anything (`any`/`interface{}`).

The param safety check verifies that, when a Safe or Unsafe parameter is constructed, the log safety level of the value
that is provided to it is compatible with the log safety level of the parameter.

The check operates on the following function calls:
* `svc1log.SafeParam` (`github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.SafeParam(key string, value interface{}) github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Param`)
  * `value` cannot be `Unsafe` or `DoNotLog` 
* `svc1log.SafeParams` (`github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.SafeParams(safe map[string]interface{}) github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Param`)
  * The values in the `safe` map cannot be `Unsafe` or `DoNotLog`
* `svc1log.UnsafeParam` (`github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.UnsafeParam(key string, value interface{}) github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Param`)
  * `value` cannot be `DoNotLog`
* `svc1log.UnsafeParams` (`github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.UnsafeParams(unsafe map[string]interface{}) github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log.Param`)
  * The values in the `unsafe` map cannot be `DoNotLog` 
* `werror.SafeParam` (`github.com/palantir/witchcraft-go-error.SafeParam(key string, val interface{}) github.com/palantir/witchcraft-go-error.Param`)
  * `val` cannot be `Unsafe` or `DoNotLog`
* `werror.SafeParams` (`github.com/palantir/witchcraft-go-error.SafeParams(vals map[string]interface{}) github.com/palantir/witchcraft-go-error.Param`)
  * The values in the `vals` map cannot be `Unsafe` or `DoNotLog`
* `werror.UnsafeParam` (`github.com/palantir/witchcraft-go-error.UnsafeParam(key string, val interface{}) github.com/palantir/witchcraft-go-error.Param`)
  * `val` cannot be `DoNotLog`
* `werror.UnsafeParams` (`github.com/palantir/witchcraft-go-error.(vals map[string]interface{}) github.com/palantir/witchcraft-go-error.Param`)
  * The values in the `vals` map cannot be `DoNotLog`
* `werror.SafeAndUnsafeParams` (`github.com/palantir/witchcraft-go-error.SafeAndUnsafeParams(safe map[string]interface{}, unsafe map[string]interface{}) github.com/palantir/witchcraft-go-error.Param`)
  * The values in the `safe` map cannot be `Unsafe` or `DoNotLog`
  * The values in the `unsafe` map cannot be `DoNotLog`

Semantically, any object can be considered to have a "safety value" that determines the safety level at which it is
acceptable to log the object. The safety levels are as follows:
* Unmarked/uncategorized
* Safe
* Unsafe
* DoNotLog

An object that is "Safe" can be provided to any parameter. An object that is "Unsafe" can only be provided to an "Unsafe"
parameter. An object that is "DoNotLog" cannot be provided to any parameter. In the current implementation of this check,
unmarked/uncategorized objects are treated as "Safe".

At a high level, this check does 2 things:
1. Determines the log safety of types and identifiers for every package based on a set of rules (struct tags, 
   comment-based annotations, and configurations)
2. Finds all invocations of the function calls to check and, using the information computed in (1), verifies that the
   log safety level of the parameter is compatible with the call

### Param value safety categorization
The safety of a parameter value is determined based on multiple factors. The primary ones are the *type* and *identifier*.

#### Specifying log safety value
The log safety value for a type or identifier can be specified in one of 3 ways:
1. Using struct tags for struct fields. The tag name is `safelogging`, and is specified as `safelogging:"{LogSafetyLevel}"` --
   for example, `safelogging:"@Safe"`, `safelogging:"@Unsafe"`, or `safelogging:"@DoNotLog"`.
2. Adding a comment to the same line as a type definition, variable or constant declaration, or function. The comment
   must be of the form `// safelogging:{SafetyLevel}` -- for example, `// safelogging:@Safe`, `// safelogging:@Unsafe`,
   or `// safelogging:@DoNotLog`.
3. Using configuration (the `safelogging.Config` struct) that is provided to the check as JSON using the `-json-config` flag

If the value provided as a parameter has log safety specified in multiple ways, it is considered to have the safety of
the "least safe" input value. The following sections outlines the specifics for how log safety can be specified.

#### Identifiers
The following identifiers can have log safety values specified:

* Variable and constant declarations
* Function definitions (standalone, interface functions, and functions defined on receivers)
* Struct fields

##### Variable and constant declarations
The log safety value for variable and constant declarations can be specified using comment-based marking. For example:

```
var PasswordVar string // safelogging:@DoNotLog
const PasswordConst = "password" // safelogging:@DoNotLog

func foo() {
    localPasswordVar := "password" // safelogging:@DoNotLog
}
```

The type safety for the identifier applies to any direct references to the identifiers. However, only *direct*
references are flagged -- the check does not track the value across assignments. For example:

```
// check will flag the following direct references
svc1log.SafeParam("testParam", PasswordVar)
svc1log.SafeParam("testParam", testpkg.PasswordConst)

// check will not flag the following
foo := PasswordVar
svc1log.SafeParam("testParam", foo)
```

##### Function definitions
The log safety value for a function can be specified using comment-based marking. Marking a function applies the log
safety value to its return value(s). Because parameter values are typically supplied individually, in practice, it only
makes sense to use function-based marking for functions that return a single value (or two values that are
`map[string]interface{}` if being used as an argument to the `SafeAndUnsafeParams` function). For example:

```
func (t TestStruct) Password() string { // safelogging:@DoNotLog
	return ""
}

type TestInterface interface {
	Password() string // safelogging:@DoNotLog
}

func Password() string { // safelogging:@DoNotLog
	return ""
}
```

The check will flag instances in which the result of the function call is provided directly to a parameter if the safety
values do not match. However, only *direct* invocations are flagged -- the check does not track the value across
assignments. For example:

```
// check will flag the following direct references
svc1log.SafeParam("testParam", Password())

// check will not flag the following
password := Password()
svc1log.SafeParam("testParam", password)

passwordFn := Password
svc1log.SafeParam("testParam", passwordFn())
```

##### Struct fields
The log safety value for a field of a struct is determined by its type and optionally by a struct tag on the field. The
safety value is always the "least safe" of all the inputs.

The log safety for a field of a struct can be specified using struct tags. The name of the struct tag is "safelogging"
and the value is one of "@Safe", "@Unsafe", or "@DoNotLog" (the "@" nomenclature comes from mirroring the Java
annotation). For example, the following is the definition for a struct with 4 fields, 3 of which have safety levels
specified using struct tags:

```
type TestStruct struct {
	UnmarkedField string
	SafeField     string `safelogging:"@Safe"`
	UnsafeField   string `safelogging:"@Unsafe"`
	DoNotLogField string `safelogging:"@DoNotLog"`
}
```

If a safety level can be determined for a struct field, any reference to it uses the safety level. For example:

```
// check will flag the following direct references
svc1log.SafeParam("testParam", TestStruct{}.UnsafeField)
```

#### Types

##### Struct
The log safety level for a struct type is determined by the safety level of its fields and optionally by a comment-based
annotation specified on the struct definition. The log safety level of the overall struct is the "least safe" safety
value of all the inputs. For example:

```
// Log safety level is "Unsafe" due to comment on definition
type TestStructOne struct { // safelogging:@Unsafe
    Name string
}

// Log safety level is "DoNotLog" because a field is at level "DoNotLog" due to struct tag
type TestStructTwo struct {
    Password string `safelogging:"@DoNotLog"`
}

// Log safety level is "DoNotLog" because a field is at level "DoNotLog" due to its type
type TestStructThree struct {
    PasswordVal Password
}

type Password string // safelogging:@DoNotLog

// Log safety level is "DoNotLog" because a field is at level "DoNotLog" (even though struct is annotated as "Unsafe",
// "DoNotLog" is the "least safe" of the inputs, so that is the level of the overall struct)  
type TestStructFour struct { // safelogging:@Unsafe
    innerStruct TestStructThree
}
```

##### Named types and aliases
The log safety level for named types and type aliases that are not structs are determined based on comment-based
annotations of the type definition and the safety of any of the underlying types of the named type. The safety level of
the named type is the "least safe" safety value of all the inputs. For example:

```
type TestUnsafeStruct struct {} // safelogging:@Unsafe

// Log safety level is "Unsafe" due to underlying type being "Unsafe"
type NamedStructType TestUnsafeStruct
type NamedStructAlias = TestUnsafeStruct

// Log safety level is "Unsafe" due to comment on definition
type NamedStringType string // safelogging:@Unsafe
type NamedStringAlias = string // safelogging:@Unsafe
```

##### Pointers and containers
If a type has a log safety level defined, then pointers to that type and containers that contain that type consider the
log safety level of the type as input, and the log safety level of the overall type is the least safe level of all the
inputs. For example:

```
type UnsafeType string // safelogging:@Unsafe

// All the following are "Unsafe" based on type
var (
    pointerVar  *UnsafeType
    sliceVar    []UnsafeType
    mapVar      map[string]UnsafeType
    mapVar2     map[UnsafeType]string
    compoundVar map[string][]*UnsafeType
)
```
