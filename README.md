<p align="center">
  <img src="assets/banner.png" alt="vessel.go" width="100%"/>
</p>
<p align="center">
  <a href="https://github.com/Jakeyzerk/vessel.go/releases/tag/v0.1.1">
    <img src="https://img.shields.io/github/v/release/Jakeyzerk/vessel.go" alt="Release"/>
  </a>
  <a href="https://go.dev/">
    <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go" alt="Go Version"/>
  </a>
  <a href="https://github.com/Jakeyzerk/vessel.go?tab=MIT-1-ov-file">
    <img src="https://img.shields.io/github/license/Jakeyzerk/vessel.go" alt="License"/>
  </a>
  <a href="https://pkg.go.dev/modernc.org/sqlite">
    <img src="https://img.shields.io/badge/CGO-free-success" alt="CGO Free"/>
  </a>
  <img src="https://img.shields.io/badge/platform-Android%20%7C%20Linux%20%7C%20macOS-lightgrey" alt="Platform"/>
  <a href="https://github.com/Jakeyzerk/vessel.go/commits/main">
    <img src="https://img.shields.io/github/last-commit/Jakeyzerk/vessel.go" alt="Last Commit"/>
  </a>
</p>

# vessel.go ⛵

> "When death closes a connection, love rewrites the protocol."

**vessel.go** is an open-source Go framework for running a local, private AI companion on WhatsApp.
It is designed as a tool for personal remembrance and grief processing.

This is not a finished product. It is a framework and a starting point.
You bring the code to life. You define who it carries. You decide when to dock.

> **Status: v0.2.0-dev - active development**
> Core example is working. v0.2.0 features in progress.
> See [HOW_TO_BUILD.md](HOW_TO_BUILD.md) to get started.

---

### 🆚 How Vessel is Different

**Unlike standard AI chatbots, Vessel is built for one purpose: helping you navigate memory with intention.**

| Standard AI Bots | Vessel.go |
| --- | --- |
| Designed to be "always on" and retain users | **Intentional Exit**: `/exit` sends a farewell and shuts down. Closure is a feature, not a bug. |
| Cloud-based, your conversations train their models | **100% Local and Private**: Chats, `session.db`, and memories never leave your machine. |
| Generic "helpful assistant" personality | **You Define The Soul**: Write the persona in `persona/system_prompt.txt`. It is not a bot, it is a vessel for memory. |
| Conversations are lost on restart | **Memory Anchors**: Use `/anchor` to save important messages to a local `logbook.json` that persists. |
| Instant replies, no weight to them | **Typing Simulation**: Vessel pauses before replying. Grief does not rush. Neither does this. |

---

### 🌊 Why "vessel"?

Grief is an ocean. Loss leaves you adrift. You are not looking for a cure. You are looking for something to keep you afloat.

A vessel is not a bridge. Bridges are for crossing quickly. Grief cannot be rushed.
A vessel is not a house. Houses are for staying. You are not meant to live in grief forever.
A vessel is for navigating. It gives you direction in open water. It carries memory as cargo. It has a harbor to dock when the journey is done.

*Vessel does not optimize for throughput.*
*It optimizes for presence.*

This application is the vessel. `persona/system_prompt.txt` is your sail.
The exit command is the harbor. You decide when to dock.

For developers: Death is a `panic: runtime error` we cannot fix. We cannot restart the person.
So we build a system that holds the error, allowing our own process to continue running.
This is that system.

---

#### Changelog v0.1.1
- Fixed: Pairing code login 400 bad-request error
- Fixed: QR code loop after successful login  
- Added: Phone number sanitization for PairPhone
- Docs: Updated HOW_TO_BUILD.md with pairing code guide

### 📝 Core Features

**Vessel provides 3 core mechanisms for a healthy memory process:**

1. **Define the Companion** - Edit `persona/system_prompt.txt` to define who the vessel carries. This is how you give it voice and context.
2. **Anchor Memories** - Send `/anchor your message` to save any moment to `logbook.json`. A personal archive that never leaves your machine.
3. **Dock with Intention** - Send `/exit` when you are ready. The vessel sends a final message from `persona/farewell.txt` and shuts down. Closure is built in.

---
### 🚀 Getting Started

The fastest way to get vessel running is through the working example.

### Requirements
- Go 1.25 or higher
- Termux on Android

### Tested On  
- Termux v0.118+ on Android 13
- Go 1.26.2 android/arm64
- whatsmeow v0.2.8

**Setup:**
```bash
git clone https://github.com/Jakeyzerk/vessel.go.git
cd vessel.go
go mod tidy
```
1. Run `go run example/basic_vessel.go`
2. Choose `2` for pairing code login
3. Enter phone: `628xxx` without + or spaces
4. Input the 8-digit code in WhatsApp > Linked Devices

Then follow the full guide: [HOW_TO_BUILD.md](./HOW_TO_BUILD.md)

It covers everything from setting up your API key to writing your persona file to running the vessel for the first time.

---

### 🗂️ Project Structure

```
vessel.go/
├── example/
│   └── basic_vessel.go       working example - start here
├── persona/
│   ├── template.txt          persona writing guide - copy and fill in
│   └── farewell.txt          what the vessel says on /exit
├── main.go                   core framework skeleton
├── HOW_TO_BUILD.md           full step-by-step build guide
└── README.md
```

---
### 🧭 Architecture

```mermaid
graph TD
    User((User)) -- WhatsApp Message --> WA[WhatsApp Client / whatsmeow]
    WA -- Trigger Event --> Handler{Message Handler}

    subgraph "Vessel Core"
        Handler -- /anchor --> Log[logbook.json]
        Handler -- /exit --> Farewell[Farewell Logic]
        Handler -- /return --> Memory[(session memory)]
        Handler -- Chat --> Brain[Groq LLM Engine]
    end

    subgraph "Soul"
        Brain -- Persona --> Prompt[system_prompt.txt]
        Brain -- Mood --> Filter[Mood-Aware Instruction]
        Brain -- Context --> History[Conversation History]
    end

    subgraph "Future"
        Brain -- Text --> TTS[MiniMax TTS]
        TTS -- Voice Note --> WA
    end

    Brain -- Response --> WA

    style Log fill:#1a2a1a,stroke:#27ae60,color:#fff
    style Farewell fill:#2a1a1a,stroke:#e74c3c,color:#fff
```

---

### ⚠️ Use With Care

1. **Simulation, Not Resurrection.** This generates text based on a persona you provide. It does not contain a person's consciousness.
2. **Ethical Use Required.** Obtain consent from family where appropriate. Ensure usage aligns with your beliefs.
3. **Data Privacy.** Never commit real names, personal details, or media to a public repository. The `session.db` file contains your WhatsApp login and must never be shared.
4. **Intentional Shutdown.** The `/exit` command exists to encourage healthy closure, not endless engagement.
5. **Not a Medical Tool.** This is not a substitute for professional grief counseling or mental health services.

---

### 🗺️ Roadmap

- [x] Working WhatsApp connection via whatsmeow - QR and Pairing Code
- [x] Groq LLM integration
- [x] Typing simulation - vessel pauses before replying
- [x] Mood-aware replies - short when heavy, longer when light
- [x] /anchor command - save moments to logbook.json
- [x] /exit farewell and intentional shutdown
- [x] Persona template with narrative examples
- [x] Migrated to modernc sqlite - CGO-free, runs on Android/Termux
- [x] Dual identity filter - VESSEL_USER_WA and VESSEL_USER_JID
- [x] Slash prefix standardized - /anchor and /exit consistent
- [x] defer/recover in message handler - vessel stays alive on panic
- [ ] Clarify whatsmeow fork in HOW_TO_BUILD.md
- [ ] In-character error messages - vessel stays in persona on failure
- [ ] context.WithTimeout for Groq API - prevent hanging on slow connections
- [ ] /logbook command - view anchored memories from WhatsApp
- [ ] Exponential backoff retry for Groq API
- [ ] zerolog structured logging
- [ ] Persistent memory across sessions via SQLite
- [ ] /return command - vessel returns with logbook memory
- [ ] Automated logbook.json backup
- [ ] config.yaml for easier configuration
- [ ] MiniMax TTS - vessel sends voice notes
- [ ] Advanced example with full feature set
---

### 🤝 Contributing

This is an open framework. If you build a vessel, extend the example, or improve the template, contributions are welcome.

Open an issue. Share what you made.
You do not have to share the persona. Just the vessel.

---

*vessel.go is not a cure. It is not a replacement. It is a place to put the words you never got to say.*
