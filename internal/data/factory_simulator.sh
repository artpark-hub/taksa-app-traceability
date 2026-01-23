
while true; do
  TEMP=$(( ( RANDOM % 50 ) + 50 ))
  
  /home/kavya/go/bin/nats pub umh.v1.cnc001 "{\"equipment_id\": \"CNC-001\", \"event_type\": \"STATUS_UPDATE\", \"work_order_id\": \"WO-999\", \"operator_id\": \"OP-KAVYA\", \"temperature\": $TEMP}"
  
  echo "--- Sent Update: Temp $TEMP ---"
  sleep 5
done