-- ============================================================================
-- PHASE 2: Traceability & Genealogy Engine — Schema Enhancements
-- ============================================================================
-- Migration: 002_traceability_features.sql
-- Description: Adds 6 new tables for lot management, genealogy, work orders,
--              operators, and equipment telemetry.
-- Depends on: schema.sql (001)
-- Idempotent: YES — uses CREATE TABLE IF NOT EXISTS
-- ============================================================================

-- ============================================================================
-- TABLE 1: material_definition
-- Purpose: Catalog of material types (e.g., "Epoxy Resin", "PCB Board")
-- Used by: Feature 1 (Lot Lifecycle), Feature 4-6 (Tracing)
-- ============================================================================
CREATE TABLE IF NOT EXISTS material_definition (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,           -- e.g., 'Epoxy Resin Grade A'
    material_type   VARCHAR(50) NOT NULL             -- 'RAW_MATERIAL', 'WORK_IN_PROGRESS', 'FINISHED_GOOD', 'PACKAGING'
                    CHECK (material_type IN ('RAW_MATERIAL', 'WORK_IN_PROGRESS', 'FINISHED_GOOD', 'PACKAGING')),
    unit_of_measure VARCHAR(50),                     -- default UoM (e.g., 'kg', 'liters', 'units')
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- TABLE 2: material_lot
-- Purpose: Individual lot instances with status tracking
-- Used by: Feature 1 (Lot Lifecycle), Feature 3-6 (Genealogy & Tracing)
-- ============================================================================
CREATE TABLE IF NOT EXISTS material_lot (
    lot_id                VARCHAR(100) PRIMARY KEY,      -- e.g., 'RM-001', 'WIP-003', 'FG-010'
    material_definition_id INTEGER NOT NULL REFERENCES material_definition(id) ON DELETE RESTRICT,
    status                VARCHAR(50) NOT NULL DEFAULT 'available'
                          CHECK (status IN ('available', 'in_process', 'completed', 'quarantined', 'released', 'shipped', 'scrapped')),
    quantity              NUMERIC(15,4),                  -- quantity in this lot
    unit_of_measure       VARCHAR(50),                    -- UoM for this specific lot
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index: Filter lots by status and material type
CREATE INDEX IF NOT EXISTS idx_material_lot_status ON material_lot(status);
CREATE INDEX IF NOT EXISTS idx_material_lot_material_def ON material_lot(material_definition_id);

-- ============================================================================
-- TABLE 3: operator
-- Purpose: Operator master record with shift tracking
-- Used by: Feature 2 (Work Orders), Feature 7 (Process History)
-- ============================================================================
CREATE TABLE IF NOT EXISTS operator (
    operator_id     VARCHAR(100) PRIMARY KEY,         -- e.g., 'OP-KAVYA-01'
    name            VARCHAR(255) NOT NULL,
    role            VARCHAR(100),                     -- e.g., 'Machine Operator', 'Supervisor'
    shift           VARCHAR(50),                      -- e.g., 'DAY_SHIFT', 'NIGHT_SHIFT', 'MORNING_SHIFT'
    status          VARCHAR(50) NOT NULL DEFAULT 'active'
                    CHECK (status IN ('active', 'inactive')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- TABLE 4: work_order
-- Purpose: Production instructions — consumes inputs, produces outputs
-- Used by: Feature 2 (Work Order Execution), Feature 3 (Genealogy)
-- ============================================================================
CREATE TABLE IF NOT EXISTS work_order (
    work_order_id   VARCHAR(100) PRIMARY KEY,         -- e.g., 'WO-2024-001'
    description     TEXT,
    status          VARCHAR(50) NOT NULL DEFAULT 'planned'
                    CHECK (status IN ('planned', 'in_progress', 'completed', 'closed')),
    equipment_id    VARCHAR(50) REFERENCES equipment_master(id) ON DELETE SET NULL,    -- primary equipment
    operator_id     VARCHAR(100) REFERENCES operator(operator_id) ON DELETE SET NULL,  -- assigned operator
    output_lot_id   VARCHAR(100),                     -- the lot this WO is producing (set on creation/start)
    planned_quantity NUMERIC(15,4),
    actual_quantity  NUMERIC(15,4),
    unit_of_measure VARCHAR(50),
    planned_start   TIMESTAMPTZ,
    planned_end     TIMESTAMPTZ,
    actual_start    TIMESTAMPTZ,
    actual_end      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index: Filter work orders by status and equipment
CREATE INDEX IF NOT EXISTS idx_work_order_status ON work_order(status);
CREATE INDEX IF NOT EXISTS idx_work_order_equipment ON work_order(equipment_id);

-- ============================================================================
-- TABLE 5: lot_genealogy
-- Purpose: Parent-child links between lots (THE CORE OF TRACEABILITY)
-- Used by: Feature 3 (Registration), Feature 4-6 (Backward/Forward/Full Trace)
--
-- Data model:
--   parent_lot_id (input consumed) --> child_lot_id (output produced)
--   Example: RM-001 (parent) --> WIP-001 (child) via WO-2024-001
--
-- This table forms a DIRECTED ACYCLIC GRAPH (DAG) of lot relationships.
-- Recursive CTEs walk this graph for tracing.
-- ============================================================================
CREATE TABLE IF NOT EXISTS lot_genealogy (
    id                  SERIAL PRIMARY KEY,
    parent_lot_id       VARCHAR(100) NOT NULL REFERENCES material_lot(lot_id) ON DELETE RESTRICT,
    child_lot_id        VARCHAR(100) NOT NULL REFERENCES material_lot(lot_id) ON DELETE RESTRICT,
    work_order_id       VARCHAR(100) REFERENCES work_order(work_order_id) ON DELETE SET NULL,
    equipment_id        VARCHAR(50) REFERENCES equipment_master(id) ON DELETE SET NULL,
    quantity_consumed   NUMERIC(15,4),               -- how much of parent lot was consumed
    quantity_produced   NUMERIC(15,4),               -- how much of child lot was produced
    event_time          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Prevent duplicate links
    CONSTRAINT uq_genealogy_link UNIQUE (parent_lot_id, child_lot_id, work_order_id)
);

-- Critical indexes for recursive CTE performance
CREATE INDEX IF NOT EXISTS idx_genealogy_parent ON lot_genealogy(parent_lot_id);
CREATE INDEX IF NOT EXISTS idx_genealogy_child ON lot_genealogy(child_lot_id);
CREATE INDEX IF NOT EXISTS idx_genealogy_work_order ON lot_genealogy(work_order_id);

-- ============================================================================
-- TABLE 6: equipment_telemetry
-- Purpose: Real-time time-series process parameter readings
-- Used by: Feature 7 (Equipment Process History for a Lot)
--
-- This is a TimescaleDB HYPERTABLE — optimized for:
--   - High-frequency INSERTs (telemetry from PLCs/SCADA)
--   - Time-range queries with time_bucket aggregation
--   - Efficient JOINs within time windows
--
-- Difference from equipment_property:
--   equipment_property  = current static metadata (e.g., "Max RPM: 3000")
--   equipment_telemetry = time-series readings  (e.g., "RPM at 14:32:05 = 2847")
-- ============================================================================
CREATE TABLE IF NOT EXISTS equipment_telemetry (
    id              SERIAL,
    equipment_id    VARCHAR(50) NOT NULL REFERENCES equipment_master(id) ON DELETE CASCADE,
    parameter_name  VARCHAR(100) NOT NULL,            -- e.g., 'temperature', 'pressure', 'speed', 'vibration'
    value           NUMERIC(15,4) NOT NULL,           -- the reading value
    unit_of_measure VARCHAR(50),                      -- e.g., '°C', 'bar', 'RPM'
    recorded_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, recorded_at)
);

-- Convert to TimescaleDB hypertable
SELECT create_hypertable('equipment_telemetry', 'recorded_at', if_not_exists => TRUE);

-- Index: Query telemetry by equipment + parameter within a time window
CREATE INDEX IF NOT EXISTS idx_telemetry_equip_param ON equipment_telemetry(equipment_id, parameter_name, recorded_at DESC);

-- ============================================================================
-- GRANTS: Extend AI reader permissions to new tables
-- ============================================================================
GRANT SELECT ON ALL TABLES IN SCHEMA public TO taksa_ai_reader;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO taksa_ai_reader;
