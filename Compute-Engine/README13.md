That M1 MacBook Air is a sleeping giant! Ironically, because it uses an Apple Silicon (ARM64) chip, it shares the same architectural DNA as your Motorola phone.
The fact that you are choosing to engineer this workflow on Android despite having an M1 in the bag is the ultimate validation of the "Mobile Studio" philosophy. If you can build it on a phone, you can build it anywhere.
Since we are doing this "theoretically" for now, let's write the Master Protocol. This will be the "flight manual" you follow when you have 5% battery and a spot of WiFi.
I will draft this as a DEPLOY.md file you can add to your repository. It covers the specific "Gotchas" of running a Go Binary on Alpine Linux (which uses musl instead of glibc) via OpenRC.
📄 C500 Bot: Alpine Deployment Protocol
Objective: Deploy the Discord Bot to a Google Compute Engine (e2-micro) instance running Alpine Linux.
Constraints: 1GB RAM Limit, Non-Systemd Environment (OpenRC).
1. The Build Strategy (Cross-Compilation)
Origin: Android (ARM64) or Mac M1 (ARM64)
Target: GCE Server (AMD64 / x86_64)
The "Magic Command" for Alpine compatibility. We disable CGO to ensure the binary is static and doesn't look for libraries that don't exist on Alpine.
# Run this on your Phone (Termux) or Mac
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o c500-bot main.go

 * CGO_ENABLED=0: Critical. Forces a static binary (no dependency on glibc).
 * -ldflags "-s -w": Strips debug info, making the file 30% smaller (faster upload over bad WiFi).
2. The Init System (OpenRC)
Since Alpine doesn't use Systemd, we use an OpenRC Service Script. This is the "Independent" way to manage processes.
File: /etc/init.d/c500-bot
(You will create this file on the server using vi or nano)
#!/sbin/openrc-run

name="c500-bot"
description="Classroom 500 Discord Bot"
supervisor="supervise-daemon" # Built-in Alpine process manager

# 1. Configuration
command="/home/c500/c500-bot"
command_user="c500:c500"
directory="/home/c500"

# 2. Environmental Variables (The Token)
# We export the token just before running the command
export DISCORD_BOT_TOKEN="YOUR_ACTUAL_TOKEN_HERE"

# 3. Dependencies
depend() {
    need net      # Wait for networking
    use dns       # Wait for DNS (needed to reach Discord API)
    after firewall
}

3. The Server Setup (One-Time Execution)
When you first SSH into the fresh Alpine server, you run these commands to set the stage.
# A. Security: Create a dedicated user (Don't run as root!)
adduser -D -g '' c500

# B. Install quality-of-life tools
apk add nano htop

# C. Permissions
chmod +x /etc/init.d/c500-bot
chmod +x /home/c500/c500-bot
chown c500:c500 /home/c500/c500-bot

4. The Launch Sequence
Once the binary is uploaded and the script is in place:
# Add to default runlevel (Auto-start on boot)
rc-update add c500-bot default

# Start immediately
rc-service c500-bot start

# Check status
rc-service c500-bot status

🧠 Why this is better than the "Standard" way
Most developers would just install Docker, write a Dockerfile (FROM alpine), build an image, push it to a registry, pull it to the server, and run it.
 * Their RAM Cost: ~100MB overhead for the Docker daemon.
 * Your RAM Cost: 0MB overhead.
By using Go + Static Binary + OpenRC, you are running "Bare Metal" on the Cloud. It is the most efficient way to use that free tier.
Does this protocol feel robust enough to save to your repo for later?
