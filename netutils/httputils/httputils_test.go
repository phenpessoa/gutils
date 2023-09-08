package httputils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestChiLogger(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080", nil)
	ChiLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("test"))
	})).ServeHTTP(w, req)

	if logs.Len() != 1 {
		t.Errorf("\nlen of logs is not 1, it is %d", logs.Len())
		t.FailNow()
	}

	msg := logs.All()[0].Message
	if !strings.Contains(msg, "[CHI]") {
		t.Error("\n[CHI] not found in the logged message")
	}

	if !strings.Contains(msg, green+" 200 "+reset) {
		t.Error("\nstatus code not found in the logged message")
	}

	if !strings.Contains(msg, cyan+" POST    "+reset) {
		t.Error("\nmethod not found in the logged message")
	}
}

func TestReadUserIP(t *testing.T) {
	for _, tc := range []struct {
		name string
		req  *http.Request
		want string
	}{
		{
			"invalid ip addres",
			&http.Request{
				RemoteAddr: "invalid ip addres",
			},
			"",
		},
		{
			"valid ip addres",
			&http.Request{
				RemoteAddr: "127.0.0.1:8080",
			},
			"127.0.0.1",
		},
		{
			"fail to parse ip",
			&http.Request{
				RemoteAddr: "127.0:1",
			},
			"",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := ReadUserIP(tc.req); got != tc.want {
				t.Errorf("\ntest '%s' failed\nwant: %v\ngot: %v", tc.name, tc.want, got)
			}
		})
	}
}

func TestFindMethodColor(t *testing.T) {
	expectedResults := map[string]string{
		http.MethodConnect: reset,
		http.MethodDelete:  red,
		http.MethodGet:     blue,
		http.MethodHead:    magenta,
		http.MethodOptions: white,
		http.MethodPatch:   green,
		http.MethodPost:    cyan,
		http.MethodPut:     yellow,
		http.MethodTrace:   reset,
	}

	for method, expected := range expectedResults {
		if got := findMethodColor(method); got != expected {
			t.Errorf("\nmethod '%s' failed\nwant: %v\ngot: %v", method, expected, got)
		}
	}
}

func TestFindStatusCodeColor(t *testing.T) {
	for _, tc := range []struct {
		code int
		want string
	}{
		{http.StatusContinue, red},
		{http.StatusSwitchingProtocols, red},
		{http.StatusEarlyHints, red},
		{http.StatusOK, green},
		{http.StatusCreated, green},
		{http.StatusAccepted, green},
		{http.StatusNonAuthoritativeInfo, green},
		{http.StatusNoContent, green},
		{http.StatusResetContent, green},
		{http.StatusPartialContent, green},
		{http.StatusIMUsed, green},
		{http.StatusMultipleChoices, white},
		{http.StatusMovedPermanently, white},
		{http.StatusFound, white},
		{http.StatusSeeOther, white},
		{http.StatusNotModified, white},
		{http.StatusTemporaryRedirect, white},
		{http.StatusPermanentRedirect, white},
		{http.StatusBadRequest, yellow},
		{http.StatusUnauthorized, yellow},
		{http.StatusPaymentRequired, yellow},
		{http.StatusForbidden, yellow},
		{http.StatusNotFound, yellow},
		{http.StatusMethodNotAllowed, yellow},
		{http.StatusNotAcceptable, yellow},
		{http.StatusProxyAuthRequired, yellow},
		{http.StatusRequestTimeout, yellow},
		{http.StatusConflict, yellow},
		{http.StatusGone, yellow},
		{http.StatusLengthRequired, yellow},
		{http.StatusPreconditionFailed, yellow},
		{413, yellow}, // Payload too large
		{414, yellow}, // URI too log
		{http.StatusUnsupportedMediaType, yellow},
		{http.StatusRequestedRangeNotSatisfiable, yellow},
		{http.StatusExpectationFailed, yellow},
		{418, yellow}, // I'm a teapot
		{http.StatusUnprocessableEntity, yellow},
		{http.StatusTooEarly, yellow},
		{http.StatusUpgradeRequired, yellow},
		{http.StatusPreconditionRequired, yellow},
		{http.StatusTooManyRequests, yellow},
		{http.StatusRequestHeaderFieldsTooLarge, yellow},
		{http.StatusUnavailableForLegalReasons, yellow},
		{http.StatusInternalServerError, red},
		{http.StatusNotImplemented, red},
		{http.StatusBadGateway, red},
		{http.StatusServiceUnavailable, red},
		{http.StatusGatewayTimeout, red},
		{http.StatusHTTPVersionNotSupported, red},
		{http.StatusVariantAlsoNegotiates, red},
		{http.StatusInsufficientStorage, red},
		{http.StatusLoopDetected, red},
		{http.StatusNotExtended, red},
		{http.StatusNetworkAuthenticationRequired, red},
	} {
		if got := findStatusCodeColor(tc.code); tc.want != got {
			t.Errorf("\nstatus code '%d' failed\nwant: %v\ngot: %v", tc.code, tc.want, got)
		}
	}
}
