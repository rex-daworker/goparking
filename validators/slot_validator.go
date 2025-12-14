package validators

import (
	"errors"
	"strings"
)

// ValidateSlotInput checks distance and status before DB writes.
func ValidateSlotInput(id string, distance int, status string, deviceID string, sensorStatus string) error {
    if id == "" {
        return errors.New("slot ID required")
    }
    if deviceID == "" {
        return errors.New("device ID required")
    }
    if sensorStatus != "active" && sensorStatus != "inactive" {
        return errors.New("sensor status must be active or inactive")
    }
    if distance < 0 {
        return errors.New("distance must be non-negative")
    }
    return nil
}

