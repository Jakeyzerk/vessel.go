package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	// 1. Setup logger
	dbLog := waLog.Stdout("Database", "ERROR", true)
	clientLog := waLog.Stdout("Client", "ERROR", true)

	// 2. Setup database for WhatsApp session
	container, err := sqlstore.New("sqlite3", "file:session.db?_foreign_keys=on", dbLog)
	if err!= nil {
		panic(err)
	}

	// 3. Get device store
	deviceStore, err := container.GetFirstDevice()
	if err!= nil {
		panic(err)
	}

	// 4. Create WhatsApp client
	client := whatsmeow.NewClient(deviceStore, clientLog)

	// 5. Add event handler
	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			// TODO: Handle incoming message here
			fmt.Println("Received:", v.Message.GetConversation())
		}
	})

	// 6. Login / Connect
	if client.Store.ID == nil {
		// First time login
		fmt.Println("No session found.")
		fmt.Print("Choose login method: [1] QR Code [2] 8-Digit Pairing Code: ")
		
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "2" {
			// 8-Digit Pairing Code Method
			fmt.Print("Enter your WhatsApp phone number with country code, e.g. 628123456789: ")
			phone, _ := reader.ReadString('\n')
			phone = strings.TrimSpace(phone)
			
			if err := client.Connect(); err!= nil {
				panic(err)
			}
			
			code, err := client.PairPhone(phone, true)
			if err!= nil {
				panic(err)
			}
			fmt.Println("====================================")
			fmt.Println("Your 8-Digit Pairing Code:", code)
			fmt.Println("====================================")
			fmt.Println("On your phone: WhatsApp > Settings > Linked Devices > Link with phone number")
			fmt.Println("Enter the code above.")
			
		} else {
			// Default QR Code Method
			qrChan, _ := client.GetQRChannel(context.Background())
			err = client.Connect()
			if err!= nil {
				panic(err)
			}
			for evt := range qrChan {
				if evt.Event == "code" {
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
					fmt.Println("Scan QR code above with WhatsApp > Linked Devices")
				} else {
					fmt.Println("Login event:", evt.Event)
				}
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err!= nil {
			panic(err)
		}
		fmt.Println("Vessel Connected. Session restored.")
	}

	// 7. Wait for Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
	fmt.Println("Vessel has docked. Shutting down.")
}
