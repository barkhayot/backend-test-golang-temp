# Backend Golang Test

- This backend application is written in Go for testing and will be deprecated after after review.
- Reference for task is available [here](https://github.com/kolikosoft/backend-test-golang)

## What the Application Has & Does
- Fetches tradable/untradable item data from Skinport API
- Manages user balances and tracks balance history
- Allows charging user balances and retrieving balance history
- Runs database migrations for user and balance history tables
- Exposes REST API endpoints for items, balance charge, and balance history
- Uses Fiber web framework for HTTP server

## REST API Endpoints
- `GET /items`: List Skinport items
- `POST /balance/:id/charge`: Charge user balance
- `GET /balance/:id/history`: Get user balance history

## Configuration
Configure via environment variables (see `.env.example`):
- `FETCH_INTERVAL_IN_MINUTE`: Item fetch interval (minutes)
- `FETCH_ON_LAUNCH`: Fetch items on startup (`true`/`false`)
- `APP_PORT`: HTTP server port (e.g., `8080`)
- `PGDB`: Database connection string
- `RUN_MIGRATIONS`: Run DB migrations on startup (`true`/`false`)
- ‼️ Important: It will create new tables in the database for users and balance history.

## Running the Application
1. **Clone the repository**
2. **Set environment variables** (copy `.env.example` to `.env` and fill in values, or export them directly)
3. **Install dependencies:**
   ```sh
   go mod tidy
   ```
4. **Run using Makefile:**
   ```sh
   make run
   ```
   - To run tests: `make test`
   - To format code: `make fmt`

The server will start and listen on the configured port. Access `http://localhost:8080/items` (or your chosen port) to get the items data.


## Example of `curl` Requests
- List items:
  ```sh
  curl http://localhost:8080/items
  ```

- Charge user balance:
  ```sh
  curl -X POST http://localhost:8080/balance/1/charge \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.2
  }'
  ```

- Get user balance history:
  ```sh
  curl http://localhost:8080/balance/1/history
  ```
