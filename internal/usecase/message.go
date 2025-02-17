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
	auto   = "* Automated response *"
	signal = "[Signal](https://signal.org/download/)"
)

// HandleMessage processes incoming messages
func HandleMessage(client *whatsmeow.Client, msg *events.Message) {
	// Ignore messages from the bot itself and prevent infinite loops
	if msg.Info.IsFromMe || strings.Contains(msg.Message.String(), auto) {
		return
	}

	// Check if this is a group message where the bot is tagged
	mentioned := false
	if msg.Info.IsGroup {
		for _, mentionedJID := range msg.Message.ExtendedTextMessage.GetContextInfo().GetMentionedJid() {
			if mentionedJID == client.Store.ID.ToNonAD().String() {
				mentioned = true
				break
			}
		}
	}

	// If it's a group message but the bot isn't tagged, ignore it
	if msg.Info.IsGroup && !mentioned {
		return
	}

	// Prepare the reply
	reply := fmt.Sprintf("%s\nHoi ik gebruik geen WhatsApp meer, je kan mij bereiken op %s.", auto, signal)

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
