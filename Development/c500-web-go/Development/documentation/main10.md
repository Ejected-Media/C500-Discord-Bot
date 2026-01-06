// ... imports and main() setup ...

	// ... after connecting to DB ...
	log.Println("Successfully connected to PostgreSQL database!")

	// --- Initialize Data Models ---
	productModel := &models.ProductModel{DB: db}
    // NEW: Init OrderModel
	orderModel := &models.OrderModel{DB: db}

	// ... parse templates ...

	// --- Initialize Handlers ---
    // ... other handlers ...
	orderHandler := handlers.NewOrderHandler(tmpl)

	// CHANGE: Update Webhook handler constructor to pass orderModel
	webhookHandler := handlers.NewWebhookHandler(productModel, orderModel)

	// ... router setup and server start remain the same ...
  
