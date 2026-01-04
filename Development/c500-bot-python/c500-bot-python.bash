c500-bot-python/
├── .env                  # Secrets (Local dev only. NOT committed to Git)
├── .gitignore            # Ignores .env, __pycache__, venv/
├── Dockerfile            # Instructions for Cloud Build to package the bot
├── requirements.txt      # Python dependencies list
├── config.py             # Loads environment variables securely
├── main.py               # The entry point. Starts the bot, loads extensions.
│
├── cogs/                 # The command modules (the features)
│   ├── __init__.py
│   ├── seller_commands.py  # /c500 setup, /c500 drop
│   └── fulfillment.py      # /c500 live, /c500 ship
│
├── ui/                   # Visual components (Modals, Buttons, Embeds)
│   ├── __init__.py
│   ├── modals.py         # The "Create Drop" pop-up form
│   ├── views.py          # Persistent buttons like "Buy Now"
│   └── embeds.py         # Central file for defining the "Cozy Pastel" aesthetic
│
└── services/             # External communication layer
    ├── __init__.py
    └── core_api_client.py # The Async HTTP client that talks to c500-core-go
