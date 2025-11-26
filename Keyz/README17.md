This is an exciting moment. We have successfully moved from a high-level concept—"selling mechanical keyboards on Discord"—to a fully fleshed-out, scalable, and defensible business model with a unique brand identity.
We have defined the Business, the Brand, the Legal, and the Data Structure.
The next logical step is to package all of this into a single, coherent document that you can hand to a developer (or use yourself) to start building the Alpha version of the bot.
Let's assemble the C500 Technical Briefing. This is the instruction manual for turning our conversations into code.
C500 Marketplace: Technical Briefing Document
Version: 1.0 (Alpha Build)
Objective: Build a federated Discord bot that manages inventory, processes Stripe payments via escrow, and integrates with Twitch for vendor verification.
1. The 30-Second Context
C500 is not a single server. It is a bot that lives in the private Discord servers of high-end keyboard builders ("The Collective").
 * The bot allows builders to post items for sale in their own channels via embeds.
 * Buyers click a button on the embed to pay C500 directly via Stripe.
 * C500 holds the funds.
 * The builder must verify the build (via Twitch stream or tracking number) to trigger the fund release to their Stripe Connect account.
2. The Tech Stack (Confirmed)
 * Bot Language: Python (using discord.py library).
 * Database: Google Firestore (NoSQL).
 * Payments: Stripe Connect (Express Accounts, using Destination Charges).
 * Live Streaming: Twitch API (EventSub/Webhooks for real-time status).
 * Hosting: (TBD, e.g., Heroku, AWS, DigitalOcean).
3. Required Environment Variables (.env)
The developer will need these keys to make it work.

```
# Discord
DISCORD_BOT_TOKEN=your_bot_token_here
DISCORD_CLIENT_ID=your_app_id
DISCORD_CLIENT_SECRET=your_secret

# Stripe (Test Mode first!)
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...

# Twitch
TWITCH_CLIENT_ID=your_twitch_id
TWITCH_CLIENT_SECRET=your_twitch_secret
```

4. The Database Schema
(Refer to the previously agreed-upon Firestore JSON structure for collections: users, builders, inventory, and orders.)
5. Critical Logic Flows (The "Hard Parts")
Flow A: The Builder Onboarding (Stripe Connect)
 * Trigger: Builder types /c500 setup payment in their private staff channel.
 * Bot Action:
   * Calls Stripe API to generate an Express Account Link.
   * DMs the link to the Builder as an ephemeral message: "Click here to set up how you get paid. Do not share this link."
 * Success State: Stripe sends a webhook to your server confirming account creation. The Bot updates the builders Firestore document with the new stripe_connect_id.
Flow B: The Transaction (Escrow)
 * Trigger: User clicks "BUY NOW" button on an inventory embed.
 * Bot Action:
   * Checks inventory doc: Is status still "available"?
   * Checks users doc: Does user meet Reputation Tier requirements for this price?
   * If pass: Calls Stripe API to create a Checkout Session.
     * Important: Must set payment_intent_data[on_behalf_of] to the Builder's Connect ID.
     * Important: Must set payment_intent_data[transfer_data][destination] to the Builder's Connect ID.
   * Sends ephemeral reply with the Stripe payment link.
Flow C: The Twitch Verification Trigger
 * Trigger: Builder types /c500 live [Order_ID].
 * Bot Action:
   * Looks up the Builder's twitch_username in Firestore.
   * Pings Twitch API: Is this user currently streaming?
   * If YES:
     * Updates orders doc status to fulfillment_status: "building".
     * DMs the buyer linked to that Order ID with the stream link.
   * If NO: Returns error message to builder: "You aren't live on Twitch right now."
6. UI/UX Style Guide ("The Cozy Aesthetic")
All Discord Embeds sent by the bot must adhere to this style.
 * Color Hex Codes:
   * Primary/Buy Buttons: #FFD1DC (Sakura Milk)
   * Secondary/Info: #AEEEEE (Dreamy Sky)
   * Success/Verified: #98FF98 (Mint Matcha)
   * Error/Sold Out: #E6E6FA (Lavender Mist - used instead of harsh red)
 * Tone of Voice:
   * Friendly, encouraging, professional but not corporate.
   * Bad: "Transaction executed successfully."
   * Good: "Got it! The funds are secured. Time to get cozy while [Builder] gets to work."
What's next after this briefing?
Once this document is handed off to begin development, the next logical step for us (the business side) is to plan the Day 1 Launch Operations.
Would you like to map out that Launch Day Checklist (e.g., the exact sequence of events for announcing the platform, onboarding the first 5 builders, and managing the first live drop)?
