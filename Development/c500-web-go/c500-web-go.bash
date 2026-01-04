c500-web-go/
├── go.mod                  # Dependencies
├── Dockerfile              # Deployment instructions
├── main.go                 # The entry point. Sets up Gin, loads templates, starts server.
│
├── assets/                 # Public Static Files (served directly)
│   ├── css/
│   │   └── main.css        # The base C500 theme CSS
│   └── images/
│       └── logo.png
│
├── templates/              # Server-Side HTML Templates
│   ├── layouts/            # Base templates holding the common structure
│   │   └── base.html       # Contains <html>, <head>, nav, footer
│   ├── pages/              # Individual page content
│   │   ├── home.html       # Landing page content
│   │   ├── success.html    # "Thanks for your order!"
│   │   └── profile.html    # The complex one: renders custom builder content
│   └── partials/           # Reusable components (e.g., navbar, footer)
│       ├── nav.html
│       └── footer.html
│
└── internal/
    ├── config/             # Env vars (PORT, CORE_API_URL, INTERNAL_API_KEY)
    │   └── config.go
    ├── clients/            # Communicates with c500-core-go
    │   └── core_client.go  # HTTP client to fetch builder data
    └── handlers/           # The controllers that render HTML
        ├── static_handlers.go # Home, Success, Cancel pages
        └── profile_handler.go # The complex handler for builder pages
