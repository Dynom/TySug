package server

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func serviceHandler(l *logrus.Logger, sr ServiceRegistry, validators []Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[6:]
		if name == "" {
			l.Info("no list name defined")
			w.WriteHeader(400)
			return
		}

		if !sr.HasServiceForList(name) {
			l.Infof("list '%s' not defined", name)
			w.WriteHeader(404)
			return
		}

		svc := sr.GetServiceForList(name)
		hf := createRequestHandler(
			l,
			svc,
			validators,
		)

		hf(w, r)
	}
}

func createRequestIDHandler(h http.Handler) http.HandlerFunc {

	// #nosec G404 -- Ignoring linter errors for using a weak random implementation
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	instanceID := strconv.Itoa(int(rnd.Int31()))

	var requestCounter int
	var lock = sync.Mutex{}
	return func(w http.ResponseWriter, r *http.Request) {
		lock.Lock()
		requestCounter++
		requestID := instanceID + "-" + strconv.Itoa(requestCounter)
		lock.Unlock()

		ctx := context.WithValue(r.Context(), CtxRequestID, requestID)
		w.Header().Set(HeaderRequestID, requestID)

		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
