package service

import (
	"errors"
	"strings"

	"frapuccino/internal/dal"
	"frapuccino/models"
)

// InventoryService defines methods for handling inventory operations.
type InventoryService interface {
	ServiceGetInvItem() ([]models.InventoryItem, error)         // Retrieves all inventory items.
	ServicePostInv(content models.InventoryItem) error          // Adds new inventory items.
	ServiceGetInvID(id int) (models.InventoryItem, error)       // Retrieves a single inventory item by ID.
	ServicePutInvID(id int, newEdit models.InventoryItem) error // Updates an existing inventory item by ID.
	// EditInvStructure(EditableStructure models.InventoryItem, newEdit models.InventoryItem) (models.InventoryItem, error) // Edits specific fields of an inventory item.
	ServiceInvDelete(id int) error // Deletes an inventory item by ID.
}

// invService implements the InventoryService interface using InventoryRepository.
type invService struct {
	invRepo dal.InventoryRepository
}

// NewInvService creates and returns a new instance of invService.
func NewInvService(invRepo dal.InventoryRepository) InventoryService {
	return &invService{invRepo: invRepo}
}

// ServicePostInv adds new inventory items to the inventory if they pass validation and don't already exist.
func (s *invService) ServicePostInv(content models.InventoryItem) error {
	if check, err := s.CheckInvPost(content); !check {
		return err // Return error if the new item fails validation.
	}

	exists, err := s.invRepo.CheckIfNameExists(content.Name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Such Item already exists")
	}

	err = s.invRepo.AddItems(content) // Save the updated inventory.
	if err != nil {
		return err
	}

	return nil
}

// ServiceGetInvItem retrieves all inventory items from storage.
func (s *invService) ServiceGetInvItem() ([]models.InventoryItem, error) {
	return s.invRepo.ReadJSONInv()
}

func (s *invService) ServiceGetInvID(id int) (models.InventoryItem, error) {
	checker := false
	newGetInvID := models.InventoryItem{}
	jsonfileinv, err := s.invRepo.ReadJSONInv()
	if err != nil {
		return newGetInvID, err
	}
	for _, value := range jsonfileinv {
		if value.IngredientID == id {
			checker = true
			newGetInvID = value
		}
	}
	if !checker {
		return newGetInvID, errors.New("ID not found")
	}
	return newGetInvID, nil
}

// ServicePutInvID updates an existing inventory item identified by ID with new data.
func (s *invService) ServicePutInvID(id int, newEdit models.InventoryItem) error {
	exists, err := s.invRepo.CheckIfExists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Such ID doesn't exist")
	}

	if err := s.invRepo.UpdateItem(id, newEdit); err != nil {
		return err
	}
	return nil
}

func (s *invService) ServiceInvDelete(id int) error {
	exists, err := s.invRepo.CheckIfExists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Such ID doesn't exist")
	}
	s.invRepo.DeleteItem(id)
	return nil
}

func (r *invService) CheckInvPost(newinv models.InventoryItem) (bool, error) {
	newInvName := strings.TrimSpace(newinv.Name)
	if newInvName == "" {
		return false, errors.New("Missing name")
	}
	if newinv.Quantity < 0 {
		return false, errors.New("Quantity cannot be negative")
	}
	newInvUnit := strings.TrimSpace(newinv.Unit)
	if newInvUnit == "" {
		return false, errors.New("Missing Unit")
	}

	return true, nil
}
