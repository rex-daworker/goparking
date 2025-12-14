package models

type Event struct {
    ID        string `json:"id"`
    Type      string `json:"type"`
    Message   string `json:"message"`
    Timestamp int64  `json:"timestamp"`
}
