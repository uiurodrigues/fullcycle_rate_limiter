# Rate Limiter


## Installation

- Execute the command [docker-compose up -d] to ensure that the environment is working
- Execute the command [go run cmd/main.go]

## Explaining the configuration
* ENABLE_RATE_LIMIT_BY_IP
    - Enables tha rate limit to work on IP when its value is [true]

* ENABLE_RATE_LIMIT_BY_TOKEN
    - Enables tha rate limit to work on Tokens when its value is [true]

* MAX_REQUESTS_BY_IP
    - Sets the max value of requests received by second by one IP (Example: [10])

* BLOCK_DURATION_IP 
    - Sets the time in minutes that one IP will be blocked to execute news requests (Example: [1])

* BLOCK_DURATION_TOKEN
    - Sets the time in minutes that one TOKEN will be blocked to execute news requests (Example: [1])

* TOKEN_LIMIT_LIST 
    - Sets the max value of requests received by second by one specific Token (Example: [{ "ABC123": 2, "ABC456": 4, "DEF122": 5}])

* STORAGE_TYPE
    - Sets the type of storage that the application will use
        - Set it to [redis] if you want redis
        - Any other value (or none) the application will use a in memory storage


