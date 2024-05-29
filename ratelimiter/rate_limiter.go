package ratelimiter

import (
	"fullcycle_rate_limiter/internal/ports"
	"log"
	"sync"
	"time"
)

type RateLimiter struct {
	mu         *sync.RWMutex
	Configs    ports.Conf
	repository ports.Repository
}

func NewRateLimiter(configs ports.Conf, repository ports.Repository) *RateLimiter {
	r := &RateLimiter{
		mu:         &sync.RWMutex{},
		Configs:    configs,
		repository: repository,
	}
	return r
}

func (r *RateLimiter) IsAllowed(key string, maxRequests, blockDuration int64) bool {
	isBlocked, err := r.repository.IsBlocked(key)
	if err != nil {
		log.Println("Error checking if user is blocked: ", err)
		return false
	}
	if isBlocked {
		return false
	}

	userCounter, exists, err := r.repository.Get(key)
	if err != nil {
		log.Println("error getting user counter: ", err)
		return false
	}
	if !exists {
		err = r.repository.Set(key, 1, time.Duration(1)*time.Second)
		if err != nil {
			log.Println("error setting user counter: ", err)
			return false
		}
		return true
	}
	if userCounter >= maxRequests {
		_ = r.repository.Block(key, time.Duration(blockDuration)*time.Minute)
		return false
	}
	_, _ = r.repository.Incr(key)

	return true
}
