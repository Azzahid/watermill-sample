## How To Use (Go Compiler)
- Install golang
- Make sure kafka is ready
- Run the code below
```
# Build
go mod tidy
go mod vendor
go build -o publisher ./publisher/publisher.go
go build -o subscriber ./subscriber/subscriber.go

# Set ENV
export KAFKA_BROKERS=kafka:29092,kafka:29093

# Run publisher
./publisher

# Run subscriber
./subscriber
```

## How To Use (Docker)
- Install docker
- Run the code below
```
go mod tidy
go mod vendor
docker compose up -d
```