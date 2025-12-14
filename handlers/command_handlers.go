package handlers

import (
    "encoding/json"
    "net/http"

    "goparking/models"
    "goparking/store"
)

type CommandHandler struct {
    Store *store.Store
}

// GET /api/command (Arduino polls this)
func (h *CommandHandler) GetCommand(w http.ResponseWriter, r *http.Request) {
    cmd, err := h.Store.GetCommand()
    if err != nil {
        http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cmd)
}

// POST /api/command (manual control from dashboard)
func (h *CommandHandler) SetCommand(w http.ResponseWriter, r *http.Request) {
    var dto models.Command
    if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Defaults
    if dto.Threshold <= 0 {
        dto.Threshold = 30
    }
    if dto.Action == "" {
        dto.Action = "none"
    }

    if err := h.Store.SetCommand(dto); err != nil {
        http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(dto)
}
