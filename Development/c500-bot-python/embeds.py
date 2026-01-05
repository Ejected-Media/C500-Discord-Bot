import discord
from datetime import datetime
from typing import Optional

# ==========================================
# C500 Brand Color Palette (Pastel/Cozy Aesthetic)
# ==========================================
COLOR_PRIMARY = 0xFFB7C5  # Sakura Pink (Used for main listing embeds)
COLOR_SECONDARY = 0xAEC6CF # Pastel Blue (Used for info/thinking states)
COLOR_SUCCESS = 0x77DD77  # Pastel Green (Used for confirmations)
COLOR_ERROR = 0xFF6961    # Pastel Red (Used for errors, but soft, not alarming)
COLOR_GOLD = 0xFDFD96     # Pastel Yellow (Used for special alerts/live status)

# Standard Footer Text
FOOTER_TEXT = "C500 Collective • Cozy Builds & Community"
FOOTER_ICON_URL = "https://c500.store/static/images/bot-icon-small.png" # Placeholder URL

# ==========================================
# Specialized Embed Builders
# ==========================================

def create_marketplace_drop_embed(
    title: str,
    description: str,
    price_str: str,
    seller_name: str,
    drop_id: str,
    image_url: Optional[str] = None,
    drop_type: str = "rts"
) -> discord.Embed:
    """
    Creates the main, highly visible embed for a new item listing in the #marketplace channel.
    This is the "storefront window" for an item.
    """
    # Choose an emoji based on drop type
    type_emoji = "📦" if drop_type == "rts" else "🎨"
    type_label = "Ready-to-Ship" if drop_type == "rts" else "Commission Slot"

    embed = discord.Embed(
        title=f"✨ New Drop: {title}",
        description=description,
        color=COLOR_PRIMARY,
        timestamp=datetime.utcnow()
    )

    # Prominent Price Field
    embed.add_field(
        name="🏷️ Price",
        value=f"**{price_str}**",
        inline=True
    )

    # Type Field
    embed.add_field(
        name=f"{type_emoji} Type",
        value=type_label,
        inline=True
    )

    # Drop ID (small, for reference)
    embed.add_field(
        name="🆔 Drop ID",
        value=f"`{drop_id}`",
        inline=False # new line
    )

    # The hero image of the keyboard
    if image_url and image_url.startswith("http"):
        embed.set_image(url=image_url)

    # Set the seller in the author slot at the top
    embed.set_author(name=f"Listed by {seller_name}", icon_url=FOOTER_ICON_URL)

    # Consistent footer
    embed.set_footer(text=FOOTER_TEXT, icon_url=FOOTER_ICON_URL)

    return embed


def create_order_confirmation_dm(
    order_id: str,
    item_title: str,
    price_str: str,
    expected_action: str
) -> discord.Embed:
    """
    Sent privately to a buyer after a successful Stripe payment.
    """
    embed = discord.Embed(
        title="🎉 Order Confirmed!",
        description=f"Hurray! Your purchase of **{item_title}** is confirmed.",
        color=COLOR_SUCCESS,
        timestamp=datetime.utcnow()
    )

    embed.add_field(name="Total Paid", value=price_str, inline=True)
    embed.add_field(name="Order #", value=f"`{order_id}`", inline=True)

    # Inform them what happens next (e.g., "Wait for shipping" or "Watch streamer")
    embed.add_field(
        name="What's Next?",
        value=expected_action,
        inline=False
    )

    embed.set_footer(text="Thank you for supporting C500 builders!")
    return embed

# ==========================================
# Generic Utility Embeds
# ==========================================

def success_embed(message: str) -> discord.Embed:
    """A simple, reusable green success message."""
    return discord.Embed(
        description=f"✅ {message}",
        color=COLOR_SUCCESS
    )

def error_embed(message: str, title: str = "Oops!") -> discord.Embed:
    """A simple, reusable soft red error message."""
    #We use a soft title to keep the anxiety low.
    return discord.Embed(
        title=f"⚠️ {title}",
        description=message,
        color=COLOR_ERROR
    )

def thinking_embed(message: str = "Processing...") -> discord.Embed:
    """Used for ephemeral deferrals while waiting for the Core API."""
    return discord.Embed(
        description=f"🔄 {message}",
        color=COLOR_SECONDARY
  )
  
