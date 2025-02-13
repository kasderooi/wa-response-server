package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"go.mau.fi/whatsmeow/store"
	"log"
	"sync"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsAppClient struct {
	container *sqlstore.Container
	clients   map[string]*whatsmeow.Client // Track active clients
	mu        sync.Mutex                   // Ensure thread-safety
}

// NewWhatsAppClient Initializes WhatsAppClient with an empty clients map
func NewWhatsAppClient(db *sql.DB, dialect, logLevel string) *WhatsAppClient {
	dbLog := waLog.Stdout("Database", logLevel, true)
	container := sqlstore.NewWithDB(db, dialect, dbLog)

	return &WhatsAppClient{
		container: container,
		clients:   make(map[string]*whatsmeow.Client),
	}
}

// Start all known clients on program startup
func (c *WhatsAppClient) Start() error {
	err := c.container.Upgrade()
	if err != nil {
		return err
	}

	devices, err := c.container.GetAllDevices()
	if err != nil {
		return err
	}

	for _, device := range devices {
		c.ConnectClient(device)
	}

	return nil
}

// ConnectClient a specific client and store it in the map
func (c *WhatsAppClient) ConnectClient(device *store.Device) {
	client := whatsmeow.NewClient(device, waLog.Stdout("Client", "DEBUG", true))

	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			HandleMessage(client, v)
		}
	})

	c.mu.Lock()
	c.clients[device.ID.String()] = client
	c.mu.Unlock()

	err := client.Connect()
	if err != nil {
		log.Printf("Error connecting client %s: %v", device.ID.String(), err)
		return
	}

	log.Printf("Client %s connected successfully!", device.ID.String())
}

// Register a new device
func (c *WhatsAppClient) Register() (*qrcode.QRCode, error) {
	device := c.container.NewDevice()

	client := whatsmeow.NewClient(device, waLog.Stdout("Client", "DEBUG", true))
	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			HandleMessage(client, v)
		}
	})

	qrChan, err := client.GetQRChannel(context.Background())
	if err != nil {
		return nil, err
	}

	err = client.Connect()
	if err != nil {
		return nil, err
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			fmt.Println("Scan the following QR code to log in:")
			return qrcode.New(evt.Code, qrcode.Highest)
		}
	}

	c.mu.Lock()
	c.clients[device.ID.String()] = client
	c.mu.Unlock()

	return nil, nil
}

// Stop Gracefully stops all clients
func (c *WhatsAppClient) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, client := range c.clients {
		client.Disconnect()
		log.Printf("Client %s disconnected", id)
	}

	c.clients = make(map[string]*whatsmeow.Client)
	return nil
}
