That is a massive and critical insight. You are absolutely right—in many parts of the world (including India, Brazil, and parts of Southeast Asia), the "Primary Computer" is an Android phone, not a laptop.
If Classroom 500 ignores them, we ignore a huge chunk of the market. If we optimize for them, we unlock a user base that expensive bootcamps can't touch.
Your research on Spck Editor is spot on. It is effectively "VS Code for Android" without the Microsoft branding. It supports Git natively (which is rare on mobile), meaning a student on a $50 Android phone can push code to the same GitHub repo as a student on a $3,000 MacBook.
📱 The "Android-First" Strategy
We need to formalize this into the curriculum. We aren't just "allowing" mobile; we are treating it as a first-class citizen.
Here is how we update the Independent Stack to be fully inclusive:
| Component | The Desktop Path | The Mobile Path (Android) | Why? |
|---|---|---|---|
| Editor | Zed | Spck Editor | Spck has a built-in local server and Git client. |
| Terminal | iTerm2 / Warp | Termux | Termux is a full Linux environment. It allows students to run npm or go on their phone. |
| Browser | Firefox Dev Edition | Firefox Nightly (Android) | Firefox Mobile allows remote debugging over USB, which is huge. |
| Collaboration | Zed Channels | Discord Mobile Screen Share | We can't do dual-cursor on mobile yet, but we can do "Over the Shoulder" mentoring. |
1. Updating the Curriculum (Level 00)
We need to split Level 00 into two tracks.
In CURRICULUM.md:
## 🟢 Level 00: The Setup (Choose Your Path)

### 🖥️ Path A: The Desktop Laboratory
*For students with Mac, Linux, or Windows machines.*
* **Editor:** Zed (Recommended) or Sublime Text.
* **Browser:** Firefox Developer Edition.
* **Terminal:** iTerm2 (Mac) or standard terminal.

### 📱 Path B: The Mobile Studio
*For students building the web from their pocket.*
* **Editor:** Spck Editor (Download from Play Store).
    * *Why:* It has a built-in "Preview" server and Git integration.
* **Browser:** Firefox for Android.
* **Keyboard:** Hacker's Keyboard (App) - *Adds Arrow keys and CTRL key to your touchscreen.*

2. Updating the Mentor Portal
This is crucial. A mentor who only knows how to debug in Chrome DevTools on a Mac might be useless to a kid using Spck on a Samsung Galaxy. We need to match them correctly.
Update to main.go (HTML Form):
<h3>Your Environment</h3>
<p style="color:#b9bbbe; font-size:0.8em;">
    C500 supports both Desktop and Mobile learners.
    Which environments are you comfortable teaching in?
</p>

<div style="margin-bottom:20px;">
    <input type="checkbox" name="env_desktop" value="true" checked>
    <label>Desktop (Mac/PC/Linux)</label>
    <br>
    <input type="checkbox" name="env_mobile" value="true">
    <label>Mobile (Android/Spck/Termux)</label>
</div>

Update to MentorApplication Struct:
type MentorApplication struct {
    // ... existing fields ...
    EnvDesktop bool `firestore:"env_desktop"`
    EnvMobile  bool `firestore:"env_mobile"`
}

3. The "Responsive Design" Twist
This actually gives C500 a unique pedagogical angle.
 * Desktop Students usually struggle to make their sites look good on mobile.
 * Mobile Students will struggle to make their sites look good on desktop.
We can create a "Peer Pairing" program:
> "Level 03 Project: The Cross-Device Swap."
> You pair a Desktop Student with a Mobile Student. They review each other's code. The Mobile student checks if the Desktop student's site works on a phone, and vice versa.
> 
This turns the hardware limitation into a learning superpower.
Does this "Two-Track" approach feel like the right way to honor your research?
