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
	var tw *TeeResponseWriter
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		tw = NewTeeResponseWriter(rw, buf)
		tw.Header().Set("content-type", wantContentType)
		tw.WriteHeader(wantStatus)
		_, _ = tw.Write([]byte(wantBody))
	}))
	defer srv.Close()
	resp, err := srv.Client().Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != wantStatus {
		t.Errorf("expected status code=%d; but got=%d", wantStatus, resp.StatusCode)
	}
	if resp.StatusCode != tw.StatusCode() {
		t.Errorf("StatusCode(): expected status code=%d; but got=%d", resp.StatusCode, tw.StatusCode())
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
	var tw *TeeResponseWriter
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		tw = NewTeeResponseWriter(rw, buf)
		tw.Header().Set("content-type", wantContentType)
		_, _ = tw.Write([]byte(wantBody))
	}))
	defer srv.Close()
	resp, err := srv.Client().Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != wantStatus {
		t.Errorf("expected status code=%d; but got=%d", wantStatus, resp.StatusCode)
	}
	if resp.StatusCode != tw.StatusCode() {
		t.Errorf("StatusCode(): expected status code=%d; but got=%d", resp.StatusCode, tw.StatusCode())
	}
	if ct := resp.Header.Get("content-type"); ct != wantContentType {
		t.Errorf("expected content-type header=%s; but got=%s", wantContentType, ct)
	}
	if body := buf.String(); body != wantBody {
		t.Errorf("expected response body=%s; but got=%s", wantBody, body)
	}
}
