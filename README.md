# Activity Tracking Service

## Description

The service provides an HTTP API to ingest events and fetch activity data or aggregated statistics per user. A background cron job runs every 4 hours to aggregate activity events into time-bucketed statistics stored in PostgreSQL. A minimal embedded React web client is available for viewing events directly in the browser.

## Run Instructions
**Prerequisites**
  * Go 1.25+ (for local execution)
  * Docker & Docker Compose
  * PostgreSQL 17+ (for local execution)
  * golang-migrate CLI (optional, for local migrations)

1. **Docker Run (Recommended)**
Clone the repository:

```Bash
git clone https://github.com/tokyohardrock/inforce_task.git
cd inforce_task
```
  
Start the application:

```Bash
docker compose build --no-cache app
docker compose up -d
```
Access the service:
  * Web Client: http://localhost:9090/
  * API Base URL: http://localhost:9090

Stop the stack:

```Bash
docker compose down
```

2. **Local Run**
Start PostgreSQL instance and run migrations:

```bash
make migrate_up
```

Run the Go application:

```bash
make run
```

## Sample Requests
1. Create Event
POST /events

```bash
curl -X POST http://localhost:9090/events \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
    "user_id": 42,
    "action": "view",
    "action_object_id": 123,
    "action_object_type": "page",
    "timestamp": "2026-08-07T14:00:00Z",
    "metadata": {
      "path": "/home"
    }
  }'
```

2. Get Events
GET /events?user_id={id}&start_time={startTime}&end_time={endTime}

```bash
curl -X GET "http://localhost:9090/events?user_id=42&start_time=2026-08-07T12:00:00%2B03:00&end_time=2026-08-07T16:00:00%2B03:00"
```

3. Get User Stats
GET /stats/{user}?start_time={startTime}&end_time={endTime}

```bash
curl -X GET "http://localhost:9090/stats/42?start_time=2026-08-07T12:00:00%2B03:00&end_time=2026-08-07T16:00:00%2B03:00"
```

## Notes

**React Web Client:** Built using React 18, Babel Standalone, and UMD bundles loaded via CDN. Embedded directly into the Go binary (//go:embed index.html) and served at GET /. No Node.js build step or NPM packages are required to run the frontend.
