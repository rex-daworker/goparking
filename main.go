package main

import (
	"log"
	"net/http"
	"time"

	"goparking/handlers"
	"goparking/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. Init SQLite store
	st := store.NewStore("parking.db")

	// 2. Handlers
	slotHandler := &handlers.SlotHandler{Store: st}
	cmdHandler := &handlers.CommandHandler{Store: st}

	// 3. Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health
	r.Get("/health", handlers.Health)

	// API
	r.Route("/api", func(r chi.Router) {

		// Slots CRUD (admin/testing)
		r.Get("/slots", slotHandler.ListSlots)
		r.Post("/slots", slotHandler.CreateSlot)
		r.Get("/slots/{id}", slotHandler.GetSlot)
		r.Put("/slots/{id}", slotHandler.UpdateSlot)
		r.Delete("/slots/{id}", slotHandler.DeleteSlot)

		// Device → backend
		r.Post("/device/update", slotHandler.DeviceUpdate)

		// Command control (cloud ↔ device)
		r.Get("/command", cmdHandler.GetCommand)  // Arduino polls this
		r.Post("/command", cmdHandler.SetCommand) // control panel / manual
	})

	addr := ":8080"
	log.Printf("Go backend running on http://localhost%v\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
		log.Println()
	}
}
