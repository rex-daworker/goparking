package validators

import (
	"errors"
	"strings"
)

// ValidateSlotInput checks distance and status before DB writes.
func ValidateSlotInput(id string, distance int, status string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id is required")
	}
	if distance < 0 || distance > 500 {
		return errors.New("distance must be between 0 and 500 cm")
	}
	status = strings.ToLower(status)
	if status != "free" && status != "occupied" && status != "unknown" {
		return errors.New(`status must be "free", "occupied" or "unknown"`)
	}
	return nil
}
