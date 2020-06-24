package fetcher

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/Integrity-178B/url-fetcher/internal/pkg/log"
)

type worker struct {
	client   *http.Client
	urls     <-chan string
	contents chan<- URLContent
	logger   *log.Logger
}

func newWorker(client *http.Client, urls <-chan string, contents chan<- URLContent, num int) *worker {
	return &worker{
		client:   client,
		urls:     urls,
		contents: contents,
		logger:   log.NewLogger(fmt.Sprintf("[fetch worker %d] ", num)),
	}
}

func (w *worker) process(ctx context.Context, cancel context.CancelFunc, group *sync.WaitGroup) {
	defer group.Done()

	w.logger.Printf("started")
	for {
		select {
		case url, ok := <-w.urls:
			if !ok {
				w.logger.Print("done processing. stopped")
				return
			}

			w.logger.Printf("fetching: %s", url)
			content, err := w.fetch(ctx, url)
			if err != nil {
				w.logger.Printf("fetch url: %s. error: %s", url, err)
				cancel()
				return
			}

			w.logger.Printf("done fetching: %s", url)

			w.contents <- URLContent{
				URL:     url,
				Content: content,
			}
		case <-ctx.Done():
			w.logger.Print("context done. stopped")
			return
		}
	}
}

func (w *worker) fetch(ctx context.Context, url string) (content string, err error) {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}

	response, err := w.client.Do(request)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	content = string(body)

	return
}
