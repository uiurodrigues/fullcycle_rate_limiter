FROM golang:1.21 as build
WORKDIR /app/cmd
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o rate_limiter
COPY ./cmd/.env /app

FROM scratch
WORKDIR /app
COPY --from=build /app/ .
COPY ./cmd/.env /app

ENTRYPOINT ["./cmd/rate_limiter"]