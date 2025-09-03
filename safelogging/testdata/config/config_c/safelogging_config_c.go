package config_c

import (
	"net/http"

	"github.com/palantir/conjure-go-runtime/v2/conjure-go-client/httpclient"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

func ParamTests() {
	svc1log.SafeParam("testParam", http.Header{})
	svc1log.SafeParam("testParam", httpclient.BasicAuth{}.Password)
}
