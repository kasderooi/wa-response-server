package controllers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"whatsapp-bot/internal/usecase"
)

type WhatsAppController struct {
	client *usecase.WhatsAppClient
	logger *logrus.Logger
}

func NewWhatsAppController(client *usecase.WhatsAppClient, logger *logrus.Logger) *WhatsAppController {
	return &WhatsAppController{client: client, logger: logger}
}

func (c *WhatsAppController) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	c.logger.WithField("req", r.Host).Info("received register request")
	qr, err := c.client.Register()
	if err != nil {
		c.logger.WithField("req", r.Host).Errorf("error register %v", err)
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	err = qr.Write(256, w)
	if err != nil {
		c.logger.WithField("req", r.Host).Errorf("error writing %v", err)
		http.Error(w, "Failed to write QR code", http.StatusInternalServerError)
		return
	}
	c.logger.WithField("req", r.Host).Info("handled request")
}
