// In main.go (Mental map for next steps)

// 1. Init Repo Layer (Already done)
productRepo := &models.ProductModel{DB: db}

// 2. Init Service Layer (New step)
// We inject the repo into the service.
dropService := services.NewDropService(productRepo)

// 3. Init Handlers (Update step)
// We inject the SERVICE into the handler, not the repo directly.
storeHandler := handlers.NewStoreHandler(tmpl, dropService)
apiHandler := handlers.NewAPIHandler(dropService)
