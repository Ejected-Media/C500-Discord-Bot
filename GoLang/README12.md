Since you recently purchased C500.tech and we have established your preferred Go + Firebase stack, I have a very specific plan to bridge your "Teacher/Mentor Onboarding" goal with your new infrastructure.
The biggest friction point in community learning communities (like Discord) is vetting mentors manually. We can automate the "intake" phase using your new domain.
Here is the proposal: The C500 Mentor Intake Portal.
1. The Workflow Concept
Instead of a manual Google Form, we build a sleek Go web app hosted on C500.tech.
 * The Hook: User visits c500.tech/mentors.
 * The Identity: They click "Login with Discord" (Using Go OAuth2). This proves they are real and captures their Discord ID.
 * The Vetting: They submit a form with:
   * Their "CSS Philosophy" (Text).
   * 3 Links to CodePens or GitHub Gists (Proof of skill).
   * Twitch Channel URL (for streaming classes).
 * The Integration: The Go server saves this to Firebase and automatically pings a private channel in your Discord server with the application details for you to review.
2. The Tech Stack for C500.tech
This fits your "Selfielunge" and "Drift Miner" preferences perfectly:
 * Backend: Go (App Engine).
 * Database: Firestore (storing applications).
 * Auth: Discord OAuth2 (Critical for linking the web user to the Discord server).
 * Frontend: Standard Go Templates (keeping it simple/fast) or a light React/Vue layer.
3. Why start here?
It solves the immediate problem ("How do I onboard people?") by building a piece of permanent infrastructure on your new domain. It validates the domain purchase immediately.
Would you like to start by writing the Go code to handle the "Login with Discord" OAuth flow? (This is the trickiest part but the most powerful).
