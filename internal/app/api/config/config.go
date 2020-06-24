package config

import (
	"sync"
	"time"

	api "github.com/Integrity-178B/url-fetcher/internal/app/api/server"
	"github.com/Integrity-178B/url-fetcher/internal/pkg/fetcher"
	"github.com/Integrity-178B/url-fetcher/internal/pkg/server"
)

var (
	conf Config
	once sync.Once
)

// Config keeps application config
type Config struct {
	Server       server.Config
	FetchHandler api.FetchHandlerConfig
	Fetcher      fetcher.Config
}

// Init initializes application config
// (with hardcoded data, but here configs should be parsed from envs, flag, config files, etc.)
func Init() {
	once.Do(func() {
		conf = Config{
			Server: server.Config{
				Host: "0.0.0.0",
				Port: "6666",
			},
			FetchHandler: api.FetchHandlerConfig{
				MaxRequests: 100,
				MaxUrls:     20,
			},
			Fetcher: fetcher.Config{
				ProcessTimeout:      10 * time.Second,
				URLFetchTimeout:     time.Second,
				MaxFetchConcurrency: 4,
			},
		}
	})
}

// Get returns application configuration
func Get() *Config {
	return &conf
}
