package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "goparking/models"
    "goparking/store"
)

type EventHandler struct {
    Store *store.Store
}

// POST /api/events
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
    var dto models.Event
    if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    if dto.ID == "" {
        http.Error(w, "Event ID is required", http.StatusBadRequest)
        return
    }
    if dto.Type == "" {
        dto.Type = "info"
    }
    dto.Timestamp = time.Now().UnixMilli()

    if err := h.Store.AddEvent(dto); err != nil {
        http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    log.Println("Event created:", dto.Type, dto.Message)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(dto)
}

// GET /api/events
func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
    events, err := h.Store.ListEvents()
    if err != nil {
        http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(events)
}
