package main

import "net/http"

// NewNpkgdServer returns an http.ServeMux that handles all the requests.
func NewNpkgdServer(upstream string) *http.ServeMux {
	mux := http.NewServeMux()
	proxy := NewUpstreamProxy(upstream)
	mux.Handle("/", proxy)
	return mux
}
