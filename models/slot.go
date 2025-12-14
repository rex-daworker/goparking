package models

// ParkingSlot represents one physical parking space monitored by a sensor.
type ParkingSlot struct {
    ID         string `json:"id"`
    Distance   int    `json:"distance"`
    Status     string `json:"status"`
    LastUpdate int64  `json:"lastUpdate"`
}

// Command is what the cloud sends to the device (gate control, etc.).
type Command struct {
    Action    string `json:"action"`
    Threshold int    `json:"threshold"`
}
