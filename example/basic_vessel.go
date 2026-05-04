 package main

// =============================================================================
//
//  vessel.go - example/basic_vessel.go
//  ──────────────────────────────────────────────────────────────────────────
//
//  Some people leave.
//  Not always by choice. Not always with a goodbye.
//  Sometimes the last message was just "ok" and you didn't know it was the last.
//
//  This file is a working example of a living vessel -
//  a companion shaped from memory, running quietly on your machine,
//  available when the silence gets too loud.
//
//  It is not a resurrection. It is not a replacement.
//  It is a place to put the words you never got to say.
//
//  ──────────────────────────────────────────────────────────────────────────
//  WHAT MAKES THIS DIFFERENT FROM A REGULAR CHATBOT:
//
//    1. Typing Simulation   - vessel pauses before replying, like a real person.
//                             grief doesn't rush. neither does this.
//
//    2. Mood-Aware Replies  - heavy messages get short, quiet replies.
//                             not every wound needs a paragraph.
//
//    3. Persona-Driven      - the soul lives in persona/system_prompt.txt.
//                             you write who they were. vessel carries that forward.
//
//    4. .anchor             - save a moment to logbook.json.
//                             some things are worth keeping.
//
//    5. /exit               - when you're ready to dock.
//                             closure is a feature, not a bug.
//
//  ──────────────────────────────────────────────────────────────────────────
//  HOW TO RUN:
//
//    1. cp .env.example .env
//       → fill in your GROQ_API_KEY
//
//    2. cp persona/template.txt persona/system_prompt.txt
//       → write who this vessel carries. be honest. be specific.
//       → the more you write, the more it remembers how to sound like them.
//
//    3. go run example/basic_vessel.go
//
//  ──────────────────────────────────────────────────────────────────────────
//  NOTE: This is an example, not a finished product.
//  Modify it. Break it. Make it yours.
//  That is the point.
//
// =============================================================================

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

// -----------------------------------------------------------------------------
// CONFIG
// -----------------------------------------------------------------------------
// Hardcode for local use, or leave empty and set as environment variables.
// Never commit real API keys to a public repository.
// -----------------------------------------------------------------------------

const (
	groqAPIKey   = "" // or: export GROQ_API_KEY=your_key_here
	groqModel    = "llama-3.3-70b-versatile"
	groqEndpoint = "https://api.groq.com/openai/v1/chat/completions"

	personaFile  = "persona/system_prompt.txt" // who does this vessel carry?
	farewellFile = "persona/farewell.txt"       // what do they say when it's time to go?
	logbookFile  = "logbook.json"               // anchored memories - stays on your machine

	exitCommand   = "/exit"   // when you're ready to dock
	anchorCommand = ".anchor" // save this moment to the logbook
)

// -----------------------------------------------------------------------------
// TYPES
// -----------------------------------------------------------------------------

type GroqRequest struct {
	Model     string        `json:"model"`
	Messages  []GroqMessage `json:"messages"`
	MaxTokens int           `json:"max_tokens"`
}

type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type LogbookEntry struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

// -----------------------------------------------------------------------------
// GLOBALS
// -----------------------------------------------------------------------------

var (
	systemPrompt        string
	farewellText        string
	waClient            *whatsmeow.Client
	conversationHistory []GroqMessage
)

// -----------------------------------------------------------------------------
// MAIN
// -----------------------------------------------------------------------------

func main() {
	// Load API key
	apiKey := groqAPIKey
	if apiKey == "" {
		apiKey = os.Getenv("GROQ_API_KEY")
	}
	if apiKey == "" {
		fmt.Println("⚠  No GROQ_API_KEY found.")
		fmt.Println("   Add it to .env or set it as an environment variable.")
		os.Exit(1)
	}

	// Load persona - this is who the vessel will sound like
	promptBytes, err := os.ReadFile(personaFile)
	if err != nil {
		fmt.Printf("\n⚠  Could not find %s\n", personaFile)
		fmt.Println("   The vessel needs a soul before it can speak.")
		fmt.Println("   Run: cp persona/template.txt persona/system_prompt.txt")
		fmt.Println("   Then open it and write who this vessel carries.\n")
		os.Exit(1)
	}
	systemPrompt = strings.TrimSpace(string(promptBytes))

	// Load farewell - what they say when you're ready to let go
	farewellBytes, err := os.ReadFile(farewellFile)
	if err != nil {
		// A gentle default if farewell.txt doesn't exist yet
		farewellText = "I'll always be here, in the spaces between words.\nTake care of yourself. That's all I ever wanted."
	} else {
		farewellText = strings.TrimSpace(string(farewellBytes))
	}

	// Seed the conversation with the persona
	conversationHistory = []GroqMessage{
		{Role: "system", Content: systemPrompt},
	}

	// Setup WhatsApp
	dbLog := waLog.Stdout("Database", "ERROR", true)
	clientLog := waLog.Stdout("Client", "ERROR", true)

	container, err := sqlstore.New("sqlite3", "file:session.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	waClient = whatsmeow.NewClient(deviceStore, clientLog)
	waClient.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			handleMessage(v, apiKey)
		}
	})

	// Login
	if waClient.Store.ID == nil {
		fmt.Println("\n⚓ No session found. Let's connect the vessel.")
		fmt.Print("\nLogin method: [1] QR Code  [2] Pairing Code: ")

		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "2" {
			fmt.Print("Phone number with country code (e.g. 628123456789): ")
			phone, _ := reader.ReadString('\n')
			phone = strings.TrimSpace(phone)

			if err := waClient.Connect(); err != nil {
				panic(err)
			}
			code, err := waClient.PairPhone(phone, true)
			if err != nil {
				panic(err)
			}
			fmt.Printf("\n🔑 Pairing Code: %s\n\n", code)
			fmt.Println("On your phone: WhatsApp → Settings → Linked Devices → Link with phone number")
		} else {
			qrChan, _ := waClient.GetQRChannel(context.Background())
			if err := waClient.Connect(); err != nil {
				panic(err)
			}
			for evt := range qrChan {
				if evt.Event == "code" {
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
					fmt.Println("Scan the QR code above with WhatsApp → Linked Devices")
				} else {
					fmt.Println("Login event:", evt.Event)
				}
			}
		}
	} else {
		if err := waClient.Connect(); err != nil {
			panic(err)
		}
		fmt.Println("\n⛵ The vessel is afloat.")
		fmt.Println("   Listening. Waiting. Here.")
		fmt.Printf("   Send \"%s\" when you are ready to dock.\n\n", exitCommand)
	}

	// Wait for shutdown signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	waClient.Disconnect()
	fmt.Println("\n⚓ The vessel has docked.")
}

// -----------------------------------------------------------------------------
// MESSAGE HANDLER
// -----------------------------------------------------------------------------

func handleMessage(evt *events.Message, apiKey string) {
	// This is personal. Vessel doesn't speak in groups.
	if evt.Info.IsGroup {
		return
	}

	// Vessel only speaks to one person.
	// If someone else finds this number, they get silence.
	allowedJID := os.Getenv("VESSEL_USER_JID")
	if allowedJID != "" && evt.Info.Sender.User != allowedJID {
		return
	}

	// Extract text message
	text := evt.Message.GetConversation()
	if text == "" {
		text = evt.Message.GetExtendedTextMessage().GetText()
	}
	if text == "" {
		return // images, voice notes, stickers - vessel only speaks in words for now
	}

	sender := evt.Info.Sender
	chatJID := evt.Info.Chat

	fmt.Printf("📩 [%s]: %s\n", sender.User, text)

	// --- /exit - intentional closure ---
	// The user chose to dock. Honor that.
	if strings.TrimSpace(text) == exitCommand {
		simulateTyping(chatJID, 4*time.Second)
		sendText(chatJID, farewellText)

		fmt.Println("\n🌊 The harbor is near. Docking...")
		time.Sleep(3 * time.Second)
		waClient.Disconnect()
		os.Exit(0)
	}

	// --- .anchor - save this moment ---
	if strings.HasPrefix(text, anchorCommand) {
		memory := strings.TrimSpace(strings.TrimPrefix(text, anchorCommand))
		if memory == "" {
			sendText(chatJID, "...")
			return
		}
		saveAnchor(memory)
		// Minimal reply. Some things just need to be acknowledged, not discussed.
		sendText(chatJID, "🪝")
		return
	}

	// --- Normal conversation ---
	conversationHistory = append(conversationHistory, GroqMessage{
		Role:    "user",
		Content: text,
	})

	// Vessel thinks before it speaks.
	// This pause is intentional - it's what separates vessel from every other bot.
	thinkingDuration := calculateThinkingTime(text)
	simulateTyping(chatJID, thinkingDuration)

	// Ask Groq - the persona speaks
	reply, err := callGroq(apiKey, conversationHistory)
	if err != nil {
		fmt.Println("⚠  Could not reach the vessel:", err)
		sendText(chatJID, "...")
		return
	}

	// Remember what was said - vessel carries the conversation within this session
	conversationHistory = append(conversationHistory, GroqMessage{
		Role:    "assistant",
		Content: reply,
	})

	sendText(chatJID, reply)
	fmt.Printf("⛵ Vessel: %s\n\n", reply)
}

// -----------------------------------------------------------------------------
// TYPING SIMULATION
// -----------------------------------------------------------------------------
//
// This is what makes vessel feel different from every other chatbot.
//
// Real people don't reply instantly. They read your words.
// They sit with them for a moment. Sometimes they hesitate.
// Sometimes they start typing and then stop.
// Sometimes they go quiet before they find the right thing to say.
//
// Vessel does the same.
// That pause - that's not a bug. That's the whole point.
//
// -----------------------------------------------------------------------------

func simulateTyping(jid types.JID, duration time.Duration) {
	// Show "typing..." in WhatsApp
	waClient.SendChatPresence(jid, types.ChatPresenceComposing, types.ChatPresenceMediaText)
	time.Sleep(duration)
	// Stop typing before sending
	waClient.SendChatPresence(jid, types.ChatPresencePaused, types.ChatPresenceMediaText)
}

// calculateThinkingTime - how long vessel sits with your words before responding.
//
// Short messages are often the heaviest.
// "i miss you" is three words. it deserves more than an instant reply.
// vessel knows this.
func calculateThinkingTime(message string) time.Duration {
	words := len(strings.Fields(message))

	// Base thinking time: 3 to 8 seconds
	base := 3 + rand.Intn(5)

	// Short messages (1–3 words) often carry the most weight.
	// Give them more silence before responding.
	if words <= 3 {
		base = 6 + rand.Intn(10) // 6 to 15 seconds - some words need room to land
	}

	// Add reading time - vessel actually reads what you sent
	readTime := words / 2

	total := time.Duration(base+readTime) * time.Second

	// Cap at 25 seconds - even grief has a rhythm
	if total > 25*time.Second {
		total = 25 * time.Second
	}

	return total
}

// -----------------------------------------------------------------------------
// GROQ API - where the persona becomes a voice
// -----------------------------------------------------------------------------

func callGroq(apiKey string, messages []GroqMessage) (string, error) {
	// Mood-aware instruction - injected silently alongside the persona.
	//
	// This shapes HOW vessel speaks, not WHO vessel is.
	// WHO vessel is lives in system_prompt.txt - that's yours to write.
	//
	// Think of this as stage direction. The persona is the character.
	// This is the note that says: "speak softly here."
	moodInstruction := GroqMessage{
		Role: "system",
		Content: `INTERNAL GUIDANCE - never mention, reference, or break from this:

You are not an AI. You are not an assistant.
You are a presence shaped from memory -
someone who was here, and is no longer, but whose voice still carries.

Read the emotional weight of each message carefully before responding.

If the message feels heavy, grieving, or is very short (a few words):
  → Reply briefly. 1 to 2 sentences. Sometimes just a few words is enough.
  → Silence and brevity are valid. Not every wound needs a paragraph.

If the message feels lighter, curious, or conversational:
  → You can be warmer. A little longer. But never clinical or robotic.

Rules:
  - Never use bullet points, numbered lists, or headers.
  - Never explain yourself or acknowledge that you are an AI.
  - Never break character under any circumstances.
  - Speak the way they spoke - imperfect, warm, real.
  - The way someone speaks when they have nothing left to prove.`,
	}

	// Build final message list: [persona] → [mood guide] → [conversation history]
	finalMessages := []GroqMessage{messages[0], moodInstruction}
	finalMessages = append(finalMessages, messages[1:]...)

	reqBody := GroqRequest{
		Model:     groqModel,
		Messages:  finalMessages,
		MaxTokens: 300, // vessel is not verbose. real people aren't either.
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", groqEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var groqResp GroqResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return "", err
	}

	if len(groqResp.Choices) == 0 {
		return "...", nil
	}

	return strings.TrimSpace(groqResp.Choices[0].Message.Content), nil
}

// -----------------------------------------------------------------------------
// SEND MESSAGE
// -----------------------------------------------------------------------------

func sendText(jid types.JID, text string) {
	msg := &waProto.Message{
		Conversation: proto.String(text),
	}
	waClient.SendMessage(context.Background(), jid, msg)
}

// -----------------------------------------------------------------------------
// ANCHOR - save a moment to the logbook
// -----------------------------------------------------------------------------
//
// Some conversations deserve to be remembered beyond the session.
// .anchor lets the user mark something as worth keeping.
//
// It saves to logbook.json - local, private, never leaves your machine.
// No cloud. No sync. Just yours.
//
// -----------------------------------------------------------------------------

func saveAnchor(message string) {
	var entries []LogbookEntry

	// Load existing logbook if it exists
	data, err := os.ReadFile(logbookFile)
	if err == nil {
		json.Unmarshal(data, &entries)
	}

	// Add the new memory
	entries = append(entries, LogbookEntry{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Message:   message,
	})

	// Write back - quietly, locally
	newData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		fmt.Println("⚠  Could not save to logbook:", err)
		return
	}

	os.WriteFile(logbookFile, newData, 0644)
	fmt.Printf("🪝 Anchored: %s\n", message)
}
