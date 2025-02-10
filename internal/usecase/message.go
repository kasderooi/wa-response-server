package usecase

import (
	"context"
	"fmt"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"strings"
)

const (
	auto = "* Automated response *"
)

// HandleMessage processes incoming messages
func HandleMessage(client *whatsmeow.Client, msg *events.Message) {
	// Ignore messages from the bot itself and group messages
	if msg.Info.IsFromMe || msg.Info.IsGroup || strings.Contains(msg.Message.String(), auto) {
		return
	}

	// Prepare a reply
	reply := auto + "\nIm not using Whats-App anymore, you can reach me on Signal or by sms"

	// Send the reply
	_, err := client.SendMessage(context.Background(), msg.Info.Sender, &waE2E.Message{
		Conversation: &reply,
	})
	if err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
	} else {
		fmt.Printf("Replied to %s\n", msg.Info.Sender.User)
	}
}
