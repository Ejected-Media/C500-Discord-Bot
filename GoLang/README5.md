That is a totally fair and historically deep cut. The "Browser Wars" of the 90s (Netscape Navigator vs. Internet Explorer) were brutal, and if your family was on the Netscape/Apple side of that trench, avoiding Microsoft's ecosystem makes perfect sense.
If we want to capture that ScreenHero vibe (collaborative, dual-cursor editing) without touching VS Code, we have an incredible alternative that fits your Apple/Independent ethos perfectly: Zed.
The Alternative: Zed (The Spiritual Successor to the TextMate Era)
Zed is likely the perfect fit for C500 for a few reasons:
 * High Performance: It is written in Rust and is incredibly fast (honoring the efficient code spirit of the 90s).
 * Built-in Collaboration: It has a feature literally called "Channels" that allows multiple people to edit the same file in real-time. It is the closest thing to the old ScreenHero experience available today, and it is not owned by Microsoft.
 * Mac-First: It started as a Mac-exclusive project (though it's expanding), so it feels native to the ecosystem you prefer.
Updated Protocol: The "Independent Stack"
Let's adjust the Level 00: Setup and the Mentorship protocol to reflect this independence.
| Feature | The Microsoft Way (Avoid) | The C500 / "Independent" Way |
|---|---|---|
| Editor | VS Code | Zed (Collaborative) or Sublime Text (Classic) |
| Collaboration | Live Share | Zed Channels or JetBrains Code With Me |
| Browser | Edge / Chrome | Firefox (The spiritual successor to Netscape) |
| Terminal | Powershell | iTerm2 or Warp |
Updating the Mentor Portal
We should tweak the formHTML one last time. Instead of asking if they can do "VS Code Live Share," we should ask if they can do "Collaborative Pairing."
Updated Checkbox Logic:
<div style="display:flex; align-items:center; margin-bottom:20px;">
    <input type="checkbox" name="can_do_live" value="true" style="width:auto; margin:0 10px 0 0;">
    <label style="margin:0;">
        I am comfortable doing live "Pair Programming" sessions 
        (using Zed, JetBrains, or Screen Share)
    </label>
</div>

Updating Level 00 (Curriculum)
In your CURRICULUM.md, Level 00 becomes a stance statement:
 * Level 00: The Setup -> "Choosing your Weapon."
   * We recommend Zed for its speed and collaboration tools.
   * We use Firefox Developer Edition for its superior CSS Grid inspector (and Netscape heritage).
Does swapping VS Code for Zed (or allowing the student to choose any independent editor) sound like the right direction for the brand?
