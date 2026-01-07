package services

import (
	"fmt"

	"c500-web-go/models"
)

// DropRepository defines the interface that the database layer must implement.
// By defining this interface here, the Service says "I don't care *how* you get the data,
// as long as you have these methods." This is crucial for decoupling and testing.
// Currently, models.ProductModel satisfies this interface.
type DropRepository interface {
	All() ([]models.Product, error)
	Get(id int) (*models.Product, error)
	// Future methods: Create(p *models.Product), UpdateInventory(id int, qty int), etc.
}

// DropService defines the publicly available methods for interacting with drops business logic.
// Handlers will talk to this interface.
type DropService interface {
	ListActiveDrops() ([]models.Product, error)
	GetDropDetails(id int) (*models.Product, error)
}

// dropServiceImpl is the concrete implementation of the DropService interface.
// It holds a reference to the repository layer.
type dropServiceImpl struct {
	repo DropRepository
}

// NewDropService is the constructor. It takes a repository as a dependency
// and returns the Service interface.
func NewDropService(repo DropRepository) DropService {
	return &dropServiceImpl{
		repo: repo,
	}
}

// --- Business Logic Methods ---

// ListActiveDrops retrieves all products that are currently available for sale.
func (s *dropServiceImpl) ListActiveDrops() ([]models.Product, error) {
	// Right now, this is a simple pass-through to the repository.
	// BUT, this is where future business logic belongs.

	// Example of future logic:
	// products, err := s.repo.All()
	// if err != nil { return nil, err }
	// activeProducts := []models.Product{}
	// for _, p := range products {
	//    if p.InventoryCount > 0 && p.EndDate.After(time.Now()) {
	//        activeProducts = append(activeProducts, p)
	//    }
	// }
	// return activeProducts, nil

	// Current simple implementation:
	fmt.Println("Service layer: Fetching all active drops...")
	return s.repo.All()
}

// GetDropDetails retrieves a specific drop by ID, potentially adding extra context.
func (s *dropServiceImpl) GetDropDetails(id int) (*models.Product, error) {
	fmt.Printf("Service layer: Fetching details for drop ID %d...\n", id)

	// Business logic check example: Ensure ID is valid before hitting DB
	if id <= 0 {
		// We could define custom service errors here
		return nil, fmt.Errorf("invalid drop ID provided: %d", id)
	}

	// Call repository layer
	product, err := s.repo.Get(id)
	if err != nil {
		// We might want to wrap DB errors with application-specific context
		return nil, fmt.Errorf("service failed to retrieve drop: %w", err)
	}

	return product, nil
}
