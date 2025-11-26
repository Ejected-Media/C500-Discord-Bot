This is the critical "Defense Layer" of your business. If you handle this poorly, you either become a haven for scammers (killing buyer trust) or a totalitarian regime (killing builder trust).
Because C500 is a Federated Marketplace (many servers, one economy), your Bad Actor Protocol must be centralized but transparent.
Here is the C500 Bad Actor Protocol designed to neutralize threats while protecting the innocent.
1. The Classification of "Bad Actors"
Not all offenses are equal. You need a tiered response system so you don't use a nuclear weapon on a shoplifter.


| Tier | Offense Type | Examples | Penalty |
|---|---|---|---|
| Tier 1 | The "Flake" (Time Waster) | Winning a raffle/drop but not paying; ignoring invoices for >24 hours. | Soft Ban (30 Days)
Cannot enter drops or buy items. |
| Tier 2 | The "Rude" (Community Violation) | Harassing a builder; toxic behavior in Discord; spamming support. | Server Ban (Local)
Banned from that specific Builder's server only. |
| Tier 3 | The "Fraud" (Nuclear) | Filing a Chargeback (Friendly Fraud); Using a stolen credit card; Lying about not receiving the item. | Global Blacklist (Permanent)

Banned from ALL C500 servers instantly. IP + Hardware ID flag. |
2. The "Global Ban" Architecture
This is the technical enforcement mechanism. Since your bot lives in 50+ different servers, it acts as a Global Security Grid.
The Workflow:
 * The Trigger: A Builder (Keyz) or System Admin detects a Tier 3 offense (e.g., a Chargeback).
 * The Command: Keyz types /c500 report @Alex [Reason: Chargeback Fraud] [Evidence Link].
 * The Verification: A C500 Admin (You) reviews the evidence. If valid, you issue the Kill Switch.
 * The Execution:
   * Database Update: User Alex is marked status: BANNED in the central database.
   * The Purge: The Bot instantly kicks/bans Alex from every single server where the C500 bot is installed.
   * The DM: Before the kick, the Bot sends a final DM: "You have been globally blacklisted from the C500 Marketplace for Fraud. To appeal, visit [c500.gg/appeals]."
3. Fighting "Friendly Fraud" (Chargebacks)
"Friendly Fraud" is when a user buys a keyboard, receives it, and then tells their bank "I didn't buy this." This is the #1 killer of niche marketplaces.
Your Twitch + Discord model gives you a "Silver Bullet" against this.
The "Evidence Pack" Generator:
When a chargeback hits Stripe, your Bot automatically compiles a PDF dossier to send to the bank. This wins disputes.
The Dossier Includes:
 * The Chat Log: Screenshots of the Discord ticket where the user discussed the build specs (proving intent).
 * The "Twitch Clip": A timestamped link to the VOD where the builder said "Okay, building this for Alex!" and Alex replied in chat "Hype!" (proving they were watching).
 * The Tracking: The FedEx/USPS delivery confirmation to the billing address.
 * The Terms: A highlighted copy of the ToS they clicked "Agree" to, specifically the "Custom Commission (No Refund)" clause.
Strategy: Banks usually side with the buyer unless you prove the buyer participated in the creation of the product. The Twitch clip is your "smoking gun."
4. The "Innocent User" Protection (Appeals Court)
You will ban an innocent person eventually (e.g., a hacked account). You need a way to fix it without Builders dealing with the drama.
The "Appeals Court" Logic:
 * Separation of Powers: Builders (Streamers) should not handle appeals. They are too busy and biased. Appeals go to C500 HQ (You or a dedicated mod).
 * The Loop:
   * Banned User fills out a form on your website (not Discord, since they are banned).
   * User proves identity (e.g., "My account was hacked, here is the ticket where I recovered it from Discord Support").
   * If approved, C500 HQ runs /c500 unban @User.
   * The Bot sends an invite link to the user to rejoin the servers they were kicked from.
5. The "Pre-Crime" Prevention (Stripe Radar)
Don't wait for them to steal. Stop them at the door.
Configure Stripe Radar with these custom rules for your platform:
 * Block if: IP Country does not match Card Country (High risk of stolen card).
 * Block if: CVC Check fails.
 * Review if: User attempts 3+ different cards in 10 minutes (Card testing).
 * Review if: Billing Address is a "Freight Forwarder" (Common for international scams).
Next Step
We have the Sales Flow, Trust Tiers, and Defense Protocol.
Would you like to draft the "Builder Agreement" (Terms of Service)? This is the contract streamers must sign to join C500, defining their responsibility to ship on time so they don't get banned by you.
