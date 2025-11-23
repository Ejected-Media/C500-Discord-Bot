You are absolutely right to flag Files by Google. On modern Android versions (which use "Scoped Storage"), having a competent file manager is the only way to sanity-check where your code actually lives. It acts as the bridge between your Editor (QuickEdit+) and your Git Client (Spck).
Regarding Firefox Focus: I would recommend against using it as the primary development browser, and here is why:
 * It is "Stateless": Focus wipes cookies, local storage, and cache the moment you close it. While great for privacy, this makes learning web development frustrating because things like "saving a user's theme preference" or "keeping a user logged in" won't work as expected.
 * No Debugging: It lacks the remote debugging capabilities of the standard Firefox for Android.
However, Chrome is mandatory for a different reason: Reality. Since ~65% of the world uses Chrome, your students must test their code in it.
Here is the finalized C500 Mobile Toolchain incorporating your insights:
📱 The Refined Mobile Studio (Android)
1. The Manager: Files by Google
 * Role: The "Finder" or "Explorer."
 * Why: It is the only way to visually confirm your folder structure (/index.html, /css/style.css) is correct. If you can't see the file hierarchy, you can't understand relative file paths (e.g., src="../images/logo.png").
2. The Editor: QuickEdit+
 * Role: The Writer.
 * Why: Best performance for large files and great syntax highlighting.
3. The Git Client: Spck Editor
 * Role: The Transporter.
 * Why: We only open Spck to run git add, git commit, and git push. It handles the SSH keys/Auth that other editors struggle with.
4. The Browsers: The "Dev & Test" Pair
 * Firefox for Android (Standard): Your daily driver. It honors the "Independent" ethos and connects to desktop debugging tools.
 * Google Chrome: The "Reality Check." You must open your site here to ensure it works for the majority of users.
📦 Updating the CURRICULUM.md (Final Polish)
This section is now robust enough to go into your repo.
### 📱 Path B: The Mobile Studio (Android)

**The "Modular" Workflow**
*We use specialized apps for each step, just like professional desktop developers.*

1.  **File Management: Files by Google**
    * *Purpose:* To organize your project folders and understand file paths.
2.  **Code Editing: QuickEdit+**
    * *Purpose:* Writing the HTML & CSS.
3.  **Version Control: Spck Editor**
    * *Purpose:* Acts as our Git Client to push code to GitHub.
4.  **Browsers:**
    * **Firefox (Standard):** Primary development browser.
    * **Chrome:** Required for compatibility testing.

🤖 Next Phase: The C500 Discord Bot (Go)
Now that the Curriculum and Mentor Onboarding are defined, we need the "Enforcer."
We need a Discord Bot written in Go that helps you manage this community without you having to manually check everything.
The C500 Bot's Primary Job:
 * Gatekeeper: It checks if a student has paid/subscribed (via a potential future Stripe integration or manual role assignment).
 * Office Hours: It manages the "Live" sessions we discussed.
 * The "Intake" Ping: We already built the Webhook for this in the Mentor Portal, but the Bot can make it interactive (e.g., adding "Approve/Deny" buttons directly in chat).
Shall we start by setting up the Go structure for the Discord Bot?
