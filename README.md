# Taksa MES Traceability

The Traceability application for Taksa MES provides comprehensive capabilities to track and trace manufacturing units throughout their entire lifecycle. It serves as a source of truth for the history, location, and status of every material and product on the shop floor.

## Features

### End-to-End Unit Tracking
Gain complete visibility into your production lines by tracking manufacturing units from raw material intake to finished goods. The system assigns unique identifiers to each unit, ensuring precise tracking across all processing stages.

### Genealogy & Hierarchy
Understand the composition and history of your products with robust genealogy features.
*   **Parent-Child Relationships**: Track how individual components (child units) are assembled into larger assemblies (parent units).
*   **Backward & Forward Traceability**: Quickly trace a finished product back to its specific batch of raw materials, or identify all products affected by a specific lot of components.

### Real-time Visibility
Monitor the pulse of your operations with real-time updates.
*   **Status Tracking**: Instantly view the current status (e.g., In-Progress, On-Hold, Completed, Scrapped) of any unit.
*   **Location Tracking**: Know exactly where materials are located within your facility at any given time.

## Getting Started

### Prerequisites
*   Taksa Platform Core Services
*   Docker & Docker Compose

### Installation
Clone the repository:
```bash
git clone https://github.com/artpark-hub/taksa-app-traceability.git
cd taksa-app-traceability
```

Refer to the central [Taksa Deployments](https://github.com/artpark-hub/taksa-deployments) repository for full system deployment instructions.

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details on how to submit pull requests, report issues, and suggest improvements.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.