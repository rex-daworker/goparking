package validators

import "errors"

// ValidateSlotInput checks slot input fields for correctness.
func ValidateSlotInput(id string, distance int, status string, deviceID string, sensorStatus string) error {
    if id == "" {
        return errors.New("slot ID required")
    }
    if distance < 0 {
        return errors.New("distance must be non-negative")
    }
    if status != "free" && status != "occupied" && status != "unknown" {
        return errors.New("status must be free, occupied, or unknown")
    }
    if deviceID == "" {
        return errors.New("device ID required")
    }
    if sensorStatus != "active" && sensorStatus != "inactive" {
        return errors.New("sensor status must be active or inactive")
    }
    return nil
}
