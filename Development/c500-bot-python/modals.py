import discord
from discord import ui
from services.core_api_client import CoreAPIClient

class CreateDropModal(ui.Modal, title='Create C500 Drop'):
    # Define the fields
    item_title = ui.TextInput(label='Item Title', placeholder='e.g., Snowy TKL Keyboard', max_length=100)
    price = ui.TextInput(label='Price ($)', placeholder='450.00')
    # ... other fields like Description, Image URL ...

    # The magic function that runs when they click "Submit" on the form
    async def on_submit(self, interaction: discord.Interaction):
        # 1. Acknowledge immediately so Discord doesn't time out
        await interaction.response.defer(ephemeral=True, thinking=True)

        # 2. Gather data needed for backend
        drop_payload = {
            "seller_discord_id": str(interaction.user.id),
            "title": self.item_title.value,
            "price": float(self.price.value),
            # ...
        }

        # 3. Send to Core API using the client service
        api_client = interaction.client.core_api  # Access the shared client from main.py
        try:
            response = await api_client.submit_drop(drop_payload)
            await interaction.followup.send("✅ Drop created successfully! Posting embed...", ephemeral=True)
            # (Normally you would now call a function to post the public embed)
        except Exception as e:
            await interaction.followup.send(f"❌ Error creating drop: {e}", ephemeral=True)
          
