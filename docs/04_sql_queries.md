# SQL Queries Guide — Taksa Traceability

> This guide explains every SQL query used by the application, how each one works, and what the result looks like.  
> **File location:** `database/queries/traceability_queries.sql`  
> **Note:** These queries are documentation only. The app executes them inline in `internal/data/traceability.go` via GORM's `.Raw()` method. `$1, $2, ...` are PostgreSQL parameter placeholders.

---

## Section 1 — Material & Work Order Queries

---

### Query 1.1 — Get Single Material Lot Detail

**What it does:** Fetches everything about one specific material lot, including what material it is.

**Parameters:** `$1` = lot_id (e.g. `"RM-LI-2026-001"`)

```sql
SELECT
    ml.lot_id,
    ml.status,
    ml.quantity,
    ml.unit_of_measure,
    ml.created_at,
    ml.updated_at,
    md.id           AS material_definition_id,
    md.name         AS material_name,
    md.material_type,
    md.description  AS material_description
FROM material_lot ml
JOIN material_definition md ON md.id = ml.material_definition_id
WHERE ml.lot_id = $1;
```

**How it works:** Joins `material_lot` with `material_definition` to pull the lot's physical data alongside the definition of what kind of material it is.

**Example result:**
| lot_id | status | quantity | unit_of_measure | material_name | material_type |
|---|---|---|---|---|---|
| RM-LI-2026-001 | available | 500.0000 | kg | Lithium Carbonate Grade A | RAW_MATERIAL |

---

### Query 1.2 — List Material Lots (with optional filters)

**What it does:** Returns a summary list of lots, optionally filtered by status and/or material type.

**Parameters:** `$1` = status filter (pass `''` to skip) · `$2` = material_type filter (pass `''` to skip)

```sql
SELECT
    ml.lot_id,
    ml.status,
    ml.quantity,
    ml.unit_of_measure,
    ml.created_at,
    md.name         AS material_name,
    md.material_type
FROM material_lot ml
JOIN material_definition md ON md.id = ml.material_definition_id
WHERE ($1 = '' OR ml.status = $1)
  AND ($2 = '' OR md.material_type = $2)
ORDER BY ml.created_at DESC;
```

**How it works:** The `($1 = '' OR ...)` pattern makes each filter optional — if you pass an empty string it matches all rows. Results are sorted newest first.

**Example result (filter: status = 'available'):**
| lot_id | status | quantity | material_name | material_type |
|---|---|---|---|---|
| RM-LI-2026-042 | available | 1000 | Lithium Carbonate | RAW_MATERIAL |
| RM-CARBON-001 | available | 250 | Graphite | RAW_MATERIAL |

---

### Query 1.3 — Get Work Order Detail

**What it does:** Returns full work order info plus the equipment class name, operator name, and operator shift.

**Parameters:** `$1` = work_order_id

```sql
SELECT
    wo.work_order_id, wo.description, wo.status,
    wo.planned_quantity, wo.actual_quantity, wo.unit_of_measure,
    wo.planned_start, wo.planned_end, wo.actual_start, wo.actual_end,
    wo.output_lot_id, wo.equipment_id,
    ec.class_name       AS equipment_class_name,
    wo.operator_id,
    op.name             AS operator_name,
    op.shift            AS operator_shift
FROM work_order wo
LEFT JOIN equipment_master em ON em.id = wo.equipment_id
LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
LEFT JOIN operator op ON op.operator_id = wo.operator_id
WHERE wo.work_order_id = $1;
```

**How it works:** Three LEFT JOINs enrich the work order with human-readable names for the machine type and operator. Uses LEFT JOIN so a WO missing an equipment or operator assignment doesn't disappear from results.

**Example result:**
| work_order_id | status | equipment_id | equipment_class_name | operator_name | operator_shift | actual_quantity |
|---|---|---|---|---|---|---|
| WO-2026-0001 | completed | MIXER-001 | Hosokawa Mixer | Kavya R. | DAY_SHIFT | 247.5 |

---

### Query 1.4 — Get Input Lots for a Work Order

**What it does:** Returns all parent (input) lots that were consumed to produce this work order's output.

**Parameters:** `$1` = work_order_id

```sql
SELECT
    lg.parent_lot_id    AS input_lot_id,
    md.name             AS material_name,
    md.material_type,
    lg.quantity_consumed,
    ml.unit_of_measure
FROM lot_genealogy lg
JOIN material_lot ml ON ml.lot_id = lg.parent_lot_id
JOIN material_definition md ON md.id = ml.material_definition_id
WHERE lg.work_order_id = $1;
```

**How it works:** Looks up genealogy records for this work order and joins to material info of each parent lot.

**Example result:**
| input_lot_id | material_name | material_type | quantity_consumed | unit_of_measure |
|---|---|---|---|---|
| RM-LI-2026-001 | Lithium Carbonate | RAW_MATERIAL | 200.0000 | kg |
| RM-CARBON-001 | Graphite | RAW_MATERIAL | 50.0000 | kg |

---

## Section 2 — Traceability Queries

---

### Query 2.1 — Backward Trace (Recursive)

**What it does:** Starting from a given lot, recursively walks backwards through `lot_genealogy` to find every ancestor lot (inputs of inputs of inputs...).

**Parameters:** `$1` = the lot_id to trace backwards from (the finished product)

```sql
WITH RECURSIVE backward_trace AS (
    -- Base case: direct parents of the target lot
    SELECT lg.parent_lot_id, lg.child_lot_id, lg.work_order_id,
           lg.equipment_id, lg.quantity_consumed, lg.quantity_produced,
           lg.event_time, 1 AS depth
    FROM lot_genealogy lg
    WHERE lg.child_lot_id = $1

    UNION ALL

    -- Recursive step: parents of parents (keep going up the tree)
    SELECT lg.parent_lot_id, lg.child_lot_id, lg.work_order_id,
           lg.equipment_id, lg.quantity_consumed, lg.quantity_produced,
           lg.event_time, bt.depth + 1
    FROM lot_genealogy lg
    JOIN backward_trace bt ON lg.child_lot_id = bt.parent_lot_id
)
SELECT
    bt.depth,
    bt.parent_lot_id       AS lot_id,
    bt.child_lot_id        AS consumed_by_lot_id,
    md.name AS material_name, md.material_type,
    ml.status AS lot_status, ml.quantity, ml.unit_of_measure,
    bt.quantity_consumed,
    bt.work_order_id, bt.equipment_id,
    ec.class_name AS equipment_class_name,
    wo.operator_id, op.name AS operator_name,
    bt.event_time
FROM backward_trace bt
JOIN material_lot ml ON ml.lot_id = bt.parent_lot_id
JOIN material_definition md ON md.id = ml.material_definition_id
LEFT JOIN equipment_master em ON em.id = bt.equipment_id
LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
LEFT JOIN work_order wo ON wo.work_order_id = bt.work_order_id
LEFT JOIN operator op ON op.operator_id = wo.operator_id
ORDER BY bt.depth, bt.parent_lot_id;
```

**How it works — the recursive CTE:**
1. **Base case:** Find all rows in `lot_genealogy` where `child_lot_id = target_lot`. These are the **direct parents** (depth 1).
2. **Recursive step:** Take each parent found, and find THEIR parents (where `child_lot_id = current parent`). Increment depth by 1.
3. PostgreSQL keeps repeating step 2 until no more parents are found.
4. The final SELECT enriches each result row with material and operator info.

**Example result (tracing BATTERY-PACK-001 backward):**
| depth | lot_id | consumed_by_lot_id | material_name | quantity_consumed | equipment_id | operator_name |
|---|---|---|---|---|---|---|
| 1 | WIP-CELL-MODULE-017 | BATTERY-PACK-001 | Cell Module | 10 | ASSEMBLER-01 | Ravi K. |
| 2 | WIP-SLURRY-003 | WIP-CELL-MODULE-017 | Anode Slurry | 250 | MIXER-001 | Kavya R. |
| 3 | RM-LI-2026-042 | WIP-SLURRY-003 | Lithium Carbonate | 200 | MIXER-001 | Kavya R. |

---

### Query 2.2 — Forward Trace (Recursive)

**What it does:** Starting from a given lot, recursively walks forward through `lot_genealogy` to find every descendant lot (what products this input became part of).

**Parameters:** `$1` = the lot_id to trace forward from (usually a raw material)

```sql
WITH RECURSIVE forward_trace AS (
    -- Base case: direct children of the target lot
    SELECT lg.parent_lot_id, lg.child_lot_id, lg.work_order_id,
           lg.equipment_id, lg.quantity_consumed, lg.quantity_produced,
           lg.event_time, 1 AS depth
    FROM lot_genealogy lg
    WHERE lg.parent_lot_id = $1

    UNION ALL

    -- Recursive step: children of children
    SELECT lg.parent_lot_id, lg.child_lot_id, lg.work_order_id,
           lg.equipment_id, lg.quantity_consumed, lg.quantity_produced,
           lg.event_time, ft.depth + 1
    FROM lot_genealogy lg
    JOIN forward_trace ft ON lg.parent_lot_id = ft.child_lot_id
)
SELECT
    ft.depth,
    ft.child_lot_id        AS lot_id,
    ft.parent_lot_id       AS produced_from_lot_id,
    md.name AS material_name, md.material_type,
    ml.status AS lot_status, ml.quantity, ml.unit_of_measure,
    ft.quantity_produced,
    ft.work_order_id, ft.equipment_id,
    ec.class_name AS equipment_class_name,
    wo.operator_id, op.name AS operator_name,
    ft.event_time
FROM forward_trace ft
JOIN material_lot ml ON ml.lot_id = ft.child_lot_id
JOIN material_definition md ON md.id = ml.material_definition_id
LEFT JOIN equipment_master em ON em.id = ft.equipment_id
LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
LEFT JOIN work_order wo ON wo.work_order_id = ft.work_order_id
LEFT JOIN operator op ON op.operator_id = wo.operator_id
ORDER BY ft.depth, ft.child_lot_id;
```

**How it works:** Mirror of backward trace. Instead of following `child → parent`, it follows `parent → child`. Critical for product recall: "Which finished goods contain this recalled raw material?"

**Example result (tracing RM-LI-2026-042 forward):**
| depth | lot_id | material_name | quantity_produced | equipment_id |
|---|---|---|---|---|
| 1 | WIP-SLURRY-003 | Anode Slurry | 250 | MIXER-001 |
| 2 | WIP-CELL-MODULE-017 | Cell Module | 10 | COATER-01 |
| 3 | BATTERY-PACK-A42 | Battery Pack 100kWh | 1 | ASSEMBLER-01 |
| 3 | BATTERY-PACK-A43 | Battery Pack 100kWh | 1 | ASSEMBLER-01 |

---

### Query 2.3 — Full Genealogy Nodes

**What it does:** Collects every lot in the complete family tree of a given lot — ancestors, the lot itself, and all descendants — as a flat list of nodes.

**Parameters:** `$1` = center lot_id (used 3 times: for backward CTE, forward CTE, and as the self-node)

```sql
WITH RECURSIVE
backward AS (
    SELECT lg.parent_lot_id AS lot_id FROM lot_genealogy lg WHERE lg.child_lot_id = $1
    UNION
    SELECT lg.parent_lot_id FROM lot_genealogy lg JOIN backward b ON lg.child_lot_id = b.lot_id
),
forward AS (
    SELECT lg.child_lot_id AS lot_id FROM lot_genealogy lg WHERE lg.parent_lot_id = $1
    UNION
    SELECT lg.child_lot_id FROM lot_genealogy lg JOIN forward f ON lg.parent_lot_id = f.lot_id
),
all_lot_ids AS (
    SELECT lot_id FROM backward
    UNION SELECT $1 AS lot_id        -- include the queried lot itself
    UNION SELECT lot_id FROM forward
)
SELECT
    ml.lot_id, md.name AS material_name, md.material_type,
    ml.status, ml.quantity, ml.unit_of_measure
FROM all_lot_ids a
JOIN material_lot ml ON ml.lot_id = a.lot_id
JOIN material_definition md ON md.id = ml.material_definition_id
ORDER BY ml.created_at;
```

**How it works:** Runs both a backward CTE and a forward CTE simultaneously, then takes the UNION of both plus the center lot itself to form the complete set of all related lots. The final SELECT fetches the actual lot and material data for each ID in that set.

**Example result (center = WIP-SLURRY-003):**
| lot_id | material_name | material_type | status | quantity |
|---|---|---|---|---|
| RM-LI-2026-042 | Lithium Carbonate | RAW_MATERIAL | completed | 500 |
| WIP-SLURRY-003 | Anode Slurry | WORK_IN_PROGRESS | released | 250 |
| WIP-CELL-MODULE-017 | Cell Module | WORK_IN_PROGRESS | completed | 10 |
| BATTERY-PACK-A42 | Battery Pack 100kWh | FINISHED_GOOD | shipped | 1 |

---

### Query 2.4 — Full Genealogy Edges

**What it does:** Collects all the transformation links (edges) connecting lots in the complete family tree.

**Parameters:** `$1` = center lot_id (used 3 times, same as Query 2.3)

```sql
WITH RECURSIVE
backward AS ( ... same as above ... ),
forward  AS ( ... same as above ... ),
all_lot_ids AS ( SELECT lot_id FROM backward UNION SELECT $1 UNION SELECT lot_id FROM forward )
SELECT
    lg.parent_lot_id AS source, lg.child_lot_id AS target,
    lg.work_order_id, lg.equipment_id,
    lg.quantity_consumed, lg.quantity_produced, lg.event_time
FROM lot_genealogy lg
WHERE lg.parent_lot_id IN (SELECT lot_id FROM all_lot_ids)
   OR lg.child_lot_id  IN (SELECT lot_id FROM all_lot_ids)
ORDER BY lg.event_time;
```

**How it works:** Same family tree collection as Query 2.3, but instead of returning lot details it returns all the `lot_genealogy` rows that connect any two lots in the tree.

**Example result:**
| source | target | work_order_id | equipment_id | quantity_consumed | quantity_produced |
|---|---|---|---|---|---|
| RM-LI-2026-042 | WIP-SLURRY-003 | WO-2026-0001 | MIXER-001 | 200 | 250 |
| WIP-SLURRY-003 | WIP-CELL-MODULE-017 | WO-2026-0002 | COATER-01 | 250 | 10 |
| WIP-CELL-MODULE-017 | BATTERY-PACK-A42 | WO-2026-0003 | ASSEMBLER-01 | 10 | 1 |

> Nodes + Edges together give a complete directed graph that can be visualized as a production lineage diagram.

---

### Query 2.5 — Equipment Process History (Statistics)

**What it does:** For a given lot, finds the time window each machine was processing it (based on traceability log), then aggregates the sensor telemetry recorded during that window.

**Parameters:** `$1` = lot_id

```sql
WITH lot_processing_window AS (
    -- Step 1: Find first and last event on each machine for this lot
    SELECT tl.equipment_id,
           MIN(tl.event_time) AS process_start,
           MAX(tl.event_time) AS process_end
    FROM traceability_log tl
    WHERE tl.material_lot_id = $1
    GROUP BY tl.equipment_id
)
-- Step 2: Get all telemetry readings in that window, compute stats
SELECT
    et.equipment_id, ec.class_name AS equipment_class_name,
    et.parameter_name, et.unit_of_measure,
    MIN(et.value) AS min_value, MAX(et.value) AS max_value,
    AVG(et.value) AS avg_value, COUNT(et.value) AS reading_count,
    lpw.process_start, lpw.process_end
FROM lot_processing_window lpw
JOIN equipment_telemetry et
    ON et.equipment_id = lpw.equipment_id
   AND et.recorded_at >= lpw.process_start
   AND et.recorded_at <= lpw.process_end
JOIN equipment_master em ON em.id = et.equipment_id
JOIN equipment_class ec ON ec.id = em.equipment_class_id
GROUP BY et.equipment_id, ec.class_name, et.parameter_name,
         et.unit_of_measure, lpw.process_start, lpw.process_end
ORDER BY et.equipment_id, et.parameter_name;
```

**How it works:**
1. **CTE `lot_processing_window`:** Scans `traceability_log` to find the earliest and latest event on each machine for this lot. This gives us the exact time window the machine was working on this lot.
2. **Main query:** Looks up all `equipment_telemetry` readings from each machine within that window and computes MIN, MAX, AVG, COUNT per sensor parameter.

**Example result (lot = WIP-SLURRY-003):**
| equipment_id | parameter_name | min_value | max_value | avg_value | reading_count | process_start | process_end |
|---|---|---|---|---|---|---|---|
| MIXER-001 | temperature | 72.10 | 91.30 | 78.50 | 432 | 08:05:00 | 09:48:00 |
| MIXER-001 | rpm | 2100 | 2600 | 2400 | 432 | 08:05:00 | 09:48:00 |

---

### Query 2.6 — Equipment Process History (Raw Readings)

**What it does:** Same lot-processing-window logic as Query 2.5, but returns every individual sensor reading instead of aggregated stats.

**Parameters:** `$1` = lot_id

**How it works:** Same CTE, same JOIN, but no GROUP BY — returns one row per telemetry reading.

**Use case:** When you need the full time series (for plotting a chart or investigating a specific spike).

**Example result:**
| equipment_id | parameter_name | value | unit_of_measure | recorded_at |
|---|---|---|---|---|
| MIXER-001 | temperature | 72.10 | °C | 08:05:05 |
| MIXER-001 | temperature | 73.40 | °C | 08:05:10 |
| MIXER-001 | temperature | 74.90 | °C | 08:05:15 |
| ... | ... | ... | ... | ... |

---

## Section 3 — Analytics Queries (S014)

---

### Query 3.1 — Machine Performance KPIs

**What it does:** Computes production KPIs for a single machine over a time window.

**Parameters:** `$1, $2` = to/from for utilization calc · `$3, $4` = from/to for WO filter · `$5` = equipment_id

```sql
SELECT
    em.id AS equipment_id, ec.class_name AS equipment_class_name, em.operational_status,
    COUNT(DISTINCT wo.work_order_id) AS total_work_orders,
    COALESCE(SUM(wo.actual_quantity), 0) AS total_units_produced,
    COALESCE(AVG(EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0), 0) AS avg_cycle_time_hours,
    COALESCE(SUM(EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0), 0) AS total_active_hours,
    CASE WHEN EXTRACT(EPOCH FROM ($1::timestamptz - $2::timestamptz)) > 0 THEN
        100.0 * COALESCE(SUM(EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0), 0)
        / (EXTRACT(EPOCH FROM ($1::timestamptz - $2::timestamptz)) / 3600.0)
    ELSE 0 END AS utilization_pct
FROM equipment_master em
LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
LEFT JOIN work_order wo ON wo.equipment_id = em.id
    AND wo.actual_start >= $3::timestamptz AND wo.actual_end <= $4::timestamptz
    AND wo.status = 'completed'
WHERE em.id = $5
GROUP BY em.id, ec.class_name, em.operational_status;
```

**How the utilization formula works:**
- `EXTRACT(EPOCH FROM (to - from))` converts a time interval to seconds
- Divide by 3600 to get hours
- Total active hours = sum of all (actual_end - actual_start) for completed WOs
- Utilization % = (total active hours / total window hours) × 100

**Example:** Machine ran 5.55 hours over a 24-hour window → utilization = 5.55/24 × 100 = **23.1%**

---

### Query 3.2 — Machine Event Summary

**What it does:** Counts how often each event type was logged for a machine in a time window.

**Parameters:** `$1` = equipment_id · `$2` = from_time · `$3` = to_time

```sql
SELECT event_type, COUNT(*) AS count
FROM traceability_log
WHERE equipment_id = $1
  AND event_time >= $2::timestamptz AND event_time <= $3::timestamptz
GROUP BY event_type
ORDER BY count DESC;
```

**How it works:** Simple GROUP BY on `event_type` from the time-partitioned `traceability_log` hypertable. TimescaleDB makes this fast even for millions of events.

**Example result:**
| event_type | count |
|---|---|
| start | 12 |
| qc_pass | 10 |
| stop | 12 |
| error_jam | 2 |

---

### Query 3.3 — Compare Machine Performance

**What it does:** Runs the same KPI calculation as Query 3.1, but for multiple machines at once, plus adds an error event count per machine.

**Parameters:** `$1, $2` = to/from for utilization · `$3, $4` = from/to for WO filter · `$5, $6` = from/to for error filter · `$7, $8, ...` = equipment IDs (dynamic IN clause)

```sql
SELECT
    em.id AS equipment_id, ec.class_name, em.operational_status,
    COUNT(DISTINCT wo.work_order_id) AS total_work_orders,
    COALESCE(SUM(wo.actual_quantity), 0) AS total_units_produced,
    COALESCE(AVG(EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0), 0) AS avg_cycle_time_hours,
    CASE WHEN EXTRACT(EPOCH FROM ($1::timestamptz - $2::timestamptz)) > 0 THEN
        100.0 * COALESCE(SUM(...), 0) / (EXTRACT(EPOCH FROM ($1::timestamptz - $2::timestamptz)) / 3600.0)
    ELSE 0 END AS utilization_pct,
    COALESCE(err.error_count, 0) AS error_count
FROM equipment_master em
LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
LEFT JOIN work_order wo ON wo.equipment_id = em.id
    AND wo.actual_start >= $3::timestamptz AND wo.actual_end <= $4::timestamptz
    AND wo.status = 'completed'
LEFT JOIN (
    SELECT equipment_id, COUNT(*) AS error_count FROM traceability_log
    WHERE event_time >= $5::timestamptz AND event_time <= $6::timestamptz
      AND event_type LIKE 'error%'
    GROUP BY equipment_id
) err ON err.equipment_id = em.id
WHERE em.id IN ($7, $8, ...)  -- built dynamically
GROUP BY em.id, ec.class_name, em.operational_status, err.error_count
ORDER BY total_units_produced DESC;
```

**How it works:** Same KPI aggregation as Query 3.1, but without the `WHERE em.id = $5` restriction and with an added subquery that counts error events per machine. The `IN (...)` clause is built dynamically in Go based on how many equipment IDs are passed.

**Example result:**
| equipment_id | total_units_produced | utilization_pct | error_count |
|---|---|---|---|
| OVEN-001 | 480 | 62.5 | 0 |
| MIXER-001 | 285 | 23.1 | 2 |
| ROBOT-001 | 180 | 18.8 | 5 |

---

### Query 3.4 — Production Trends

**What it does:** Buckets completed work order output by time period (day/week/month) and machine — the data series for output charts.

**Parameters:** `$1` = from_time · `$2` = to_time · `$3` = equipment_id (optional, added dynamically)  
**Granularity:** Injected via `fmt.Sprintf` in Go — `'day'`, `'week'`, or `'month'`

```sql
SELECT
    DATE_TRUNC('day', actual_end) AS period,
    equipment_id,
    COALESCE(SUM(actual_quantity), 0) AS units_produced,
    COUNT(*) AS work_orders_completed,
    COALESCE(AVG(EXTRACT(EPOCH FROM (actual_end - actual_start)) / 3600.0), 0) AS avg_cycle_time_hours
FROM work_order
WHERE actual_end >= $1::timestamptz AND actual_end <= $2::timestamptz
  AND status = 'completed'
  -- AND equipment_id = $3  (added when equipment_id filter is provided)
GROUP BY DATE_TRUNC('day', actual_end), equipment_id
ORDER BY period, equipment_id;
```

**How it works:** `DATE_TRUNC('day', actual_end)` rounds each WO's end timestamp down to the start of its day (or week/month). Grouping by this truncated timestamp and equipment_id gives one row per (period, machine) combination. The result is chronologically ordered and ready to plot.

**Example result (daily, all machines, Feb 24–26):**
| period | equipment_id | units_produced | work_orders_completed | avg_cycle_time_hours |
|---|---|---|---|---|
| 2026-02-24 | MIXER-001 | 247.5 | 3 | 1.83 |
| 2026-02-24 | OVEN-001 | 480.0 | 2 | 3.20 |
| 2026-02-25 | MIXER-001 | 250.0 | 3 | 1.90 |
| 2026-02-25 | ROBOT-001 | 180.0 | 4 | 0.75 |

---

### Query 3.5a — Dashboard: Production KPIs

**What it does:** Single-row summary of total production output and machine activity for a time window.

**Parameters:** `$1` = from_time · `$2` = to_time

```sql
SELECT
    COALESCE(SUM(actual_quantity), 0) AS total_units_produced,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) AS total_work_orders_completed,
    COUNT(DISTINCT CASE WHEN status IN ('completed','in_progress') THEN equipment_id END) AS active_equipment_count
FROM work_order
WHERE created_at >= $1::timestamptz AND created_at <= $2::timestamptz;
```

**How it works:** Uses conditional aggregation (`CASE WHEN` inside `COUNT`) to calculate multiple KPIs in a single pass over the table. `COUNT(DISTINCT ...)` ensures each machine is counted once even if it had many work orders.

**Example result:**
| total_units_produced | total_work_orders_completed | active_equipment_count |
|---|---|---|
| 4820.5 | 18 | 5 |

---

### Query 3.5b — Dashboard: Quality KPIs

**What it does:** Counts lots by QC outcome (released vs. quarantined) to compute quality pass rate.

**Parameters:** `$1` = from_time · `$2` = to_time

```sql
SELECT
    COUNT(CASE WHEN status = 'released'    THEN 1 END) AS lots_released,
    COUNT(CASE WHEN status = 'quarantined' THEN 1 END) AS lots_quarantined
FROM material_lot
WHERE created_at >= $1::timestamptz AND created_at <= $2::timestamptz;
```

**How it works:** Same conditional aggregation pattern. The quality rate is then calculated in Go: `released / (released + quarantined) × 100`.

**Example result:**
| lots_released | lots_quarantined |
|---|---|
| 12 | 2 |

→ Quality rate = 12 / (12 + 2) × 100 = **85.7%**

---

### Query 3.5c — Dashboard: Error Events

**What it does:** Counts all error-type events across all machines in the window.

**Parameters:** `$1` = from_time · `$2` = to_time

```sql
SELECT COUNT(*) AS total_errors
FROM traceability_log
WHERE event_time >= $1::timestamptz AND event_time <= $2::timestamptz
  AND event_type LIKE 'error%';
```

**How it works:** Simple COUNT with a LIKE filter. The `LIKE 'error%'` matches any event type starting with "error" (e.g. `error_jam`, `error_sensor`, `error_overtemp`).

**Example result:**
| total_errors |
|---|
| 3 |

---

### Query 3.5d — Dashboard: Top 5 Performers

**What it does:** Ranks all machines by units produced in the window and returns the top 5.

**Parameters:** `$1, $2` = to/from for utilization calc · `$3, $4` = from/to for WO filter

```sql
SELECT
    em.id AS equipment_id, ec.class_name, em.operational_status,
    COUNT(DISTINCT wo.work_order_id) AS total_work_orders,
    COALESCE(SUM(wo.actual_quantity), 0) AS total_units_produced,
    COALESCE(AVG(EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0), 0) AS avg_cycle_time_hours,
    CASE WHEN EXTRACT(EPOCH FROM ($1::timestamptz - $2::timestamptz)) > 0 THEN
        100.0 * COALESCE(SUM(...), 0)
        / (EXTRACT(EPOCH FROM ($1::timestamptz - $2::timestamptz)) / 3600.0)
    ELSE 0 END AS utilization_pct
FROM equipment_master em
LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
LEFT JOIN work_order wo ON wo.equipment_id = em.id
    AND wo.actual_start >= $3::timestamptz AND wo.actual_end <= $4::timestamptz
    AND wo.status = 'completed'
GROUP BY em.id, ec.class_name, em.operational_status
ORDER BY total_units_produced DESC
LIMIT 5;
```

**How it works:** Same KPI calculation as Query 3.1 but across ALL machines (no WHERE on equipment_id), sorted by `total_units_produced DESC`, and limited to top 5 rows.

**Example result:**
| equipment_id | total_units_produced | avg_cycle_time_hours | utilization_pct |
|---|---|---|---|
| OVEN-001 | 1920 | 3.20 | 80.0 |
| ROBOT-001 | 1440 | 0.75 | 60.0 |
| MIXER-001 | 2400 | 1.85 | 39.0 |

---

## Quick Reference Summary

| Query | Input | Output |
|---|---|---|
| 1.1 Get Lot Detail | lot_id | Full lot + material info |
| 1.2 List Lots | status, material_type (optional) | List of lot summaries |
| 1.3 Get Work Order | work_order_id | WO + equipment class + operator |
| 1.4 Get WO Input Lots | work_order_id | All consumed parent lots |
| 2.1 Backward Trace | lot_id | All ancestor lots (recursive, with depth) |
| 2.2 Forward Trace | lot_id | All descendant lots (recursive, with depth) |
| 2.3 Full Genealogy Nodes | lot_id | All lots in the family tree |
| 2.4 Full Genealogy Edges | lot_id | All transformation links in the family tree |
| 2.5 Equipment History Stats | lot_id | Sensor min/max/avg per parameter per machine |
| 2.6 Equipment History Readings | lot_id | All raw sensor readings per machine |
| 3.1 Machine Performance | equipment_id, from, to | KPIs for one machine |
| 3.2 Machine Event Summary | equipment_id, from, to | Event counts per type |
| 3.3 Compare Machines | [equipment_ids], from, to | KPIs for many machines side-by-side |
| 3.4 Production Trends | from, to, granularity | Time-bucketed output per machine |
| 3.5a Dashboard Production | from, to | Total units, WO count, active machines |
| 3.5b Dashboard Quality | from, to | Released vs. quarantined lots |
| 3.5c Dashboard Errors | from, to | Total error event count |
| 3.5d Dashboard Top 5 | from, to | Top 5 machines by output |
