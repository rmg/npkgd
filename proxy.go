package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Director is a simple request rewriter for httputil.ReverseProxy
type Director func(*http.Request)

// NewUpstreamProxy returns an http request handler that simply proxies requests to a given upstream server.
func NewUpstreamProxy(upstreamRoot string) *httputil.ReverseProxy {
	target, _ := url.Parse(upstreamRoot)
	proxy := &httputil.ReverseProxy{
		Director:  RequestRewritingDirector(target),
		Transport: NewResponseRewritingRoundTripper(target),
	}
	return proxy
}

// RequestRewritingDirector fills in the expected Host header for the upstream
func RequestRewritingDirector(target *url.URL) Director {
	return func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = strings.TrimRight(target.Path, "/") + "/" + strings.TrimLeft(req.URL.Path, "/")
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		req.Host = target.Host
	}
}

// ResponseRewritingRoundTripper is an http agent that rewrites responses before passing them back
type ResponseRewritingRoundTripper struct {
	Upstream *url.URL
}

// NewResponseRewritingRoundTripper completes an upstream request but rewrites the response
func NewResponseRewritingRoundTripper(upstream *url.URL) http.RoundTripper {
	roundTripper := &ResponseRewritingRoundTripper{upstream}
	return roundTripper
}

// RoundTrip makes a request to an upstream and then rewrites parts of the response as
// necessary before passing it along
func (r ResponseRewritingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// TODO: check if req can be fulfilled from our cache
	res, err := http.DefaultTransport.RoundTrip(req)
	// TODO: do something with res
	return res, err
}
