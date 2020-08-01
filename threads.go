package main

import (
	"net/http"
	"sync"
)

// HTTPAction - A http action
type HTTPAction func(http.ResponseWriter, *http.Request)

var httpLock sync.Mutex

// AsThreadSafeHTTPAction - Makes a http action thread safe
func AsThreadSafeHTTPAction(a HTTPAction) HTTPAction {
	return func(w http.ResponseWriter, r *http.Request) {
		httpLock.Lock()
		defer httpLock.Unlock()

		a(w, r)
	}
}
