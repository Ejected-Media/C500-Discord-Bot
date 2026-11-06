import discord
from discord.ext import commands
from discord import app_commands

class MentorMatching(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    # --- Prefix Command: !findmentor ---
    @commands.command(name="findmentor", help="Find available C500 mentors or get a link to the booking system.")
    async def findmentor_prefix(self, ctx, *, topic: str = None):
        embed = self._build_mentor_embed(topic)
        view = MentorBookingView()
        await ctx.send(embed=embed, view=view)

    # --- Slash Command: /findmentor ---
    @app_commands.command(name="findmentor", description="Find available C500 mentors for your tech stack or session booking.")
    @app_commands.describe(topic="The specific technology or topic (e.g., CSS Subgrid, Rust, UI Layout)")
    async def findmentor_slash(self, interaction: discord.Interaction, topic: str = None):
        embed = self._build_mentor_embed(topic)
        view = MentorBookingView()
        await interaction.response.send_message(embed=embed, view=view, ephemeral=False)

    def _build_mentor_embed(self, topic: str = None) -> discord.Embed:
        embed = discord.Embed(
            title="🎓 C500 Mentor Matching (Initial Stage)",
            description="Connect with live technical mentors or schedule a dedicated 1-on-1 pairing session on the C500 platform.",
            color=discord.Color.from_rgb(114, 137, 218) # Discord Blurple / theme aligned
        )
        
        if topic:
            embed.add_field(
                name=f"🔍 Filtered Topic: `{topic}`",
                value="Showing active mentors matching your search query. To book a verified time slot, use the portal link below.",
                inline=False
            )
            # Simulated matching result
            embed.add_field(
                name="🟢 Available Now",
                value="• **@glox_mentor_01** — Native CSS Subgrid & Container Queries\n• **@glox_mentor_02** — Rust CLI & WebAssembly Tooling",
                inline=False
            )
        else:
            embed.add_field(
                name="💡 Quick Tip",
                value="Run `!findmentor <topic>` to narrow down mentors by subject (e.g., `!findmentor CSS Subgrid`).",
                inline=False
            )
            embed.add_field(
                name="🟢 Currently Online Mentors",
                value="• **@glox_mentor_01** (Frontend / C500 Workspace)\n• **@glox_mentor_02** (Systems / Rust / WASM)\n• **@glox_mentor_03** (UI Design & Accessibility)",
                inline=False
            )

        embed.set_footer(text="Classroom 500 Mentorship Desk • Powered by GloxBot")
        return embed

class MentorBookingView(discord.ui.View):
    def __init__(self):
        super().__init__(timeout=None)
        # Persistent button linking to the C500 platform mentor booking system
        self.add_item(discord.ui.Button(
            label="Open C500 Mentor Booking Portal",
            style=discord.ButtonStyle.link,
            url="https://github.com/C500Students", # Update to your production booking URL endpoint
            emoji="🔗"
        ))

async def setup(bot):
    await bot.add_cog(MentorMatching(bot))
