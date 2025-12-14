package main

import (
    "net/http"
    "goparking/store"
    "goparking/handlers"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    dbPath := "parking.db"
    store := store.NewStore(dbPath)

    // Slot and Command handlers
    slotHandler := &handlers.SlotHandler{Store: store}
    commandHandler := &handlers.CommandHandler{Store: store}

    //  NEW: Event handler
    eventHandler := &handlers.EventHandler{Store: store}

    // Routes
    r.Post("/api/device/update", slotHandler.CreateSlot)
    r.Get("/api/slots", slotHandler.ListSlots)
    r.Post("/api/command", commandHandler.SetCommand)
    r.Get("/api/command", commandHandler.GetCommand)

    //  NEW: Event routes
    r.Post("/api/events", eventHandler.CreateEvent)
    r.Get("/api/events", eventHandler.ListEvents)

    http.ListenAndServe(":8080", r)
}
