import discord
from discord.ui import View, Button
import aiohttp
import logging
from utils import embeds # Import our cozy embed builders

# Set up logging
logger = logging.getLogger(__name__)

# The address of our Go Core microservice backend
CORE_API_URL = "http://localhost:8080/api/v1"

class BuyNowView(View):
    """
    A persistent view that attaches to a marketplace drop embed.
    It contains the "Buy Now" button.
    We set timeout=None so the button works forever (or until the bot restarts).
    """
    def __init__(self, drop_id: str, bot_session: aiohttp.ClientSession):
        super().__init__(timeout=None)
        self.drop_id = drop_id
        self.bot_session = bot_session

        # Add the button dynamically so we can pass the drop_id to it
        self.add_item(BuyButton(drop_id, bot_session))

class BuyButton(Button):
    """
    The actual clickable "Buy Now" button component.
    """
    def __init__(self, drop_id: str, session: aiohttp.ClientSession):
        super().__init__(
            style=discord.ButtonStyle.green,
            label="Buy Now <C500>",
            emoji="💸",
            # A unique custom_id helps with persistence debugging if needed
            custom_id=f"buy_btn:{drop_id}"
        )
        self.drop_id = drop_id
        self.session = session

    async def callback(self, interaction: discord.Interaction):
        """
        This method triggers whenever the button is clicked.
        It handles the handshake with the Core API to get a checkout link.
        """
        # 1. Immediate Acknowledge: Show the user the bot is working.
        # ephemeral=True is CRITICAL here. Only the clicker sees this response.
        # We use our cozy "thinking" embed.
        await interaction.response.send_message(
            embed=embeds.thinking_embed("Contacting secure checkout..."),
            ephemeral=True
        )

        buyer_id = str(interaction.user.id)
        api_endpoint = f"{CORE_API_URL}/checkout/session"

        # Prepare payload for Core API
        payload = {
            "drop_id": self.drop_id,
            "buyer_discord_id": buyer_id
        }

        try:
            # 2. Call Core API to generate Stripe session
            async with self.session.post(api_endpoint, json=payload, timeout=10) as response:

                if response.status == 200:
                    # Success! Parse the response to get the URL.
                    data = await response.json()
                    checkout_url = data.get("url")

                    if not checkout_url:
                         raise ValueError("API response missing checkout URL")

                    # 3. Create a new view with a Link Button pointing to Stripe
                    # We don't need a custom class for simple link buttons.
                    link_view = View()
                    link_view.add_item(Button(
                        label="👉 Proceed to Secure Checkout",
                        url=checkout_url,
                        style=discord.ButtonStyle.link
                    ))

                    # 4. Update the ephemeral message with the link.
                    await interaction.edit_original_response(
                        content="**Click below to complete your purchase securely on Stripe.**\n*This link expires in 30 minutes.*",
                        embed=None, # Remove the thinking embed
                        view=link_view
                    )

                elif response.status == 409:
                    # Conflict: Item already pending or sold.
                    await interaction.edit_original_response(
                        embed=embeds.error_embed("Sorry, this item is currently pending purchase by someone else!", title="Too Late!")
                    )
                elif response.status == 404:
                     # Drop ID doesn't exist in DB.
                    await interaction.edit_original_response(
                        embed=embeds.error_embed("This listing appears to be invalid or expired.")
                    )
                else:
                    # Generic server error
                    logger.error(f"Checkout API Error for drop {self.drop_id}: Status {response.status}")
                    await interaction.edit_original_response(
                        embed=embeds.error_embed("An internal error occurred preparing checkout. Please try again.")
                    )

        except aiohttp.ClientError as e:
            logger.error(f"Network error contacting Core API during checkout attempt: {e}")
            await interaction.edit_original_response(
                 embed=embeds.error_embed("Could not reach the checkout server. Please try again later.", title="Connection Error")
            )
        except Exception as e:
             logger.exception(f"Unexpected error in buy button callback: {e}")
             await interaction.edit_original_response(
                 embed=embeds.error_embed("An unexpected error occurred.")
      )
          
