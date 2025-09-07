# B3 Ingestor

Main features:

1. **Data ingestion** (CLI)
2. **Data querying** (REST API)

---

## How to run

### 0. Place ZIP files

Before running `docker compose up`, place all ZIP files in the `data` directory at the root of the project.

### 1. Using Docker Compose

```bash
docker compose up --build
```

* The `postgres` container initializes the database.
* The `b3-ingestor` container can run data ingestion or the API.

Batch size and temporary directory options are read from `config.yaml`.

---

### 2. Run the REST API

After running `docker compose up`, the API should be available at `http://localhost:8080`.

#### Endpoint:

```
GET /trade_summary?ticker=PETR4&data_start=2025-09-01
```

* **ticker**: required
* **data\_start**: optional (YYYY-MM-DD). If omitted, defaults to the last 7 days until yesterday.

Example JSON response:

```json
{
  "ticker": "PETR4",
  "max_range_value": 20.50,
  "max_daily_volume": 150000
}
```

---

## Tests

Tests can be run with a mocked database using `go-sqlmock`:

```bash
go test ./...
```
