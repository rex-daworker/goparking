package models

// Kept in slot.go for your original structure; file provided for clarity.
// If you prefer one file, you can remove this and keep Command in slot.go.
type Command struct {
    Action    string `json:"action"`
    Threshold int    `json:"threshold"`
}
