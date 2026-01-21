package data

import (
	"encoding/json"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nats-io/nats.go"
)

// UMHEvent defines the JSON structure we expect to receive from the factory.
// We map JSON keys (like "machine_id") to our Go fields.
type UMHEvent struct {
	EquipmentID   string `json:"equipment_id"`
	EventType     string `json:"event_type"`
	WorkOrderID   string `json:"work_order_id"`
	MaterialLotID string `json:"material_lot_id"`
	OperatorID    string `json:"operator_id"`
}

// StartConsumer subscribes to the NATS stream and saves data to Postgres
func (d *Data) StartConsumer(logger log.Logger) {
	l := log.NewHelper(logger)

	// 1. Get JetStream Context
	js, err := d.nats.JetStream()
	if err != nil {
		l.Errorf("Failed to get JetStream context: %v", err)
		return
	}

	// 2. Define the Callback Logic
	callback := func(msg *nats.Msg) {
		// A. Parse the JSON
		var event UMHEvent
		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			l.Errorf("Failed to parse JSON: %v | Data: %s", err, string(msg.Data))
			// We Ack here because retrying garbage data won't fix it.
			msg.Ack()
			return
		}

		// B. Map to Database Model (LogORM is defined in traceability.go)
		logEntry := LogORM{
			EquipmentID:   event.EquipmentID,
			EventType:     event.EventType,
			WorkOrderID:   event.WorkOrderID,
			MaterialLotID: event.MaterialLotID,
			OperatorID:    event.OperatorID,
			EventTime:     time.Now(), // We use the current server time
		}

		// C. Insert into Postgres
		// We use d.db directly because we are inside the 'data' package
		result := d.db.Create(&logEntry)
		if result.Error != nil {
			l.Errorf("Failed to insert into DB: %v", result.Error)
			// IMPORTANT: We do NOT msg.Ack() here.
			// NATS will see we failed and redeliver the message automatically.
			return
		}

		l.Infof("✅ SAVED to DB: Equipment=%s | Event=%s", logEntry.EquipmentID, logEntry.EventType)

		// D. Acknowledge success to NATS (remove from queue)
		msg.Ack()
	}

	// 3. Subscribe
	_, err = js.Subscribe("umh.>", callback, nats.BindStream("UMH_DATA"))
	if err != nil {
		l.Errorf("Failed to subscribe: %v", err)
	} else {
		l.Info("Successfully subscribed to UMH_DATA stream! Ready to save to DB...")
	}
}
