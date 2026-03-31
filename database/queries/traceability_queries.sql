-- ============================================================================
--  Core Traceability SQL Queries
-- ============================================================================

-- Get full details of a single material lot including its material definition info.
-- Output: lot_id, status, quantity, UOM, timestamps, material name/type/description.
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


-- List all material lots with optional filters on status and material type.
-- Output: summary row per lot (lot_id, status, quantity, UOM, created date, material name/type).
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


-- Fetch full details of a single work order including the assigned equipment class and operator.
-- Output: all WO fields + equipment_class_name, operator_name, operator_shift.
SELECT
    wo.work_order_id,
    wo.description,
    wo.status,
    wo.planned_quantity,
    wo.actual_quantity,
    wo.unit_of_measure,
    wo.planned_start,
    wo.planned_end,
    wo.actual_start,
    wo.actual_end,
    wo.output_lot_id,
    wo.equipment_id,
    ec.class_name       AS equipment_class_name,
    wo.operator_id,
    op.name             AS operator_name,
    op.shift            AS operator_shift
FROM work_order wo
LEFT JOIN equipment_master em ON em.id = wo.equipment_id
LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
LEFT JOIN operator op ON op.operator_id = wo.operator_id
WHERE wo.work_order_id = $1;


-- Fetch all input (parent) lots consumed by a specific work order via genealogy links.
-- Output: parent lot_id, material name/type, quantity consumed, and UOM.
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


-- Recursively trace all upstream (ancestor) lots that contributed to producing a given lot.
-- Output: one row per ancestor lot with depth level, material info, equipment, operator, and event time.
WITH RECURSIVE backward_trace AS (
    SELECT
        lg.parent_lot_id,
        lg.child_lot_id,
        lg.work_order_id,
        lg.equipment_id,
        lg.quantity_consumed,
        lg.quantity_produced,
        lg.event_time,
        1 AS depth
    FROM lot_genealogy lg
    WHERE lg.child_lot_id = $1

    UNION ALL

    SELECT
        lg.parent_lot_id,
        lg.child_lot_id,
        lg.work_order_id,
        lg.equipment_id,
        lg.quantity_consumed,
        lg.quantity_produced,
        lg.event_time,
        bt.depth + 1
    FROM lot_genealogy lg
    JOIN backward_trace bt ON lg.child_lot_id = bt.parent_lot_id
)
SELECT
    bt.depth,
    bt.parent_lot_id       AS lot_id,
    bt.child_lot_id        AS consumed_by_lot_id,
    md.name                AS material_name,
    md.material_type,
    ml.status              AS lot_status,
    ml.quantity,
    ml.unit_of_measure,
    bt.quantity_consumed,
    bt.work_order_id,
    bt.equipment_id,
    ec.class_name          AS equipment_class_name,
    wo.operator_id,
    op.name                AS operator_name,
    bt.event_time
FROM backward_trace bt
JOIN material_lot ml ON ml.lot_id = bt.parent_lot_id
JOIN material_definition md ON md.id = ml.material_definition_id
LEFT JOIN equipment_master em ON em.id = bt.equipment_id
LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
LEFT JOIN work_order wo ON wo.work_order_id = bt.work_order_id
LEFT JOIN operator op ON op.operator_id = wo.operator_id
ORDER BY bt.depth, bt.parent_lot_id;


-- Recursively trace all downstream (descendant) lots produced from a given lot.
-- Output: one row per descendant lot with depth level, material info, equipment, operator, and event time.
WITH RECURSIVE forward_trace AS (
    SELECT
        lg.parent_lot_id,
        lg.child_lot_id,
        lg.work_order_id,
        lg.equipment_id,
        lg.quantity_consumed,
        lg.quantity_produced,
        lg.event_time,
        1 AS depth
    FROM lot_genealogy lg
    WHERE lg.parent_lot_id = $1

    UNION ALL

    SELECT
        lg.parent_lot_id,
        lg.child_lot_id,
        lg.work_order_id,
        lg.equipment_id,
        lg.quantity_consumed,
        lg.quantity_produced,
        lg.event_time,
        ft.depth + 1
    FROM lot_genealogy lg
    JOIN forward_trace ft ON lg.parent_lot_id = ft.child_lot_id
)
SELECT
    ft.depth,
    ft.child_lot_id        AS lot_id,
    ft.parent_lot_id       AS produced_from_lot_id,
    md.name                AS material_name,
    md.material_type,
    ml.status              AS lot_status,
    ml.quantity,
    ml.unit_of_measure,
    ft.quantity_produced,
    ft.work_order_id,
    ft.equipment_id,
    ec.class_name          AS equipment_class_name,
    wo.operator_id,
    op.name                AS operator_name,
    ft.event_time
FROM forward_trace ft
JOIN material_lot ml ON ml.lot_id = ft.child_lot_id
JOIN material_definition md ON md.id = ml.material_definition_id
LEFT JOIN equipment_master em ON em.id = ft.equipment_id
LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
LEFT JOIN work_order wo ON wo.work_order_id = ft.work_order_id
LEFT JOIN operator op ON op.operator_id = wo.operator_id
ORDER BY ft.depth, ft.child_lot_id;


-- Collect ALL unique lots in the complete genealogy tree of a given lot (ancestors + self + descendants).
-- Output: one node row per lot — lot_id, material name/type, status, quantity, UOM.
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
    UNION SELECT $1 AS lot_id
    UNION SELECT lot_id FROM forward
)
SELECT
    ml.lot_id,
    md.name            AS material_name,
    md.material_type,
    ml.status,
    ml.quantity,
    ml.unit_of_measure
FROM all_lot_ids a
JOIN material_lot ml ON ml.lot_id = a.lot_id
JOIN material_definition md ON md.id = ml.material_definition_id
ORDER BY ml.created_at;


-- Collect ALL genealogy links (edges) that connect lots in the complete family tree of a given lot.
-- Output: one edge row per parent→child transformation with work order, equipment, and quantities.
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
    UNION SELECT $1 AS lot_id
    UNION SELECT lot_id FROM forward
)
SELECT
    lg.parent_lot_id   AS source,
    lg.child_lot_id    AS target,
    lg.work_order_id,
    lg.equipment_id,
    lg.quantity_consumed,
    lg.quantity_produced,
    lg.event_time
FROM lot_genealogy lg
WHERE lg.parent_lot_id IN (SELECT lot_id FROM all_lot_ids)
   OR lg.child_lot_id IN (SELECT lot_id FROM all_lot_ids)
ORDER BY lg.event_time;


-- Aggregate sensor/telemetry readings per equipment during the time window a lot was being processed.
-- Output: per-parameter statistics (min, max, avg, count) for each machine that touched the lot.
WITH lot_processing_window AS (
    SELECT
        tl.equipment_id,
        MIN(tl.event_time) AS process_start,
        MAX(tl.event_time) AS process_end
    FROM traceability_log tl
    WHERE tl.material_lot_id = $1
    GROUP BY tl.equipment_id
)
SELECT
    et.equipment_id,
    ec.class_name          AS equipment_class_name,
    et.parameter_name,
    et.unit_of_measure,
    MIN(et.value)          AS min_value,
    MAX(et.value)          AS max_value,
    AVG(et.value)          AS avg_value,
    COUNT(et.value)        AS reading_count,
    lpw.process_start,
    lpw.process_end
FROM lot_processing_window lpw
JOIN equipment_telemetry et
    ON et.equipment_id = lpw.equipment_id
   AND et.recorded_at >= lpw.process_start
   AND et.recorded_at <= lpw.process_end
JOIN equipment_master em ON em.id = et.equipment_id
JOIN equipment_class ec ON ec.id = em.equipment_class_id
GROUP BY et.equipment_id, ec.class_name, et.parameter_name, et.unit_of_measure, lpw.process_start, lpw.process_end
ORDER BY et.equipment_id, et.parameter_name;


-- Fetch individual raw telemetry readings for every sensor on machines that processed a given lot.
-- Output: one row per reading — equipment_id, parameter name, value, UOM, and timestamp.
WITH lot_processing_window AS (
    SELECT
        tl.equipment_id,
        MIN(tl.event_time) AS process_start,
        MAX(tl.event_time) AS process_end
    FROM traceability_log tl
    WHERE tl.material_lot_id = $1
    GROUP BY tl.equipment_id
)
SELECT
    et.equipment_id,
    et.parameter_name,
    et.value,
    et.unit_of_measure,
    et.recorded_at
FROM lot_processing_window lpw
JOIN equipment_telemetry et
    ON et.equipment_id = lpw.equipment_id
   AND et.recorded_at >= lpw.process_start
   AND et.recorded_at <= lpw.process_end
ORDER BY et.equipment_id, et.parameter_name, et.recorded_at;


-- Historical Analytics Queries

-- ---- 1. Machine Performance ------------------------------------------------
-- Compute production KPIs for a single machine over a given time window.
-- Output: total work orders, units produced, average cycle time, total active hours, and utilization % for the machine.
-- $1 = to_time, $2 = from_time (x2 for utilization calc), $3 = from_time, $4 = to_time (WO filter), $5 = equipment_id
SELECT
    em.id                                           AS equipment_id,
    ec.class_name                                   AS equipment_class_name,
    em.operational_status,
    COUNT(DISTINCT wo.work_order_id)                AS total_work_orders,
    COALESCE(SUM(wo.actual_quantity), 0)            AS total_units_produced,
    COALESCE(AVG(
        EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0
    ), 0)                                           AS avg_cycle_time_hours,
    COALESCE(SUM(
        EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0
    ), 0)                                           AS total_active_hours,
    CASE
        WHEN EXTRACT(EPOCH FROM ($1::timestamptz - $2::timestamptz)) > 0 THEN
            100.0 * COALESCE(SUM(
                EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0
            ), 0)
            / (EXTRACT(EPOCH FROM ($1::timestamptz - $2::timestamptz)) / 3600.0)
        ELSE 0
    END                                             AS utilization_pct
FROM equipment_master em
LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
LEFT JOIN work_order wo
    ON wo.equipment_id = em.id
   AND wo.actual_start >= $3::timestamptz
   AND wo.actual_end   <= $4::timestamptz
   AND wo.status = 'completed'
WHERE em.id = $5
GROUP BY em.id, ec.class_name, em.operational_status;


-- ---- 2. Machine Event Summary (for a single machine in a time window) ------
-- Count how many times each event type occurred on a machine within a time range.
-- Output: event_type and its count, ordered from most to least frequent.
-- $1 = equipment_id, $2 = from_time, $3 = to_time
SELECT
    event_type,
    COUNT(*) AS count
FROM traceability_log
WHERE equipment_id  = $1
  AND event_time   >= $2::timestamptz
  AND event_time   <= $3::timestamptz
GROUP BY event_type
ORDER BY count DESC;


-- ---- 3. Compare Machine Performance ----------------------------------------
-- Side-by-side KPI comparison for a list of machines over the same time window.
-- Output: one row per machine with units produced, cycle time, utilization %, and error event count.
-- $1 = to_time, $2 = from_time (x2 for utilization), $3 = from_time, $4 = to_time (WO filter)
-- $5 = from_time, $6 = to_time (error filter), then IN($7, $8, ...) for equipment IDs
SELECT
    em.id                                               AS equipment_id,
    ec.class_name                                       AS equipment_class_name,
    em.operational_status,
    COUNT(DISTINCT wo.work_order_id)                    AS total_work_orders,
    COALESCE(SUM(wo.actual_quantity), 0)                AS total_units_produced,
    COALESCE(AVG(
        EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0
    ), 0)                                               AS avg_cycle_time_hours,
    CASE
        WHEN EXTRACT(EPOCH FROM ($1::timestamptz - $2::timestamptz)) > 0 THEN
            100.0 * COALESCE(SUM(
                EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0
            ), 0)
            / (EXTRACT(EPOCH FROM ($1::timestamptz - $2::timestamptz)) / 3600.0)
        ELSE 0
    END                                                 AS utilization_pct,
    COALESCE(err.error_count, 0)                        AS error_count
FROM equipment_master em
LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
LEFT JOIN work_order wo
    ON wo.equipment_id = em.id
   AND wo.actual_start >= $3::timestamptz
   AND wo.actual_end   <= $4::timestamptz
   AND wo.status = 'completed'
LEFT JOIN (
    SELECT equipment_id, COUNT(*) AS error_count
    FROM traceability_log
    WHERE event_time >= $5::timestamptz
      AND event_time <= $6::timestamptz
      AND event_type LIKE 'error%'
    GROUP BY equipment_id
) err ON err.equipment_id = em.id
WHERE em.id IN ($7 /*, $8, $9 … dynamic */)
GROUP BY em.id, ec.class_name, em.operational_status, err.error_count
ORDER BY total_units_produced DESC;


-- ---- 4. Production Trends --------------------------------------------------
-- Bucket completed work order output by time period (daily / weekly / monthly) and machine.
-- Output: one row per (period, equipment) with units produced, WO count, and average cycle time.
-- granularity = 'day' | 'week' | 'month'  (injected via fmt.Sprintf in Go)
-- $1 = from_time, $2 = to_time  (optional $3 = equipment_id)
SELECT
    DATE_TRUNC('day', actual_end)                   AS period,
    equipment_id,
    COALESCE(SUM(actual_quantity), 0)               AS units_produced,
    COUNT(*)                                        AS work_orders_completed,
    COALESCE(AVG(
        EXTRACT(EPOCH FROM (actual_end - actual_start)) / 3600.0
    ), 0)                                           AS avg_cycle_time_hours
FROM work_order
WHERE actual_end >= $1::timestamptz
  AND actual_end <= $2::timestamptz
  AND status = 'completed'
  /* AND equipment_id = $3   ← added dynamically when equipment_id is provided */
GROUP BY DATE_TRUNC('day', actual_end), equipment_id
ORDER BY period, equipment_id;


-- ---- 5a. Dashboard Summary – Production KPIs --------------------------------
-- High-level production counters for all work orders created in the time window.
-- Output: total units produced, number of completed WOs, and count of distinct active machines.
-- $1 = from_time, $2 = to_time
SELECT
    COALESCE(SUM(actual_quantity), 0)                                       AS total_units_produced,
    COUNT(CASE WHEN status = 'completed' THEN 1 END)                        AS total_work_orders_completed,
    COUNT(DISTINCT CASE WHEN status IN ('completed','in_progress')
                         THEN equipment_id END)                             AS active_equipment_count
FROM work_order
WHERE created_at >= $1::timestamptz
  AND created_at <= $2::timestamptz;


-- ---- 5b. Dashboard Summary – Quality KPIs -----------------------------------
-- Count lots by final disposition to calculate overall quality pass rate.
-- Output: lots_released (pass) and lots_quarantined (fail) created in the window.
-- $1 = from_time, $2 = to_time
SELECT
    COUNT(CASE WHEN status = 'released'    THEN 1 END)  AS lots_released,
    COUNT(CASE WHEN status = 'quarantined' THEN 1 END)  AS lots_quarantined
FROM material_lot
WHERE created_at >= $1::timestamptz
  AND created_at <= $2::timestamptz;


-- ---- 5c. Dashboard Summary – Error Events -----------------------------------
-- Count total error-type events logged across all machines in the time window.
-- Output: single count (total_errors) used as an anomaly indicator on the dashboard.
-- $1 = from_time, $2 = to_time
SELECT COUNT(*) AS total_errors
FROM traceability_log
WHERE event_time >= $1::timestamptz
  AND event_time <= $2::timestamptz
  AND event_type LIKE 'error%';


-- ---- 5d. Dashboard Summary – Top 5 Performers ------------------------------
-- Rank all machines by total units produced in the window and return the top 5.
-- Output: per-machine KPIs (units, cycle time, utilization %) sorted best-first, limited to 5 rows.
-- $1 = to_time, $2 = from_time (x2 utilization), $3 = from_time, $4 = to_time (WO filter)
SELECT
    em.id                                               AS equipment_id,
    ec.class_name                                       AS equipment_class_name,
    em.operational_status,
    COUNT(DISTINCT wo.work_order_id)                    AS total_work_orders,
    COALESCE(SUM(wo.actual_quantity), 0)                AS total_units_produced,
    COALESCE(AVG(
        EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0
    ), 0)                                               AS avg_cycle_time_hours,
    CASE
        WHEN EXTRACT(EPOCH FROM ($1::timestamptz - $2::timestamptz)) > 0 THEN
            100.0 * COALESCE(SUM(
                EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0
            ), 0)
            / (EXTRACT(EPOCH FROM ($1::timestamptz - $2::timestamptz)) / 3600.0)
        ELSE 0
    END                                                 AS utilization_pct
FROM equipment_master em
LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
LEFT JOIN work_order wo
    ON wo.equipment_id = em.id
   AND wo.actual_start >= $3::timestamptz
   AND wo.actual_end   <= $4::timestamptz
   AND wo.status = 'completed'
GROUP BY em.id, ec.class_name, em.operational_status
ORDER BY total_units_produced DESC
LIMIT 5;
