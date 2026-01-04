c500-core-go/
├── go.mod                  # The Go dependencies definition file
├── go.sum                  # Checksums ensuring dependency integrity
├── Dockerfile              # Instructions for Cloud Run deployment
├── main.go                 # The entry point. Wires everything together and starts the server.
│
├── config/                 # Handles environment variables securely
│   └── config.go
│
└── internal/               # The application logic (private to this service)
    ├── domain/             # The core data structures (Structs)
    │   ├── builder.go      # Defines what a "Builder" user is
    │   └── drop.go         # Defines what an Item Listing looks like
    │
    ├── database/           # The Data Access Layer (Talks to Firestore)
    │   └── firestore.go    # Handles reading/writing documents
    │
    ├── service/            # The Business Logic Layer ("The Brain")
    │   ├── builder_service.go # Logic for onboarding, Stripe connection
    │   └── drop_service.go    # Logic for validating and creating drops
    │
    ├── transport/          # The HTTP Layer (Talks to the outside world)
    │   ├── router.go       # Sets up Gin URLs and middleware security
    │   ├── middleware/     # Security checks (e.g., ensuring callers are authorized)
    │   │   └── auth.go
    │   └── handlers/       # The specific API endpoints
    │       ├── builder_handler.go
    │       ├── drop_handler.go
    │       └── webhook_handler.go # Handles incoming Stripe events
    │
    └── integrations/       # Clients for external APIs
        └── stripe/
            └── client.go   # Wrapper around the official Stripe Go SDK
