# Features Guide — Taksa Traceability

> This guide explains **what you can do** with the Taksa Traceability app, using real factory scenarios to show what each feature means in practice.

---

## What Is Traceability?

Traceability means being able to answer the question: **"What happened to this product, and how was it made?"**

If a battery pack fails in the field, you need to know:
- What raw materials went into it?
- Which machine processed it?
- Who was operating that machine?
- Were there any quality failures or sensor anomalies during production?
- Are there other products made from the same batch of raw material?

Taksa Traceability records all of this automatically as production happens, and then lets you query it instantly.

---

## Feature 1 — Plant Hierarchy Management

### What it is
Model your entire factory structure in the system, from company level all the way down to individual machines.

### The 5 levels
```
Enterprise  →  Site  →  Area  →  Production Line  →  Production Unit  →  Equipment
Ola Electric   Bengaluru   Cell Hall   Cell Line A      Mixing Station     MIXER-001
```

### What you can do
- Register a new factory location (site) under your company
- Add halls / zones (areas) to a site
- Add production lines inside each hall
- Add individual workstations (production units) on each line
- Attach physical machines to each workstation

### Real factory scenario
> Tata Motors sets up their new Chennai plant. They create:  
> **Tata Motors** → **Chennai EV Plant** → **Battery Bay** → **Module Assembly Line 2** → **Laser Welding Station** → **LASER-WELDER-03**  
> Now every work order and event on LASER-WELDER-03 is automatically linked to Chennai EV Plant.

---

## Feature 2 — Equipment Registry & Capabilities

### What it is
A complete digital record of every machine on the floor, including what type it is, where it is, and what it can do.

### What you can do
- Register machines with their physical serial number and asset tag
- Classify machines by type (e.g. "Kuka Robot v2", "Hosokawa Mixer")
- Track machine location (which workstation it's installed in)
- Track operational status: `active`, `idle`, `maintenance`, `decommissioned`
- Record technical capabilities (max speed, max torque, rated capacity)
- Record dynamic runtime properties (calibration status, last maintenance date)
- Model parent-child equipment (e.g. a robot arm as a sub-component of a robot cell)

### Real factory scenario
> **MIXER-001** is registered as:
> - Type: Hosokawa Planetary Mixer v2.1
> - Serial: SN-AX-99012
> - Location: Mixing Station 1, Cell Line A, Bengaluru
> - Status: Active
> - Capability: max_speed = 3000 RPM
> - Property: last_calibrated = "2026-01-15"
>
> If it goes in for maintenance, the status is updated to `maintenance`. The analytics dashboard will show it wasn't available during that window.

---

## Feature 3 — Material Management

### What it is
A two-level system: first define what materials exist (catalogue), then track each physical batch (lot).

### Material Types
| Type | Meaning | Example |
|---|---|---|
| `RAW_MATERIAL` | Comes in from a supplier | Lithium Carbonate, Graphite, NMP solvent |
| `WORK_IN_PROGRESS` | Made on the floor, not yet finished | Anode Slurry, Cell Module, Sub-assembly |
| `FINISHED_GOOD` | Final saleable product | Battery Pack 100kWh, EV Powerwall |
| `PACKAGING` | Materials used for packing | Bubble wrap, Pallet, Carton |

### Lot Statuses
| Status | Meaning |
|---|---|
| `available` | Ready to be used in production |
| `in_process` | Currently being worked on |
| `completed` | Work is done, awaiting QC decision |
| `quarantined` | Held for investigation (potential quality issue) |
| `released` | QC approved, can ship or use downstream |
| `shipped` | Sent to customer or warehouse |
| `scrapped` | Discarded due to quality failure |

### What you can do
- Add materials to the catalogue (one-time setup)
- Register incoming deliveries as lots with quantity and UOM
- Update a lot's status as it moves through the production process
- Filter lots by status or material type

### Real factory scenario
> A truck delivers 1000 kg of Lithium Carbonate. The receiving team creates:
> **Lot RM-LI-2026-042** = 1000 kg, status = `available`
>
> When it goes into the mixer: status → `in_process`  
> After QC approves the output: status → `released`  
> After it ships: status → `shipped`  
>
> If a defect is found mid-process: status → `quarantined` — the whole lot is held and investigators can trace exactly what was produced from it.

---

## Feature 4 — Operator Management

### What it is
Register every person who operates equipment and link them to work orders and events.

### What you can do
- Register operators with their name, role, and shift
- Track active vs. inactive operators
- Filter operators by shift for scheduling
- Every work order and event is linked to an operator — so you know **who was responsible** for each production step

### Real factory scenario
> If a batch of cells comes out with contamination, quality investigators can trace:
> - Which work order produced them
> - Which operator was running the machine at the time
> - Which shift they were on
> - What other lots that operator handled that day

---

## Feature 5 — Work Order Management

### What it is
A production job record — links together the machine, the operator, the input materials, and the output lot.

### Work Order Lifecycle
```
planned  →  in_progress  →  completed  →  closed
```

### What a work order captures
- Which machine performs the work
- Which operator runs it
- Planned vs. actual start/end time
- Planned vs. actual quantity produced
- Output lot ID (what was produced)
- Input lots (via genealogy links)

### What you can do
- Create a work order before production starts
- Update it to `in_progress` when work begins
- Record actual quantities and timestamps on completion
- View a full work order report including equipment class, operator name, and all input materials consumed

### Real factory scenario
> **WO-2026-0051**: Mix 500 kg Lithium Carbonate into 250 litres of Anode Slurry
> - Machine: MIXER-001
> - Operator: Arjun (Day Shift)
> - Planned: 08:00 – 10:00
> - Actual: 08:05 – 09:48 
> - Actual output: 247.5 litres
>
> This work order is the permanent production record for that batch.

---

## Feature 6 — Lot Genealogy (Core Traceability)

### What it is
**The most important feature.** Records the transformation relationships: which input lots were consumed to produce which output lot. This is the chain that makes full traceability possible.

### How it works
Every time a work order is completed, you register a genealogy link:
```
Input Lot (parent)  →  [Work Order + Machine]  →  Output Lot (child)
RM-LI-2026-042         WO-2026-0051  MIXER-001     WIP-SL-001
```

Multiple inputs → one output is supported (e.g. several raw materials mixed together).

### What you can do
- Register parent→child transformation links
- Each link records quantity consumed and quantity produced
- Build chains across multiple processing steps

### Real factory scenario
> **Step 1:** RM-LI-001 + RM-CARBON-001 → WIP-SLURRY-001 (mixing)  
> **Step 2:** WIP-SLURRY-001 → WIP-CELL-MODULE-001 (coating + pressing)  
> **Step 3:** WIP-CELL-MODULE-001 → FG-BATTERY-PACK-001 (assembly)
>
> The genealogy table now has 3 links. Any trace query on FG-BATTERY-PACK-001 can follow these links all the way back to the original raw materials.

---

## Feature 7 — Backward Trace (Where did this come from?)

### What it is
Given any lot — even a finished product — walk backwards through all genealogy links to find every upstream input that contributed to producing it.

### What you get
A list of all ancestor lots with:
- Depth (how many steps back)
- Material name and type
- Quantity consumed at each step
- Machine used at each step
- Operator responsible at each step
- Timestamp of each transformation

### Real factory scenario
> A customer reports that **BATTERY-PACK-A42** is failing.
>
> Backward trace shows:
> - Depth 1: WIP-CELL-MODULE-017 (10 cells, from ASSEMBLER-01, Ravi)
> - Depth 2: WIP-SLURRY-003 (250 litres, from MIXER-001, Kavya)
> - Depth 3: RM-LI-2026-042 (200 kg Lithium Carbonate, from supplier Sigma Aldrich)
>
> The investigation team can now check whether RM-LI-2026-042 was a bad batch, and query all other battery packs that used the same lot.

---

## Feature 8 — Forward Trace (Where did this end up?)

### What it is
Given an input lot — usually a raw material or WIP — trace forward through all genealogy links to find every downstream product it became part of.

### The most critical use case: **Product Recall**

If a raw material lot is found to be contaminated or out-of-spec, you need to immediately know every finished product that contains it — so you can recall or quarantine them.

### Real factory scenario
> Supplier notifies: batch **RM-LI-2026-042** of Lithium Carbonate may contain iron contamination.
>
> Forward trace on RM-LI-2026-042 shows:
> - WIP-SLURRY-003 (250L of slurry made from this lot)
>   - WIP-CELL-MODULE-017 (10 cells from this slurry)
>     - BATTERY-PACK-A42 ← **recall this**
>     - BATTERY-PACK-A43 ← **recall this**
>   - WIP-CELL-MODULE-018 (10 cells from this slurry)
>     - BATTERY-PACK-A44 ← **recall this**
>
> In minutes, quality team has the full recall list instead of manually searching through spreadsheets.

---

## Feature 9 — Full Genealogy Tree

### What it is
Returns the complete family tree of a lot — all ancestors, the lot itself, and all descendants — as a graph of **nodes** (lots) and **edges** (transformations).

### Ideal for
- Visual genealogy diagrams (connect to a graph visualization library)
- Exporting complete production history for regulatory compliance
- Audits and quality reports

### What you get
- **Nodes:** Every lot in the tree with its material type, status, and quantity
- **Edges:** Every transformation link with machine, work order, quantities, and timestamp

---

## Feature 10 — Equipment Process History

### What it is
For a given material lot, retrieve all the sensor data that was recorded on the machines that processed it. Links lot processing windows to machine telemetry readings.

### What you get
- Which machines processed the lot and during what time window
- Per-parameter statistics: min, max, average, reading count
- All individual raw sensor readings within that window

### Real factory scenario
> QC suspects WIP-SLURRY-003 was over-heated during mixing.
>
> Equipment history shows:
> - MIXER-001 processed this lot from 08:05 to 09:48
> - Temperature readings during that window:
>   - Min: 72.1°C, Max: 91.3°C (**exceeded spec of 85°C**), Avg: 78.5°C
>   - 432 readings recorded
>
> Exact evidence for the quality investigation, with no manual data collection needed.

---

## Feature 11 — Machine Performance Analytics

### What it is
KPI dashboard for any individual machine over any date range.

### What you get
| KPI | What it means |
|---|---|
| Total Work Orders | How many production jobs ran on this machine |
| Total Units Produced | Actual output volume |
| Avg Cycle Time | How long each job took on average (hours) |
| Total Active Hours | Total hours the machine was actually running |
| Utilization % | Active hours ÷ available hours in the window × 100 |
| Event Summary | Count of each event type (starts, stops, errors, qc_pass, etc.) |

### Real factory scenario
> Plant manager wants to know: "How busy was MIXER-001 last week?"
>
> Machine Performance for MIXER-001, Feb 24–28:
> - 12 work orders completed
> - 2,400 litres of slurry produced
> - Avg cycle time: 1.85 hours
> - Utilization: 39% (ran 18.5 hours out of 48 available hours)
> - Events: 12 starts, 10 qc_pass, 2 error_jam

---

## Feature 12 — Machine Comparison

### What it is
Compare multiple machines side by side over the same time window to identify best/worst performers.

### What you get
Same KPIs as Machine Performance but for all selected machines in one response, sorted by units produced.

### Real factory scenario
> Floor manager asks: "Which of our 3 mixers is most efficient?"
>
> Compare MIXER-001, MIXER-002, MIXER-003 over February:
> | Machine | Units Produced | Utilization % | Errors |
> |---|---|---|---|
> | MIXER-002 | 3,200 L | 66.7% | 0 |
> | MIXER-001 | 2,400 L | 39.0% | 2 |
> | MIXER-003 | 1,800 L | 37.5% | 5 |
>
> → MIXER-003 has the most errors and lowest output — schedule for maintenance review.

---

## Feature 13 — Production Trends

### What it is
Time-series output data bucketed by day, week, or month — ready for charting in dashboards.

### What you get
For each time bucket and each machine: units produced, work orders completed, average cycle time.

### Real factory scenario
> Analytics team building a monthly report wants to show production volume over Q1 2026.
>
> Production Trends with `granularity=monthly`, Jan–Mar 2026:
> - Jan: 18,400 units across all machines
> - Feb: 21,200 units (peak month)
> - Mar: 19,800 units
>
> Can also filter to a single machine to see its output trend over time.

---

## Feature 14 — Factory Dashboard

### What it is
A single API call that returns all top-level KPIs for a factory management dashboard.

### What you get
| KPI | Meaning |
|---|---|
| Total Units Produced | Sum of all actual_quantity across completed WOs |
| Work Orders Completed | Count of completed WOs in the window |
| Active Equipment Count | Distinct machines that had work in the window |
| Lots Released | Lots that passed QC |
| Lots Quarantined | Lots held for quality investigation |
| Quality Rate % | Released ÷ (Released + Quarantined) × 100 |
| Total Error Events | Count of `error*` events in traceability_log |
| Top 5 Performers | The 5 machines with highest output, with their utilization % |

### Real factory scenario
> CEO opens the factory dashboard for February 2026:
> - 4,820 units produced
> - 18 work orders completed
> - 5 machines active
> - Quality rate: 85.7% (12 released, 2 quarantined)
> - 3 error events recorded
> - Top performer: OVEN-001 at 80% utilization, 1,920 units

---

## Summary: What a Factory Manager Can Answer

| Question | Feature Used |
|---|---|
| Which raw materials went into this defective product? | Backward Trace |
| Which finished products contain this recalled raw material? | Forward Trace |
| Who was operating the machine when this batch was made? | Work Order Detail / Backward Trace |
| Were the machine temperatures in spec during processing? | Equipment Process History |
| Which machine had the most downtime this month? | Machine Performance / Comparison |
| Which machine produced the most output this week? | Machine Comparison / Dashboard |
| Is our quality rate improving or declining? | Dashboard / Production Trends |
| What is the full production lineage of this batch for regulatory audit? | Full Genealogy Tree |
| What work orders are currently in progress on Line A? | Work Order List (filter by status) |
| How many error events happened on MIXER-001 last week? | Machine Performance (Event Summary) |

---

Note :
> *Location is implicit (derived from work order → equipment → production hierarchy) rather than a standalone real-time location endpoint. All the data to answer "where is this lot right now?" is fully present in the system.
