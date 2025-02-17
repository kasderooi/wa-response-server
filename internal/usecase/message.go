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
	base = "%s/nHoi ik gebruik geen WhatsApp meer, je kan mij bereiken op Signal. \n https://signal.org/download/"
	auto = "· _Automated response_ ·"
)

// HandleMessage processes incoming messages
func HandleMessage(client *whatsmeow.Client, msg *events.Message) {
	if !shouldReply(client, msg) {
		return
	}

	reply := fmt.Sprintf(base, auto)
	_, err := client.SendMessage(context.Background(), msg.Info.Chat, &waE2E.Message{
		Conversation: &reply,
	})
	if err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
	} else {
		fmt.Printf("Replied to %s\n", msg.Info.Sender.User)
	}
}

func shouldReply(client *whatsmeow.Client, msg *events.Message) bool {
	if msg.Info.IsFromMe || strings.Contains(msg.Message.String(), auto) {
		return false
	}
	if !msg.Info.IsGroup {
		return true
	}
	for _, mentionedJID := range msg.Message.ExtendedTextMessage.GetContextInfo().GetMentionedJID() {
		if mentionedJID == client.Store.ID.ToNonAD().String() {
			return true
		}
	}
	return false
}
