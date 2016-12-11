package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewNpkgdServer(t *testing.T) {
	EXPECTED := []byte("I am upstream!")
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", EXPECTED)
	}))
	defer upstream.Close()
	handler := NewNpkgdServer(upstream.URL)
	ts := httptest.NewServer(handler)
	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(greeting, EXPECTED) {
		t.Errorf("got response '%v'; want '%v'", greeting, EXPECTED)
	}
}
