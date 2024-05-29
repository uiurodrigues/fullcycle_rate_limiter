package main

import (
	"fullcycle_rate_limiter/cmd/middleware"
	"fullcycle_rate_limiter/cmd/ratelimiter"
	"fullcycle_rate_limiter/configs"
	"fullcycle_rate_limiter/internal/database"
	"fullcycle_rate_limiter/internal/ports"
	"log"
	"net/http"

	"github.com/go-redis/redis"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db := getDB(configs.GetStorageType())

	mux := http.NewServeMux()

	mux.HandleFunc("/", middleware.Limit(helloWorldHandler, *ratelimiter.NewRateLimiter(configs, db)))

	log.Println("listening on port 8080")

	http.ListenAndServe(":8080", mux)
}

func getDB(storageType string) ports.Repository {
	if storageType == "redis" {
		return database.NewRedisDB(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	} else {
		return database.NewLocalDB()
	}
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
