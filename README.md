# saas-integration-demo
Backend is a SaaS service for integrating AMO CRM and Telegram, which allows you to automate customer service via messenger.

## Tech Stack
- **Go** — core language
- **Gin** — lightweight HTTP router
- **pgx** — PostgreSQL driver
- **Redis** — high-performance caching
- **Clean Architecture** — modular, maintainable design
- **Zap** — logging

## Structure
```
cmd/ — application entry point
env/ - application config 

internal/
  boot/
    di — dependency injection (wiring of dependencies)
    app — application initialization and startup

  constants — project constants
  data — data models / DTOs
  domain — business logic (core layer)

  external/
    adapters — adapters for external services (CRM, Telegram, etc.)

  infrastructure — infrastructure layer (database, cache, configuration)
  transport — HTTP/gRPC handlers (incoming requests)
  utils — utility functions
```

## Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Make sure to set up your `.env` file with DB, Cache credentials and amo crm API keys

### Configuration
```yaml
amo:
  channel_id: 
  channel_id_oauth: 
  channel_secret: 
  channel_secret_key: 
  redirect_url: 
  api_connect_method: 
  api_disconnect_method: 
  domain: 
  referer: 
  signature_secret: 
  source_url: /api/v4/sources
  pipeline_url: /api/v4/leads/pipelines

telegram:
  base_url: https://api.telegram.org/bot
  callback_route: 

common:
  db:
    name: 
    password: 
    username: 
    address: 

  cache:
    url: 
    password: 
    db_num: 0
    max_retries: 3
    dial_timeout: 5s
    timeout: 300s

  server:
    read_timeout: 5s
    write_timeout: 30s
    idle_timeout: 10s

  service:
    auth:
      access_ttl: 15m
      refresh_ttl: 720h

## Run Locally
```

### How to run
```
git clone https://github.com/LemeshkinIvan/saas-integration-demo.git
cd saas-integration-demo
go mod tidy
cd cmd
go run main.go
```

