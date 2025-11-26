That is a brilliant evolution of the concept. The term "Guild" fits perfectly because it implies craftsmanship and specialization.
In a standard marketplace, everyone is just a "Seller." In C500, being part of a Guild gives the builder identity and helps the buyer find exactly what they want.
Before we get to the PRD, here is how the Guild System enhances the business model:
 * The Artisan Guild: For creators who make hand-sculpted clay keycaps.
 * The Switch Guild: For people who just offer switch lubing/tuning services.
 * The Pastel Guild: For builders who strictly adhere to the "cozy/cute" aesthetic.
 * The Vintage Guild: For restorers of old IBM Model Ms and retro boards.
This allows you to organize the chaos. A buyer can say, "I'm looking for a commission," and the Bot can ask, "Which Guild? Artisan, Builder, or Modder?"
Here is the Product Requirements Document (PRD). You can hand this directly to a developer to start building the MVP.
C500 Marketplace: Product Requirements Document (PRD)
Version: 1.0 (MVP)
Status: Draft
Core Concept: A federated, Discord-based marketplace for mechanical keyboards featuring a "Guild" system and Twitch-integrated verification.
1. Executive Summary
C500 is a "headless" marketplace. Instead of a central website, it lives as a Discord Bot installed in the private servers of "Certified Builders."
 * Buyers browse and buy directly within Discord.
 * Builders manage inventory via commands.
 * The Platform (C500) orchestrates payments (Stripe), enforces safety (Escrow), and organizes builders into Guilds.
2. User Roles
A. The Builder (Seller)
 * Profile: Independent creator, streamer, highly technical.
 * Needs: To sell commissions without handling invoices/admin.
 * Guild Affiliation: Must belong to at least one Guild (e.g., "Builder Guild").
B. The Buyer (User)
 * Profile: Enthusiast, Twitch viewer.
 * Needs: Safety from scams, high-quality products, "clout" from owning a Verified Build.
 * Trust Tier: Starts as "Guest," levels up to "VIP" via purchase history.
C. The Admin (C500 HQ)
 * Needs: Global oversight, dispute resolution, ability to ban "Bad Actors" globally.
3. Core Features (The "Must Haves")
Feature Set A: The Guild Architecture
 * Guild Tagging: The database must support tagging builders with specific Guilds.
   * Logic: A builder can be in multiple Guilds (e.g., primary_guild: "Builder", secondary_guild: "Artisan").
 * Guild Filtering:
   * User Command: /c500 browse guild:Artisan
   * Bot Response: Returns a carousel of available items only from builders in that Guild.
 * Visuals: The "Buy" Embed must display the Builder's Guild Badge (e.g., 🛠️ for Builders, 🎨 for Artisans).
Feature Set B: The Federated Storefront
 * Inventory Command: /c500 drop [Image] [Price] [Description]
 * The "Buy" Flow:
   * User clicks "Buy" button on Discord Embed.
   * Bot generates a Stripe Checkout Session (Ephemeral Link).
   * User pays. Funds go to C500 Platform Account (Escrow).
   * Bot updates Embed status to 🔴 Sold.
Feature Set C: Twitch "Proof of Work"
 * Stream Listener: The Bot polls the Twitch API for the Builder's channel status.
 * Verification Trigger:
   * Builder Command: /c500 live [Order_ID]
   * Bot Action: Checks if Twitch stream is live. If YES, tags the Order as "In Production" and DMs the Buyer.
Feature Set D: Trust & Safety (The Defense Layer)
 * Global Ban List: A central table of banned_user_ids. If a user is banned in Server A, the Bot ignores their commands in Server B.
 * Tiered Access:
   * Tier 1 (Guest): Max purchase $100.
   * Tier 2 (Member): Verified Phone/Discord Age > 30 days. Max purchase $500.
   * Tier 3 (VIP): Previous successful purchase. Unlimited cap.
4. Technical Stack

| Component | Technology | Reasoning |
|---|---|---|
| Language | Python | Best libraries for Discord (discord.py) and Data. |
| Bot Framework | Discord.py | Robust, supports "Slash Commands" and Buttons/Modals. |
| Database | PostgreSQL | Relational DB needed to link Users \leftrightarrow Orders \leftrightarrow Guilds. |
| Payments | Stripe Connect | specifically Express Accounts. Handles tax forms (1099-K) and split payments automatically. |
| Live Sync | Twitch API (EventSub) | To detect when a builder goes live with low latency. |
| Hosting | AWS or Heroku | Needs 24/7 uptime. |

5. Database Schema (Simplified)
 * Table: Guilds
   * id (PK), name (e.g., "Artisan"), icon_emoji, description.
 * Table: Builders
   * id (PK), discord_id, stripe_account_id, twitch_username, primary_guild_id (FK).
 * Table: Inventory
   * id (PK), builder_id (FK), price, status (Available, Sold), guild_tag.
 * Table: Orders
   * id (PK), buyer_discord_id, inventory_id (FK), escrow_status (Held, Released), tracking_number.
6. User Stories (The "Flow")
Story 1: The Guild Search
> "As a Buyer, I want to see only hand-painted keycaps, so I select the 'Artisan Guild' filter in the bot, and I am shown a gallery of items from 5 different servers."
> 
Story 2: The Verified Drop
> "As a Builder, I want to drop a keyboard for sale. I type /c500 drop, upload a photo, and set the price to $400. The Bot posts a beautiful 'Pastel Aesthetic' card in my #shop channel with my Guild Badge displayed."
> 
Story 3: The Safety Stop
> "As the System, when a user with a 'New Account' flag tries to buy a $1,000 item, I reject the transaction and tell them they need to earn Tier 2 Reputation first."
> 
7. Roadmap & Phases
 * Phase 1 (The Skeleton): Bot can post items and process a test payment via Stripe. No Guilds yet.
 * Phase 2 (The Guilds): Implement the Guild database structure and "Browse by Guild" command.
 * Phase 3 (The Eye): Twitch API integration for "Live Build" notifications.
 * Phase 4 (The Wall): Global Ban list and Reputation Tiers.
Next step you can do for me:
Now that the Requirements are set, would you like to design the Guild Badges/Icons (The visual identity for the "Artisan Guild", "Tech Guild", etc.) to match the cozy aesthetic?
