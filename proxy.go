package main

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
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
	cache    Cache
}

// NewResponseRewritingRoundTripper completes an upstream request but rewrites the response
func NewResponseRewritingRoundTripper(upstream *url.URL) http.RoundTripper {
	roundTripper := &ResponseRewritingRoundTripper{
		Upstream: upstream,
		cache:    NewCache(),
	}
	return roundTripper
}

// Cache is a simple concurrent map wrapper
type Cache struct {
	cache map[string][]byte
	lock  sync.RWMutex
}

// NewCache returns an initialized Cache
func NewCache() Cache {
	return Cache{cache: make(map[string][]byte)}
	// return cache
}

// Get retrieves a value from the Cache
func (c *Cache) Get(key string) []byte {
	c.lock.RLock()
	value := c.cache[key]
	c.lock.RUnlock()
	return value
}

// Put stores a new value in the Cache
func (c *Cache) Put(key string, res []byte) []byte {
	c.lock.Lock()
	c.cache[key] = res
	c.lock.Unlock()
	return res
}

// RoundTrip makes a request to an upstream and then rewrites parts of the response as
// necessary before passing it along
func (r *ResponseRewritingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	fromCache := r.cache.Get(req.URL.Path)
	if fromCache == nil {
		newRes, err := http.DefaultTransport.RoundTrip(req)
		if err == nil {
			// TODO: rewrite parts of the response that reference URLs we should proxy
			toCache, err := httputil.DumpResponse(newRes, true)
			if err == nil {
				r.cache.Put(req.URL.Path, toCache)
			}
		}
	}
	fromCache = r.cache.Get(req.URL.Path)
	reader := bufio.NewReader(bytes.NewReader(fromCache))
	return http.ReadResponse(reader, req)
}
