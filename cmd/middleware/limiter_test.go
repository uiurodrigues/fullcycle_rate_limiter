package middleware

import (
	"encoding/json"
	"fullcycle_rate_limiter/cmd/ratelimiter"
	"fullcycle_rate_limiter/internal/database"
	"fullcycle_rate_limiter/internal/fixtures"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockForIp10From200(t *testing.T) {
	configs := fixtures.MockedConfig{
		EnableRateLimitByIp: "true",
		MaxRequestsByIP:     "10",
		BlockDurationIP:     "1",
	}
	qtRequests := 200
	qtParallelThreads := 10

	mux := http.NewServeMux()
	mux.HandleFunc("/", Limit(helloWorldHandler, *ratelimiter.NewRateLimiter(&configs, database.NewLocalDB())))

	parallelThreads := make(chan int, qtParallelThreads)
	defer close(parallelThreads)

	wg := &sync.WaitGroup{}
	wg.Add(qtRequests)

	var successes atomic.Int32
	var errors atomic.Int32
	for i := 0; i < qtRequests; i++ {
		parallelThreads <- 1
		go func() {
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "0.0.0.1:8000"
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			if rr.Code == http.StatusOK {
				successes.Add(1)
			} else {
				errors.Add(1)
			}
			wg.Done()
			<-parallelThreads
		}()
	}
	wg.Wait()

	assert.Equal(t, 10, int(successes.Load()))
	assert.Equal(t, 190, int(errors.Load()))
}

func TestBlockForToken10From200(t *testing.T) {
	tokenList := map[string]int64{"ABC123": 10}
	jsonMap, _ := json.Marshal(tokenList)

	configs := &fixtures.MockedConfig{
		EnableRateLimitByToken: "true",
		BlockDurationToken:     "1",
		TokenLimitList:         string(jsonMap),
	}
	_ = configs.LoadTokenLimitList()

	qtRequests := 200
	qtParallelThreads := 10

	mux := http.NewServeMux()
	mux.HandleFunc("/", Limit(helloWorldHandler, *ratelimiter.NewRateLimiter(configs, database.NewLocalDB())))

	parallelThreads := make(chan int, qtParallelThreads)
	defer close(parallelThreads)

	wg := &sync.WaitGroup{}
	wg.Add(qtRequests)

	var successes atomic.Int32
	var errors atomic.Int32
	for i := 0; i < qtRequests; i++ {
		parallelThreads <- 1
		go func() {
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "0.0.0.1:8000"
			req.Header.Set("API_TOKEN", "ABC123")
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			if rr.Code == http.StatusOK {
				successes.Add(1)
			} else {
				errors.Add(1)
			}
			wg.Done()
			<-parallelThreads
		}()
	}
	wg.Wait()

	assert.Equal(t, 10, int(successes.Load()))
	assert.Equal(t, 190, int(errors.Load()))
}

func TestBlockForTokenHigherThanIp(t *testing.T) {
	tokenList := map[string]int64{"ABC123": 10}
	jsonMap, _ := json.Marshal(tokenList)

	configs := &fixtures.MockedConfig{
		EnableRateLimitByToken: "true",
		EnableRateLimitByIp:    "true",
		BlockDurationToken:     "1",
		BlockDurationIP:        "1",
		MaxRequestsByIP:        "1",
		TokenLimitList:         string(jsonMap),
	}
	_ = configs.LoadTokenLimitList()

	qtRequests := 200
	qtParallelThreads := 1

	mux := http.NewServeMux()
	mux.HandleFunc("/", Limit(helloWorldHandler, *ratelimiter.NewRateLimiter(configs, database.NewLocalDB())))

	parallelThreads := make(chan int, qtParallelThreads)
	defer close(parallelThreads)

	wg := &sync.WaitGroup{}
	wg.Add(qtRequests)

	var successes atomic.Int32
	var errors atomic.Int32
	for i := 0; i < qtRequests; i++ {
		parallelThreads <- 1
		go func() {
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "0.0.0.1:8000"
			req.Header.Set("API_TOKEN", "ABC123")
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			if rr.Code == http.StatusOK {
				successes.Add(1)
			} else {
				errors.Add(1)
			}
			wg.Done()
			<-parallelThreads
		}()
	}
	wg.Wait()

	assert.Equal(t, 10, int(successes.Load()))
	assert.Equal(t, 190, int(errors.Load()))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
