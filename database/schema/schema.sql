CREATE TABLE IF NOT EXISTS enterprise (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS site (
    id SERIAL PRIMARY KEY,
    enterprise_id INTEGER REFERENCES enterprise(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    description TEXT
);

CREATE TABLE IF NOT EXISTS area (
    id SERIAL PRIMARY KEY,
    site_id INTEGER REFERENCES site(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS production_line (
    id SERIAL PRIMARY KEY,
    area_id INTEGER REFERENCES area(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS production_unit (
    id SERIAL PRIMARY KEY,
    production_line_id INTEGER REFERENCES production_line(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS equipment_class (
    id SERIAL PRIMARY KEY,
    class_name VARCHAR(255) NOT NULL,
    version VARCHAR(50),
    description TEXT
);

CREATE TABLE IF NOT EXISTS equipment_master (
    id VARCHAR(50) PRIMARY KEY, 
    physical_asset_id VARCHAR(100),
    production_unit_id INTEGER REFERENCES production_unit(id) ON DELETE CASCADE,
    equipment_class_id INTEGER REFERENCES equipment_class(id) ON DELETE RESTRICT,
    operational_status VARCHAR(50),
    parent_equipment_id VARCHAR(50) REFERENCES equipment_master(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS equipment_capability (
    id SERIAL PRIMARY KEY,
    equipment_id VARCHAR(50) REFERENCES equipment_master(id) ON DELETE CASCADE,
    capability_name VARCHAR(100) NOT NULL,
    value VARCHAR(255),
    uom VARCHAR(50),
    description TEXT
);

CREATE TABLE IF NOT EXISTS equipment_property (
    id SERIAL PRIMARY KEY,
    equipment_id VARCHAR(50) REFERENCES equipment_master(id) ON DELETE CASCADE,
    property_name VARCHAR(100) NOT NULL,
    current_value VARCHAR(255),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS traceability_log (
    id SERIAL PRIMARY KEY,
    equipment_id VARCHAR(50) REFERENCES equipment_master(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL,
    work_order_id VARCHAR(100),
    material_lot_id VARCHAR(100),
    operator_id VARCHAR(100),
    event_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);