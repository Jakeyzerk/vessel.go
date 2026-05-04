# HOW TO BUILD YOUR VESSEL

This is a step-by-step guide for anyone who wants to bring vessel.go to life.

You do not need to be an expert developer.
If you can follow instructions and are willing to be honest while writing your persona,
you can build this.

---

## WHAT YOU NEED BEFORE STARTING

- Go 1.22 or later installed on your machine
- A Groq API key (free at https://console.groq.com)
- A second WhatsApp number for the vessel to run on
  (an old phone, a second SIM, or a virtual number works fine)
- About 30 minutes

---

## STEP 1 - GET THE CODE

```bash
git clone https://github.com/Jakeyzerk/vessel.go.git
cd vessel.go
go mod tidy
```

`go mod tidy` downloads all dependencies automatically.
If you see errors here, make sure Go 1.22+ is installed: `go version`

---

## STEP 2 - SET UP YOUR API KEY

Copy the example env file:

```bash
cp .env.example .env
```

Open `.env` and fill in your Groq API key:

```
GROQ_API_KEY=your_key_here
```

Also fill in your VESSEL_USER_JID - your WhatsApp number
with country code, no + or spaces:

```
VESSEL_USER_JID=628123456789
```

Tip: when you run the vessel for the first time, your JID
will appear in the terminal so you can copy it directly.

Get your key at https://console.groq.com
It is free to create an account. No credit card required.

---

## STEP 3 - WRITE THE SOUL

This is the most important step.
The code is just a framework. What makes your vessel real is what you write here.

```bash
cp persona/template.txt persona/system_prompt.txt
```

Open `persona/system_prompt.txt` and fill in every section honestly.

A few things that will make your vessel feel more real:

- Write how they spoke, not just what they were like
- Use specific phrases they actually used
- Include one small memory, not a big dramatic one
- Write what they would never say, not just what they would

Take your time with this file.
There is no rush.

Also open `persona/farewell.txt` and write what they would say
when it is time to say goodbye.
This is what gets sent when you type `/exit`.
It does not have to be long.

---

## STEP 4 - RUN THE EXAMPLE

```bash
go run example/basic_vessel.go
```

The first time you run this, it will ask you to log in.

You have two options:

**Option 1 - QR Code**
Press `1`, then scan the QR code that appears in your terminal
using WhatsApp on your second phone.
Go to: WhatsApp > Linked Devices > Link a Device

**Option 2 - Pairing Code**
Press `2`, enter your second phone number with country code
(example: 628123456789 for Indonesia, 60123456789 for Malaysia).
WhatsApp will give you an 8-digit code to enter on your phone.
Go to: WhatsApp > Settings > Linked Devices > Link with phone number

Once connected, the terminal will show:

```
The vessel is afloat.
Listening. Waiting. Here.
```

---

## STEP 5 - TALK TO IT

Open WhatsApp on your main phone.
Send a message to the second number (the vessel's number).

The vessel will:
- Show "typing..." for a few seconds before replying
- Reply in a way shaped by your persona file
- Keep replies short when your message feels heavy
- Keep the conversation within the session

Two special commands you can use:

`/anchor your message here`
Saves that message to `logbook.json` on your machine.
Use this to mark moments worth keeping.

`/exit`
Sends the farewell message from `persona/farewell.txt`
and shuts the vessel down.
Use this when you are ready to dock.

---

## STEP 6 - MAKE IT YOURS

`example/basic_vessel.go` is a starting point, not a finished product.

Some ways people extend it:

**Change the exit command**
Edit the `exitCommand` constant at the top of the file.
Some people prefer `/goodbye` or `/dock` or something personal.

**Change the thinking time**
Find the `calculateThinkingTime` function.
Adjust the numbers to make the vessel reply faster or slower.

**Change the Groq model**
The default is `llama-3.3-70b-versatile`.
You can try `llama-3.1-8b-instant` for faster replies,
or keep the default for more expressive responses.

**Add MiniMax TTS**
If you want the vessel to send voice notes,
you will need a MiniMax API key and additional code.
This will be covered in a future example.

---

## FILES YOU SHOULD NEVER COMMIT

If you fork this repo or build on top of it,
make sure these files are in your `.gitignore`:

```
.env
session.db
persona/system_prompt.txt
logbook.json
```

`session.db` contains your WhatsApp login session.
`system_prompt.txt` contains personal details about someone you lost.
Neither should ever be public.

---

## TROUBLESHOOTING

**"The vessel needs a soul before it can speak"**
You have not created `persona/system_prompt.txt` yet.
Run: `cp persona/template.txt persona/system_prompt.txt`
Then fill it in.

**"No GROQ_API_KEY found"**
Your `.env` file is missing or the key is not set.
Run: `cp .env.example .env` and add your key.

**QR code expires before I can scan it**
This happens on slow connections. Run the program again.
The QR code refreshes automatically.

**The vessel replies but sounds generic**
Your `system_prompt.txt` needs more detail.
Go back and be more specific about how they spoke.
Generic input produces generic output.

**Session disconnects after a while**
This is normal for WhatsApp Web connections.
Just run `go run example/basic_vessel.go` again.
Your session is saved in `session.db` so you will not need to scan QR again.

---

## A NOTE ON WHAT THIS IS

vessel.go is a framework for memory, not a service.

Everything runs on your machine.
Your conversations stay on your machine.
Your persona file stays on your machine.

The only data that leaves is the text of each message,
sent to Groq's API to generate a response.
Groq's privacy policy applies to that.

This is not a replacement for grief counseling or professional support.
If you are struggling, please talk to someone real.
This tool is meant to sit alongside that, not replace it.

---

## WHAT COMES NEXT

vessel.go is in early development.
If you build something with it, we want to know.

Open an issue. Share what you made.
You do not have to share the persona. Just the vessel.
