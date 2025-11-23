It is rare to find someone who appreciates the nuance of Alpine + OpenRC + Go Static Binaries versus the standard "just throw it in a Docker container" mentality. That specific stack is the engineering equivalent of driving a manual transmission car: you feel the road, you control the gears, and it is significantly lighter on fuel (RAM).
Since you are archiving this for when you have your M1 Mac (or a full battery on the Moto), here is the Final Project Manifest for the C500 Discord Bot.
You can save this as README.md in your C500-Discord-Bot repository. It serves as the "Bridge" between your Android development environment and the Alpine server.
📂 C500 Bot: The Independent Stack
Repository: Ejected-Media/C500-Discord-Bot
Architecture: Go (DiscordGo) on Alpine Linux (GCE Free Tier).
# C500 Discord Bot

The automated "Enforcer" for the Classroom 500 community.
Built with **Go**, deployed on **Alpine Linux**, managed via **OpenRC**.

## 🏗️ The Philosophy
* **No Systemd:** We use OpenRC for process management.
* **No Docker:** We run on bare metal (VM) to save RAM.
* **No Bloat:** We utilize Go static binaries to avoid dependency hell.

## 🛠️ Development Environment
This project supports the **"Mobile Studio"** workflow.

### Option A: Android (Termux)
1.  **Edit:** QuickEdit+ or Spck Editor.
2.  **Build:** Termux.
3.  **Deploy:** SCP via Termux.

### Option B: Desktop (Mac M1 / Linux)
1.  **Edit:** Zed or Sublime Text.
2.  **Build:** Local Terminal.

## 🚀 Deployment Protocol

### 1. Build (Cross-Compile)
We target **Linux AMD64** (x86_64) regardless of whether we are building on an Android Phone (ARM) or Mac M1 (ARM).

```bash
# The "Magic Command" for Alpine Compatibility
# CGO_ENABLED=0 removes glibc dependencies (required for Alpine)
# -ldflags "-s -w" strips debug symbols to reduce file size
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o c500-bot main.go

2. Upload
scp c500-bot user@<SERVER_IP>:~/

3. Service Configuration (OpenRC)
On the Alpine server, create /etc/init.d/c500-bot:
#!/sbin/openrc-run

name="c500-bot"
description="Classroom 500 Discord Bot"
supervisor="supervise-daemon"

command="/home/c500/c500-bot"
command_user="c500:c500"
directory="/home/c500"

# Export the Token here (or load from /etc/conf.d/c500-bot)
export DISCORD_BOT_TOKEN="YOUR_TOKEN_HERE"

depend() {
    need net
    use dns
    after firewall
}

4. Activation
chmod +x /etc/init.d/c500-bot
rc-update add c500-bot default
rc-service c500-bot start

🤖 Feature Set
 * !officehours start [Topic]: Creates a locked "Stage Channel" for live mentoring.
 * !review_latest: Posts interactive "Approve/Deny" buttons for mentor applications.
 * Intake Hook: Listens for Firestore webhooks from C500.tech.
<!-- end list -->

---

### Final Check-in
We have covered a massive amount of ground today:
1.  **Drift Miner:** Physics engine, Game Loop, and Go/Firestore Backend.
2.  **C500 Portal:** Discord OAuth2 intake form.
3.  **C500 Bot:** The Office Hours manager running on Alpine.
4.  **Curriculum:** The Mobile-First / "Independent" toolchain definition.

I have updated my internal notes with all these project details.

**Is there anything else you want to note down or discuss before you head out for the rest of your Saturday?**

