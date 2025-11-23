This changes the logistics significantly, but it actually fits the "Independent Mobile Studio" theme perfectly.
Building a custom FreeBSD disk image (multi-gigabyte file) on an Android phone and uploading it to the cloud is technically possible but painful. You would be fighting ARM (phone) vs. x86 (server) architecture mismatches and storage limits.
Since you are on a Motorola, let's pivot to the most efficient Android-Native Workflow to get your C500 Bot running on Google Compute Engine today.
We will use Termux as your local powerhouse and Alpine Linux on the server (since it's available as a community image and doesn't require the custom build process of FreeBSD).
📱 The Android "Command Center" Setup
We are going to turn your Motorola into a Go development machine.
1. Setup Termux (The Terminal)
If you haven't already, install Termux (preferably from F-Droid, as the Play Store version is outdated).
Open it and set up your Go environment:
pkg update && pkg upgrade
pkg install golang git openssh

2. The Code Workflow (Files by Google + Termux)
You mentioned using Files by Google as the middle layer. Here is the loop:
 * Write Code: Use QuickEdit+ to write your main.go. Save it to a folder, e.g., Internal Storage/C500-Bot.
 * Access in Termux: Termux needs permission to see that folder.
   termux-setup-storage
cd ~/storage/shared/C500-Bot

 * Compile on Phone:
   Since your phone is ARM but the GCE server is x86 (AMD64), we use Go's cross-compilation magic:
   # Build a linux/amd64 binary named 'c500-bot'
GOOS=linux GOARCH=amd64 go build -o c500-bot main.go

   You now have a server-ready binary sitting on your phone.
☁️ The Server Setup (From the Phone)
Since we can't easily build a custom FreeBSD image on the phone, we will use Alpine Linux (The "Gentoo Lite" option). It honors your wish for No Systemd (it uses OpenRC) and runs beautifully on the Free Tier.
1. Provision the Server (via Browser)
 * Open Firefox Android.
 * Go to console.cloud.google.com.
 * Go to Compute Engine -> VM Instances -> Create Instance.
 * Name: c500-bot
 * Region: us-west1 (or whichever has free tier available).
 * Machine Type: e2-micro.
 * Boot Disk:
   * Click "Change".
   * Go to Public Images.
   * Note: If Alpine isn't in the default list, search for "alpinelinux" projects or just pick Debian 12 for today to get started.
   * (If you pick Debian, we can strip it down later. If you can select a "Container Optimized OS" or similar, stick to Debian for a generic Linux env).
2. Connecting via Termux
Once the VM is running, get its External IP.
In Termux:
 * Generate an SSH key (if you haven't): ssh-keygen -t ed25519
 * Add your public key to the VM metadata (You can do this in the Google Cloud Console in Firefox).
 * SSH in: ssh user@<EXTERNAL_IP>
3. Deploying the Bot
Now, upload the binary you built on your phone to the server.
In Termux (from your C500-Bot folder):
scp -i ~/.ssh/id_ed25519 c500-bot user@<EXTERNAL_IP>:~/

🚀 Keeping it Alive (The OpenRC Way)
If you managed to get an Alpine image running (or if you are on Debian and installed OpenRC), here is how you launch it.
If you are stuck on Debian (Systemd) for now because of the Android constraint, here is the "Independent" way to run it without writing a unit file, using a tool called screen (which survives disconnects).
In the SSH Session:
# 1. Install screen
sudo apt-get install screen

# 2. Start a detached session
screen -S bot

# 3. Run the bot
./c500-bot -t "YOUR_DISCORD_TOKEN"

# 4. Detach (Press CTRL+A, then D)

Now you can close Termux, and the bot keeps running.
Summary of the Android Plan
 * Code: QuickEdit+
 * Build: Termux (GOOS=linux GOARCH=amd64 go build)
 * Upload: Termux (scp)
 * Run: GCE VM (screen or Init System)
This gets you online right now using only your Motorola. We can revisit the "Custom FreeBSD Image" project when you have access to a machine that can run QEMU efficiently.
Ready to try compiling that "Hello World" bot on Termux?
