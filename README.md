## URL Fetcher  
Url fetcher is dead simple service for fetching contents of urls list   

### API  
  `POST /` - fetch contents of urls list

*Example request:* 
```
curl http://localhost:6666 -X POST -H "Content-Type: application/json" --data '["https://google.com", "https://yandex.ru", "https://habr.com", "https://onliner.by", "https://www.linkedin.com"]'
```
*Example response:*
```
[{"url":"https://google.com","content":"..."}, {"url":"https://google.com","content":"..."}, {"url":"https://yandex.ru","content":"..."}, {"url":"https://habr.com","content":"..."}, {"url":"https://onliner.by","content":"..."}] 

```

### Configuration  
Envs, flag, config files are currently not supported. Configuration is hardcoded.

 - `Server.Host` - server host
 - `Server.Port` - server port
 - `FetchHandler.MaxRequests` - maximum number of concurrent http requests to fetch urls
 - `FetchHandler.MaxUrls` - maximum number of urls to fetch per request
 - `Fetcher.ProcessTimeout` - timeout for urls fetch request processing
 - `Fetcher.URLFetchTimeout` - timeout to fetch single url
 - `Fetcher.MaxFetchConcurrency` - maximum number of concurrent url fetches per request

### How to Run
 1. Build binary for your OS and architecture
 2. Run the binary

*Example for MacOS:*
```
make build && ./build/darwin-amd64/api 
```
*Example for Linux:*
```
GOOS=linux make build && ./build/linux-amd64/api 
```
