This is the hardest part of any ed-tech project: Validation. The code is easy; getting people to pay is hard.
The difference this time—and why I think C500 has a better shot—is the "Discord First" architecture. You aren't just selling a video course (which anyone can pirate); you are selling access to a mentor in a specific channel.
That makes the Mentor Portal we just built critical. You can't scale the business if you are the only teacher.
Let's lock in this curriculum structure into your C500-Tech repo so it becomes the "Source of Truth" for both your website and your business model.
1. The CURRICULUM.md File
I have formatted your "Desk/Office" metaphor into a structure that looks like a product roadmap. You can drop this file into the root of your C500-Tech repo.
# Classroom 500: CSS Layout Curriculum

**Philosophy:** We don't just "teach CSS." We build a virtual workspace.
Every concept is tied to a physical object in a developer's office (The Desk, The Chair, The Monitor).

---

## 🟢 Level 00: The Setup (Free Tier / Lead Magnet)
*Target Audience: Absolute Beginners*
* **Theme:** "Entering the Office"
* **Goal:** Configure a professional Dev Environment.
* **Projects:**
    * Installing VS Code & Extensions.
    * Setting up Git & GitHub.
    * "Hello World" in an `index.html` file.
* **Outcome:** Student has a working local server and a GitHub repo.

## 🟡 Level 01: Typography & Aesthetics
* **Theme:** "The Brand"
* **Metaphor:** Choosing the paint color and furniture style.
* **CSS Properties:**
    * `color` vs `background-color`
    * `font-family` (Web Safe vs Google Fonts)
    * `font-weight` & `font-style`
    * `text-align`
* **Project:** Stylizing a "Nameplate" for the desk.

## 🟠 Level 02: The Box Model (The Core)
* **Theme:** "The Desk Structure"
* **Metaphor:** Every object is a box. Some boxes have padding (cushions), some have borders (frames).
* **CSS Properties:**
    * `width` & `height`
    * `padding` (The space *inside* the box)
    * `border` (The edge of the box)
    * `margin` (The space *outside* the box)
    * `box-sizing: border-box` (The magic fix)
* **Project:** Building the rectangular "Desk" and "Monitor" shapes.

## 🔴 Level 03: Semantics & Flow
* **Theme:** "Organization"
* **Metaphor:** Everything has a proper place.
* **Concepts:**
    * Block vs Inline elements.
    * The Document Flow (Gravity).
    * Semantic HTML (`<header>`, `<nav>`, `<main>`, `<footer>`).
* **Project:** structuring the HTML to match a real room layout.

---

## 🔮 Future Levels (Monetization Upsells)
* **Level 04:** Flexbox (Alignment & Distribution)
* **Level 05:** CSS Grid (2D Layouts)
* **Level 06:** Responsive Design (Mobile Adaptation)

2. Integrating this into the Go Code
To make this real, let's update the Mentor Application Form (main.go). We need to ask applicants: "Which level are you qualified to teach?"
This helps you segment your mentors. You might have a Junior Mentor who can handle Level 01, leaving the expensive Level 05 (Grid) teaching to you or Senior Mentors.
Update the formHTML const in your main.go:
<label>Which Curriculum Level can you teach?</label>
    <select name="teaching_level" required>
        <option value="" disabled selected>Select a Level...</option>
        <option value="lvl00">Lvl 00: Setup & Git (Beginner Support)</option>
        <option value="lvl01">Lvl 01: Text & Color</option>
        <option value="lvl02">Lvl 02: Box Model (Intermediate)</option>
        <option value="lvl03">Lvl 03: Semantics & Flow</option>
        <option value="all">All of the Above (Senior Mentor)</option>
    </select>


Update the MentorApplication Struct:
type MentorApplication struct {
    // ... existing fields ...
    TeachingLevel string `firestore:"teaching_level"` // NEW
}

Update the handleSubmit function:
// ... inside handleSubmit ...
app := MentorApplication{
    // ... existing fields ...
    TeachingLevel: r.FormValue("teaching_level"),
}

3. The Business Pivot (How to get paying users)
Since we are "starting semi-fresh," here is a strategy to avoid the fatigue of past attempts:
The "Tuition" Model:
 * Lvl 00 is Free: Use it to get people into your Discord. It's low effort for you but high value for them.
 * Lvl 01-03 is "Tuition": They don't pay for the videos (content is commodity). They pay for "Code Review".
 * The Mechanism:
   * Student watches your video/guide.
   * Student pushes code to GitHub.
   * Student pays $X/month.
   * C500 Bot (which we can build later) checks their role.
   * If they have the "Paid Student" role, they can post their PR link in the #code-review channel.
   * A Mentor (recruited via your new portal) reviews it.
This way, you aren't selling content; you're selling feedback. That is what people actually pay for.
Does that align with how you want to run C500 this time?
