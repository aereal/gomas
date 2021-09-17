package httputil

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTeeResponseWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	wantStatus := http.StatusCreated
	wantBody := `{"ok":true}`
	wantContentType := "application/json"
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		w := NewTeeResponseWriter(rw, buf)
		w.Header().Set("content-type", wantContentType)
		w.WriteHeader(wantStatus)
		_, _ = w.Write([]byte(wantBody))
	}))
	defer srv.Close()
	resp, err := srv.Client().Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != wantStatus {
		t.Errorf("expected status code=%d; but got=%d", wantStatus, resp.StatusCode)
	}
	if ct := resp.Header.Get("content-type"); ct != wantContentType {
		t.Errorf("expected content-type header=%s; but got=%s", wantContentType, ct)
	}
	if body := buf.String(); body != wantBody {
		t.Errorf("expected response body=%s; but got=%s", wantBody, body)
	}
}

func TestTeeResponseWriter_implicitWriteHeader(t *testing.T) {
	buf := new(bytes.Buffer)
	wantStatus := http.StatusOK
	wantBody := `{"ok":true}`
	wantContentType := "application/json"
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		w := NewTeeResponseWriter(rw, buf)
		w.Header().Set("content-type", wantContentType)
		_, _ = w.Write([]byte(wantBody))
	}))
	defer srv.Close()
	resp, err := srv.Client().Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != wantStatus {
		t.Errorf("expected status code=%d; but got=%d", wantStatus, resp.StatusCode)
	}
	if ct := resp.Header.Get("content-type"); ct != wantContentType {
		t.Errorf("expected content-type header=%s; but got=%s", wantContentType, ct)
	}
	if body := buf.String(); body != wantBody {
		t.Errorf("expected response body=%s; but got=%s", wantBody, body)
	}
}
