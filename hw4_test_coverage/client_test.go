package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestCase struct {
	Request SearchRequest
	IsError bool
}

func SearchServerBeta(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("AccessToken")
	if token != "qwer" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	order := r.FormValue("order_field")
	if order != "age" && order != "name" && order != "id" {
		w.WriteHeader(http.StatusBadRequest)
		switch order {
		case "unknown_json_error":
			w.Write([]byte(`{"error": "1111"}`))
		case "json_error":
			w.Write([]byte(`{"error": "broken json`))
		case "else":
			w.Write([]byte(`{"error": "NotErrorBadOrderField"}`))
		default:
			w.Write([]byte(`{"error": "ErrorBadOrderField"}`))
		}
		return
	}
	if order == "age" {
		w.Write([]byte(`[{"name":"bob {"name": "ja"`))
		return
	}
	if order == "bad_json" {
		w.Write([]byte(`[{"name":"bob"}, {"name": "ja"`))
		return
	}
	w.Write([]byte(`[{"name":"bob"}, {"name": "ja"}]`))
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func TimeoutError(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
}

func TestWrongUrl(t *testing.T) {
	c := &SearchClient{
		URL:         "qwer",
		AccessToken: "qwer",
	}
	_, err := c.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("Wrong URL")
	}
}

func TestTimeoutError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(TimeoutError))
	c := &SearchClient{
		URL:         ts.URL,
		AccessToken: "qwer",
	}
	_, err := c.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("TimeoutError")
	}
}

func TestBadOrderField(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerBeta))
	sc := &SearchClient{
		AccessToken: "qwer",
		URL:         ts.URL,
	}
	cases := []TestCase{
		{SearchRequest{OrderField: "json_error"},true},
		{SearchRequest{OrderField: "default"},true},
		{SearchRequest{OrderField: "else"},true},
		{SearchRequest{OrderField: "age"},true},
		{SearchRequest{OrderField: "unknown_json_error"},true},
		{SearchRequest{OrderField: "bad_json"},true},
	}
	for _, testCase := range cases {
		_, err := sc.FindUsers(SearchRequest{OrderField: testCase.Request.OrderField})
		if err != nil && !testCase.IsError {
			t.Errorf("Unexpected error")
		}
		if err == nil && testCase.IsError {
			t.Errorf("Bad Order Field")
		}
	}
	ts.Close()
}

func TestInternalServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(InternalServerError))
	sc := &SearchClient{
		AccessToken: "qwer",
		URL:         ts.URL,
	}
	_, err := sc.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("InternalServerError")
	}
	ts.Close()
}

func TestBadToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerBeta))
	sc := &SearchClient{
		AccessToken: "not qwer",
		URL:         ts.URL,
	}
	_, err := sc.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("Token error")
	}
	ts.Close()
}

func TestBadLimit(t *testing.T) {
	cases := []TestCase{
		{SearchRequest{Limit: -1, OrderField: "name"},true},
		{SearchRequest{Limit: 26, OrderField: "name"},false},
		{SearchRequest{Offset: -1, OrderField: "name"},true},
		{SearchRequest{Limit: 1, OrderField: "name"},false},
	}
	ts := httptest.NewServer(http.HandlerFunc(SearchServerBeta))
	for _, testCase := range cases {
		sc := &SearchClient{
			AccessToken: "qwer",
			URL:         ts.URL,
		}
		_, err := sc.FindUsers(testCase.Request)
		if err != nil && !testCase.IsError {
			t.Errorf("unexpected error: %#v", err)
		}
		if err == nil && testCase.IsError {
			t.Errorf("expected error, got nil")
		}
	}
	ts.Close()
}
