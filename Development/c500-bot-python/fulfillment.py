import discord
from discord import app_commands
from discord.ext import commands
import aiohttp
import os
import logging

# Set up basic logging
logger = logging.getLogger(__name__)

# In a real app, grab this from os.getenv("CORE_API_BASE_URL")
# This is the address of our Go Core microservice (e.g., running in Cloud Run)
CORE_API_URL = "http://localhost:8080/api/v1"

class FulfillmentCog(commands.Cog):
    """
    Handles commands related to order fulfillment by sellers.
    These commands are restricted to users who are verified sellers.
    """

    def __init__(self, bot: commands.Bot):
        self.bot = bot
        # We use the bot's shared aiohttp session for efficient API calls
        self.session = bot.http_session

    # We group these commands under '/fulfill' so in Discord it looks like:
    # /fulfill ship ...
    # /fulfill live ...
    fulfill_group = app_commands.Group(name="fulfill", description="Commands to complete orders and release funds.")

    async def cog_check(self, interaction: discord.Interaction) -> bool:
        """
        Global check for this cog: Is the user allowed to perform fulfillment?
        In a real app, we would check a cached role or ask the Core API if
        this user is a 'verified_builder'.
        For now, we'll assume true to show the logic flow.
        """
        # Example check:
        # Use discord.utils.get(interaction.user.roles, name="Verified Builder")
        return True

    # =========================================
    # Command: /fulfill ship
    # =========================================
    @fulfill_group.command(name="ship", description="Mark an RTS order as shipped with tracking.")
    @app_commands.describe(
        order_id="The Order ID provided when the item sold",
        tracking_number="The tracking number from the carrier",
        carrier="e.g., UPS, USPS, DHL"
    )
    async def fulfill_ship(
        self,
        interaction: discord.Interaction,
        order_id: str,
        tracking_number: str,
        carrier: str
    ):
        # 1. Acknowledge the interaction so Discord doesn't timeout.
        # Ephemeral=True means only the user running the command sees the response.
        await interaction.response.defer(ephemeral=True, thinking=True)

        # 2. Prepare the payload for the Go Core API endpoint.
        # We assume the API expects POST /api/v1/orders/{id}/fulfill/ship
        api_endpoint = f"{CORE_API_URL}/orders/{order_id}/fulfill/ship"
        payload = {
            "tracking_number": tracking_number,
            "carrier": carrier,
            # Send the Discord ID so Core API ensures the command runner owns the order
            "seller_discord_id": str(interaction.user.id)
        }

        try:
            # 3. Send request to the Go Core API
            async with self.session.post(api_endpoint, json=payload) as response:
                if response.status == 200:
                    # Success! The Core API handled the database and Stripe update.
                    await interaction.followup.send(
                        f"✅ **Success!** Order `{order_id}` marked as shipped via {carrier}.\n"
                        f"Tracking: `{tracking_number}`.\n"
                        "Funds will be released shortly once tracking activates."
                    )
                elif response.status == 404:
                     await interaction.followup.send(f"❌ Error: Order ID `{order_id}` not found or does not belong to you.")
                elif response.status == 400:
                     # Bad request, maybe the order is already shipped
                     data = await response.json()
                     await interaction.followup.send(f"⚠️ Cannot ship: {data.get('error', 'Invalid request')}")
                else:
                    logger.error(f"API Error shipping order {order_id}: Status {response.status}")
                    await interaction.followup.send("🔥 An internal API error occurred. Please contact support.")

        except aiohttp.ClientError as e:
            logger.error(f"Network error connecting to Core API: {e}")
            await interaction.followup.send("📡 Could not connect to the C500 Core API. Try again later.")


    # =========================================
    # Command: /fulfill live
    # =========================================
    @fulfill_group.command(name="live", description="Mark a commission order as completed via live stream.")
    @app_commands.describe(
        order_id="The Order ID being built",
        vod_url="Link to the Twitch/YouTube VOD or clip proving completion"
    )
    async def fulfill_live(
        self,
        interaction: discord.Interaction,
        order_id: str,
        vod_url: str
    ):
        await interaction.response.defer(ephemeral=True, thinking=True)

        # Basic validation on the bot side before bothering the API
        if "twitch.tv" not in vod_url and "youtube.com" not in vod_url:
             await interaction.followup.send("⚠️ Please provide a valid Twitch or YouTube URL.")
             return

        api_endpoint = f"{CORE_API_URL}/orders/{order_id}/fulfill/live"
        payload = {
            "vod_url": vod_url,
            "seller_discord_id": str(interaction.user.id)
        }

        try:
            async with self.session.post(api_endpoint, json=payload) as response:
                if response.status == 200:
                    await interaction.followup.send(
                        f"🎥 **Success!** Order `{order_id}` marked as verified LIVE.\n"
                        f"VOD linked: <{vod_url}>\n"
                        "Funds have been released to your balance!"
                    )
                elif response.status == 404:
                     await interaction.followup.send(f"❌ Error: Order ID `{order_id}` not found.")
                else:
                     # Handle other API errors generically for now
                     data = await response.json()
                     await interaction.followup.send(f"⚠️ API Error: {data.get('error', 'Unknown error')}")

        except aiohttp.ClientError as e:
             logger.error(f"Network error: {e}")
             await interaction.followup.send("📡 Connection error to Core API.")


# Standard setup function for discord.py cogs
async def setup(bot: commands.Bot):
    await bot.add_cog(FulfillmentCog(bot))
          
