package main

import (
	"context"
	"log"
	"os"

	"github.com/alecthomas/repr"

	unifi "github.com/silverlyra/unifi-go/api"
	"github.com/silverlyra/unifi-go/api/protect"
)

func main() {
	ctx := context.Background()

	username := os.Getenv("UNIFI_USERNAME")
	password := os.Getenv("UNIFI_PASSWORD")

	baseClient, err := unifi.New("https://reverie.localdomain", unifi.Login{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatalf("Failed to create Unifi client: %v\n", err)
	}

	client, err := protect.New(baseClient)
	if err != nil {
		log.Fatalf("Failed to create Unifi Protect client: %v\n", err)
	}

	receiver, err := client.ReceiveUpdates(ctx)
	if err != nil {
		log.Fatalf("Failed to receive Unifi Protect updates: %v\n", err)
	}

	repr.Print(client.Bootstrap)

	for {
		update := <-receiver.C
		log.Printf("Got %s %s", update.Action.ModelKey, update.Action.Action)
		repr.Print(update.Payload)
	}
}
