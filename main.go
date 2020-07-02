package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
)

type config struct {
	AuthKey string `env:"AuthKey,required"`
	KeyID   string `env:"KeyID,required"`
	TeamID  string `env:"TEAM_ID,required"`
}

func getClient(authKeyFileName string, keyID string, teamID string) *apns2.Client {
	authKey, err := token.AuthKeyFromFile(authKeyFileName)
	if err != nil {
		log.Fatal("token error:", err)
	}

	token := &token.Token{
		AuthKey: authKey,
		// KeyID from developer account
		KeyID: keyID,
		// TeamID from developer account
		TeamID: teamID,
	}

	client := apns2.NewTokenClient(token)
	return client
}

func getNotification(deviceToken string, topic string, payload string, expiration int, priority int, pushType apns2.EPushType) *apns2.Notification {
	notification := &apns2.Notification{
		DeviceToken: deviceToken,
		Topic:       topic,
		Payload:     []byte(payload),
		Priority:    priority,
		PushType:    pushType,
	}

	if expiration == 1 {
		notification.Expiration = time.Now()
	}

	return notification
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	client := getClient("./cert/"+cfg.AuthKey, cfg.KeyID, cfg.TeamID)
	client.Host = apns2.HostDevelopment

	notification := getNotification("123", "topic", "payload", 1, 5, "action")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := client.PushWithContext(ctx, notification)
	defer cancel()

	if err != nil {
		log.Println("There was an error", err)
		return
	}

	if res.Sent() {
		log.Println("Sent:", res.ApnsID)
	} else {
		fmt.Printf("Not Sent: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	}
}
