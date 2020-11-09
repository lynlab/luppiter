# Luppiter

> Restful APIs for Luppiter services.
> For details, see documentations on [Luppiter Console](https://console.luppiter.dev).

## Development

### Prerequisites

- Go 1.15 (or greater)
- Docker Compose

### Start Database

```sh
# Start database
docker-compose up -d

# Run migration
go run ./app/migrate up
```

### Run Server

```sh
go run ./app/api
```
