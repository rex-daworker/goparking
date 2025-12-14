package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"goparking/models"
	"goparking/store"
	"goparking/validators"

	"github.com/go-chi/chi/v5"
)

type SlotHandler struct {
	Store *store.Store
}

// ---------- DTO ----------
type SlotCreateUpdateDTO struct {
	ID       string `json:"id"`
	Distance int    `json:"distance"`
	Status   string `json:"status"`
}

// ============================================================
// 🔹 ADMIN / TESTING CRUD ENDPOINTS
// ============================================================

// POST /api/slots
func (h *SlotHandler) CreateSlot(w http.ResponseWriter, r *http.Request) {
	var dto SlotCreateUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateSlotInput(dto.ID, dto.Distance, dto.Status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slot := models.ParkingSlot{
		ID:         dto.ID,
		Distance:   dto.Distance,
		Status:     dto.Status,
		LastUpdate: time.Now().UnixMilli(),
	}

	if err := h.Store.UpsertSlot(slot); err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(slot)
}

// GET /api/slots
func (h *SlotHandler) ListSlots(w http.ResponseWriter, r *http.Request) {
	slots, err := h.Store.ListSlots()
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(slots)
}

// GET /api/slots/{id}
func (h *SlotHandler) GetSlot(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	slot, err := h.Store.GetSlot(id)
	if err != nil {
		http.Error(w, "Slot not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(slot)
}

// PUT /api/slots/{id}
func (h *SlotHandler) UpdateSlot(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var dto SlotCreateUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if dto.ID == "" {
		dto.ID = id
	}

	if err := validators.ValidateSlotInput(dto.ID, dto.Distance, dto.Status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slot := models.ParkingSlot{
		ID:         dto.ID,
		Distance:   dto.Distance,
		Status:     dto.Status,
		LastUpdate: time.Now().UnixMilli(),
	}

	if err := h.Store.UpsertSlot(slot); err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(slot)
}

// DELETE /api/slots/{id}
func (h *SlotHandler) DeleteSlot(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.Store.DeleteSlot(id); err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================
// 🔹 DEVICE ENDPOINT (Arduino / ESP32 Sensor Update)
// ============================================================

// POST /api/device/update
func (h *SlotHandler) DeviceUpdate(w http.ResponseWriter, r *http.Request) {
	var dto SlotCreateUpdateDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if dto.ID == "" {
		http.Error(w, "Slot ID is required", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateSlotInput(dto.ID, dto.Distance, dto.Status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slot := models.ParkingSlot{
		ID:         dto.ID,
		Distance:   dto.Distance,
		Status:     dto.Status,
		LastUpdate: time.Now().UnixMilli(),
	}

	if err := h.Store.UpsertSlot(slot); err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// history
	occupied := 0
	if dto.Status == "occupied" {
		occupied = 1
	}
	_ = h.Store.AddHistoryPoint(dto.ID, occupied)

	// simple auto-gate logic
	if dto.Status == "occupied" {
		_ = h.Store.SetCommandAction("close")
	} else if dto.Status == "free" {
		_ = h.Store.SetCommandAction("open")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
}
