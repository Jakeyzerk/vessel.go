# vessel.go ⛵
### Compile Memory Into Companionship

> "When death closes a connection, love rewrites the protocol."

**vessel.go** is an open-source Go framework for running a local, private AI companion on WhatsApp. It is designed as a tool for personal remembrance and grief processing.

This is not a service. It is a self-hosted application. All data is stored locally. No conversation data is shared with third parties except for API calls required for language and voice generation.

---

### 🆚 How Vessel is Different
**Unlike standard AI chatbots, Vessel is built for one purpose: helping you navigate memory with intention.**

| Standard AI Bots | Vessel.go |
| --- | --- |
| Designed to be "always on" and retain users | **Intentional Exit**: `/exit` command sends a farewell and shuts down. Closure is a feature, not a bug. |
| Cloud-based, your conversations train their models | **100% Local & Private**: Chats, `session.db`, and memories never leave your machine. |
| Generic "helpful assistant" personality | **You Define The Soul**: Craft the persona in `persona/system_prompt.txt`. It's not a bot, it's a vessel for memory. |
| Conversations are lost on restart | **Memory Anchors**: Use `.anchor` to save important messages to a local `logbook.json` that persists. |

---

### 🌊 Why "vessel"?

Grief is an ocean. Loss leaves you adrift. You are not looking for a cure. You are looking for something to keep you afloat.

A vessel is not a bridge. Bridges are for crossing quickly. Grief cannot be rushed.
A vessel is not a house. Houses are for staying. You are not meant to live in grief forever.
A vessel is for navigating. It gives you direction in open water. It carries memory as cargo. It has a harbor to dock when the journey is done.

This application is the vessel. `config.yaml` is your map. `persona/system_prompt.txt` is your sail.
The exit command is the harbor. You decide when to dock.

For developers: Death is a `panic: runtime error` we cannot fix. We cannot restart the person. So we build a system that holds the error, allowing our own process to continue running. This is that system.

---

### 📝 Core Features
**Vessel provides 3 core mechanisms for a healthy memory process:**

1. **Define the Companion**: Edit `persona/system_prompt.txt` to define who the vessel embodies. This is how you give it voice and context.
2. **Anchor Memories**: Send `.anchor your message` to save any chat to `logbook.json`. Build a personal archive of moments that matter.
3. **Dock with Intention**: Send `/exit` when you are ready. The vessel will send a final message from `persona/farewell.txt` and shut down permanently. This encourages completion over endless engagement.

---

### ⚠️ Important: Use With Care

This software deals with sensitive subject matter. Please review before use:

1.  **Simulation, Not Resurrection.** This app generates text based on a persona you provide. It does not contain a person's consciousness.
2.  **Ethical Use Required.** Obtain consent from family where appropriate. Ensure usage aligns with your beliefs.
3.  **Data Privacy.** **Never commit real names, personal details, or media to a public repository.** Use placeholders in all configs. The `session.db` file contains your WhatsApp login and must never be shared.
4.  **Intentional Shutdown.** The bot includes a configurable exit command, default `/exit`. Sending this command triggers a farewell and permanently terminates the app. This feature exists to encourage healthy closure. Change it in `config.yaml`.
5.  **Not a Medical Tool.** This is not a substitute for professional grief counseling or mental health services.

---

### 🛠️ Technical Architecture

| Component | Technology | Description |
| --- | --- | --- |
| **Language** | `Go 1.22+` | Core app logic with concurrency via Goroutines |
| **Messaging** | `whatsmeow` | WhatsApp Web API client library |
| **LLM** | `Groq API` | Language model for generating responses |
| **TTS** | `MiniMax API` | Optional text-to-speech for voice notes |
| **Persona** | `persona/system_prompt.txt` | User-defined system prompt for the AI's personality |
| **Memory** | `logbook.json` | Local file for anchored messages [+] BARU |
| **Session** | `SQLite` | Local database for storing WhatsApp session |

### 🚀 Installation & Setup

#### 1. Prerequisites
- Go 1.22 or later
- API keys from Groq and MiniMax

#### 2. Clone & Install
```bash
git clone https://github.com/Jakeyzerk/vessel.go.git
cd vessel.go
go mod tidy


#### 3. Setup Persona
cp persona/template.txt persona/system_prompt.txt
# Edit persona/system_prompt.txt to define your companion's name, traits, and backstory.
# IMPORTANT: Never commit your personal system_prompt.txt to a public repo.
# You can also edit persona/farewell.txt to customize the /exit message.
