package models

// Event logs system actions or alerts (e.g., gate closed, sensor error).
type Event struct {
    ID        string `json:"id"`
    Type      string `json:"type"`
    Message   string `json:"message"`
    Timestamp int64  `json:"timestamp"`
}
