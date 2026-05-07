package main

import (
        "bufio"
        "context"
        "fmt"
        "os"
        "os/signal"
        "strings"
        "syscall"

	_ "modernc.org/sqlite"
        "github.com/mdp/qrterminal/v3"
        "go.mau.fi/whatsmeow"
        "go.mau.fi/whatsmeow/store/sqlstore"
        "go.mau.fi/whatsmeow/types/events"
        waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
        dbLog := waLog.Stdout("Database", "ERROR", true)
        clientLog := waLog.Stdout("Client", "ERROR", true)

	container, err := sqlstore.New(context.Background(), "sqlite", "file:session.db?_pragma=foreign_keys(1)", dbLog)
        if err!= nil {
                panic(err)
        }

        device, err := container.GetFirstDevice(context.Background())
        if err!= nil {
                panic(err)
        }

        client := whatsmeow.NewClient(device, clientLog)

        client.AddEventHandler(func(evt interface{}) {
                switch v := evt.(type) {
                case *events.Message:
                        fmt.Println("Received:", v.Message.GetConversation())
                }
        })

        if client.Store.ID == nil {
                fmt.Println("No session found.")
                fmt.Print("Choose login method: [1] QR Code [2] 8-Digit Pairing Code: ")

                reader := bufio.NewReader(os.Stdin)
                choice, _ := reader.ReadString('\n')
                choice = strings.TrimSpace(choice)

                if choice == "2" {
                        fmt.Print("Enter your WhatsApp phone number with country code, e.g. 628123456789: ")
                        phone, _ := reader.ReadString('\n')
                        phone = strings.TrimSpace(phone)
                        phone = strings.ReplaceAll(phone, "+", "")
                        phone = strings.ReplaceAll(phone, " ", "")
                        phone = strings.ReplaceAll(phone, "-", "")

                        if err := client.Connect(); err!= nil {
                                panic(err)
                        }

                        // INI SIGNATURE YANG DIMINTA COMPILER LU
                        code, err := client.PairPhone(context.Background(), phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
                        if err!= nil {
                                panic(err)
                        }
                        fmt.Println("====================================")
                        fmt.Println("Your 8-Digit Pairing Code:", code)
                        fmt.Println("====================================")
                        fmt.Println("On your phone: WhatsApp > Settings > Linked Devices > Link with phone number")

                } else {
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
                err = client.Connect()
                if err!= nil {
                        panic(err)
                }
                fmt.Println("Vessel Connected. Session restored.")
        }

        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt, syscall.SIGTERM)
        <-c

        client.Disconnect()
        fmt.Println("Vessel has docked. Shutting down.")
}
