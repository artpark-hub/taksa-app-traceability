# API Reference — Taksa Traceability

> **Base URL:** `http://localhost:8000`  
> **Protocol:** HTTP/JSON (REST) · gRPC also available on port `9000`  
> **Auth:** None (add auth middleware for production)

All request bodies are JSON. All responses are JSON. Timestamps are RFC3339 format (e.g. `"2026-02-25T10:00:00Z"`).

---

## 1. Enterprise Management

### `POST /api/v1/traceability/enterprises`
Create a new enterprise.

**Request body:**
```json
{ "name": "Ola Electric", "description": "EV manufacturer" }
```
**Response:** `{ "id": 1 }`

---

### `GET /api/v1/traceability/enterprises`
List all enterprises.

**Response:** `{ "enterprises": [{ "id": 1, "name": "Ola Electric", "description": "..." }] }`

---

### `PATCH /api/v1/traceability/enterprises/{id}`
Update an enterprise's name or description.

**Request body:** `{ "name": "Ola Electric Ltd", "description": "Updated desc" }`  
**Response:** `{ "success": true }`

---

### `DELETE /api/v1/traceability/enterprises/{id}`
Delete an enterprise and cascade-delete all its sites, areas, lines, units, and equipment.

**Response:** `{ "success": true }`

---

## 2. Site Management

### `POST /api/v1/traceability/sites`
Create a factory site under an enterprise.

**Request body:**
```json
{ "enterprise_id": 1, "name": "Bengaluru Gigafactory", "location": "Bengaluru, India", "description": "Main plant" }
```
**Response:** `{ "id": 1 }`

---

### `GET /api/v1/traceability/sites?enterprise_id=1`
List all sites for an enterprise.

**Response:** `{ "sites": [{ "id": 1, "name": "Bengaluru Gigafactory", "location": "...", "description": "..." }] }`

---

### `PATCH /api/v1/traceability/sites/{id}`
Update site location.

**Request body:** `{ "location": "Electronic City, Bengaluru" }`  
**Response:** `{ "success": true }`

---

### `DELETE /api/v1/traceability/sites/{id}`
Delete a site and all its sub-resources.

**Response:** `{ "success": true }`

---

## 3. Area Management

### `POST /api/v1/traceability/areas`
Create a production area inside a site.

**Request body:**
```json
{ "site_id": 1, "name": "Cell Assembly Hall", "description": "Battery cell production zone" }
```
**Response:** `{ "id": 1 }`

---

### `GET /api/v1/traceability/areas?site_id=1`
List all areas in a site.

**Response:** `{ "areas": [{ "id": 1, "name": "Cell Assembly Hall", "description": "..." }] }`

---

### `PATCH /api/v1/traceability/areas/{id}`
Update area name/description.

**Response:** `{ "success": true }`

---

### `DELETE /api/v1/traceability/areas/{id}`
Delete an area.

**Response:** `{ "success": true }`

---

## 4. Production Line Management

### `POST /api/v1/traceability/lines`
Create a production line inside an area.

**Request body:**
```json
{ "area_id": 1, "name": "Cell Line A", "description": "Main cell stacking line" }
```
**Response:** `{ "id": 1 }`

---

### `GET /api/v1/traceability/lines?area_id=1`
List all lines in an area.

---

### `PATCH /api/v1/traceability/lines/{id}` · `DELETE /api/v1/traceability/lines/{id}`
Update or delete a line.

---

## 5. Production Unit Management

### `POST /api/v1/traceability/production-units`
Create a production slot/station on a line.

**Request body:**
```json
{ "production_line_id": 1, "name": "Mixing Station 1", "description": "Slurry mixing slot" }
```
**Response:** `{ "id": 1 }`

---

### `GET /api/v1/traceability/production-units?line_id=1`
List all units on a line.

---

### `PATCH /api/v1/traceability/production-units/{id}` · `DELETE /api/v1/traceability/production-units/{id}`
Update or delete a unit.

---

## 6. Equipment Class Management

### `POST /api/v1/traceability/equipment-classes`
Register a type/model of machine.

**Request body:**
```json
{ "class_name": "Hosokawa Mixer", "version": "v2.1", "description": "Planetary mixer for electrode slurry" }
```
**Response:** `{ "id": 1 }`

---

### `GET /api/v1/traceability/equipment-classes`
List all equipment classes.

---

### `PATCH /api/v1/traceability/equipment-classes/{id}` · `DELETE /api/v1/traceability/equipment-classes/{id}`
Update or delete a class.

---

## 7. Equipment Master Management

### `POST /api/v1/traceability/equipment`
Register a physical machine on the floor.

**Request body:**
```json
{
  "id": "MIXER-001",
  "physical_asset_id": "SN-AX-99012",
  "production_unit_id": 1,
  "equipment_class_id": 1,
  "operational_status": "active",
  "parent_equipment_id": ""
}
```
**Response:** `{ "id": "MIXER-001" }`

---

### `GET /api/v1/traceability/equipment?line_id=1`
List all machines on a line (or filter by `parent_equipment_id`).

---

### `PATCH /api/v1/traceability/equipment/{id}`
Update machine status or move it to a new production unit.

**Request body:** `{ "operational_status": "maintenance", "production_unit_id": 2 }`  
**Response:** `{ "success": true }`

---

### `DELETE /api/v1/traceability/equipment/{id}`
Remove a machine from the registry.

---

## 8. Equipment Capability Management

### `POST /api/v1/traceability/equipment/{equipment_id}/capabilities`
Add a technical capability to a machine.

**Request body:**
```json
{ "capability_name": "max_speed", "value": "3000", "uom": "RPM", "description": "Maximum mixing speed" }
```
**Response:** `{ "id": 1 }`

---

### `GET /api/v1/traceability/equipment/{equipment_id}/capabilities`
List all capabilities of a machine.

---

### `PATCH /api/v1/traceability/equipment/{equipment_id}/capabilities/{id}` · `DELETE ...`
Update a capability value or remove it.

---

## 9. Equipment Property Management

### `POST /api/v1/traceability/equipment/{equipment_id}/properties`
Set a dynamic runtime property on a machine.

**Request body:**
```json
{ "property_name": "calibration_status", "current_value": "valid" }
```
**Response:** `{ "id": 1 }`

---

### `GET /api/v1/traceability/equipment/{equipment_id}/properties`
List all properties of a machine.

---

### `PATCH /api/v1/traceability/equipment/{equipment_id}/properties/{id}` · `DELETE ...`
Update a property value or remove it.

---

## 10. Traceability Log

### `POST /api/v1/traceability/logs`
Log an event that occurred on a machine.

**Request body:**
```json
{
  "equipment_id": "MIXER-001",
  "event_type": "qc_pass",
  "work_order_id": "WO-2026-0001",
  "material_lot_id": "WIP-SL-001",
  "operator_id": "OP-KAVYA-01"
}
```
**Response:** `{ "id": 42 }`

---

### `GET /api/v1/traceability/logs?work_order_id=WO-2026-0001`
List all events logged for a work order.

**Response:**
```json
{
  "logs": [
    { "equipmentId": "MIXER-001", "eventType": "start", "eventTime": "2026-02-25T10:00:00Z", "operatorId": "OP-KAVYA-01" }
  ]
}
```

---

## 11. Material Definition Management

### `POST /api/v1/traceability/material-definitions`
Add a material type to the catalogue.

**Request body:**
```json
{ "name": "Lithium Carbonate Grade A", "material_type": "RAW_MATERIAL", "unit_of_measure": "kg", "description": "High-purity Li2CO3" }
```
**Response:** `{ "id": 1 }`

> **Valid material_type values:** `RAW_MATERIAL` · `WORK_IN_PROGRESS` · `FINISHED_GOOD` · `PACKAGING`

---

### `GET /api/v1/traceability/material-definitions?material_type=RAW_MATERIAL`
List material definitions, optionally filtered by type.

---

## 12. Material Lot Management

### `POST /api/v1/traceability/lots`
Register a physical batch of material.

**Request body:**
```json
{ "lot_id": "RM-LI-2026-001", "material_definition_id": 1, "quantity": 500, "unit_of_measure": "kg" }
```
**Response:** `{ "lotId": "RM-LI-2026-001" }`

---

### `GET /api/v1/traceability/lots/{lot_id}`
Get full details of a lot including its material definition info.

**Response:**
```json
{
  "lotId": "RM-LI-2026-001",
  "status": "available",
  "quantity": 500,
  "unitOfMeasure": "kg",
  "createdAt": "2026-02-25T08:00:00Z",
  "materialName": "Lithium Carbonate Grade A",
  "materialType": "RAW_MATERIAL"
}
```

---

### `GET /api/v1/traceability/lots?status=available&material_type=RAW_MATERIAL`
List lots with optional filters.

---

### `PATCH /api/v1/traceability/lots/{lot_id}/status`
Update a lot's status (e.g. after QC decision).

**Request body:** `{ "status": "released" }`  
**Response:** `{ "success": true }`

> **Valid status values:** `available` · `in_process` · `completed` · `quarantined` · `released` · `shipped` · `scrapped`

---

## 13. Operator Management

### `POST /api/v1/traceability/operators`
Register a factory operator.

**Request body:**
```json
{ "operator_id": "OP-KAVYA-01", "name": "Kavya R.", "role": "Quality Control", "shift": "DAY_SHIFT" }
```
**Response:** `{ "operatorId": "OP-KAVYA-01" }`

---

### `GET /api/v1/traceability/operators?shift=DAY_SHIFT&status=active`
List operators with optional filters.

---

### `PATCH /api/v1/traceability/operators/{operator_id}`
Update operator role, shift, or status.

**Request body:** `{ "shift": "NIGHT_SHIFT", "status": "active" }`  
**Response:** `{ "success": true }`

---

## 14. Work Order Management

### `POST /api/v1/traceability/work-orders`
Create a production work order.

**Request body:**
```json
{
  "work_order_id": "WO-2026-0001",
  "description": "Anode slurry mixing batch 1",
  "equipment_id": "MIXER-001",
  "operator_id": "OP-KAVYA-01",
  "output_lot_id": "WIP-SL-001",
  "planned_quantity": 100,
  "unit_of_measure": "liters",
  "planned_start": "2026-02-25T10:00:00Z",
  "planned_end": "2026-02-25T12:00:00Z"
}
```
**Response:** `{ "workOrderId": "WO-2026-0001" }`

---

### `GET /api/v1/traceability/work-orders/{work_order_id}`
Get full work order detail including equipment class, operator name, and all input lots consumed.

**Response:**
```json
{
  "workOrderId": "WO-2026-0001",
  "status": "completed",
  "equipmentId": "MIXER-001",
  "equipmentClassName": "Hosokawa Mixer",
  "operatorName": "Kavya R.",
  "operatorShift": "DAY_SHIFT",
  "actualQuantity": 98.5,
  "inputLots": [
    { "lotId": "RM-LI-2026-001", "materialName": "Lithium Carbonate", "quantityConsumed": 200 }
  ]
}
```

---

### `GET /api/v1/traceability/work-orders?status=completed&equipment_id=MIXER-001`
List work orders with optional filters.

---

### `PATCH /api/v1/traceability/work-orders/{work_order_id}/status`
Update work order status (start it, complete it, etc.).

**Request body:**
```json
{
  "status": "completed",
  "actual_quantity": 98.5,
  "actual_start": "2026-02-25T10:05:00Z",
  "actual_end": "2026-02-25T11:55:00Z"
}
```
**Response:** `{ "success": true }`

---

## 15. Lot Genealogy Registration

### `POST /api/v1/traceability/genealogy`
Record that a parent lot was consumed to produce a child lot.

**Request body:**
```json
{
  "parent_lot_id": "RM-LI-2026-001",
  "child_lot_id": "WIP-SL-001",
  "work_order_id": "WO-2026-0001",
  "equipment_id": "MIXER-001",
  "quantity_consumed": 200,
  "quantity_produced": 100
}
```
**Response:** `{ "id": 1 }`

---

## 16. Traceability Queries

### `GET /api/v1/traceability/trace/backward/{lot_id}`
**Trace backward** — find every upstream input lot (and sub-inputs) that contributed to producing this lot.

**Example:** `GET /api/v1/traceability/trace/backward/BATTERY-PACK-001`

**Response:**
```json
{
  "queriedLotId": "BATTERY-PACK-001",
  "traceDirection": "backward",
  "nodes": [
    {
      "depth": 1,
      "lotId": "WIP-CELL-MODULE-001",
      "relatedLotId": "BATTERY-PACK-001",
      "materialName": "Cell Module",
      "materialType": "WORK_IN_PROGRESS",
      "lotStatus": "completed",
      "quantity": 10,
      "quantityUsed": 10,
      "workOrderId": "WO-2026-0003",
      "equipmentId": "ASSEMBLER-001",
      "equipmentClassName": "Pack Assembler",
      "operatorName": "Ravi K.",
      "eventTime": "2026-02-25T14:00:00Z"
    }
  ]
}
```

---

### `GET /api/v1/traceability/trace/forward/{lot_id}`
**Trace forward** — find every downstream lot produced from this input lot (where did this raw material end up?).

**Response:** Same structure as backward trace but `traceDirection: "forward"`.

---

### `GET /api/v1/traceability/trace/full/{lot_id}`
**Full genealogy tree** — returns the complete family tree as nodes (lots) + edges (links).

**Response:**
```json
{
  "centerLotId": "WIP-CELL-001",
  "nodes": [
    { "lotId": "RM-LI-001", "materialName": "Lithium Carbonate", "materialType": "RAW_MATERIAL", "status": "completed", "quantity": 500 },
    { "lotId": "WIP-CELL-001", "materialName": "Cell Module", "materialType": "WORK_IN_PROGRESS", "status": "released", "quantity": 100 }
  ],
  "edges": [
    { "sourceLotId": "RM-LI-001", "targetLotId": "WIP-CELL-001", "workOrderId": "WO-2026-0001", "equipmentId": "MIXER-001", "quantityConsumed": 200, "quantityProduced": 100 }
  ]
}
```

---

### `GET /api/v1/traceability/trace/equipment-history/{lot_id}`
Show all sensor readings and telemetry summaries for machines that processed this lot.

**Response:**
```json
{
  "lotId": "WIP-SL-001",
  "parameters": [
    { "equipmentId": "MIXER-001", "parameterName": "temperature", "minValue": 72.1, "maxValue": 85.3, "avgValue": 78.5, "readingCount": 360 }
  ],
  "readings": [
    { "equipmentId": "MIXER-001", "parameterName": "temperature", "value": 78.5, "unitOfMeasure": "°C", "recordedAt": "2026-02-25T10:00:05Z" }
  ]
}
```

---

## 17. Historical Analytics

### `GET /api/v1/traceability/analytics/machine-performance`
KPI report for a single machine over a date range.

**Query parameters:**
| Param | Required | Example |
|---|---|---|
| `equipment_id` | Yes | `MIXER-001` |
| `from_time` | Yes | `2026-02-25T00:00:00Z` |
| `to_time` | Yes | `2026-02-26T00:00:00Z` |

**Response:**
```json
{
  "equipmentId": "MIXER-001",
  "equipmentClassName": "Hosokawa Mixer",
  "operationalStatus": "active",
  "fromTime": "2026-02-25T00:00:00Z",
  "toTime": "2026-02-26T00:00:00Z",
  "totalWorkOrders": 3,
  "totalUnitsProduced": 285.5,
  "avgCycleTimeHours": 1.85,
  "totalActiveHours": 5.55,
  "utilizationPct": 23.1,
  "eventSummary": [
    { "eventType": "start", "count": 3 },
    { "eventType": "qc_pass", "count": 3 }
  ]
}
```

---

### `POST /api/v1/traceability/analytics/machine-comparison`
Compare multiple machines side by side over the same time window.

**Request body:**
```json
{
  "equipment_ids": ["MIXER-001", "OVEN-001", "ROBOT-001"],
  "from_time": "2026-02-25T00:00:00Z",
  "to_time": "2026-02-26T00:00:00Z"
}
```

**Response:**
```json
{
  "fromTime": "2026-02-25T00:00:00Z",
  "toTime": "2026-02-26T00:00:00Z",
  "machines": [
    { "equipmentId": "OVEN-001", "totalUnitsProduced": 480, "utilizationPct": 62.5, "errorCount": 0 },
    { "equipmentId": "MIXER-001", "totalUnitsProduced": 285, "utilizationPct": 23.1, "errorCount": 1 }
  ]
}
```

---

### `GET /api/v1/traceability/analytics/production-trends`
Time-bucketed output chart data — use for charts/graphs in dashboards.

**Query parameters:**
| Param | Required | Example |
|---|---|---|
| `from_time` | Yes | `2026-02-01T00:00:00Z` |
| `to_time` | Yes | `2026-02-28T00:00:00Z` |
| `granularity` | No (default: `daily`) | `daily` / `weekly` / `monthly` |
| `equipment_id` | No | `MIXER-001` |

**Response:**
```json
{
  "granularity": "daily",
  "dataPoints": [
    { "period": "2026-02-25T00:00:00Z", "equipmentId": "MIXER-001", "unitsProduced": 285.5, "workOrdersCompleted": 3, "avgCycleTimeHours": 1.85 },
    { "period": "2026-02-25T00:00:00Z", "equipmentId": "OVEN-001", "unitsProduced": 480, "workOrdersCompleted": 2, "avgCycleTimeHours": 3.2 }
  ]
}
```

---

### `GET /api/v1/traceability/analytics/dashboard`
Top-level factory dashboard summary.

**Query parameters:**
| Param | Required | Example |
|---|---|---|
| `from_time` | Yes | `2026-02-01T00:00:00Z` |
| `to_time` | Yes | `2026-02-28T00:00:00Z` |

**Response:**
```json
{
  "totalUnitsProduced": 4820.5,
  "totalWorkOrdersCompleted": 18,
  "activeEquipmentCount": 5,
  "lotsReleased": 12,
  "lotsQuarantined": 2,
  "qualityRatePct": 85.71,
  "totalErrorEvents": 3,
  "topPerformers": [
    { "equipmentId": "OVEN-001", "totalUnitsProduced": 1920, "utilizationPct": 80.0 },
    { "equipmentId": "ROBOT-001", "totalUnitsProduced": 1440, "utilizationPct": 60.0 }
  ]
}
```

---

## Quick Reference Table

| # | Method | Endpoint | What it does |
|---|---|---|---|
| 1 | POST | `/enterprises` | Create enterprise |
| 2 | GET | `/enterprises` | List enterprises |
| 3 | POST | `/sites` | Create site |
| 4 | GET | `/sites` | List sites |
| 5 | POST | `/areas` | Create area |
| 6 | POST | `/lines` | Create production line |
| 7 | POST | `/production-units` | Create production unit |
| 8 | POST | `/equipment-classes` | Register equipment class |
| 9 | POST | `/equipment` | Register physical machine |
| 10 | POST | `/equipment/{id}/capabilities` | Add machine capability |
| 11 | POST | `/equipment/{id}/properties` | Set machine property |
| 12 | POST | `/logs` | Log a machine event |
| 13 | GET | `/logs` | List events for a WO |
| 14 | POST | `/material-definitions` | Add material to catalogue |
| 15 | POST | `/lots` | Register a material lot |
| 16 | GET | `/lots/{lot_id}` | Get lot detail |
| 17 | PATCH | `/lots/{lot_id}/status` | Update lot status |
| 18 | POST | `/operators` | Register operator |
| 19 | POST | `/work-orders` | Create work order |
| 20 | GET | `/work-orders/{id}` | Get WO detail with input lots |
| 21 | PATCH | `/work-orders/{id}/status` | Update WO status |
| 22 | POST | `/genealogy` | Record parent→child lot link |
| 23 | GET | `/trace/backward/{lot_id}` | Trace all upstream inputs |
| 24 | GET | `/trace/forward/{lot_id}` | Trace all downstream outputs |
| 25 | GET | `/trace/full/{lot_id}` | Full genealogy tree (nodes+edges) |
| 26 | GET | `/trace/equipment-history/{lot_id}` | Sensor history for a lot |
| 27 | GET | `/analytics/machine-performance` | Single machine KPIs |
| 28 | POST | `/analytics/machine-comparison` | Side-by-side machine KPIs |
| 29 | GET | `/analytics/production-trends` | Time-bucketed output trends |
| 30 | GET | `/analytics/dashboard` | Factory dashboard summary |

> All endpoints are prefixed with `/api/v1/traceability/`
