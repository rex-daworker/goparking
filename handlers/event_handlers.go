package handlers

import (
    "encoding/json"
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
    dto.Timestamp = time.Now().UnixMilli()
    if err := h.Store.AddEvent(dto); err != nil {
        http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(dto)
}

// GET /api/events
func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
    events, err := h.Store.ListEvents()
    if err != nil {
        http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(events)
    
}
