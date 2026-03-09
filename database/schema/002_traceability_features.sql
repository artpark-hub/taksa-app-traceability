CREATE TABLE IF NOT EXISTS material_definition (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,           
    material_type   VARCHAR(50) NOT NULL             
                    CHECK (material_type IN ('RAW_MATERIAL', 'WORK_IN_PROGRESS', 'FINISHED_GOOD', 'PACKAGING')),
    unit_of_measure VARCHAR(50),                     
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS material_lot (
    lot_id                VARCHAR(100) PRIMARY KEY,      
    material_definition_id INTEGER NOT NULL REFERENCES material_definition(id) ON DELETE RESTRICT,
    status                VARCHAR(50) NOT NULL DEFAULT 'available'
                          CHECK (status IN ('available', 'in_process', 'completed', 'quarantined', 'released', 'shipped', 'scrapped')),
    quantity              NUMERIC(15,4),                  
    unit_of_measure       VARCHAR(50),                    
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_material_lot_status ON material_lot(status);
CREATE INDEX IF NOT EXISTS idx_material_lot_material_def ON material_lot(material_definition_id);

CREATE TABLE IF NOT EXISTS operator (
    operator_id     VARCHAR(100) PRIMARY KEY,         
    name            VARCHAR(255) NOT NULL,
    role            VARCHAR(100),                     
    shift           VARCHAR(50),                      
    status          VARCHAR(50) NOT NULL DEFAULT 'active'
                    CHECK (status IN ('active', 'inactive')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS work_order (
    work_order_id   VARCHAR(100) PRIMARY KEY,         
    description     TEXT,
    status          VARCHAR(50) NOT NULL DEFAULT 'planned'
                    CHECK (status IN ('planned', 'in_progress', 'completed', 'closed')),
    equipment_id    VARCHAR(50) REFERENCES equipment_master(id) ON DELETE SET NULL,    
    operator_id     VARCHAR(100) REFERENCES operator(operator_id) ON DELETE SET NULL,  
    output_lot_id   VARCHAR(100),                     
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


CREATE INDEX IF NOT EXISTS idx_work_order_status ON work_order(status);
CREATE INDEX IF NOT EXISTS idx_work_order_equipment ON work_order(equipment_id);


CREATE TABLE IF NOT EXISTS lot_genealogy (
    id                  SERIAL PRIMARY KEY,
    parent_lot_id       VARCHAR(100) NOT NULL REFERENCES material_lot(lot_id) ON DELETE RESTRICT,
    child_lot_id        VARCHAR(100) NOT NULL REFERENCES material_lot(lot_id) ON DELETE RESTRICT,
    work_order_id       VARCHAR(100) REFERENCES work_order(work_order_id) ON DELETE SET NULL,
    equipment_id        VARCHAR(50) REFERENCES equipment_master(id) ON DELETE SET NULL,
    quantity_consumed   NUMERIC(15,4),              
    quantity_produced   NUMERIC(15,4),               
    event_time          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    
    CONSTRAINT uq_genealogy_link UNIQUE (parent_lot_id, child_lot_id, work_order_id)
);


CREATE INDEX IF NOT EXISTS idx_genealogy_parent ON lot_genealogy(parent_lot_id);
CREATE INDEX IF NOT EXISTS idx_genealogy_child ON lot_genealogy(child_lot_id);
CREATE INDEX IF NOT EXISTS idx_genealogy_work_order ON lot_genealogy(work_order_id);


CREATE TABLE IF NOT EXISTS equipment_telemetry (
    id              SERIAL,
    equipment_id    VARCHAR(50) NOT NULL REFERENCES equipment_master(id) ON DELETE CASCADE,
    parameter_name  VARCHAR(100) NOT NULL,            
    value           NUMERIC(15,4) NOT NULL,          
    unit_of_measure VARCHAR(50),                      
    recorded_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, recorded_at)
);


SELECT create_hypertable('equipment_telemetry', 'recorded_at', if_not_exists => TRUE);


CREATE INDEX IF NOT EXISTS idx_telemetry_equip_param ON equipment_telemetry(equipment_id, parameter_name, recorded_at DESC);

GRANT SELECT ON ALL TABLES IN SCHEMA public TO taksa_ai_reader;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO taksa_ai_reader;