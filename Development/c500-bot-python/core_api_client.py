import aiohttp
import config

class CoreAPIClient:
    """Handles all communication with the c500-core-go microservice."""
    
    def __init__(self):
        self.base_url = config.CORE_API_URL
        # We use a shared session for performance
        self.session = aiohttp.ClientSession() 

    async def close(self):
        await self.session.close()

    # --- Seller Endpoints ---

    async def get_stripe_onboarding_link(self, discord_user_id: str):
        """Calls backend to generate a Stripe Express link."""
        endpoint = f"{self.base_url}/api/internal/seller/onboard"
        payload = {"discord_id": discord_user_id}
        
        async with self.session.post(endpoint, json=payload) as resp:
            if resp.status != 200:
                # Handle errors (e.g., user already registered, API down)
                resp.raise_for_status()
            data = await resp.json()
            return data.get("url")

    async def submit_drop(self, drop_data: dict):
        """Sends the data collected from the Discord Modal to the Core."""
        endpoint = f"{self.base_url}/api/internal/drops/create"
        # ... similar async POST logic ...
      
