# internal/cogs/admin.py

import discord
from discord.ext import commands
import io
import textwrap
import traceback
import contextlib

class Admin(commands.Cog):
    """
    High-level administrative commands for managing the bot instance.
    CRITICAL: These commands are restricted to the BOT OWNER only.
    """

    def __init__(self, bot):
        self.bot = bot
        self._last_result = None # Used for the eval command history

    async def cog_check(self, ctx):
        """
        This check runs before ANY command in this cog.
        It ensures only the bot application owner can run these commands.
        """
        is_owner = await self.bot.is_owner(ctx.author)
        if not is_owner:
            raise commands.NotOwner("You must be the bot owner to use this cog.")
        return True

    # ===========================
    # BOT MANAGEMENT
    # ===========================

    @commands.command(aliases=['logout', 'die'])
    async def shutdown(self, ctx):
        """Safely shuts down the bot connection."""
        await ctx.send("👋 Shutting down agent processes. Goodbye.")
        await self.bot.close()

    @commands.command(name='setstatus')
    async def set_status(self, ctx, activity_type: str, *, message: str):
        """
        Changes the bot's presence.
        Usage: !setstatus <playing|watching|listening> <message>
        Example: !setstatus watching over orders
        """
        activity_type = activity_type.lower()
        
        if activity_type == 'playing':
            activity = discord.Game(name=message)
        elif activity_type == 'watching':
            activity = discord.Activity(type=discord.ActivityType.watching, name=message)
        elif activity_type == 'listening':
            activity = discord.Activity(type=discord.ActivityType.listening, name=message)
        else:
            await ctx.send("❌ Invalid activity type. Use: `playing`, `watching`, or `listening`.")
            return

        await self.bot.change_presence(activity=activity)
        await ctx.send(f"✅ Status updated to: **{activity_type.capitalize()} {message}**")

    @commands.command()
    async def say(self, ctx, channel: discord.TextChannel, *, message: str):
        """
        Makes the bot send a message to a specific channel.
        Usage: !say #channel-name Hello world
        """
        try:
            await channel.send(message)
            await ctx.message.add_reaction("✅")
        except discord.Forbidden:
            await ctx.send(f"❌ I don't have permission to speak in {channel.mention}.")
        except Exception as e:
            await ctx.send(f"❌ Error: {e}")

    # ===========================
    # COG (EXTENSION) MANAGEMENT
    # ===========================
    # These allow you to update code without restarting the bot.
    # Paths must be dot-separated, e.g., "internal.cogs.general"

    @commands.command()
    async def load(self, ctx, extension_path: str):
        """Loads a new extension."""
        try:
            await self.bot.load_extension(extension_path)
            await ctx.send(f"✅ Loaded extension: `{extension_path}`")
        except Exception as e:
            await ctx.send(f"❌ Failed to load `{extension_path}`\n```py\n{traceback.format_exc()}```")

    @commands.command()
    async def unload(self, ctx, extension_path: str):
        """Unloads an existing extension."""
        # Prevent unloading the admin cog itself, lest you lock yourself out.
        if "admin" in extension_path:
            await ctx.send("⚠️ Request denied. Cannot unload the Admin cog.")
            return

        try:
            await self.bot.unload_extension(extension_path)
            await ctx.send(f"✅ Unloaded extension: `{extension_path}`")
        except Exception as e:
            await ctx.send(f"❌ Failed to unload `{extension_path}`\n```py\n{traceback.format_exc()}```")

    @commands.command(aliases=['rl'])
    async def reload(self, ctx, extension_path: str):
        """Reloads an extension (useful after editing code)."""
        try:
            await self.bot.reload_extension(extension_path)
            await ctx.send(f"🔄 Reloaded extension: `{extension_path}`")
        except Exception as e:
            await ctx.send(f"❌ Failed to reload `{extension_path}`\n```py\n{traceback.format_exc()}```")

    # ===========================
    # DEBUGGING / EVAL
    # ===========================

    def _cleanup_code(self, content):
        """Automatically removes code blocks from discord messages."""
        # remove ```py\n```
        if content.startswith('```') and content.endswith('```'):
            return '\n'.join(content.split('\n')[1:-1])
        # remove `foo`
        return content.strip('` \n')

    @commands.command(name='eval', aliases=['debug', 'run', 'exec'])
    async def _eval(self, ctx, *, body: str):
        """
        Evaluates arbitrary Python code.
        EXTREMELY DANGEROUS. Only the bot owner can use this.
        Environment variables available: ctx, bot, channel, author, guild, message, _ (last result)
        """
        
        env = {
            'bot': self.bot,
            'ctx': ctx,
            'channel': ctx.channel,
            'author': ctx.author,
            'guild': ctx.guild,
            'message': ctx.message,
            '_': self._last_result
        }

        env.update(globals())

        body = self._cleanup_code(body)
        stdout = io.StringIO()

        to_compile = f'async def func():\n{textwrap.indent(body, "  ")}'

        try:
            exec(to_compile, env)
        except Exception as e:
            return await ctx.send(f'```py\n{e.__class__.__name__}: {e}\n```')

        func = env['func']
        try:
            with contextlib.redirect_stdout(stdout):
                ret = await func()
        except Exception:
            value = stdout.getvalue()
            await ctx.send(f'```py\n{value}{traceback.format_exc()}\n```')
        else:
            value = stdout.getvalue()
            try:
                await ctx.message.add_reaction('✅')
            except:
                pass

            if ret is None:
                if value:
                    await ctx.send(f'```py\n{value}\n```')
            else:
                self._last_result = ret
                await ctx.send(f'```py\n{value}{ret}\n```')

async def setup(bot):
    await bot.add_cog(Admin(bot))
      
