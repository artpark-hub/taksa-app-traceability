import os

# Base Directory
BASE_DIR = os.path.expanduser("~/artpark-hub/taksa-app-traceability/tests/api")
os.makedirs(BASE_DIR, exist_ok=True)

# Create bruno.json
with open(os.path.join(BASE_DIR, "bruno.json"), "w") as f:
    f.write('{\n  "version": "1",\n  "name": "Taksa Traceability",\n  "type": "collection"\n}\n')

# Create Environment
env_dir = os.path.join(BASE_DIR, "environments")
os.makedirs(env_dir, exist_ok=True)
with open(os.path.join(env_dir, "Dev VM.bru"), "w") as f:
    f.write('vars {\n  baseUrl: http://localhost:8000\n}\n')

# Helper to create a .bru file
def create_bru(folder, filename, name, method, url, body=None, assertions=None, seq=1):
    folder_path = os.path.join(BASE_DIR, folder)
    os.makedirs(folder_path, exist_ok=True)
    
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
        file_content += f"\nbody:json {{\n  {body}\n}}\n"
    
    if assertions:
        file_content += "\nassert {\n"
        for a in assertions:
            file_content += f"  {a}\n"
        file_content += "}\n"

    with open(os.path.join(folder_path, filename), "w") as f:
        f.write(file_content)

# --- 1. Enterprise ---
create_bru("1_Enterprise", "1. Create Enterprise.bru", "1. Create Enterprise", "post", "{{baseUrl}}/api/v1/traceability/enterprises", 
           '{\n    "name": "Global Motors",\n    "description": "Global Headquarters"\n  }', 
           ["res.status: eq 200", "res.body.id: isNumber"], 1)

create_bru("1_Enterprise", "2. List Enterprises.bru", "2. List Enterprises", "get", "{{baseUrl}}/api/v1/traceability/enterprises", 
           None, ["res.status: eq 200", "res.body.enterprises: isArray"], 2)

create_bru("1_Enterprise", "3. Update Enterprise.bru", "3. Update Enterprise", "patch", "{{baseUrl}}/api/v1/traceability/enterprises/1", 
           '{\n    "name": "Global Motors Inc."\n  }', 
           ["res.status: eq 200", "res.body.success: isTruthy"], 3)

# --- 2. Site ---
create_bru("2_Site", "1. Create Site.bru", "1. Create Site", "post", "{{baseUrl}}/api/v1/traceability/sites", 
           '{\n    "enterprise_id": 1,\n    "name": "Detroit Plant",\n    "location": "Michigan, USA"\n  }', 
           ["res.status: eq 200", "res.body.id: isNumber"], 1)

create_bru("2_Site", "2. List Sites.bru", "2. List Sites", "get", "{{baseUrl}}/api/v1/traceability/sites?enterprise_id=1", 
           None, ["res.status: eq 200", "res.body.sites: isArray"], 2)

create_bru("2_Site", "3. Update Site.bru", "3. Update Site", "patch", "{{baseUrl}}/api/v1/traceability/sites/1", 
           '{\n    "location": "Ohio, USA"\n  }', 
           ["res.status: eq 200", "res.body.success: isTruthy"], 3)

# --- 3. Area ---
create_bru("3_Area", "1. Create Area.bru", "1. Create Area", "post", "{{baseUrl}}/api/v1/traceability/areas", 
           '{\n    "site_id": 1,\n    "name": "Body Shop",\n    "description": "Welding Area"\n  }', 
           ["res.status: eq 200", "res.body.id: isNumber"], 1)

create_bru("3_Area", "2. List Areas.bru", "2. List Areas", "get", "{{baseUrl}}/api/v1/traceability/areas?site_id=1", 
           None, ["res.status: eq 200", "res.body.areas: isArray"], 2)

create_bru("3_Area", "3. Update Area.bru", "3. Update Area", "patch", "{{baseUrl}}/api/v1/traceability/areas/1", 
           '{\n    "name": "Paint Shop"\n  }', 
           ["res.status: eq 200", "res.body.success: isTruthy"], 3)

# --- 4. Line ---
create_bru("4_Line", "1. Create Line.bru", "1. Create Line", "post", "{{baseUrl}}/api/v1/traceability/lines", 
           '{\n    "area_id": 1,\n    "name": "Line A",\n    "description": "Main Line"\n  }', 
           ["res.status: eq 200", "res.body.id: isNumber"], 1)

create_bru("4_Line", "2. List Lines.bru", "2. List Lines", "get", "{{baseUrl}}/api/v1/traceability/lines?area_id=1", 
           None, ["res.status: eq 200", "res.body.lines: isArray"], 2)

create_bru("4_Line", "3. Update Line.bru", "3. Update Line", "patch", "{{baseUrl}}/api/v1/traceability/lines/1", 
           '{\n    "name": "Line B"\n  }', 
           ["res.status: eq 200", "res.body.success: isTruthy"], 3)

# --- 5. Unit ---
create_bru("5_Unit", "1. Create Unit.bru", "1. Create Unit", "post", "{{baseUrl}}/api/v1/traceability/production-units", 
           '{\n    "production_line_id": 1,\n    "name": "Station-01",\n    "description": "Robot Slot"\n  }', 
           ["res.status: eq 200", "res.body.id: isNumber"], 1)

create_bru("5_Unit", "2. List Units.bru", "2. List Units", "get", "{{baseUrl}}/api/v1/traceability/production-units?line_id=1", 
           None, ["res.status: eq 200", "res.body.units: isArray"], 2)

create_bru("5_Unit", "3. Update Unit.bru", "3. Update Unit", "patch", "{{baseUrl}}/api/v1/traceability/production-units/1", 
           '{\n    "name": "Station-01-B"\n  }', 
           ["res.status: eq 200", "res.body.success: isTruthy"], 3)

# --- 6. Equipment ---
create_bru("6_Equipment", "1. Create Class.bru", "1. Create Class", "post", "{{baseUrl}}/api/v1/traceability/equipment-classes", 
           '{\n    "class_name": "Kuka Robot",\n    "version": "v2.0"\n  }', 
           ["res.status: eq 200", "res.body.id: isNumber"], 1)

create_bru("6_Equipment", "2. Register Equip Main.bru", "2. Register Equip Main", "post", "{{baseUrl}}/api/v1/traceability/equipment", 
           '{\n    "id": "ROBOT-99",\n    "physical_asset_id": "SN-9988",\n    "production_unit_id": 1,\n    "equipment_class_id": 1,\n    "operational_status": "Idle"\n  }', 
           ["res.status: eq 200", 'res.body.id: eq "ROBOT-99"'], 2)

create_bru("6_Equipment", "3. Register Equip Child.bru", "3. Register Equip Child", "post", "{{baseUrl}}/api/v1/traceability/equipment", 
           '{\n    "id": "SENSOR-01",\n    "physical_asset_id": "SN-SENS-001",\n    "production_unit_id": 1,\n    "equipment_class_id": 1,\n    "operational_status": "Running",\n    "parent_equipment_id": "ROBOT-99"\n  }', 
           ["res.status: eq 200", 'res.body.id: eq "SENSOR-01"'], 3)

create_bru("6_Equipment", "4. List Equip.bru", "4. List Equip", "get", "{{baseUrl}}/api/v1/traceability/equipment?line_id=1", 
           None, ["res.status: eq 200", "res.body.equipment: isArray"], 4)

create_bru("6_Equipment", "5. Update Equip.bru", "5. Update Equip", "patch", "{{baseUrl}}/api/v1/traceability/equipment/ROBOT-99", 
           '{\n    "operational_status": "Maintenance"\n  }', 
           ["res.status: eq 200", "res.body.success: isTruthy"], 5)

create_bru("6_Equipment", "6. Update Class.bru", "6. Update Class", "patch", "{{baseUrl}}/api/v1/traceability/equipment-classes/1", 
           '{\n    "version": "v2.1"\n  }', 
           ["res.status: eq 200", "res.body.success: isTruthy"], 6)

# --- 7. Capabilities ---
create_bru("7_Capabilities", "1. Add Capability.bru", "1. Add Capability", "post", "{{baseUrl}}/api/v1/traceability/equipment/ROBOT-99/capabilities", 
           '{\n    "capability_name": "Max Load",\n    "value": "100",\n    "uom": "kg"\n  }', 
           ["res.status: eq 200", "res.body.id: isNumber"], 1)

create_bru("7_Capabilities", "2. List Capabilities.bru", "2. List Capabilities", "get", "{{baseUrl}}/api/v1/traceability/equipment/ROBOT-99/capabilities", 
           None, ["res.status: eq 200", "res.body.capabilities: isArray"], 2)

create_bru("7_Capabilities", "3. Update Capability.bru", "3. Update Capability", "patch", "{{baseUrl}}/api/v1/traceability/equipment/ROBOT-99/capabilities/1", 
           '{\n    "value": "150"\n  }', 
           ["res.status: eq 200", "res.body.success: isTruthy"], 3)

# --- 8. Properties ---
create_bru("8_Properties", "1. Set Property.bru", "1. Set Property", "post", "{{baseUrl}}/api/v1/traceability/equipment/ROBOT-99/properties", 
           '{\n    "property_name": "Temp Setpoint",\n    "current_value": "200"\n  }', 
           ["res.status: eq 200", "res.body.id: isNumber"], 1)

create_bru("8_Properties", "2. List Properties.bru", "2. List Properties", "get", "{{baseUrl}}/api/v1/traceability/equipment/ROBOT-99/properties", 
           None, ["res.status: eq 200", "res.body.properties: isArray"], 2)

create_bru("8_Properties", "3. Update Property.bru", "3. Update Property", "patch", "{{baseUrl}}/api/v1/traceability/equipment/ROBOT-99/properties/1", 
           '{\n    "current_value": "210"\n  }', 
           ["res.status: eq 200", "res.body.success: isTruthy"], 3)

# --- 9. Logs ---
create_bru("9_Logs", "1. Log Event.bru", "1. Log Event", "post", "{{baseUrl}}/api/v1/traceability/logs", 
           '{\n    "equipment_id": "ROBOT-99",\n    "event_type": "Production Start",\n    "work_order_id": "WO-2025-001"\n  }', 
           ["res.status: eq 200", "res.body.id: isNumber"], 1)

create_bru("9_Logs", "2. View Logs.bru", "2. View Logs", "get", "{{baseUrl}}/api/v1/traceability/logs?work_order_id=WO-2025-001", 
           None, ["res.status: eq 200", "res.body.logs: isArray"], 2)

# --- 99. Cleanup ---
create_bru("99_Cleanup", "1. Delete Property.bru", "1. Delete Property", "delete", "{{baseUrl}}/api/v1/traceability/equipment/ROBOT-99/properties/1", 
           None, ["res.status: eq 200", "res.body.success: isTruthy"], 1)

create_bru("99_Cleanup", "2. Delete Capability.bru", "2. Delete Capability", "delete", "{{baseUrl}}/api/v1/traceability/equipment/ROBOT-99/capabilities/1", 
           None, ["res.status: eq 200", "res.body.success: isTruthy"], 2)

create_bru("99_Cleanup", "3. Delete Child Equip.bru", "3. Delete Child Equip", "delete", "{{baseUrl}}/api/v1/traceability/equipment/SENSOR-01", 
           None, ["res.status: eq 200", "res.body.success: isTruthy"], 3)

create_bru("99_Cleanup", "4. Delete Main Equip.bru", "4. Delete Main Equip", "delete", "{{baseUrl}}/api/v1/traceability/equipment/ROBOT-99", 
           None, ["res.status: eq 200", "res.body.success: isTruthy"], 4)

create_bru("99_Cleanup", "5. Delete Class.bru", "5. Delete Class", "delete", "{{baseUrl}}/api/v1/traceability/equipment-classes/1", 
           None, ["res.status: eq 200", "res.body.success: isTruthy"], 5)

create_bru("99_Cleanup", "6. Delete Unit.bru", "6. Delete Unit", "delete", "{{baseUrl}}/api/v1/traceability/production-units/1", 
           None, ["res.status: eq 200", "res.body.success: isTruthy"], 6)

create_bru("99_Cleanup", "7. Delete Line.bru", "7. Delete Line", "delete", "{{baseUrl}}/api/v1/traceability/lines/1", 
           None, ["res.status: eq 200", "res.body.success: isTruthy"], 7)

create_bru("99_Cleanup", "8. Delete Area.bru", "8. Delete Area", "delete", "{{baseUrl}}/api/v1/traceability/areas/1", 
           None, ["res.status: eq 200", "res.body.success: isTruthy"], 8)

create_bru("99_Cleanup", "9. Delete Site.bru", "9. Delete Site", "delete", "{{baseUrl}}/api/v1/traceability/sites/1", 
           None, ["res.status: eq 200", "res.body.success: isTruthy"], 9)

create_bru("99_Cleanup", "10. Delete Enterprise.bru", "10. Delete Enterprise", "delete", "{{baseUrl}}/api/v1/traceability/enterprises/1", 
           None, ["res.status: eq 200", "res.body.success: isTruthy"], 10)

print("✅ All Bruno files generated successfully with CORRECT syntax!")
