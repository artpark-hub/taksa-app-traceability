-- ============================================================================
--  Core Traceability SQL Queries
-- ============================================================================

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
