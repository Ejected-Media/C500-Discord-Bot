This is a phenomenal idea. It introduces an Affiliate Sales Layer directly into the C500 ecosystem.
By allowing a Builder to partner with a dedicated Marketer/Streamer, you are supercharging sales. The Builder focuses on craftsmanship (Instagram aesthetics), and the Marketer focuses on hype and audience reach (Twitch energy).
This changes the financial model from a 2-Way Split to a 3-Way Split.
Here is the system design for the C500 Builder + Marketer Partnership Program.

1. The New Financial Model: The 3-Way Split
Let's take a $500.00 sale as an example.


| Role | Share % | Amount | Responsibility |
|---|---|---|---|
| C500 Platform | 10% (Fixed) | $50.00 | Technology, Escrow, Safety. |
| The "Talent Pool" | 90% (Remaining) | $450.00 | The product and the promotion. |


Now, we split that remaining $450.00. You need a flexible system where the Builder and Marketer agree on their split beforehand. Let's say they agree on a 70/30 split of their pool.
 * Builder (70% of pool): $450.00 * 0.70 = **$315.00**
 * Marketer (30% of pool): $450.00 * 0.30 = **$135.00**
The C500 Bot will handle this math automatically.

2. Technical Implementation (Stripe & Database)
This requires a slight upgrade to our existing architecture.

A. The Stripe Connect Update
We are already using Stripe Express. This is perfect.
 * The Marketer must also onboard with Stripe Express, just like a Builder, so they have a stripe_connect_id.
 * When the payout trigger happens (Builder ships the item), the C500 backend will execute two separate transfers instead of one.

B. The Database Schema Update (Firestore)
We need to link the Marketer to the inventory item and track their payout.
1. Update builders collection to a broader partners collection
We can add a role field to distinguish them.

```
// Collection: partners (formerly builders)
// Document ID: "DISCORD_USER_ID_STRING"
{
  "role": "marketer", // or "builder"
  "stripe_connect_id": "acct_98765abc", // The Marketer's payout account
  "twitch_username": "HypeStreamerLive",
  // ... other fields
}
```


2. Update inventory collection
When a builder creates a drop, they "tag" their marketing partner.

```
// Collection: inventory
{
  "builder_id": "BUILDER_ID_STRING",
  "marketer_id": "MARKETER_ID_STRING", // <-- NEW FIELD (Optional)
  "marketer_split_percentage": 30,       // <-- NEW FIELD (e.g., 30%)
  "title": "Collab Build: Snowy TKL",
  "price_cents": 50000,
  // ...
}
```


3. Update orders collection
We need to record exactly how much everyone is getting paid for financial reporting.

```
// Collection: orders
{
  "total_paid_cents": 50000,
  "platform_fee_cents": 5000,   // C500's 10%
  "builder_payout_cents": 31500, // Builder's share
  "marketer_payout_cents": 13500, // Marketer's share
  // ...
}
```

3. The Bot Workflow (User Experience)
How do they set this up in Discord? It needs to be simple.

Step 1: Marketer Onboarding
 * Marketer: Joins the C500 server and types /c500 register marketer.
 * Bot: DMs them the Stripe Connect onboarding link. Once complete, they are ready to be tagged.

Step 2: Creating the Partnership Drop
 * Builder: Goes to their private channel and types the drop command. We add optional arguments for the partner.
 * Command: /c500 drop partner:@HypeStreamer split:30
   * partner: Tags the Marketer's Discord user.
   * split: Sets the Marketer's percentage of the 90% pool (e.g., 30%).

Step 3: The Sale & Payout
 * The drop is posted. It looks the same to the buyer.
 * Buyer pays $500. C500 holds it in escrow.
 * Builder streams the build and ships it.
 * Bot Trigger: The bot calculates the split and sends $315 to the Builder and $135 to the Marketer simultaneously via Stripe.

4. Visualizing the New Flow
Here is an updated graphic showing the 3-way financial split.
This system is powerful because it's flexible. A builder could do a 50/50 split for a huge streamer, or a 90/10 split for a smaller helper. C500 is just the neutral platform facilitating their agreement.
