package httputils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestChiLogger(t *testing.T) {
	assert := assert.New(t)

	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	router := chi.NewRouter()
	router.Use(ChiLogger(logger))
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("test"))
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/test")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	chiRegex := regexp.MustCompile(`\[CHI\]`)
	assert.True(chiRegex.MatchString(logs.All()[0].Message))

	colorString := "\033\\[97;42m 200 \033\\[0m"
	colorRegex := regexp.MustCompile(colorString)
	assert.True(colorRegex.MatchString(logs.All()[0].Message))

	endpointRegex := regexp.MustCompile(`"/test"`)
	assert.True(endpointRegex.MatchString(logs.All()[0].Message))

	methodString := "\033\\[97;44m GET     \033\\[0m"
	methodRegex := regexp.MustCompile(methodString)
	assert.True(methodRegex.MatchString(logs.All()[0].Message))
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
			assert.Equal(t, tc.want, ReadUserIP(tc.req))
		})
	}
}

func TestFindMethodColor(t *testing.T) {
	assert := assert.New(t)

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
		assert.Equal(expected, findMethodColor(method))
	}
}

func TestFindStatusCodeColor(t *testing.T) {
	assert := assert.New(t)

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
		assert.Equal(tc.want, findStatusCodeColor(tc.code))
	}
}

func TestAddAllowHeader(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`ok`)) })
	r.Put("/", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`ok`)) })
	r.Delete("/", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`ok`)) })
	r.Patch("/", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`ok`)) })

	for _, tc := range []struct {
		name string
		req  *http.Request
	}{
		{
			"filled rawpath",
			func() *http.Request {
				req := httptest.NewRequest("POST", "/", nil)
				rctx := chi.NewRouteContext()
				rctx.RoutePath = ""
				req.URL.RawPath = "/"
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
				return req
			}(),
		},
		{
			"empty rawpath",
			func() *http.Request {
				req := httptest.NewRequest("POST", "/", nil)
				rctx := chi.NewRouteContext()
				rctx.RoutePath = ""
				req.URL.RawPath = ""
				req.URL.Path = ""
				rctx.RouteMethod = ""
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
				return req
			}(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			w := httptest.NewRecorder()
			AddAllowHeader(r, w, tc.req)

			assert.Equal(http.StatusMethodNotAllowed, w.Result().StatusCode)
			want := []string{"GET", "PUT", "DELETE", "PATCH"}
			assert.ElementsMatch(want, w.Result().Header.Values("Allow"))
		})
	}
}
