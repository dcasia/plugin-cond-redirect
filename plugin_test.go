package plugin_cond_redirect_test

import (
	"context"
	. "github.com/dcasia/plugin_cond_redirect"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {
	testHeaderRedirection(t)
	testCookieRedirection(t)
	testComplexRedirection(t)
}

func testHeaderRedirection(t *testing.T) {
	var cfg *Config
	err := mapstructure.Decode(map[string]interface{}{
		"rules": []map[string]interface{}{{
			"sourcePattern":      "(.*)",
			"destinationPattern": "$1/foo",
			"condition": map[string]interface{}{
				"type":    "header",
				"name":    "X-Foo",
				"pattern": "foo",
			},
		},
		},
	}, &cfg)
	if err != nil {
		t.Error(err)
	}

	handler, ctx := prepare(t, cfg)

	req, recorder := prepareCase(t, ctx, "http://localhost/bar")

	// should not redirect
	handler.ServeHTTP(recorder, req)
	assertNoRedirection(t, recorder)

	// should redirect
	recorder = httptest.NewRecorder()
	req.Header.Set("X-Foo", "foo")
	handler.ServeHTTP(recorder, req)
	assertRedirection(t, recorder, "/bar/foo")
}

func testCookieRedirection(t *testing.T) {
	var cfg *Config
	err := mapstructure.Decode(map[string]interface{}{
		"rules": []map[string]interface{}{{
			"sourcePattern":      ".*",
			"destinationPattern": "/foo",
			"condition": map[string]interface{}{
				"type":    "cookie",
				"name":    "foo",
				"pattern": ".*",
			},
		},
		},
	}, &cfg)
	if err != nil {
		t.Error(err)
	}

	handler, ctx := prepare(t, cfg)

	req, recorder := prepareCase(t, ctx, "http://localhost")

	// should not redirect
	handler.ServeHTTP(recorder, req)
	assertNoRedirection(t, recorder)

	// should redirect
	recorder = httptest.NewRecorder()
	req.AddCookie(&http.Cookie{
		Name:  "foo",
		Value: "foooo",
	})
	handler.ServeHTTP(recorder, req)
	assertRedirection(t, recorder, "/foo")
}

func testComplexRedirection(t *testing.T) {
	var cfg *Config
	err := mapstructure.Decode(map[string]interface{}{
		"rules": []map[string]interface{}{{
			"sourcePattern":      ".*",
			"destinationPattern": "/foo",
			"condition": map[string]interface{}{
				"type": "and",
				"children": []map[string]interface{}{
					{
						"type":    "cookie",
						"name":    "foo",
						"pattern": ".*",
					},
					{
						"type": "not",
						"condition": map[string]interface{}{
							"type":    "header",
							"name":    "Referer",
							"pattern": "http://foo.bar",
						},
					},
				},
			},
		},
		},
	}, &cfg)
	if err != nil {
		t.Error(err)
	}

	handler, ctx := prepare(t, cfg)

	// should not redirect
	req, recorder := prepareCase(t, ctx, "http://localhost")
	handler.ServeHTTP(recorder, req)
	assertNoRedirection(t, recorder)

	// should redirect
	req, recorder = prepareCase(t, ctx, "http://localhost")
	req.AddCookie(&http.Cookie{
		Name: "foo", Value: "foo",
	})
	handler.ServeHTTP(recorder, req)
	assertRedirection(t, recorder, "/foo")

	// should not redirect
	req, recorder = prepareCase(t, ctx, "http://localhost")
	req.AddCookie(&http.Cookie{
		Name: "foo", Value: "foo",
	})
	req.Header.Set("Referer", "http://foo.bar")
	handler.ServeHTTP(recorder, req)
	assertNoRedirection(t, recorder)
}

func prepare(t *testing.T, cfg *Config) (http.Handler, context.Context) {
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg)
	if err != nil {
		t.Fatal(err)
	}

	return handler, ctx
}

func prepareCase(t *testing.T, ctx context.Context, url string) (*http.Request, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	return req, recorder
}

func assertRedirection(t *testing.T, recorder *httptest.ResponseRecorder, location string) {
	assertStatusCode(t, recorder, 302)
	assertHeader(t, recorder, "Location", location)
}

func assertNoRedirection(t *testing.T, recorder *httptest.ResponseRecorder) {
	assertStatusCode(t, recorder, 200)
	assertHeader(t, recorder, "Location", "")
}

func assertStatusCode(t *testing.T, recorder *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if recorder.Code != expected {
		t.Errorf("Wrong status code. Expected: %d. Actual: %d.", expected, recorder.Code)
	}
}

func assertHeader(t *testing.T, recorder *httptest.ResponseRecorder, key, expected string) {
	t.Helper()

	actual := recorder.Header().Get(key)
	if actual != expected {
		t.Errorf("Wrong header. Expected: %s. Actual: %s", expected, actual)
	}
}
