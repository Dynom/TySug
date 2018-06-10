package server

import (
	"encoding/json"
	"net/http"

	"time"

	"io"
	"io/ioutil"

	"github.com/Dynom/TySug/server/service"
	"github.com/pkg/errors"
	"github.com/rs/cors"
)

const maxBodySize = 1 << 20

var (
	ErrMissingBody        = errors.New("missing body")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrInvalidRequestBody = errors.New("invalid request body")
)

type request struct {
	Input string `json:"input"`
}

type response struct {
	Result string  `json:"result"`
	Score  float64 `json:"score"`
}

func getRequestFromHttpRequest(r *http.Request) (request, error) {
	var req request

	b, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		if err == io.EOF {
			return req, ErrMissingBody
		}
		return req, ErrInvalidRequest
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		return req, ErrInvalidRequestBody
	}

	return req, nil
}

type TySugServer struct {
	server http.Server
}

func (tss TySugServer) ListenOnAndServe(addr string) error {
	s := tss.server
	s.Addr = addr

	return s.ListenAndServe()
}

func NewHTTP(tysug service.Domain, mux *http.ServeMux) TySugServer {

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut},
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		req, err := getRequestFromHttpRequest(r)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(400)
			return
		}

		var res response
		res.Result, res.Score = tysug.Rank(req.Input)

		response, err := json.Marshal(res)
		if err != nil {
			w.Write([]byte("unable to marshal result, b00m"))
			w.WriteHeader(500)
			return
		}

		w.Write(response)
	})

	server := http.Server{
		//Addr:              "0.0.0.0:1337",
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
		Handler:           defaultHeaderHandler(c.Handler(mux)),
	}

	return TySugServer{
		server: server,
	}
}

func defaultHeaderHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'none'")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin")

		h.ServeHTTP(w, req)
	}
}
