package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Integrity-178B/url-fetcher/internal/app/api/config"
	api "github.com/Integrity-178B/url-fetcher/internal/app/api/server"
	"github.com/Integrity-178B/url-fetcher/internal/pkg/fetcher"
	"github.com/Integrity-178B/url-fetcher/internal/pkg/server"
)

func init() {
	config.Init()
}

func main() {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, stopServer := context.WithCancel(context.Background())

	go func() {
		<-signals
		stopServer()
	}()

	f := fetcher.NewFetcher(&config.Get().Fetcher)
	r := api.NewRouter(&config.Get().FetchHandler, f)

	server.NewServer(&config.Get().Server, r).ListenAndServe(ctx)
}
