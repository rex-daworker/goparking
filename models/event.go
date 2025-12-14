package models

import "time"

type Event struct {
	ID         int64     `json:"id"`
	SlotID     string    `json:"slotId"`
	Status     string    `json:"status"`
	DistanceCM int       `json:"distanceCm"`
	CreatedAt  time.Time `json:"createdAt"`
}
