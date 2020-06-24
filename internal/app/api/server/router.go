package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Integrity-178B/url-fetcher/internal/pkg/fetcher"
	"github.com/Integrity-178B/url-fetcher/internal/pkg/log"
	"github.com/Integrity-178B/url-fetcher/internal/pkg/server"
)

const (
	jsonContentType = "application/json"
)

// NewRouter creates new router (http handler)
func NewRouter(conf *FetchHandlerConfig, fetcher *fetcher.Fetcher) http.Handler {
	mux := http.NewServeMux()

	handler := newFetchHandler(conf, fetcher)
	mux.Handle("/", server.MaxRequestsMiddleware(handler, conf.MaxRequests))

	return mux
}

// FetchHandlerConfig keeps fetch handler configuration
type FetchHandlerConfig struct {
	MaxRequests    int
	RequestTimeout time.Duration
	MaxUrls        int
}

// FetchHandler is handler for url fetch requests
type FetchHandler struct {
	conf    *FetchHandlerConfig
	fetcher *fetcher.Fetcher
	logger  *log.Logger
}

func newFetchHandler(conf *FetchHandlerConfig, fetcher *fetcher.Fetcher) *FetchHandler {
	return &FetchHandler{
		conf:    conf,
		fetcher: fetcher,
		logger:  log.NewLogger("[fetch handler] "),
	}
}

// ServeHTTP handles url fetch requests and implements handler interface
func (fh FetchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := server.ResponseWriter{ResponseWriter: w}

	switch r.Method {
	case http.MethodPost:
		if r.Header.Get("Content-Type") != jsonContentType {
			fh.errorResponse(rw, fmt.Errorf("Content-Type header should be %s", jsonContentType))
			return
		}

		urls := make(fetcher.Urls, 0, fh.conf.MaxUrls)
		if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
			fh.errorResponse(rw, fmt.Errorf("unable to parse json url list: %s", err))
			return
		}
		if len(urls) > fh.conf.MaxUrls {
			fh.errorResponse(rw, fmt.Errorf("url count limit exceeded (max %d urls)", fh.conf.MaxUrls))
			return
		}
		if err := urls.Validate(); err != nil {
			fh.errorResponse(rw, fmt.Errorf("url is invalid: %s", err))
			return
		}

		data, err := fh.fetcher.Fetch(r.Context(), urls)
		if err != nil {
			fh.errorResponse(rw, fmt.Errorf("failed to fetch some urls: %s", err))
			return
		}

		err = rw.WriteJSON(data)
		if err != nil {
			fh.logger.Printf("write response data: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (fh FetchHandler) errorResponse(rw server.ResponseWriter, e error) {
	fh.logger.Print(e)

	err := rw.WriteError(e)
	if err != nil {
		fh.logger.Printf("write response data: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusBadRequest)
}
