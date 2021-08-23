// Package httputil provides some utilities related to net/http.

package httputil

import "net/http"

type Middleware = func(next http.Handler) http.Handler

// ApplyIf returns new Middlware that wraps `next` handler if `cond` returns true otherwise just delegates to `next` handler.
//
// You can determine `cond` result according to incoming *http.Request's value.
func ApplyIf(cond func(r *http.Request) bool, mw Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cond(r) {
				mw(next).ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
