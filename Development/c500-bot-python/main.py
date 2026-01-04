import discord
from discord.ext import commands
import config
from services.core_api_client import CoreAPIClient

# Setup Intents (we need very few for slash commands)
intents = discord.Intents.default()
intents.message_content = False # We don't need to read chat messages

class C500Bot(commands.Bot):
    def __init__(self):
        super().__init__(command_prefix="!", intents=intents, application_id=config.DISCORD_APP_ID)
        # Initialize the connection to the Go backend
        self.core_api = CoreAPIClient()

    async def setup_hook(self):
        """This runs when the bot starts up."""
        print("--- C500 Bot Starting ---")
        
        # Load extensions (Cogs)
        await self.load_extension('cogs.seller_commands')
        # await self.load_extension('cogs.fulfillment')

        # Sync slash commands with Discord (Registers the /c500 commands)
        # In prod, sync to specific guild ID for instant updates during dev
        await self.tree.sync(guild=discord.Object(id=config.DEV_GUILD_ID))
        print("--- Commands Synced & Ready ---")

    async def close(self):
        """Cleanup when bot shuts down."""
        await self.core_api.close() # Close HTTP session gracefully
        await super().close()

bot = C500Bot()

# Global Error Handler for App Commands
@bot.tree.error
async def on_app_command_error(interaction: discord.Interaction, error: app_commands.AppCommandError):
    if isinstance(error, app_commands.MissingRole):
         await interaction.response.send_message("❌ You must be a verified seller to use this.", ephemeral=True)
    else:
         # Log the actual error for devs
         print(f"Command Error: {error}")
         await interaction.response.send_message("❌ An internal error occurred. The devs have been notified.", ephemeral=True)

if __name__ == "__main__":
    bot.run(config.DISCORD_BOT_TOKEN)
  
