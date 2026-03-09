-- ============================================================================
-- MASTER SEED: Complete Factory Simulation
-- Wipes ALL data, then seeds ALL 16 tables
-- Scenario: EV Battery Pack Manufacturing at Taksa Global, Bengaluru Plant
-- Flow: Raw Cells + Raw Electrolyte → Mixer → WIP Slurry → Oven → WIP Baked Cells
--       → Assembly Robot → Finished Battery Pack → Packaging
-- ============================================================================

-- ============================================================================
-- STEP 0: WIPE EVERYTHING (respecting FK order — children first)
-- ============================================================================
TRUNCATE TABLE lot_genealogy CASCADE;
TRUNCATE TABLE equipment_telemetry CASCADE;
TRUNCATE TABLE traceability_log CASCADE;
TRUNCATE TABLE work_order CASCADE;
TRUNCATE TABLE material_lot CASCADE;
TRUNCATE TABLE material_definition CASCADE;
TRUNCATE TABLE operator CASCADE;
TRUNCATE TABLE equipment_property CASCADE;
TRUNCATE TABLE equipment_capability CASCADE;
TRUNCATE TABLE equipment_master CASCADE;
TRUNCATE TABLE equipment_class CASCADE;
TRUNCATE TABLE production_unit CASCADE;
TRUNCATE TABLE production_line CASCADE;
TRUNCATE TABLE area CASCADE;
TRUNCATE TABLE site CASCADE;
TRUNCATE TABLE enterprise CASCADE;

-- Reset sequences
ALTER SEQUENCE enterprise_id_seq RESTART WITH 1;
ALTER SEQUENCE site_id_seq RESTART WITH 1;
ALTER SEQUENCE area_id_seq RESTART WITH 1;
ALTER SEQUENCE production_line_id_seq RESTART WITH 1;
ALTER SEQUENCE production_unit_id_seq RESTART WITH 1;
ALTER SEQUENCE equipment_class_id_seq RESTART WITH 1;
ALTER SEQUENCE equipment_capability_id_seq RESTART WITH 1;
ALTER SEQUENCE equipment_property_id_seq RESTART WITH 1;
ALTER SEQUENCE material_definition_id_seq RESTART WITH 1;
ALTER SEQUENCE lot_genealogy_id_seq RESTART WITH 1;
ALTER SEQUENCE traceability_log_id_seq RESTART WITH 1;
ALTER SEQUENCE equipment_telemetry_id_seq RESTART WITH 1;

-- ============================================================================
-- STEP 1: ISA-95 EQUIPMENT HIERARCHY
-- ============================================================================

-- Enterprises
INSERT INTO enterprise (name, description) VALUES
('Taksa Global', 'Advanced Manufacturing Corporation — EV Battery Division'),
('Artpark Research', 'Robotics & AI Innovation Hub at IISc Bengaluru');

-- Sites
INSERT INTO site (enterprise_id, name, location, description) VALUES
(1, 'Bengaluru Plant', 'Peenya Industrial Area, Bengaluru, KA', 'High-volume EV battery assembly plant'),
(1, 'Pune Facility', 'Chakan MIDC, Pune, MH', 'Component fabrication and sub-assembly'),
(2, 'Innovation Lab', 'IISc Campus, Bengaluru', 'R&D testing ground for new manufacturing processes');

-- Areas (within Bengaluru Plant)
INSERT INTO area (site_id, name, description) VALUES
(1, 'Cell Preparation', 'Raw material receiving and cell preparation zone'),
(1, 'Main Assembly', 'Battery module and pack assembly'),
(1, 'Quality Control', 'Testing, inspection, and charge-discharge validation'),
(1, 'Packaging & Dispatch', 'Final packaging and shipment staging');

-- Production Lines
INSERT INTO production_line (area_id, name, description) VALUES
(1, 'Line CP-1: Mixing', 'Electrolyte mixing and slurry preparation'),
(2, 'Line A-1: Cell Assembly', 'Automated cell stacking and welding'),
(2, 'Line A-2: Pack Assembly', 'Battery pack final assembly'),
(3, 'Line QC-1: Validation', 'Electrical and thermal validation'),
(4, 'Line PK-1: Packaging', 'Shrink wrap and boxing');

-- Production Units (stations within lines)
INSERT INTO production_unit (production_line_id, name, description) VALUES
(1, 'Mixer Station', 'Industrial planetary mixer for electrolyte'),
(1, 'Oven Station', 'Drying oven for coated cells'),
(2, 'Cell Stacking Station', 'Robotic cell stacking unit'),
(2, 'Laser Welding Station', 'Precision laser welding for bus bars'),
(3, 'Module Assembly Station', 'Final module integration'),
(4, 'Charge-Discharge Bay', 'Capacity and cycle testing'),
(4, 'Thermal Chamber', 'Temperature stress testing'),
(5, 'Packaging Station', 'Automated packaging line');

-- Equipment Classes
INSERT INTO equipment_class (class_name, version, description) VALUES
('Industrial Mixer', 'v3.1', 'Planetary mixer for viscous slurries, 200L capacity'),
('Drying Oven', 'v2.0', 'Convection oven with nitrogen atmosphere, 300°C max'),
('Robotic Arm', 'v2.4', '6-axis industrial manipulator, 20kg payload'),
('Laser Welder', 'v1.2', 'Fiber laser welder, 2kW, 0.1mm precision'),
('Conveyor Belt', 'v3.0', 'Variable speed servo-driven transport'),
('Charge-Discharge Unit', 'v4.0', '64-channel battery cycler'),
('Thermal Chamber', 'v2.1', '-40°C to +150°C environmental chamber'),
('Packaging Machine', 'v1.5', 'Automated shrink-wrap and boxing system');

-- Equipment Master (actual machines on the floor)
INSERT INTO equipment_master (id, physical_asset_id, production_unit_id, equipment_class_id, operational_status, parent_equipment_id) VALUES
('MIXER-001', 'PA-MIX-2024-001', 1, 1, 'active', NULL),
('OVEN-001',  'PA-OVN-2024-001', 2, 2, 'active', NULL),
('ROBOT-001', 'PA-ROB-2024-001', 3, 3, 'active', NULL),
('ROBOT-002', 'PA-ROB-2024-002', 3, 3, 'maintenance', NULL),
('LASER-001', 'PA-LAS-2024-001', 4, 4, 'active', NULL),
('CONV-001',  'PA-CNV-2024-001', 5, 5, 'active', NULL),
('CYCLER-001','PA-CYC-2024-001', 6, 6, 'active', NULL),
('THERM-001', 'PA-THR-2024-001', 7, 7, 'active', NULL),
('PKG-001',   'PA-PKG-2024-001', 8, 8, 'active', NULL),
-- Sub-components of ROBOT-001
('ROBOT-001-GRIPPER', 'PA-GRP-2024-001', 3, 3, 'active', 'ROBOT-001'),
('ROBOT-001-VISION',  'PA-VIS-2024-001', 3, 3, 'active', 'ROBOT-001');

-- Equipment Capabilities
INSERT INTO equipment_capability (equipment_id, capability_name, value, uom, description) VALUES
('MIXER-001', 'max_capacity', '200', 'liters', 'Maximum batch volume'),
('MIXER-001', 'max_speed', '1500', 'RPM', 'Maximum impeller speed'),
('OVEN-001',  'max_temperature', '300', '°C', 'Maximum operating temperature'),
('OVEN-001',  'atmosphere', 'nitrogen', NULL, 'Inert atmosphere capability'),
('ROBOT-001', 'payload', '20', 'kg', 'Maximum lifting capacity'),
('ROBOT-001', 'reach', '1800', 'mm', 'Maximum arm reach'),
('LASER-001', 'power', '2000', 'W', 'Maximum laser power'),
('LASER-001', 'precision', '0.1', 'mm', 'Welding spot precision'),
('CYCLER-001','channels', '64', 'units', 'Number of test channels'),
('THERM-001', 'temp_range', '-40 to +150', '°C', 'Operating temperature range');

-- Equipment Properties (current state)
INSERT INTO equipment_property (equipment_id, property_name, current_value) VALUES
('MIXER-001', 'motor_temperature', '42.3'),
('MIXER-001', 'impeller_speed', '800'),
('MIXER-001', 'batch_viscosity', '3200'),
('OVEN-001',  'zone1_temperature', '185.0'),
('OVEN-001',  'zone2_temperature', '192.5'),
('OVEN-001',  'nitrogen_flow_rate', '12.5'),
('ROBOT-001', 'joint1_temperature', '38.2'),
('ROBOT-001', 'cycle_count', '15847'),
('ROBOT-001', 'uptime_hours', '2340'),
('LASER-001', 'power_output', '1850'),
('LASER-001', 'focal_distance', '150.2'),
('CYCLER-001','active_channels', '48'),
('THERM-001', 'chamber_temperature', '25.0');

-- ============================================================================
-- STEP 2: OPERATORS
-- ============================================================================
INSERT INTO operator (operator_id, name, role, shift, status) VALUES
('OP-KAVYA-01',  'Kavya Sharma',    'Senior Operator',    'DAY_SHIFT',     'active'),
('OP-RAHUL-02',  'Rahul Menon',     'Machine Operator',   'DAY_SHIFT',     'active'),
('OP-SURESH-03', 'Suresh Kumar',    'Welding Specialist', 'DAY_SHIFT',     'active'),
('OP-PRIYA-04',  'Priya Nair',      'QC Inspector',       'DAY_SHIFT',     'active'),
('OP-ARJUN-05',  'Arjun Reddy',     'Machine Operator',   'NIGHT_SHIFT',   'active'),
('OP-MEERA-06',  'Meera Patel',     'Packaging Operator', 'DAY_SHIFT',     'active'),
('OP-VIKRAM-07', 'Vikram Singh',    'Shift Supervisor',   'NIGHT_SHIFT',   'active'),
('OP-DEEPA-08',  'Deepa Iyer',     'QC Inspector',       'NIGHT_SHIFT',   'active'),
('OP-RAVI-09',   'Ravi Prasad',     'Maintenance Tech',   'DAY_SHIFT',     'inactive');

-- ============================================================================
-- STEP 3: MATERIAL DEFINITIONS
-- ============================================================================
INSERT INTO material_definition (name, material_type, unit_of_measure, description) VALUES
('Lithium Cell NMC811',     'RAW_MATERIAL',       'units',  'Nickel-Manganese-Cobalt cathode cells, 3.7V nominal'),
('Electrolyte Solution EC', 'RAW_MATERIAL',       'liters', 'Ethylene carbonate based electrolyte, 1M LiPF6'),
('Copper Bus Bar',          'RAW_MATERIAL',       'units',  'Pre-cut copper bus bars for series/parallel connections'),
('Thermal Paste TG-7',      'RAW_MATERIAL',       'kg',     'Thermal interface compound, 7 W/mK conductivity'),
('Cell Slurry Mix',         'WORK_IN_PROGRESS',   'liters', 'Mixed electrolyte-cell slurry for coating'),
('Coated Cell Sheet',       'WORK_IN_PROGRESS',   'units',  'Oven-dried coated cell sheets, ready for stacking'),
('Cell Stack Assembly',     'WORK_IN_PROGRESS',   'units',  'Stacked and aligned cell assembly, pre-welding'),
('Welded Cell Module',      'WORK_IN_PROGRESS',   'units',  'Laser-welded cell module with bus bars'),
('Battery Pack BP-7200',    'FINISHED_GOOD',      'units',  'Complete 72V battery pack, tested and validated'),
('Corrugated Box CB-L',     'PACKAGING',          'units',  'Large corrugated shipping box with foam inserts'),
('Battery Pack BP-7200 Shipped', 'FINISHED_GOOD', 'units',  'Packaged battery pack ready for dispatch');

-- ============================================================================
-- STEP 4: MATERIAL LOTS (the actual batches)
-- Simulates a full production run on 2026-02-25
-- ============================================================================
INSERT INTO material_lot (lot_id, material_definition_id, status, quantity, unit_of_measure, created_at, updated_at) VALUES
-- Raw Materials received
('RM-CELL-2026-001',  1, 'completed',  500,   'units',  '2026-02-24 06:00:00+05:30', '2026-02-25 14:00:00+05:30'),
('RM-CELL-2026-002',  1, 'available',  500,   'units',  '2026-02-24 06:05:00+05:30', '2026-02-24 06:05:00+05:30'),
('RM-ELEC-2026-001',  2, 'completed',  100,   'liters', '2026-02-24 06:10:00+05:30', '2026-02-25 10:30:00+05:30'),
('RM-BUSBAR-2026-001',3, 'completed',  200,   'units',  '2026-02-24 06:15:00+05:30', '2026-02-25 12:00:00+05:30'),
('RM-PASTE-2026-001', 4, 'completed',  25,    'kg',     '2026-02-24 06:20:00+05:30', '2026-02-25 13:00:00+05:30'),
-- WIP lots created during processing
('WIP-SLURRY-2026-001', 5, 'completed', 80,   'liters', '2026-02-25 08:30:00+05:30', '2026-02-25 10:30:00+05:30'),
('WIP-COATED-2026-001', 6, 'completed', 200,  'units',  '2026-02-25 09:30:00+05:30', '2026-02-25 11:30:00+05:30'),
('WIP-STACK-2026-001',  7, 'completed', 50,   'units',  '2026-02-25 11:00:00+05:30', '2026-02-25 12:30:00+05:30'),
('WIP-MODULE-2026-001', 8, 'completed', 50,   'units',  '2026-02-25 12:00:00+05:30', '2026-02-25 13:30:00+05:30'),
-- Finished Goods
('FG-BP-2026-001',      9, 'released',  10,   'units',  '2026-02-25 13:00:00+05:30', '2026-02-25 16:00:00+05:30'),
('FG-BP-2026-002',      9, 'quarantined', 2,  'units',  '2026-02-25 13:30:00+05:30', '2026-02-25 15:30:00+05:30'),
-- Packaging
('PKG-BOX-2026-001',   10, 'completed', 12,   'units',  '2026-02-25 15:00:00+05:30', '2026-02-25 16:30:00+05:30'),
-- Shipped finished good
('FG-SHIP-2026-001',   11, 'shipped',   10,   'units',  '2026-02-25 16:00:00+05:30', '2026-02-25 17:00:00+05:30'),
-- Second production run lots (partial, for variety)
('RM-CELL-2026-003',   1, 'in_process', 500,  'units',  '2026-02-26 06:00:00+05:30', '2026-02-26 08:00:00+05:30'),
('RM-ELEC-2026-002',   2, 'in_process', 100,  'liters', '2026-02-26 06:05:00+05:30', '2026-02-26 08:00:00+05:30'),
('WIP-SLURRY-2026-002',5, 'in_process', 80,   'liters', '2026-02-26 08:30:00+05:30', '2026-02-26 08:30:00+05:30');

-- ============================================================================
-- STEP 5: WORK ORDERS
-- ============================================================================
INSERT INTO work_order (work_order_id, description, status, equipment_id, operator_id, output_lot_id, planned_quantity, actual_quantity, unit_of_measure, planned_start, planned_end, actual_start, actual_end) VALUES
-- Production Run 1 (2026-02-25, Day Shift)
('WO-2026-0001', 'Mix electrolyte with cell material to create coating slurry',
 'completed', 'MIXER-001', 'OP-KAVYA-01', 'WIP-SLURRY-2026-001',
 100, 80, 'liters',
 '2026-02-25 08:00:00+05:30', '2026-02-25 10:00:00+05:30',
 '2026-02-25 08:15:00+05:30', '2026-02-25 10:30:00+05:30'),

('WO-2026-0002', 'Dry coated cell sheets in nitrogen oven',
 'completed', 'OVEN-001', 'OP-KAVYA-01', 'WIP-COATED-2026-001',
 200, 200, 'units',
 '2026-02-25 09:00:00+05:30', '2026-02-25 11:00:00+05:30',
 '2026-02-25 09:30:00+05:30', '2026-02-25 11:30:00+05:30'),

('WO-2026-0003', 'Stack dried cells into module assemblies',
 'completed', 'ROBOT-001', 'OP-RAHUL-02', 'WIP-STACK-2026-001',
 60, 50, 'units',
 '2026-02-25 10:30:00+05:30', '2026-02-25 12:00:00+05:30',
 '2026-02-25 11:00:00+05:30', '2026-02-25 12:30:00+05:30'),

('WO-2026-0004', 'Laser weld bus bars onto cell stacks',
 'completed', 'LASER-001', 'OP-SURESH-03', 'WIP-MODULE-2026-001',
 50, 50, 'units',
 '2026-02-25 11:30:00+05:30', '2026-02-25 13:00:00+05:30',
 '2026-02-25 12:00:00+05:30', '2026-02-25 13:30:00+05:30'),

('WO-2026-0005', 'Final assembly of battery pack — modules + thermal paste + casing',
 'completed', 'CONV-001', 'OP-RAHUL-02', 'FG-BP-2026-001',
 12, 10, 'units',
 '2026-02-25 12:30:00+05:30', '2026-02-25 14:30:00+05:30',
 '2026-02-25 13:00:00+05:30', '2026-02-25 14:00:00+05:30'),

('WO-2026-0006', 'QC validation — charge-discharge cycle testing',
 'completed', 'CYCLER-001', 'OP-PRIYA-04', 'FG-BP-2026-001',
 12, 10, 'units',
 '2026-02-25 14:00:00+05:30', '2026-02-25 16:00:00+05:30',
 '2026-02-25 14:00:00+05:30', '2026-02-25 15:30:00+05:30'),

('WO-2026-0007', 'Package released battery packs for shipping',
 'completed', 'PKG-001', 'OP-MEERA-06', 'FG-SHIP-2026-001',
 10, 10, 'units',
 '2026-02-25 15:30:00+05:30', '2026-02-25 17:00:00+05:30',
 '2026-02-25 16:00:00+05:30', '2026-02-25 17:00:00+05:30'),

-- Production Run 2 (2026-02-26, Night Shift — in progress)
('WO-2026-0008', 'Mix electrolyte batch 2',
 'in_progress', 'MIXER-001', 'OP-ARJUN-05', 'WIP-SLURRY-2026-002',
 100, NULL, 'liters',
 '2026-02-26 08:00:00+05:30', '2026-02-26 10:00:00+05:30',
 '2026-02-26 08:15:00+05:30', NULL),

-- Planned future work order
('WO-2026-0009', 'Dry coated cell sheets batch 2',
 'planned', 'OVEN-001', 'OP-ARJUN-05', NULL,
 200, NULL, 'units',
 '2026-02-26 10:00:00+05:30', '2026-02-26 12:00:00+05:30',
 NULL, NULL);

-- ============================================================================
-- STEP 6: LOT GENEALOGY (the DAG of material transformations)
-- Production Run 1: RM → WIP-SLURRY → WIP-COATED → WIP-STACK → WIP-MODULE → FG → Shipped
-- ============================================================================
INSERT INTO lot_genealogy (parent_lot_id, child_lot_id, work_order_id, equipment_id, quantity_consumed, quantity_produced, event_time) VALUES
-- Stage 1: Mixing (Raw Cells + Electrolyte → Slurry)
('RM-CELL-2026-001',  'WIP-SLURRY-2026-001', 'WO-2026-0001', 'MIXER-001', 200, 80, '2026-02-25 10:30:00+05:30'),
('RM-ELEC-2026-001',  'WIP-SLURRY-2026-001', 'WO-2026-0001', 'MIXER-001', 50,  80, '2026-02-25 10:30:00+05:30'),

-- Stage 2: Oven Drying (Slurry → Coated Sheets)
('WIP-SLURRY-2026-001', 'WIP-COATED-2026-001', 'WO-2026-0002', 'OVEN-001', 80, 200, '2026-02-25 11:30:00+05:30'),

-- Stage 3: Cell Stacking (Coated Sheets → Cell Stack)
('WIP-COATED-2026-001', 'WIP-STACK-2026-001', 'WO-2026-0003', 'ROBOT-001', 200, 50, '2026-02-25 12:30:00+05:30'),

-- Stage 4: Laser Welding (Cell Stack + Bus Bars → Welded Module)
('WIP-STACK-2026-001',  'WIP-MODULE-2026-001', 'WO-2026-0004', 'LASER-001', 50, 50, '2026-02-25 13:30:00+05:30'),
('RM-BUSBAR-2026-001',  'WIP-MODULE-2026-001', 'WO-2026-0004', 'LASER-001', 50, 50, '2026-02-25 13:30:00+05:30'),

-- Stage 5: Final Assembly (Module + Thermal Paste → Battery Pack)
('WIP-MODULE-2026-001', 'FG-BP-2026-001', 'WO-2026-0005', 'CONV-001', 10, 10, '2026-02-25 14:00:00+05:30'),
('RM-PASTE-2026-001',   'FG-BP-2026-001', 'WO-2026-0005', 'CONV-001', 5,  10, '2026-02-25 14:00:00+05:30'),

-- Stage 5b: Defective units split into quarantined lot
('WIP-MODULE-2026-001', 'FG-BP-2026-002', 'WO-2026-0005', 'CONV-001', 2, 2, '2026-02-25 14:00:00+05:30'),

-- Stage 6: Packaging (Released FG + Boxes → Shipped FG)
('FG-BP-2026-001',     'FG-SHIP-2026-001', 'WO-2026-0007', 'PKG-001', 10, 10, '2026-02-25 17:00:00+05:30'),
('PKG-BOX-2026-001',   'FG-SHIP-2026-001', 'WO-2026-0007', 'PKG-001', 10, 10, '2026-02-25 17:00:00+05:30'),

-- Production Run 2 (in progress)
('RM-CELL-2026-003',  'WIP-SLURRY-2026-002', 'WO-2026-0008', 'MIXER-001', 200, NULL, '2026-02-26 08:15:00+05:30'),
('RM-ELEC-2026-002',  'WIP-SLURRY-2026-002', 'WO-2026-0008', 'MIXER-001', 50,  NULL, '2026-02-26 08:15:00+05:30');

-- ============================================================================
-- STEP 7: TRACEABILITY LOG (event stream — the timeline of everything)
-- ============================================================================
INSERT INTO traceability_log (equipment_id, event_type, work_order_id, material_lot_id, operator_id, event_time) VALUES
-- Production Run 1: Day Shift, 2026-02-25
-- Mixing
('MIXER-001', 'work_order_started',  'WO-2026-0001', 'RM-CELL-2026-001',     'OP-KAVYA-01', '2026-02-25 08:15:00+05:30'),
('MIXER-001', 'material_loaded',     'WO-2026-0001', 'RM-CELL-2026-001',     'OP-KAVYA-01', '2026-02-25 08:20:00+05:30'),
('MIXER-001', 'material_loaded',     'WO-2026-0001', 'RM-ELEC-2026-001',     'OP-KAVYA-01', '2026-02-25 08:25:00+05:30'),
('MIXER-001', 'process_started',     'WO-2026-0001', 'WIP-SLURRY-2026-001',  'OP-KAVYA-01', '2026-02-25 08:30:00+05:30'),
('MIXER-001', 'process_completed',   'WO-2026-0001', 'WIP-SLURRY-2026-001',  'OP-KAVYA-01', '2026-02-25 10:30:00+05:30'),
('MIXER-001', 'work_order_completed','WO-2026-0001', 'WIP-SLURRY-2026-001',  'OP-KAVYA-01', '2026-02-25 10:30:00+05:30'),
-- Oven
('OVEN-001',  'work_order_started',  'WO-2026-0002', 'WIP-SLURRY-2026-001',  'OP-KAVYA-01', '2026-02-25 09:30:00+05:30'),
('OVEN-001',  'material_loaded',     'WO-2026-0002', 'WIP-SLURRY-2026-001',  'OP-KAVYA-01', '2026-02-25 09:35:00+05:30'),
('OVEN-001',  'process_started',     'WO-2026-0002', 'WIP-COATED-2026-001',  'OP-KAVYA-01', '2026-02-25 09:40:00+05:30'),
('OVEN-001',  'process_completed',   'WO-2026-0002', 'WIP-COATED-2026-001',  'OP-KAVYA-01', '2026-02-25 11:30:00+05:30'),
('OVEN-001',  'work_order_completed','WO-2026-0002', 'WIP-COATED-2026-001',  'OP-KAVYA-01', '2026-02-25 11:30:00+05:30'),
-- Cell Stacking
('ROBOT-001', 'work_order_started',  'WO-2026-0003', 'WIP-COATED-2026-001',  'OP-RAHUL-02', '2026-02-25 11:00:00+05:30'),
('ROBOT-001', 'material_loaded',     'WO-2026-0003', 'WIP-COATED-2026-001',  'OP-RAHUL-02', '2026-02-25 11:05:00+05:30'),
('ROBOT-001', 'process_started',     'WO-2026-0003', 'WIP-STACK-2026-001',   'OP-RAHUL-02', '2026-02-25 11:10:00+05:30'),
('ROBOT-001', 'error_jam',           'WO-2026-0003', 'WIP-STACK-2026-001',   'OP-RAHUL-02', '2026-02-25 11:45:00+05:30'),
('ROBOT-001', 'process_resumed',     'WO-2026-0003', 'WIP-STACK-2026-001',   'OP-RAHUL-02', '2026-02-25 11:50:00+05:30'),
('ROBOT-001', 'process_completed',   'WO-2026-0003', 'WIP-STACK-2026-001',   'OP-RAHUL-02', '2026-02-25 12:30:00+05:30'),
('ROBOT-001', 'work_order_completed','WO-2026-0003', 'WIP-STACK-2026-001',   'OP-RAHUL-02', '2026-02-25 12:30:00+05:30'),
-- Laser Welding
('LASER-001', 'work_order_started',  'WO-2026-0004', 'WIP-STACK-2026-001',   'OP-SURESH-03','2026-02-25 12:00:00+05:30'),
('LASER-001', 'material_loaded',     'WO-2026-0004', 'WIP-STACK-2026-001',   'OP-SURESH-03','2026-02-25 12:05:00+05:30'),
('LASER-001', 'material_loaded',     'WO-2026-0004', 'RM-BUSBAR-2026-001',   'OP-SURESH-03','2026-02-25 12:10:00+05:30'),
('LASER-001', 'process_started',     'WO-2026-0004', 'WIP-MODULE-2026-001',  'OP-SURESH-03','2026-02-25 12:15:00+05:30'),
('LASER-001', 'process_completed',   'WO-2026-0004', 'WIP-MODULE-2026-001',  'OP-SURESH-03','2026-02-25 13:30:00+05:30'),
('LASER-001', 'work_order_completed','WO-2026-0004', 'WIP-MODULE-2026-001',  'OP-SURESH-03','2026-02-25 13:30:00+05:30'),
-- Final Assembly
('CONV-001',  'work_order_started',  'WO-2026-0005', 'WIP-MODULE-2026-001',  'OP-RAHUL-02', '2026-02-25 13:00:00+05:30'),
('CONV-001',  'material_loaded',     'WO-2026-0005', 'WIP-MODULE-2026-001',  'OP-RAHUL-02', '2026-02-25 13:05:00+05:30'),
('CONV-001',  'material_loaded',     'WO-2026-0005', 'RM-PASTE-2026-001',    'OP-RAHUL-02', '2026-02-25 13:10:00+05:30'),
('CONV-001',  'process_started',     'WO-2026-0005', 'FG-BP-2026-001',       'OP-RAHUL-02', '2026-02-25 13:15:00+05:30'),
('CONV-001',  'quality_flag',        'WO-2026-0005', 'FG-BP-2026-002',       'OP-RAHUL-02', '2026-02-25 13:45:00+05:30'),
('CONV-001',  'process_completed',   'WO-2026-0005', 'FG-BP-2026-001',       'OP-RAHUL-02', '2026-02-25 14:00:00+05:30'),
('CONV-001',  'work_order_completed','WO-2026-0005', 'FG-BP-2026-001',       'OP-RAHUL-02', '2026-02-25 14:00:00+05:30'),
-- QC Validation
('CYCLER-001','work_order_started',  'WO-2026-0006', 'FG-BP-2026-001',       'OP-PRIYA-04', '2026-02-25 14:00:00+05:30'),
('CYCLER-001','qc_test_started',     'WO-2026-0006', 'FG-BP-2026-001',       'OP-PRIYA-04', '2026-02-25 14:05:00+05:30'),
('CYCLER-001','qc_test_passed',      'WO-2026-0006', 'FG-BP-2026-001',       'OP-PRIYA-04', '2026-02-25 15:30:00+05:30'),
('CYCLER-001','work_order_completed','WO-2026-0006', 'FG-BP-2026-001',       'OP-PRIYA-04', '2026-02-25 15:30:00+05:30'),
('CYCLER-001','qc_test_started',     'WO-2026-0006', 'FG-BP-2026-002',       'OP-PRIYA-04', '2026-02-25 14:10:00+05:30'),
('CYCLER-001','qc_test_failed',      'WO-2026-0006', 'FG-BP-2026-002',       'OP-PRIYA-04', '2026-02-25 15:30:00+05:30'),
('CYCLER-001','lot_quarantined',     'WO-2026-0006', 'FG-BP-2026-002',       'OP-PRIYA-04', '2026-02-25 15:35:00+05:30'),
-- Packaging
('PKG-001',   'work_order_started',  'WO-2026-0007', 'FG-BP-2026-001',       'OP-MEERA-06', '2026-02-25 16:00:00+05:30'),
('PKG-001',   'material_loaded',     'WO-2026-0007', 'FG-BP-2026-001',       'OP-MEERA-06', '2026-02-25 16:05:00+05:30'),
('PKG-001',   'material_loaded',     'WO-2026-0007', 'PKG-BOX-2026-001',     'OP-MEERA-06', '2026-02-25 16:10:00+05:30'),
('PKG-001',   'process_completed',   'WO-2026-0007', 'FG-SHIP-2026-001',     'OP-MEERA-06', '2026-02-25 17:00:00+05:30'),
('PKG-001',   'lot_shipped',         'WO-2026-0007', 'FG-SHIP-2026-001',     'OP-MEERA-06', '2026-02-25 17:00:00+05:30'),
-- Production Run 2 (in progress)
('MIXER-001', 'work_order_started',  'WO-2026-0008', 'RM-CELL-2026-003',     'OP-ARJUN-05', '2026-02-26 08:15:00+05:30'),
('MIXER-001', 'material_loaded',     'WO-2026-0008', 'RM-CELL-2026-003',     'OP-ARJUN-05', '2026-02-26 08:20:00+05:30'),
('MIXER-001', 'material_loaded',     'WO-2026-0008', 'RM-ELEC-2026-002',     'OP-ARJUN-05', '2026-02-26 08:25:00+05:30'),
('MIXER-001', 'process_started',     'WO-2026-0008', 'WIP-SLURRY-2026-002',  'OP-ARJUN-05', '2026-02-26 08:30:00+05:30');

-- ============================================================================
-- STEP 8: EQUIPMENT TELEMETRY (time-series process parameters)
-- Simulates 15-minute interval readings during Production Run 1
-- ============================================================================

-- MIXER-001 telemetry during WO-2026-0001 (08:15 - 10:30)
INSERT INTO equipment_telemetry (equipment_id, parameter_name, value, unit_of_measure, recorded_at) VALUES
('MIXER-001', 'temperature',     25.0, '°C',  '2026-02-25 08:15:00+05:30'),
('MIXER-001', 'temperature',     32.5, '°C',  '2026-02-25 08:30:00+05:30'),
('MIXER-001', 'temperature',     41.2, '°C',  '2026-02-25 08:45:00+05:30'),
('MIXER-001', 'temperature',     45.8, '°C',  '2026-02-25 09:00:00+05:30'),
('MIXER-001', 'temperature',     48.3, '°C',  '2026-02-25 09:15:00+05:30'),
('MIXER-001', 'temperature',     50.1, '°C',  '2026-02-25 09:30:00+05:30'),
('MIXER-001', 'temperature',     49.7, '°C',  '2026-02-25 09:45:00+05:30'),
('MIXER-001', 'temperature',     48.5, '°C',  '2026-02-25 10:00:00+05:30'),
('MIXER-001', 'temperature',     46.2, '°C',  '2026-02-25 10:15:00+05:30'),
('MIXER-001', 'temperature',     42.0, '°C',  '2026-02-25 10:30:00+05:30'),
('MIXER-001', 'impeller_speed',  0,    'RPM', '2026-02-25 08:15:00+05:30'),
('MIXER-001', 'impeller_speed',  400,  'RPM', '2026-02-25 08:30:00+05:30'),
('MIXER-001', 'impeller_speed',  800,  'RPM', '2026-02-25 08:45:00+05:30'),
('MIXER-001', 'impeller_speed',  1200, 'RPM', '2026-02-25 09:00:00+05:30'),
('MIXER-001', 'impeller_speed',  1200, 'RPM', '2026-02-25 09:15:00+05:30'),
('MIXER-001', 'impeller_speed',  1200, 'RPM', '2026-02-25 09:30:00+05:30'),
('MIXER-001', 'impeller_speed',  1200, 'RPM', '2026-02-25 09:45:00+05:30'),
('MIXER-001', 'impeller_speed',  800,  'RPM', '2026-02-25 10:00:00+05:30'),
('MIXER-001', 'impeller_speed',  400,  'RPM', '2026-02-25 10:15:00+05:30'),
('MIXER-001', 'impeller_speed',  0,    'RPM', '2026-02-25 10:30:00+05:30'),
('MIXER-001', 'viscosity',       500,  'cP',  '2026-02-25 08:30:00+05:30'),
('MIXER-001', 'viscosity',       1200, 'cP',  '2026-02-25 09:00:00+05:30'),
('MIXER-001', 'viscosity',       2400, 'cP',  '2026-02-25 09:30:00+05:30'),
('MIXER-001', 'viscosity',       3100, 'cP',  '2026-02-25 10:00:00+05:30'),
('MIXER-001', 'viscosity',       3200, 'cP',  '2026-02-25 10:30:00+05:30'),

-- OVEN-001 telemetry during WO-2026-0002 (09:30 - 11:30)
('OVEN-001', 'zone1_temperature', 25.0,  '°C',  '2026-02-25 09:30:00+05:30'),
('OVEN-001', 'zone1_temperature', 95.0,  '°C',  '2026-02-25 09:45:00+05:30'),
('OVEN-001', 'zone1_temperature', 155.0, '°C',  '2026-02-25 10:00:00+05:30'),
('OVEN-001', 'zone1_temperature', 182.0, '°C',  '2026-02-25 10:15:00+05:30'),
('OVEN-001', 'zone1_temperature', 185.0, '°C',  '2026-02-25 10:30:00+05:30'),
('OVEN-001', 'zone1_temperature', 185.5, '°C',  '2026-02-25 10:45:00+05:30'),
('OVEN-001', 'zone1_temperature', 185.2, '°C',  '2026-02-25 11:00:00+05:30'),
('OVEN-001', 'zone1_temperature', 120.0, '°C',  '2026-02-25 11:15:00+05:30'),
('OVEN-001', 'zone1_temperature', 65.0,  '°C',  '2026-02-25 11:30:00+05:30'),
('OVEN-001', 'nitrogen_flow',     0,     'L/min','2026-02-25 09:30:00+05:30'),
('OVEN-001', 'nitrogen_flow',     8.5,   'L/min','2026-02-25 09:45:00+05:30'),
('OVEN-001', 'nitrogen_flow',     12.0,  'L/min','2026-02-25 10:00:00+05:30'),
('OVEN-001', 'nitrogen_flow',     12.5,  'L/min','2026-02-25 10:15:00+05:30'),
('OVEN-001', 'nitrogen_flow',     12.5,  'L/min','2026-02-25 10:30:00+05:30'),
('OVEN-001', 'nitrogen_flow',     12.5,  'L/min','2026-02-25 10:45:00+05:30'),
('OVEN-001', 'nitrogen_flow',     12.5,  'L/min','2026-02-25 11:00:00+05:30'),
('OVEN-001', 'nitrogen_flow',     6.0,   'L/min','2026-02-25 11:15:00+05:30'),
('OVEN-001', 'nitrogen_flow',     0,     'L/min','2026-02-25 11:30:00+05:30'),

-- ROBOT-001 telemetry during WO-2026-0003 (11:00 - 12:30) — includes the jam event
('ROBOT-001', 'joint1_temp',      35.0, '°C', '2026-02-25 11:00:00+05:30'),
('ROBOT-001', 'joint1_temp',      37.5, '°C', '2026-02-25 11:15:00+05:30'),
('ROBOT-001', 'joint1_temp',      39.2, '°C', '2026-02-25 11:30:00+05:30'),
('ROBOT-001', 'joint1_temp',      52.8, '°C', '2026-02-25 11:45:00+05:30'),
('ROBOT-001', 'joint1_temp',      41.0, '°C', '2026-02-25 12:00:00+05:30'),
('ROBOT-001', 'joint1_temp',      38.5, '°C', '2026-02-25 12:15:00+05:30'),
('ROBOT-001', 'joint1_temp',      37.0, '°C', '2026-02-25 12:30:00+05:30'),
('ROBOT-001', 'cycle_time',       4.2,  'sec','2026-02-25 11:00:00+05:30'),
('ROBOT-001', 'cycle_time',       4.1,  'sec','2026-02-25 11:15:00+05:30'),
('ROBOT-001', 'cycle_time',       4.3,  'sec','2026-02-25 11:30:00+05:30'),
('ROBOT-001', 'cycle_time',       99.9, 'sec','2026-02-25 11:45:00+05:30'),
('ROBOT-001', 'cycle_time',       4.5,  'sec','2026-02-25 12:00:00+05:30'),
('ROBOT-001', 'cycle_time',       4.2,  'sec','2026-02-25 12:15:00+05:30'),
('ROBOT-001', 'cycle_time',       4.1,  'sec','2026-02-25 12:30:00+05:30'),

-- LASER-001 telemetry during WO-2026-0004 (12:00 - 13:30)
('LASER-001', 'power_output',     0,     'W',  '2026-02-25 12:00:00+05:30'),
('LASER-001', 'power_output',     1500,  'W',  '2026-02-25 12:15:00+05:30'),
('LASER-001', 'power_output',     1850,  'W',  '2026-02-25 12:30:00+05:30'),
('LASER-001', 'power_output',     1870,  'W',  '2026-02-25 12:45:00+05:30'),
('LASER-001', 'power_output',     1855,  'W',  '2026-02-25 13:00:00+05:30'),
('LASER-001', 'power_output',     1840,  'W',  '2026-02-25 13:15:00+05:30'),
('LASER-001', 'power_output',     0,     'W',  '2026-02-25 13:30:00+05:30'),
('LASER-001', 'focal_distance',   150.0, 'mm', '2026-02-25 12:15:00+05:30'),
('LASER-001', 'focal_distance',   150.2, 'mm', '2026-02-25 12:30:00+05:30'),
('LASER-001', 'focal_distance',   150.1, 'mm', '2026-02-25 12:45:00+05:30'),
('LASER-001', 'focal_distance',   150.3, 'mm', '2026-02-25 13:00:00+05:30'),
('LASER-001', 'focal_distance',   150.2, 'mm', '2026-02-25 13:15:00+05:30');

-- ============================================================================
-- GRANTS
-- ============================================================================
GRANT USAGE ON SCHEMA public TO taksa_ai_reader;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO taksa_ai_reader;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO taksa_ai_reader;
