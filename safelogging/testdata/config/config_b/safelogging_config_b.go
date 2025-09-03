package config_b

import (
	"net/http"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

func ParamTests() {
	svc1log.SafeParam("testParam", http.Header{})                   // want `^SafeParam called with unsafe argument: argument is of type "net/http\.Header", which is @DoNotLog \(reason: configuration specified log safety value for type "net/http\.Header"\)$`
	svc1log.SafeParam("testParam", httpclient.BasicAuth{}.Password) // want `^SafeParam called with unsafe argument: argument is @DoNotLog \(reason: configuration specified log safety value for struct field "github\.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient\.BasicAuth\.Password"\)$`
}
