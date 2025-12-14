package main

import (
    "goparking/handlers"
    "goparking/store"
    "net/http"

    "github.com/go-chi/chi/v5"
)

func main() {
    r := chi.NewRouter()
    db := store.NewStore("parking.db")
    slotHandler := &handlers.SlotHandler{Store: db}

    // Device endpoint
    r.Post("/api/device/update", slotHandler.DeviceUpdate)

    // Admin endpoints
    r.Post("/api/slots", slotHandler.CreateSlot)
    r.Get("/api/slots", slotHandler.ListSlots)
    r.Get("/api/slots/{id}", slotHandler.GetSlot)
    r.Put("/api/slots/{id}", slotHandler.UpdateSlot)
    r.Delete("/api/slots/{id}", slotHandler.DeleteSlot)
    
    
    log.Println("Server started on http://localhost:8080")

    http.ListenAndServe(":8080", r)
}
