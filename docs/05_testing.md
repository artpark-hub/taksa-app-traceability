# Testing Guide — Taksa Traceability

> The test suite uses **Bruno** — a lightweight, Git-friendly API testing tool. Tests are stored as `.bru` files, committed to the repo, and run headlessly via the Bruno CLI.

---

## Quick Start — Run All Tests

```bash
./run_tests.sh
```

That single command:
1. Builds the Go binary
2. Stops any old running server
3. Starts the server in the background
4. Generates all `.bru` test files from the Python script
5. Runs the full test suite
6. Cleans up the binary

Expected output:
```
📊 Execution Summary
┌───────────────┬────────────────┐
│ Status        │     ✓ PASS     │
│ Requests      │ 24 (24 Passed) │
│ Assertions    │     47/47      │
└───────────────┴────────────────┘
```

---

## What Is Bruno?

Bruno is an open-source API client (alternative to Postman/Insomnia) where every test is a plain-text `.bru` file that can be committed to Git. It has a CLI (`@usebruno/cli`) for running tests in CI pipelines or scripts without opening a GUI.

Install the CLI (one-time):
```bash
npm install -g @usebruno/cli
# or use npx (no global install needed)
npx @usebruno/cli --version
```

---

## Directory Structure

```
tests/api/
├── bruno.json                    ← Collection metadata (name, version)
├── environments/
│   └── Dev VM.bru               ← Environment variables (baseUrl)
├── 01_Enterprise/
│   ├── folder.bru               ← Folder metadata
│   └── 1. Create Enterprise.bru ← Test request
├── 02_Site/
├── 03_Area/
├── 04_Line/
├── 05_Unit/
├── 06_EquipmentClass/
├── 07_Equipment/
├── 09_Materials/
├── 10_Operators/
├── 11_Lots/
├── 12_WorkOrders/
├── 13_Genealogy/
├── 15_Queries/                  ← Traceability query tests (S006)
├── 17_Analytics/                ← Historical analytics tests (S014)
└── 99_Cleanup/                  ← Teardown (deletes test data)
```

Folders are numbered so Bruno runs them in the correct dependency order (e.g., you must create an Enterprise before creating a Site inside it).

---

## The `.bru` File Format

Each test is a plain-text file with sections:

```
meta {
  name: 1. Machine Performance   ← Display name shown in output
  type: http
  seq: 1                         ← Order within the folder
}

get {                            ← HTTP method (get/post/patch/delete)
  url: {{baseUrl}}/api/v1/traceability/analytics/machine-performance?equipment_id=MIXER-001&from_time=2026-02-25T00:00:00Z&to_time=2026-02-26T00:00:00Z
  body: none
  auth: none
}

assert {
  res.status: eq 200                         ← HTTP status must be 200
  res.body.equipmentId: eq "MIXER-001"       ← Response body field check
  res.body.totalWorkOrders: isNumber         ← Type assertion
  res.body.utilizationPct: isNumber
}
```

### Available Assertion Types

| Assertion | Meaning |
|---|---|
| `eq VALUE` | Field equals VALUE |
| `isNumber` | Field is a numeric value |
| `isString` | Field is a string |
| `isArray` | Field is an array |
| `isTruthy` | Field is truthy (true, non-zero, non-empty) |
| `isDefined` | Field exists and is not undefined |
| `gt VALUE` | Greater than |
| `lt VALUE` | Less than |

---

## The Environment File

`tests/api/environments/Dev VM.bru`:
```
vars {
  baseUrl: http://localhost:8000
}
```

This defines `{{baseUrl}}` which is substituted into every test URL. To test against a different server, change this value.

---

## How Test Variables Flow Between Tests

Tests use `script:post-response` to save response values into Bruno variables, which later tests can use. This is how dependent resources are chained:

```
# 01_Enterprise/1. Create Enterprise.bru
post { url: {{baseUrl}}/api/v1/traceability/enterprises }
body:json { "name": "Global Motors" }
assert { res.body.id: isNumber }
script:post-response {
  bru.setVar("enterpriseId", res.body.id);   ← saves the new ID
}

# 02_Site/1. Create Site.bru
post { url: {{baseUrl}}/api/v1/traceability/sites }
body:json { "enterprise_id": {{enterpriseId}}, ... }  ← uses saved ID
```

### Full variable chain in the test suite:

```
enterpriseId → siteId → areaId → lineId → unitId → (equipment created with unitId)
defId + wipDefId → (lots created with these IDs)
LOT-SIL-{TS} + WIP-{TS} → genealogy link → trace queries
```

---

## How `generate_tests.py` Works

The test files are **generated dynamically** by a Python script before each run. This is done for one key reason:

> Tests create resources with unique IDs (e.g. `ROBOT-1773127072`) using a **Unix timestamp** (`TS = int(time.time())`). This prevents test runs from colliding with each other if the database already contains data from a previous run.

The script:
1. **Wipes** the `tests/api/` folder entirely
2. **Recreates** it from scratch with fresh `.bru` files
3. Embeds the current timestamp into every resource ID
4. Writes a fresh `bruno.json` collection file and `Dev VM.bru` environment file

**Important:** The `tests/api/` directory is regenerated on every `./run_tests.sh`. You should never manually edit `.bru` files there — edit `generate_tests.py` instead, as your changes will be overwritten.

The analytics tests (`17_Analytics`) are an exception — they use **fixed seed data dates** (`2026-02-25`) rather than the dynamic timestamp, because they query the pre-loaded seed data that's always in the database.

---

## How `run_tests.sh` Works — Step by Step

```bash
# Step 1: Build the Go binary
go build -o traceability ./cmd/traceability/

# Step 2: Stop any previous instance
pkill traceability 2>/dev/null || true

# Step 3: Start server in background, redirect output to server.log
./traceability -conf ./configs/config.yaml > server.log 2>&1 &

# Step 4: Wait for server to be ready
sleep 3

# Step 5: Generate fresh .bru test files
python3 generate_tests.py

# Step 6: Run Bruno CLI against the "Dev VM" environment
cd tests/api
npx @usebruno/cli run --env "Dev VM"

# Step 7: Clean up binary (server keeps running)
rm -f traceability
```

After the script completes, **the server is still running**. Check `server.log` for SQL queries and any errors.

To stop the server:
```bash
pkill traceability
```

---

## Test Coverage

### What Each Folder Tests

| Folder | Tests | Verifies |
|---|---|---|
| `01_Enterprise` | Create Enterprise | Resource creation returns 200 + numeric ID |
| `02_Site` | Create Site | Uses `{{enterpriseId}}` from previous step |
| `03_Area` | Create Area | Uses `{{siteId}}` |
| `04_Line` | Create Production Line | Uses `{{areaId}}` |
| `05_Unit` | Create Production Unit | Uses `{{lineId}}` |
| `06_EquipmentClass` | Create Equipment Class | Class registration |
| `07_Equipment` | Register Equipment | Physical machine with unique timestamped ID |
| `09_Materials` | Create 2 Definitions | RAW_MATERIAL + WORK_IN_PROGRESS types |
| `10_Operators` | Register Operator | Operator with DAY_SHIFT |
| `11_Lots` | Create Parent + Child Lot | Two lots for genealogy |
| `12_WorkOrders` | Create Work Order | Links machine, operator, output lot |
| `13_Genealogy` | Register Link | Parent→Child genealogy link |
| `15_Queries` | 4 trace queries | Backward, Full, Forward, Equipment History |
| `17_Analytics` | 4 analytics queries | Machine perf, compare, trends, dashboard |
| `99_Cleanup` | Delete Site + Enterprise | Cascade delete cleans all test data |

### Assertions Summary (47 total)

| Category | Assertions |
|---|---|
| HTTP 200 status | All 24 requests |
| Resource IDs are numbers | Enterprise, Site, Area, Line, Unit, Class, Material Definitions, Genealogy |
| Trace direction labels | backward, forward |
| Response array fields | genealogy nodes, analytics dataPoints, machines, topPerformers |
| Specific field values | equipmentId, granularity, qualityRatePct, totalUnitsProduced |
| Boolean success | Delete operations |

---

## Running Tests Manually (Without `run_tests.sh`)

If the server is already running and you just want to run tests:

```bash
# Generate fresh test files
python3 generate_tests.py

# Run the suite
cd tests/api
npx @usebruno/cli run --env "Dev VM"
```

To run only one folder:
```bash
cd tests/api
npx @usebruno/cli run 17_Analytics --env "Dev VM"
```

To run with verbose output:
```bash
npx @usebruno/cli run --env "Dev VM" --reporter-json results.json
```

---

## Adding a New Test

1. **Edit `generate_tests.py`** — add a `create_bru(...)` call for your new endpoint.
2. Run `python3 generate_tests.py` to regenerate the `.bru` files.
3. Run `cd tests/api && npx @usebruno/cli run --env "Dev VM"` to verify.

### `create_bru()` function signature

```python
create_bru(
    folder,       # e.g. "17_Analytics" — creates tests/api/17_Analytics/
    filename,     # e.g. "5. New Test.bru"
    name,         # Display name in Bruno output
    method,       # "get", "post", "patch", "delete"
    url,          # Full URL with {{baseUrl}} template
    body,         # JSON string for request body, or None for GET
    assertions,   # List of strings: ["res.status: eq 200", ...]
    seq,          # Order within folder (integer)
    vars_to_set   # Dict to save response fields: {"myVar": "id"}
)
```

### Example — Adding a new analytics test:
```python
create_bru("17_Analytics", "5. Quality Rate.bru", "5. Quality Rate", "get",
    "{{baseUrl}}/api/v1/traceability/analytics/dashboard?from_time=2026-02-01T00:00:00Z&to_time=2026-02-28T00:00:00Z",
    None,
    ["res.status: eq 200", "res.body.qualityRatePct: isNumber", "res.body.lotsReleased: isNumber"],
    5)
```

---

## Debugging Failures

### Check the server log
```bash
tail -50 server.log
```

Look for lines starting with `[ERROR]` or containing `pq:` (PostgreSQL errors).

### Run a single test manually with curl
```bash
curl -s http://localhost:8000/api/v1/traceability/analytics/machine-performance \
  ?equipment_id=MIXER-001\&from_time=2026-02-25T00:00:00Z\&to_time=2026-02-26T00:00:00Z | jq .
```

### Check if server is running
```bash
curl -s http://localhost:8000/api/v1/traceability/enterprises
```

### The analytics tests require seed data
The `17_Analytics` folder tests query against dates `2026-02-24` to `2026-02-27`. These dates correspond to the seed data loaded into PostgreSQL via `database/schema/seed.sql`. If the database was reset without re-seeding, those tests may return empty results but still pass (they assert `isNumber` / `isArray`, not specific values).

---

## CI/CD Integration

The test suite can be integrated into any CI pipeline since it's fully headless:

```yaml
# Example GitHub Actions step
- name: Run API Tests
  run: |
    go build -o traceability ./cmd/traceability/
    ./traceability -conf ./configs/config.yaml > server.log 2>&1 &
    sleep 3
    python3 generate_tests.py
    cd tests/api && npx @usebruno/cli run --env "Dev VM"
```
