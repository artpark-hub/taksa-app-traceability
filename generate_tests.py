import os
import time
import shutil

BASE_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), "tests", "api")

if os.path.exists(BASE_DIR):
    shutil.rmtree(BASE_DIR)
os.makedirs(BASE_DIR, exist_ok=True)

TS = int(time.time())

with open(os.path.join(BASE_DIR, "bruno.json"), "w") as f:
    f.write('{\n  "version": "1",\n  "name": "Taksa Traceability",\n  "type": "collection"\n}\n')

env_dir = os.path.join(BASE_DIR, "environments")
os.makedirs(env_dir, exist_ok=True)

with open(os.path.join(env_dir, "Dev VM.bru"), "w") as f:
    f.write('vars {\n  baseUrl: http://localhost:8000\n}\n')

def create_bru(folder, filename, name, method, url, body=None, assertions=None, seq=1, vars_to_set=None):
    folder_path = os.path.join(BASE_DIR, folder)
    os.makedirs(folder_path, exist_ok=True)

    # Bruno requires a folder.bru metadata file in every directory
    folder_meta = os.path.join(folder_path, "folder.bru")
    if not os.path.exists(folder_meta):
        with open(folder_meta, "w") as fm:
            fm.write(f"meta {{\n  name: {folder}\n}}\n")

    file_content = f"""meta {{
  name: {name}
  type: http
  seq: {seq}
}}

{method} {{
  url: {url}
  body: {'json' if body else 'none'}
  auth: none
}}
"""
    if body:
        # Wrap the raw JSON payload inside Bruno's body:json { } block
        file_content += "\nbody:json {\n" + body.strip() + "\n}\n"

    if assertions:
        file_content += "\nassert {\n"
        for a in assertions:
            file_content += f"  {a}\n"
        file_content += "}\n"

    if vars_to_set:
        file_content += "\nscript:post-response {\n"
        for var, path in vars_to_set.items():
            file_content += f'  bru.setVar("{var}", res.body.{path});\n'
        file_content += "}\n"

    with open(os.path.join(folder_path, filename), "w") as f:
        f.write(file_content)

# --- 01. Enterprise ---
create_bru("01_Enterprise", "1. Create Enterprise.bru", "1. Create Enterprise", "post",
    "{{baseUrl}}/api/v1/traceability/enterprises",
    '{\n    "name": "Global Motors",\n    "description": "Global Headquarters"\n  }',
    ["res.status: eq 200", "res.body.id: isNumber"], 1, {"enterpriseId": "id"})

# --- 02. Site ---
create_bru("02_Site", "1. Create Site.bru", "1. Create Site", "post",
    "{{baseUrl}}/api/v1/traceability/sites",
    '{\n    "enterprise_id": {{enterpriseId}},\n    "name": "Detroit Plant",\n    "location": "Michigan, USA"\n  }',
    ["res.status: eq 200", "res.body.id: isNumber"], 1, {"siteId": "id"})

# --- 03. Area ---
create_bru("03_Area", "1. Create Area.bru", "1. Create Area", "post",
    "{{baseUrl}}/api/v1/traceability/areas",
    '{\n    "site_id": {{siteId}},\n    "name": "Body Shop",\n    "description": "Welding Area"\n  }',
    ["res.status: eq 200", "res.body.id: isNumber"], 1, {"areaId": "id"})

# --- 04. Line ---
create_bru("04_Line", "1. Create Line.bru", "1. Create Line", "post",
    "{{baseUrl}}/api/v1/traceability/lines",
    '{\n    "area_id": {{areaId}},\n    "name": "Line A",\n    "description": "Main Line"\n  }',
    ["res.status: eq 200", "res.body.id: isNumber"], 1, {"lineId": "id"})

# --- 05. Unit ---
create_bru("05_Unit", "1. Create Unit.bru", "1. Create Unit", "post",
    "{{baseUrl}}/api/v1/traceability/production-units",
    '{\n    "production_line_id": {{lineId}},\n    "name": "Station-01",\n    "description": "Robot Slot"\n  }',
    ["res.status: eq 200", "res.body.id: isNumber"], 1, {"unitId": "id"})

# --- 06. Equipment Class ---
create_bru("06_EquipmentClass", "1. Create Class.bru", "1. Create Class", "post",
    "{{baseUrl}}/api/v1/traceability/equipment-classes",
    '{\n    "class_name": "Kuka Robot",\n    "version": "v2.0"\n  }',
    ["res.status: eq 200", "res.body.id: isNumber"], 1, {"classId": "id"})

# --- 07. Equipment Master ---
create_bru("07_Equipment", "1. Register Equipment.bru", "1. Register Equipment", "post",
    "{{baseUrl}}/api/v1/traceability/equipment",
    f'{{\n    "id": "ROBOT-{TS}",\n    "physical_asset_id": "SN-{TS}",\n    "production_unit_id": {{{{unitId}}}},\n    "equipment_class_id": {{{{classId}}}},\n    "operational_status": "active"\n  }}',
    ["res.status: eq 200"], 1)

# --- 09. Material Definitions ---
create_bru("09_Materials", "1. Create Definition.bru", "1. Create Definition", "post",
    "{{baseUrl}}/api/v1/traceability/material-definitions",
    '{\n    "name": "Raw Silicon",\n    "material_type": "RAW_MATERIAL",\n    "unit_of_measure": "kg"\n  }',
    ["res.status: eq 200", "res.body.id: isNumber"], 1, {"defId": "id"})

create_bru("09_Materials", "2. Create WIP Definition.bru", "2. Create WIP Definition", "post",
    "{{baseUrl}}/api/v1/traceability/material-definitions",
    '{\n    "name": "Silicon Slurry",\n    "material_type": "WORK_IN_PROGRESS",\n    "unit_of_measure": "liters"\n  }',
    ["res.status: eq 200", "res.body.id: isNumber"], 2, {"wipDefId": "id"})

# --- 10. Operators ---
create_bru("10_Operators", "1. Create Operator.bru", "1. Create Operator", "post",
    "{{baseUrl}}/api/v1/traceability/operators",
    f'{{\n    "operator_id": "OP-TEST-{TS}",\n    "name": "Test Operator",\n    "role": "Quality Control",\n    "shift": "DAY_SHIFT"\n  }}',
    ["res.status: eq 200"], 1)

# --- 11. Material Lots ---
create_bru("11_Lots", "1. Create Parent Lot.bru", "1. Create Parent Lot", "post",
    "{{baseUrl}}/api/v1/traceability/lots",
    f'{{\n    "lot_id": "LOT-SIL-{TS}",\n    "material_definition_id": {{{{defId}}}},\n    "quantity": 100,\n    "unit_of_measure": "kg"\n  }}',
    ["res.status: eq 200"], 1)

create_bru("11_Lots", "2. Create Child Lot.bru", "2. Create Child Lot", "post",
    "{{baseUrl}}/api/v1/traceability/lots",
    f'{{\n    "lot_id": "WIP-{TS}",\n    "material_definition_id": {{{{wipDefId}}}},\n    "quantity": 0,\n    "unit_of_measure": "liters"\n  }}',
    ["res.status: eq 200"], 2)

# --- 12. Work Orders ---
create_bru("12_WorkOrders", "1. Create Work Order.bru", "1. Create Work Order", "post",
    "{{baseUrl}}/api/v1/traceability/work-orders",
    f'{{\n    "work_order_id": "WO-TEST-{TS}",\n    "description": "Test Production Run",\n    "equipment_id": "ROBOT-{TS}",\n    "operator_id": "OP-TEST-{TS}",\n    "output_lot_id": "WIP-{TS}",\n    "planned_quantity": 10,\n    "unit_of_measure": "units",\n    "planned_start": "2026-02-27T10:00:00Z",\n    "planned_end": "2026-02-27T18:00:00Z"\n  }}',
    ["res.status: eq 200"], 1)

# --- 13. Genealogy ---
create_bru("13_Genealogy", "1. Register Link.bru", "1. Register Link", "post",
    "{{baseUrl}}/api/v1/traceability/genealogy",
    f'{{\n    "parent_lot_id": "LOT-SIL-{TS}",\n    "child_lot_id": "WIP-{TS}",\n    "work_order_id": "WO-TEST-{TS}",\n    "equipment_id": "ROBOT-{TS}",\n    "quantity_consumed": 50,\n    "quantity_produced": 10\n  }}',
    ["res.status: eq 200", "res.body.id: isNumber"], 1)

# --- 15. Queries ---
create_bru("15_Queries", "1. Trace Backward.bru", "1. Trace Backward", "get",
    f'{{{{baseUrl}}}}/api/v1/traceability/trace/backward/WIP-{TS}',
    None, ["res.status: eq 200", 'res.body.traceDirection: eq "backward"'], 1)

create_bru("15_Queries", "2. Full Genealogy.bru", "2. Full Genealogy", "get",
    f'{{{{baseUrl}}}}/api/v1/traceability/trace/full/WIP-{TS}',
    None, ["res.status: eq 200", "res.body.nodes: isArray"], 2)

# --- 99. Cleanup ---
create_bru("99_Cleanup", "9. Delete Site.bru", "9. Delete Site", "delete",
    "{{baseUrl}}/api/v1/traceability/sites/{{siteId}}",
    None, ["res.status: eq 200", "res.body.success: isTruthy"], 9)

create_bru("99_Cleanup", "10. Delete Enterprise.bru", "10. Delete Enterprise", "delete",
    "{{baseUrl}}/api/v1/traceability/enterprises/{{enterpriseId}}",
    None, ["res.status: eq 200", "res.body.success: isTruthy"], 10)

print(f"✅ All Bruno files generated! TS={TS}")