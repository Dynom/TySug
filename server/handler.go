package server

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func defaultHeaderHandler(h http.Handler) http.HandlerFunc {

	type kv struct {
		Key   string
		Value string
	}

	return func(w http.ResponseWriter, req *http.Request) {
		for _, h := range []kv{
			{Key: "Strict-Transport-Security", Value: "max-age=31536000; includeSubDomains"},
			{Key: "Content-Security-Policy", Value: "default-src 'none'"},
			{Key: "X-Frame-Options", Value: "DENY"},
			{Key: "X-XSS-Protection", Value: "1; mode=block"},
			{Key: "X-Content-Type-Options", Value: "nosniff"},
			{Key: "Referrer-Policy", Value: "strict-origin"},
		} {
			w.Header().Set(h.Key, h.Value)
		}

		h.ServeHTTP(w, req)
	}
}

func createRequestIDHandler(h http.Handler) http.HandlerFunc {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	instanceID := strconv.Itoa(int(rnd.Int31()))

	var requestCounter int
	var lock = sync.Mutex{}
	var buf = strings.Builder{}
	return func(w http.ResponseWriter, r *http.Request) {
		lock.Lock()
		requestCounter++
		buf.Reset()
		buf.WriteString(instanceID)
		buf.WriteString("-")
		buf.WriteString(strconv.Itoa(requestCounter))
		requestID := buf.String()
		lock.Unlock()

		ctx := context.WithValue(r.Context(), CtxRequestID, requestID)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
