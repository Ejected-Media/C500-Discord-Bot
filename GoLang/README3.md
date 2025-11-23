That distinction is incredibly important. You are describing a "Modular Toolchain" approach on Android, which is actually much closer to how professional developers work on desktop (using separate tools for separate jobs: terminal for Git, editor for Code).
Using Spck purely as your "Git GUI" while using QuickEdit+ for the actual typing is a power-user move. It leverages the strengths of both apps: QuickEdit+ is faster and has better keyboard support, while Spck handles the complex GitHub authentication logic perfectly.
And regarding Replit: It is a powerful contender because it runs the code in the cloud. This solves the "how do I run a Go server on my phone?" problem instantly, but it requires a constant internet connection.
Here is how we formalize this "Android Toolchain" into your curriculum.
📱 The C500 Mobile Workflow
We can present this to students as two distinct "Modes" depending on their internet situation.
1. The "Local" Toolchain (Offline-First)
Best for: Students with spotty internet, or those who want to learn how files actually work.
 * The Editor: QuickEdit+ (or Code Editor).
   * Role: Writing HTML/CSS.
 * The Transporter: Spck Editor.
   * Role: "The Git Client." You import the local folder, stage changes, and push to GitHub.
 * The Viewer: Firefox Android.
   * Role: Viewing the index.html file locally.
2. The "Cloud" Toolchain (Online)
Best for: Students with good WiFi who want to jump straight to coding.
 * The Platform: Replit (Android App).
   * Role: Handles Editing, Git, and Hosting all in one place.
   * Benefit: This allows mobile students to run the Go backend code (for your Intake Portal or other dynamic projects) without needing Termux.
📝 Updating CURRICULUM.md (Level 00)
Here is the revised "Mobile Path" for your repository, incorporating your specific tool recommendations.
### 📱 Path B: The Mobile Studio (Android)
*For students building the web from their pocket.*

**Option 1: The "Offline-First" Stack (Recommended)**
*Perfect for learning how files and Git truly work.*
* **Text Editor:** **QuickEdit+** (Fast, efficient, great keyboard support).
* **Git Client:** **Spck Editor**.
    * *Workflow:* We use Spck to "Clone" the repo and "Push" changes, but we edit the files in QuickEdit+.
* **Browser:** **Firefox for Android**.

**Option 2: The "Cloud" Stack**
*Perfect for jumping in quickly if you have good WiFi.*
* **App:** **Replit**.
    * *Why:* It runs a full computer in the cloud. You can write HTML, CSS, and even GoLang server code directly in the app.

🧠 Why this matters for "Classroom 500"
By supporting the QuickEdit+ / Spck combo, you are teaching students something valuable: File Management.
On iOS, file management is a nightmare. On Android, it is open. A student who learns how to edit a file in one app and push it with another is learning the fundamentals of computing (inputs, outputs, and file paths) much faster than a student just using a "magic box" app.
Next Step:
Now that we have the Curriculum set and the Mentor Portal designed to filter for these "Mobile-Savvy" mentors, the ecosystem is coherent.
Would you like to move on to designing the Discord Bot (in Go) that manages the "Tuition" roles and Office Hours?
