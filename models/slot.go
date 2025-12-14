package models

// ParkingSlot represents one physical parking space monitored by a sensor.
type ParkingSlot struct {
    ID         string `json:"id"`
    Distance   int    `json:"distance"`
    Status     string `json:"status"`
    LastUpdate int64  `json:"lastUpdate"`
}

