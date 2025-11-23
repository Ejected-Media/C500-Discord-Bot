This is a very common transition for App Engine veterans. I moved from the Python 2.5/2.7 era to Go as well, so I understand your mental model completely.
Here is the breakdown of why this switch matters for a Discord Bot, and exactly what OS you should choose.
1. The Conceptual Shift: "Request" vs. "Persistent"
 * App Engine (GAE): It is designed for HTTP Requests. It spins up your code when a request comes in (someone visits your site) and spins it down when traffic stops.
   * Problem for Bots: Discord Bots utilize a WebSocket connection (a 24/7 phone line to Discord's servers). If App Engine "scales down" to save money, your bot disconnects and goes offline.
 * Compute Engine (GCE): It is a Virtual Machine. It’s like leaving a Mac Mini running in a closet. It runs 24/7 until you turn it off.
   * Benefit: Your bot stays online forever.
   * Cost: Google offers a "Free Tier" (e2-micro instance) in specific regions (like us-west1, us-central1) that makes this effectively free.
2. Do you need a Container OS?
No. You typically do not want a "Container-Optimized OS" (COS) for this unless you are already deep into Docker.
 * Container-Optimized OS (COS): This is a stripped-down OS made by Google that only runs Docker containers. It has a read-only file system. It is great for Kubernetes, but annoying if you just want to SSH in and check a log file or tweak a config.
 * Standard Linux: This is what you want. It lets you treat the server like a remote computer.
3. My Recommendation: Ubuntu 24.04 LTS (Minimal)
Since you are writing in Go, you have a superpower: Static Binaries.
Go compiles your entire application into a single file. You don't need to install Python, pip, virtualenvs, or libraries on the server. You just drop that one file onto the OS and run it.
Why Ubuntu?
 * Documentation: If you get stuck, 99% of Linux tutorials assume you are using Ubuntu.
 * Systemd: It uses the standard service manager to keep your bot running (if it crashes, Systemd restarts it—just like App Engine used to do for you).
4. The Deployment "Recipe" (GAE vs. GCE)
Here is how your workflow changes.
The Old Way (App Engine):
 * You write app.yaml.
 * You run gcloud app deploy.
 * Google handles the rest.
The New Way (Compute Engine):
You only have to do the "Setup" once. After that, updating is easy.
Step A: Create the VM
 * Go to Compute Engine -> VM Instances.
 * Create Instance.
 * Region: us-west1 or us-central1 (Critical for Free Tier).
 * Machine Type: e2-micro.
 * Boot Disk: Change to Ubuntu 24.04 LTS Minimal (Standard persistent disk).
 * Firewall: You don't need HTTP/HTTPS allowed for a bot (since it makes outbound calls to Discord), but keeping SSH open is standard.
Step B: The "Go Build" Strategy
Since you are on a Mac/PC (likely x86 or ARM) and the server is Linux (x86), you cross-compile on your machine.
In your local terminal:
# Tell Go to build for Linux
GOOS=linux GOARCH=amd64 go build -o c500-bot main.go

You now have a file named c500-bot that is a Linux executable.
Step C: Uploading (The "Deploy")
You use scp (Secure Copy), which comes built-in with SSH.
# Upload the binary and your .env file
gcloud compute scp ./c500-bot ./env user@your-instance-name:~/

Step D: The "Keep it Running" Magic (Systemd)
On App Engine, app.yaml told Google how to run your app. On GCE, we create a System Service.
SSH into your VM:
gcloud compute ssh your-instance-name

Create a service file:
sudo nano /etc/systemd/system/c500-bot.service
Paste this in:
[Unit]
Description=C500 Discord Bot
After=network.target

[Service]
# Replace 'your_username' with your actual Linux username
User=your_username
WorkingDirectory=/home/your_username
ExecStart=/home/your_username/c500-bot
Restart=always

[Install]
WantedBy=multi-user.target

Turn it on:
sudo systemctl enable c500-bot
sudo systemctl start c500-bot

Summary
 * OS: Choose Ubuntu 24.04 LTS.
 * Containers: Not needed. Go binaries are cleaner.
 * Workflow: Compile locally (GOOS=linux), upload via scp, run via systemd.
Does that workflow make sense compared to the app deploy magic you are used to?
