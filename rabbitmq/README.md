## How To Use (Go Compiler)
- Install golang
- Make sure rabbitmq is ready
- Run the code below
```
# Build
go mod tidy
go mod vendor
go build -o publisher ./publisher/publisher.go
go build -o subscriber ./subscriber/subscriber.go

# Set ENV
export RABBITMQ_URI=amqp://guest:guest@rabbitmq-sample:5672/

# Run Publisher
./publisher

# Run Subscriber
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