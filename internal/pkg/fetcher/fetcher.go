package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Integrity-178B/url-fetcher/internal/pkg/log"
)

// Config keeps fetcher configuration
type Config struct {
	ProcessTimeout      time.Duration
	URLFetchTimeout     time.Duration
	MaxFetchConcurrency int
}

// Fetcher is basically http.Client that fetches provided urls content
type Fetcher struct {
	conf   *Config
	client *http.Client
	logger *log.Logger
}

// NewFetcher creates new fetcher instance
func NewFetcher(conf *Config) *Fetcher {
	return &Fetcher{
		conf:   conf,
		client: &http.Client{Timeout: conf.URLFetchTimeout},
		logger: log.NewLogger("[fetcher] "),
	}
}

// Urls is url list
type Urls []string

// Validate checks that all urls are valid
func (urs Urls) Validate() error {
	for _, ur := range urs {
		u, err := url.Parse(ur)
		if err != nil {
			return err
		}
		if u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("'%s' has invalid scheme or host", ur)
		}
	}

	return nil
}

// URLContent keeps the url and its fetched content
type URLContent struct {
	URL     string `json:"url"`
	Content string `json:"content"`
}

// Fetch fetches contents of provided list of urls. The content is simply the body of http response
func (f Fetcher) Fetch(ctx context.Context, urls []string) ([]URLContent, error) {
	ctx, cancel := context.WithTimeout(ctx, f.conf.ProcessTimeout)
	defer cancel()

	contents := make([]URLContent, 0, len(urls))

	urlsCh, group := make(chan string, len(urls)), new(sync.WaitGroup)
	contentsCh, contentsDone := make(chan URLContent, len(urls)), make(chan struct{})

	group.Add(f.conf.MaxFetchConcurrency)
	for i := 0; i < f.conf.MaxFetchConcurrency; i++ {
		worker := newWorker(f.client, urlsCh, contentsCh, i)
		go worker.process(ctx, cancel, group)
	}
	f.logger.Printf("worker pool of size %d created", f.conf.MaxFetchConcurrency)

	for _, u := range urls {
		urlsCh <- u
	}
	close(urlsCh)
	f.logger.Printf("urls are pushed to job queue. total jobs: %d", len(urls))

	go func() {
		for content := range contentsCh {
			contents = append(contents, content)
		}

		contentsDone <- struct{}{}
	}()
	f.logger.Print("started waiting for fetched content")

	group.Wait()
	close(contentsCh)
	<-contentsDone

	f.logger.Print("done collecting url contents")

	return contents, nil
}
