package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoot(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	proxy(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got http %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestProxyJPEG(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(proxy))

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./testdata/oyaki.jpg")
	}))

	orgSrvURL = origin.URL

	url := ts.URL + "/oyaki.jpg"

	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	orgRes, err := http.Get(orgSrvURL)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("HTTP status is %d, want %d", res.StatusCode, http.StatusOK)
	}

	if res.ContentLength < 0 {
		t.Errorf("Content-Length header does not exist")
	}

	if res.ContentLength >= orgRes.ContentLength {
		t.Errorf("value of Content-Length header %d is larger than origin's one %d", res.ContentLength, orgRes.ContentLength)
	}
}

func TestProxyPNG(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(proxy))

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./testdata/corn.png")
	}))

	orgSrvURL = origin.URL
	url := ts.URL + "/corn.png"

	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	orgRes, err := http.Get(orgSrvURL)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("HTTP status is %d, want %d", res.StatusCode, http.StatusOK)
	}

	if res.ContentLength < 0 {
		t.Errorf("Content-Length header does not exist")
	}

	if res.ContentLength != orgRes.ContentLength {
		t.Errorf("value of Content-Length header %d is not equal to origin's one, want %d", res.ContentLength, orgRes.ContentLength)
	}
}

func TestOriginNotExist(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(proxy))

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}))

	orgSrvURL = origin.URL

	url := ts.URL + "/none.jpg"

	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusBadGateway {
		t.Errorf("HTTP status is %d, want %d", res.StatusCode, http.StatusBadGateway)
	}
}
