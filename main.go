package main

import (
	"log"
	"net/http"

	"goparking/handlers"
	"goparking/store"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Initialize router and database
	r := chi.NewRouter()
	db := store.NewStore("parking.db")

	// Handlers
	slotHandler := &handlers.SlotHandler{Store: db}
	commandHandler := &handlers.CommandHandler{Store: db}

	// Device endpoint
	r.Post("/api/device/update", slotHandler.DeviceUpdate)

	// Slot CRUD endpoints
	r.Post("/api/slots", slotHandler.CreateSlot)
	r.Get("/api/slots", slotHandler.ListSlots)
	r.Get("/api/slots/{id}", slotHandler.GetSlot)
	r.Put("/api/slots/{id}", slotHandler.UpdateSlot)
	r.Delete("/api/slots/{id}", slotHandler.DeleteSlot)

	// Command endpoints
	r.Get("/api/command", commandHandler.GetCommand)
	r.Post("/api/command", commandHandler.SetCommand)

	log.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
