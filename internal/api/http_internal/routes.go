package http_internal

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"whatsapp-bot/internal/controllers"
)

func SetupRoutes(r *chi.Mux, h *controllers.WhatsAppController) {
	r.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("ok"))
	})
	r.Get("/register", h.RegisterDevice)
}
