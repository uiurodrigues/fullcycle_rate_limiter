package middleware

import (
	"fullcycle_rate_limiter/cmd/ratelimiter"
	"net"
	"net/http"
)

func Limit(next http.HandlerFunc, rl ratelimiter.RateLimiter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rl.Configs.IsRateLimitByTokenEnabled() {
			if !isTokenAllowed(r, rl) {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		} else {
			if rl.Configs.IsRateLimitByIPEnabled() {
				if !isIPAllowed(r, rl) {
					http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func isIPAllowed(r *http.Request, rl ratelimiter.RateLimiter) bool {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return rl.IsAllowed(ip, rl.Configs.GetMaxRequestsByIP(), rl.Configs.GetBlockDurationIP())
}

func isTokenAllowed(r *http.Request, rl ratelimiter.RateLimiter) bool {
	token := r.Header.Get("API_TOKEN")
	if token == "" {
		return true
	}
	tokenLimit, exists := rl.Configs.GetTokenLimit(token)
	if !exists {
		return false
	}
	return rl.IsAllowed(token, tokenLimit, rl.Configs.GetBlockDurationToken())
}
