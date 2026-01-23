package data

import (
	"encoding/json"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nats-io/nats.go"
)

type UMHEvent struct {
	EquipmentID   string `json:"equipment_id"`
	EventType     string `json:"event_type"`
	WorkOrderID   string `json:"work_order_id"`
	MaterialLotID string `json:"material_lot_id"`
	OperatorID    string `json:"operator_id"`
}

func (d *Data) StartConsumer(logger log.Logger) {
	l := log.NewHelper(logger)

	js, err := d.nats.JetStream()
	if err != nil {
		l.Errorf("Failed to get JetStream context: %v", err)
		return
	}

	callback := func(msg *nats.Msg) {

		var event UMHEvent
		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			l.Errorf("Failed to parse JSON: %v | Data: %s", err, string(msg.Data))

			msg.Ack()
			return
		}

		logEntry := LogORM{
			EquipmentID:   event.EquipmentID,
			EventType:     event.EventType,
			WorkOrderID:   event.WorkOrderID,
			MaterialLotID: event.MaterialLotID,
			OperatorID:    event.OperatorID,
			EventTime:     time.Now(),
		}

		result := d.db.Create(&logEntry)
		if result.Error != nil {
			l.Errorf("Failed to insert into DB: %v", result.Error)

			return
		}

		l.Infof("✅ SAVED to DB: Equipment=%s | Event=%s", logEntry.EquipmentID, logEntry.EventType)

		msg.Ack()
	}

	_, err = js.Subscribe("umh.>", callback, nats.BindStream("UMH_DATA"))
	if err != nil {
		l.Errorf("Failed to subscribe: %v", err)
	} else {
		l.Info("Successfully subscribed to UMH_DATA stream! Ready to save to DB...")
	}
}
