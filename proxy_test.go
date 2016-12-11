package main

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDirector(t *testing.T) {
	testCases := []struct {
		upstream string
		req      string
		host     string
		path     string
	}{
		{"https://upstream.com/root/", "http://localhost:1234/foo", "upstream.com", "/root/foo"},
		{"https://upstream.com/root", "http://localhost:1234/foo", "upstream.com", "/root/foo"},
		{"https://upstream.com/", "http://localhost:1234/foo", "upstream.com", "/foo"},
		{"https://upstream.com", "http://localhost:1234/foo", "upstream.com", "/foo"},
		{"https://upstream.com/", "http://localhost:1234/foo/", "upstream.com", "/foo/"},
		{"https://internal.local:1234/", "http://localhost:1234/foo/", "internal.local:1234", "/foo/"},
		// {"https://upstream.com/foo?bar=baz", "http://localhost:1234/fizz?bix=buzz", "upstream.com", "/foo/fizz?bar=baz&bix=buzz"},
	}
	for _, tc := range testCases {
		tc := tc // capture range variable
		name := fmt.Sprintf("%s => %s => (%s, %s)", tc.req, tc.upstream, tc.host, tc.path)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			parsed, _ := url.Parse(tc.upstream)
			director := RequestRewritingDirector(parsed)
			req := httptest.NewRequest("GET", tc.req, nil)
			director(req)
			// if host, ok := req.Header["Host"]; !ok || host[0] != tc.host {
			if req.Host != tc.host {
				t.Errorf("got host %s; want %s", req.Host, tc.host)
			}
			if req.URL.Path != tc.path {
				t.Errorf("got path %s; want %s", req.URL.Path, tc.path)
			}
		})
	}
}
