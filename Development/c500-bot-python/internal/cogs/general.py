# internal/cogs/general.py

import discord
from discord.ext import commands
from datetime import datetime

class General(commands.Cog):
    """
    General utility commands for the C500 Bot.
    Includes ping, uptime, server info, and user info.
    """

    def __init__(self, bot):
        self.bot = bot
        # Record the time the cog was loaded to calculate uptime.
        # In a robust bot, this might be stored in the main bot instance instead.
        self.start_time = datetime.utcnow()

    # ===========================
    # HELPER FUNCTIONS
    # ===========================
    def _format_timedelta(self, delta):
        """Formats a timedelta object into a readable string."""
        days, seconds = delta.days, delta.seconds
        hours = seconds // 3600
        minutes = (seconds % 3600) // 60
        seconds = seconds % 60
        return f"{days}d {hours}h {minutes}m {seconds}s"

    # ===========================
    # COMMANDS
    # ===========================

    @commands.command(aliases=['lat'])
    async def ping(self, ctx):
        """
        Checks the bot's responsiveness and API latency.
        Usage: !ping
        """
        # Calculate API latency by measuring message send time
        start_time = datetime.utcnow()
        message = await ctx.send("🏓 Pinging...")
        end_time = datetime.utcnow()
        api_latency = (end_time - start_time).total_seconds() * 1000

        # Get Discord WebSocket latency
        ws_latency = round(self.bot.latency * 1000)

        await message.edit(content=f"🏓 **Pong!**\nAPI Latency: `{api_latency:.0f}ms`\nWebSocket Latency: `{ws_latency}ms`")

    @commands.command(aliases=['up'])
    async def uptime(self, ctx):
        """
        Shows how long the bot's current session has been running.
        Usage: !uptime
        """
        now = datetime.utcnow()
        delta = now - self.start_time
        uptime_str = self._format_timedelta(delta)
        await ctx.send(f"⏱️ **Uptime:** {uptime_str}")

    @commands.command(aliases=['info', 'botinfo'])
    async def about(self, ctx):
        """
        Displays information about the C500 Python Bot.
        Usage: !about
        """
        embed = discord.Embed(
            title="🤖 About C500 Python Bot",
            description="A custom Python bot built using discord.py for fulfillment and community management.",
            color=discord.Color.blurple(),
            timestamp=datetime.utcnow()
        )
        
        # You can customize these fields with real data
        embed.add_field(name="Developer", value="[Your Name/Team]", inline=True)
        embed.add_field(name="Library", value=f"discord.py v{discord.__version__}", inline=True)
        embed.add_field(name="Python Version", value=platform.python_version(), inline=True)
        
        # Calculate uptime for the embed
        delta = datetime.utcnow() - self.start_time
        uptime_str = self._format_timedelta(delta)
        embed.add_field(name="Uptime", value=uptime_str, inline=True)

        embed.set_thumbnail(url=self.bot.user.avatar.url if self.bot.user.avatar else None)
        embed.set_footer(text=f"Requested by {ctx.author.display_name}")
        
        await ctx.send(embed=embed)

    @commands.guild_only()
    @commands.command(aliases=['server'])
    async def serverinfo(self, ctx):
        """
        Displays detailed information about the current server.
        Usage: !serverinfo
        """
        guild = ctx.guild
        
        embed = discord.Embed(title=f"ℹ️ Server Info: {guild.name}", color=discord.Color.gold())
        
        if guild.icon:
            embed.set_thumbnail(url=guild.icon.url)

        embed.add_field(name="Server ID", value=guild.id, inline=True)
        embed.add_field(name="Owner", value=guild.owner.mention, inline=True)
        embed.add_field(name="Members", value=str(guild.member_count), inline=True)
        
        # Count roles, text channels, and voice channels
        roles_count = len(guild.roles)
        text_channels = len(guild.text_channels)
        voice_channels = len(guild.voice_channels)

        embed.add_field(name="Roles", value=str(roles_count), inline=True)
        embed.add_field(name="Channels", value=f"📝 {text_channels} Text | 🔊 {voice_channels} Voice", inline=True)
        
        created_at = guild.created_at.strftime("%Y-%m-%d")
        embed.set_footer(text=f"Server Created on {created_at}")

        await ctx.send(embed=embed)

    @commands.guild_only()
    @commands.command(aliases=['whois', 'user'])
    async def userinfo(self, ctx, member: discord.Member = None):
        """
        Displays information about a specific user.
        If no user is specified, shows info about the command caller.
        Usage: !userinfo [@user]
        """
        member = member or ctx.author

        embed = discord.Embed(title=f"👤 User Info: {member.display_name}", color=member.color)
        
        if member.avatar:
            embed.set_thumbnail(url=member.avatar.url)

        embed.add_field(name="User ID", value=member.id, inline=True)
        embed.add_field(name="Username", value=str(member), inline=True)
        
        # Format dates
        joined_at = member.joined_at.strftime("%Y-%m-%d %H:%M:%S") if member.joined_at else "Unknown"
        created_at = member.created_at.strftime("%Y-%m-%d %H:%M:%S")
        
        embed.add_field(name="Joined Server", value=joined_at, inline=False)
        embed.add_field(name="Account Created", value=created_at, inline=False)
        
        # List key roles (excluding @everyone)
        roles = [role.mention for role in member.roles if role.name != "@everyone"]
        # Truncate if too many roles
        if len(roles) > 10:
            roles_display = ", ".join(roles[:10]) + f" and {len(roles)-10} more..."
        else:
            roles_display = ", ".join(roles) if roles else "None"

        embed.add_field(name=f"Roles [{len(roles)}]", value=roles_display, inline=False)
        
        await ctx.send(embed=embed)

# Need to import platform for the about command
import platform

async def setup(bot):
    await bot.add_cog(General(bot))
  
