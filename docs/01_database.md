# Database Guide — Taksa Traceability

> **Tech stack:** PostgreSQL 15 + TimescaleDB extension  
> **Schema files:** `database/schema/schema.sql` · `database/schema/002_traceability_features.sql`

---

## Overview

The database models a **manufacturing plant hierarchy** from the top-level company all the way down to individual sensor readings. Every physical thing that happens on the factory floor — a raw material arriving, a machine processing it, an operator completing a work order — is recorded here and linked together so you can trace any product's complete history.

---

## Complete Table Map

```
enterprise
  └── site
        └── area
              └── production_line
                    └── production_unit
                          └── equipment_master ←── equipment_class
                                ├── equipment_capability
                                ├── equipment_property
                                ├── traceability_log (TimescaleDB hypertable)
                                └── equipment_telemetry (TimescaleDB hypertable)

material_definition
  └── material_lot
        └── lot_genealogy (parent_lot ↔ child_lot via work_order + equipment)

operator
  └── work_order ──── equipment_master
                └──── output_lot_id → material_lot
```

---

## Table-by-Table Reference

### 1. `enterprise`
| Column | Type | Description |
|---|---|---|
| `id` | SERIAL PK | Auto-incremented integer ID |
| `name` | VARCHAR(255) | Company name (e.g. "Tesla Inc.") |
| `description` | TEXT | Optional notes |

**Real-life example:** Tesla Inc., Ola Electric, Tata Motors

---

### 2. `site`
| Column | Type | Description |
|---|---|---|
| `id` | SERIAL PK | Auto-incremented integer ID |
| `enterprise_id` | FK → enterprise | Which company owns this site |
| `name` | VARCHAR(255) | Site name (e.g. "Gigafactory Nevada") |
| `location` | VARCHAR(255) | Physical address or city |
| `description` | TEXT | Optional notes |

**Real-life example:** Gigafactory Nevada, Chennai Assembly Plant, Pune Battery Hub

---

### 3. `area`
| Column | Type | Description |
|---|---|---|
| `id` | SERIAL PK | Auto-incremented integer ID |
| `site_id` | FK → site | Which site this area belongs to |
| `name` | VARCHAR(255) | Area name (e.g. "Battery Assembly Hall") |
| `description` | TEXT | Optional notes |

**Real-life example:** Body Welding Zone, Paint Shop, Battery Pack Assembly

---

### 4. `production_line`
| Column | Type | Description |
|---|---|---|
| `id` | SERIAL PK | Auto-incremented integer ID |
| `area_id` | FK → area | Which area this line is inside |
| `name` | VARCHAR(255) | Line name (e.g. "Cell Module Line A") |
| `description` | TEXT | Optional notes |

**Real-life example:** Cell Stacking Line 1, Pack Assembly Line 3

---

### 5. `production_unit`
| Column | Type | Description |
|---|---|---|
| `id` | SERIAL PK | Auto-incremented integer ID |
| `production_line_id` | FK → production_line | Which line this unit sits on |
| `name` | VARCHAR(255) | Unit name (e.g. "Welding Station 4") |
| `description` | TEXT | Optional notes |

**Real-life example:** Laser Welder Station 2, Electrolyte Fill Station

---

### 6. `equipment_class`
| Column | Type | Description |
|---|---|---|
| `id` | SERIAL PK | Auto-incremented integer ID |
| `class_name` | VARCHAR(255) | Machine type (e.g. "Laser Welder") |
| `version` | VARCHAR(50) | Hardware/firmware version (e.g. "v3.1") |
| `description` | TEXT | Notes about the machine type |

**Real-life example:** Kuka Robot v2, ABB Laser Welder v1, Fanuc Welding Robot

---

### 7. `equipment_master`
The actual physical machine on the floor.

| Column | Type | Description |
|---|---|---|
| `id` | VARCHAR(50) PK | Your own asset ID (e.g. "MIXER-001") |
| `physical_asset_id` | VARCHAR(100) | Serial number / barcode on the machine |
| `production_unit_id` | FK → production_unit | Where the machine lives |
| `equipment_class_id` | FK → equipment_class | What type of machine it is |
| `operational_status` | VARCHAR(50) | `active`, `idle`, `maintenance`, `decommissioned` |
| `parent_equipment_id` | FK → equipment_master | For sub-components (e.g. a robot arm inside a cell) |

**Real-life example:** MIXER-001 is a physical Hosokawa mixer with serial "SN-AX-99012", installed in Slot 3 of Line A

---

### 8. `equipment_capability`
Optional technical capabilities of a machine.

| Column | Type | Description |
|---|---|---|
| `id` | SERIAL PK | Auto-incremented ID |
| `equipment_id` | FK → equipment_master | Which machine |
| `capability_name` | VARCHAR(100) | e.g. "max_torque", "max_speed" |
| `value` | VARCHAR(255) | e.g. "1200" |
| `uom` | VARCHAR(50) | e.g. "Nm", "RPM" |
| `description` | TEXT | Free text notes |

**Real-life example:** MIXER-001 has `max_speed = 3000 RPM`

---

### 9. `equipment_property`
Current runtime properties — things that change over time.

| Column | Type | Description |
|---|---|---|
| `id` | SERIAL PK | Auto-incremented ID |
| `equipment_id` | FK → equipment_master | Which machine |
| `property_name` | VARCHAR(100) | e.g. "maintenance_due_date" |
| `current_value` | VARCHAR(255) | Current value |
| `last_updated` | TIMESTAMP | When it was last set |

**Real-life example:** `calibration_status = "valid"`, `last_cleaned = "2026-02-20"`

---

### 10. `traceability_log` ⚡ TimescaleDB Hypertable
Every event that happens on a machine, linked to a work order / lot / operator.

| Column | Type | Description |
|---|---|---|
| `id` | SERIAL | Row ID (part of composite PK) |
| `equipment_id` | FK → equipment_master | Which machine generated the event |
| `event_type` | VARCHAR(100) | e.g. `start`, `stop`, `error_jam`, `qc_pass`, `qc_fail` |
| `work_order_id` | VARCHAR(100) | Which work order was active |
| `material_lot_id` | VARCHAR(100) | Which material lot was being processed |
| `operator_id` | VARCHAR(100) | Who was operating the machine |
| `event_time` | TIMESTAMPTZ PK | When the event occurred (used for partitioning) |

> **TimescaleDB:** This table is automatically partitioned by `event_time` into time-based "chunks". This means queries like "show all events in February" are extremely fast even with millions of rows.

**Real-life example:** At 10:32 AM, MIXER-001 logged `qc_pass` for lot `RM-CELL-001` under work order `WO-2026-0001`, operated by Kavya.

---

### 11. `material_definition`
A catalogue of all material types used in the factory.

| Column | Type | Description |
|---|---|---|
| `id` | SERIAL PK | Auto-incremented ID |
| `name` | VARCHAR(255) | Material name (e.g. "Lithium Carbonate Grade A") |
| `material_type` | VARCHAR(50) | One of: `RAW_MATERIAL`, `WORK_IN_PROGRESS`, `FINISHED_GOOD`, `PACKAGING` |
| `unit_of_measure` | VARCHAR(50) | e.g. `kg`, `liters`, `units` |
| `description` | TEXT | Optional notes |
| `created_at` | TIMESTAMPTZ | When added to catalogue |

**Real-life example:**
- `Lithium Carbonate` → RAW_MATERIAL, kg
- `Cell Module` → WORK_IN_PROGRESS, units
- `Battery Pack 100kWh` → FINISHED_GOOD, units

---

### 12. `material_lot`
A physical batch of material — a specific delivery of 500 kg of Lithium Carbonate, or a specific set of 10 finished battery packs.

| Column | Type | Description |
|---|---|---|
| `lot_id` | VARCHAR(100) PK | Your own lot ID (e.g. "RM-LI-2026-001") |
| `material_definition_id` | FK → material_definition | What material this lot is |
| `status` | VARCHAR(50) | `available`, `in_process`, `completed`, `quarantined`, `released`, `shipped`, `scrapped` |
| `quantity` | NUMERIC(15,4) | How much material is in this lot |
| `unit_of_measure` | VARCHAR(50) | kg, liters, units, etc. |
| `created_at` | TIMESTAMPTZ | When the lot entered the system |
| `updated_at` | TIMESTAMPTZ | When the status last changed |

**Real-life example:** Lot `RM-LI-2026-001` = 500 kg of Lithium Carbonate, status = `available`

---

### 13. `operator`
Every person who operates equipment on the factory floor.

| Column | Type | Description |
|---|---|---|
| `operator_id` | VARCHAR(100) PK | Your own operator ID (e.g. "OP-KAVYA-01") |
| `name` | VARCHAR(255) | Full name |
| `role` | VARCHAR(100) | e.g. "Line Operator", "Quality Inspector" |
| `shift` | VARCHAR(50) | `DAY_SHIFT`, `NIGHT_SHIFT`, `MORNING_SHIFT` |
| `status` | VARCHAR(50) | `active` or `inactive` |
| `created_at` | TIMESTAMPTZ | When registered |

---

### 14. `work_order`
A production job — "produce 100 battery cells using MIXER-001, operated by Arjun, output into lot WIP-CELL-001".

| Column | Type | Description |
|---|---|---|
| `work_order_id` | VARCHAR(100) PK | Your own WO ID (e.g. "WO-2026-0001") |
| `description` | TEXT | What this job is |
| `status` | VARCHAR(50) | `planned`, `in_progress`, `completed`, `closed` |
| `equipment_id` | FK → equipment_master | Which machine does the work |
| `operator_id` | FK → operator | Who operates it |
| `output_lot_id` | VARCHAR(100) | Which lot is produced as output |
| `planned_quantity` | NUMERIC | How many units planned |
| `actual_quantity` | NUMERIC | How many units actually produced |
| `unit_of_measure` | VARCHAR(50) | Units |
| `planned_start/end` | TIMESTAMPTZ | Scheduled window |
| `actual_start/end` | TIMESTAMPTZ | Real execution window |

---

### 15. `lot_genealogy`
**The heart of traceability.** Records which input lot(s) were consumed to produce which output lot.

| Column | Type | Description |
|---|---|---|
| `id` | SERIAL PK | Auto-incremented ID |
| `parent_lot_id` | FK → material_lot | The INPUT lot (what was consumed) |
| `child_lot_id` | FK → material_lot | The OUTPUT lot (what was produced) |
| `work_order_id` | FK → work_order | The WO that performed the transformation |
| `equipment_id` | FK → equipment_master | Which machine was used |
| `quantity_consumed` | NUMERIC | How much of the parent was used |
| `quantity_produced` | NUMERIC | How much of the child was produced |
| `event_time` | TIMESTAMPTZ | When the transformation happened |

**Unique constraint:** `(parent_lot_id, child_lot_id, work_order_id)` — no duplicate genealogy links.

**Real-life example:**  
`RM-LI-2026-001` (500 kg raw lithium) → 200 kg consumed → `WIP-CELL-001` (100 cell modules produced) via WO-2026-0001 on MIXER-001

---

### 16. `equipment_telemetry` ⚡ TimescaleDB Hypertable
Continuous sensor data from machines — temperature, pressure, speed, etc.

| Column | Type | Description |
|---|---|---|
| `id` | SERIAL | Row ID (part of composite PK) |
| `equipment_id` | FK → equipment_master | Which machine |
| `parameter_name` | VARCHAR(100) | Sensor name (e.g. "temperature", "pressure") |
| `value` | NUMERIC(15,4) | Sensor reading |
| `unit_of_measure` | VARCHAR(50) | e.g. "°C", "bar", "RPM" |
| `recorded_at` | TIMESTAMPTZ PK | Exact timestamp of reading (used for partitioning) |

> **TimescaleDB:** Same as `traceability_log` — partitioned by time for ultra-fast range queries over millions of data points.

**Real-life example:** MIXER-001 logs `temperature = 78.5 °C` every 5 seconds.

---

## Complete Real-Life Flow

Here is a **battery cell manufacturing scenario** showing how all tables connect:

```
1. SETUP (Master Data)
   ─────────────────────────────────────────
   enterprise:   "Ola Electric"
     └── site:   "Bengaluru Gigafactory"
           └── area:   "Cell Assembly Hall"
                 └── production_line:   "Cell Line 1"
                       └── production_unit:   "Mixing Station"
                             └── equipment_master:   MIXER-001 (Hosokawa v2, active)
                                   equipment_capability: max_speed = 3000 RPM
                                   equipment_property:   calibration_status = valid

2. MATERIALS (Catalogue + Lots)
   ─────────────────────────────────────────
   material_definition:
     - "Lithium Carbonate" → RAW_MATERIAL, kg
     - "Anode Slurry"      → WORK_IN_PROGRESS, liters
     - "Battery Cell 2170" → FINISHED_GOOD, units

   material_lot:
     - RM-LI-001   → 500 kg Lithium Carbonate   (status: available)
     - WIP-SL-001  → 0 liters Anode Slurry       (status: available, waiting to be filled)

3. PRODUCTION (Work Order + Genealogy)
   ─────────────────────────────────────────
   operator:       Arjun (OP-Arjun-01, DAY_SHIFT, Quality Control)

   work_order:     WO-2026-0001
                     equipment = MIXER-001
                     operator  = OP-Arjun-01
                     output    = WIP-SL-001
                     planned   = 200 kg → 100 liters slurry
                     status:   completed

   lot_genealogy:  RM-LI-001 (200 kg consumed) → WIP-SL-001 (100 liters produced)
                     via WO-2026-0001 on MIXER-001

4. EVENTS & TELEMETRY (Real-time)
   ─────────────────────────────────────────
   traceability_log:
     10:00 AM  MIXER-001  start        WO-2026-0001  RM-LI-001  OP-Arjun-01
     10:05 AM  MIXER-001  qc_check     WO-2026-0001  RM-LI-001  OP-Arjun-01
     10:30 AM  MIXER-001  qc_pass      WO-2026-0001  WIP-SL-001 OP-Arjun-01
     10:32 AM  MIXER-001  stop         WO-2026-0001  WIP-SL-001 OP-Arjun-01

   equipment_telemetry (every 5 sec during mixing):
     MIXER-001  temperature   78.5 °C    10:00:05
     MIXER-001  temperature   81.2 °C    10:00:10
     MIXER-001  rpm           2400 RPM   10:00:10
     ...

5. TRACEABILITY QUESTION: "What went into WIP-SL-001?"
   ─────────────────────────────────────────
   Backward trace via lot_genealogy:
     WIP-SL-001 was produced from:
       RM-LI-001 (200 kg Lithium Carbonate)
         processed on MIXER-001
         by OP-Arjun-01
         at 10:30 AM on 2026-02-25
```

---

## Indexes & Performance Notes

| Table | Index | Purpose |
|---|---|---|
| `material_lot` | `idx_material_lot_status` | Fast filter by lot status |
| `material_lot` | `idx_material_lot_material_def` | Fast join to material_definition |
| `work_order` | `idx_work_order_status` | Fast filter by WO status |
| `work_order` | `idx_work_order_equipment` | Fast lookup of all WOs on a machine |
| `lot_genealogy` | `idx_genealogy_parent/child/work_order` | Fast recursive trace queries |
| `equipment_telemetry` | `idx_telemetry_equip_param` | Fast time-range queries per sensor |
| `traceability_log` | TimescaleDB partition | Auto-partitioned by `event_time` |
| `equipment_telemetry` | TimescaleDB partition | Auto-partitioned by `recorded_at` |

---

## Access Control

A read-only PostgreSQL role `taksa_ai_reader` is granted `SELECT` on all tables. This is designed for an AI/analytics agent that needs to read data but never write it.
