// ... (previous code for users and drops remains above)

// Uncomment this constant now
const (
	// ...
	ordersCollection = "orders"
)

// =================================================================
// OrderRepository Implementation
// These methods fulfill the interface defined in checkout_service.go
// and future interfaces needed for fulfillment logic.
// =================================================================

// CreateOrder is called by the Webhook Handler (via the Service layer)
// the moment Stripe tells us a payment succeeded.
func (f *FirestoreClient) CreateOrder(ctx context.Context, order *domain.Order) error {
	// We use .Create to ensure we don't overwrite an existing order by accident,
	// though UUID collisions should be impossible.
	_, err := f.client.Collection(ordersCollection).Doc(order.ID).Create(ctx, order)
	if err != nil {
		return fmt.Errorf("firestore create order error: %w", err)
	}
	return nil
}

// GetOrderByID is needed for the Fulfillment flow.
// Before a seller can ship an item, we need to fetch the order to verify
// it exists and belongs to them.
func (f *FirestoreClient) GetOrderByID(ctx context.Context, orderID string) (*domain.Order, error) {
	docRef := f.client.Collection(ordersCollection).Doc(orderID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Use a generic service error so layers above don't need gRPC imports
			// You'd need to define ErrOrderNotFound in service/errors.go
			return nil, fmt.Errorf("order not found: %w", service.ErrOrderNotFound)
		}
		return nil, fmt.Errorf("firestore get order error: %w", err)
	}

	var order domain.Order
	if err := docSnap.DataTo(&order); err != nil {
		return nil, fmt.Errorf("failed to map data to order struct: %w", err)
	}
	return &order, nil
}

// UpdateOrderFulfillment is the critical step where funds are released.
// It's called when a seller provides valid tracking or a VOD link.
func (f *FirestoreClient) UpdateOrderFulfillment(ctx context.Context, orderID string, updates map[string]interface{}) error {
	// 'updates' is a map of fields to change, e.g.:
	// {
	//   "tracking_number": "1Z999...",
	//   "carrier": "UPS",
	//   "escrow_status": "released",
	//   "updated_at": time.Now()
	// }

	docRef := f.client.Collection(ordersCollection).Doc(orderID)

	// Convert the map into Firestore Update objects
	fsUpdates := make([]firestore.Update, 0, len(updates))
	for field, value := range updates {
		fsUpdates = append(fsUpdates, firestore.Update{
			Path:  field,
			Value: value,
		})
	}

	// Perform the atomic partial update.
	_, err := docRef.Update(ctx, fsUpdates)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return service.ErrOrderNotFound
		}
		return fmt.Errorf("firestore update order fulfillment error: %w", err)
	}

	return nil
}
