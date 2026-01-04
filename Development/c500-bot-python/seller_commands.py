import discord
from discord import app_commands
from discord.ext import commands
from ui.modals import CreateDropModal

class SellerCommandsCog(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    # Define the main command group '/c500'
    c500_group = app_commands.Group(name="c500", description="C500 Collective Seller Commands")

    @c500_group.command(name="drop", description="Create a new listing marketplace drop.")
    async def create_drop(self, interaction: discord.Interaction):
        # 1. Check if user has the "Verified Seller" role in the Discord server
        # (Crucial security check before showing the modal)
        
        # 2. Pop up the modal form defined in ui/modals.py
        await interaction.response.send_modal(CreateDropModal())

    @c500_group.command(name="setup", description="Connect your Stripe account to start selling.")
    async def setup_seller(self, interaction: discord.Interaction):
        await interaction.response.defer(ephemeral=True, thinking=True)
        
        # Call Core API to get link
        link = await self.bot.core_api.get_stripe_onboarding_link(str(interaction.user.id))
        
        # Send link hidden (ephemeral) so only they see it
        await interaction.followup.send(f"Click here to connect Stripe: {link}", ephemeral=True)

# Standard setup function to load the cog
async def setup(bot):
    await bot.add_cog(SellerCommandsCog(bot))
  
