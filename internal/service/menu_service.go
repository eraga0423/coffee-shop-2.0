package service

import (
	"errors"
	"strings"

	"frapuccino/internal/dal"
	"frapuccino/models"
)

type MenuService interface {
	ServiceGetMenuItem() ([]models.MenuItem, error)
	ServicePostMenu(content models.MenuItem) error
	ServiceGetMenuID(id int) (models.MenuItem, error)
	ServicePutMenuID(id int, newEdit models.MenuItem) error
	ServiceDelete(id int) error
}

type menuService struct {
	menuRepo dal.MenuRepository
}

// Initializes and returns a new instance of menuService with the provided repository
func NewMenuService(menuRepo dal.MenuRepository) MenuService {
	return &menuService{menuRepo: menuRepo}
}

// Adds new menu items to the menu, checking for duplicates and validating data
func (s *menuService) ServicePostMenu(content models.MenuItem) error {
	err := s.CheckMenu(content)
	if err != nil {
		return err
	}
	err = s.menuRepo.PostRepoMenu(content)
	if err != nil {
		return err
	}
	return nil
}

// Retrieves all menu items from the repository
func (s *menuService) ServiceGetMenuItem() ([]models.MenuItem, error) {
	return s.menuRepo.GetMenuRepo()
}

// Retrieves a specific menu item by ID, returning an error if not found
func (s *menuService) ServiceGetMenuID(id int) (models.MenuItem, error) {
	return s.menuRepo.GetMenuItemID(id)
}

// Updates a specific menu item by ID with new data provided, validating changes
func (s *menuService) ServicePutMenuID(id int, newEdit models.MenuItem) error {
	err := s.CheckMenu(newEdit)
	if err != nil {
		return err
	}
	return s.menuRepo.UpdateMenu(id, newEdit)
}

// Deletes a menu item by ID, returning an error if the ID is not found
func (s *menuService) ServiceDelete(id int) error {
	return s.menuRepo.DeleteMenuItem(id)
}

// Validates the fields of a new menu item to ensure all required fields are filled correctly
func (s *menuService) CheckMenu(newmenu models.MenuItem) error {
	if newmenu.ID == 0 {
		return errors.New("Missing ID")
	}
	newmenuName := strings.Trim(newmenu.Name, " ")
	if newmenuName == "" {
		return errors.New("Missing name")
	}
	newmenuDescription := strings.Trim(newmenu.Description, " ")
	if newmenuDescription == "" {
		return errors.New("Missing Description")
	}
	if newmenu.Price < 0.0 {
		return errors.New("Price cannot be negative")
	}
	for _, msq := range newmenu.Ingredients {
		if msq.Quantity <= 0 {
			return errors.New("Ingredients quantity cannot be 0 or negative")
		}
		if msq.IngredientID == 0 {
			return errors.New("Missing ingredients ID")
		}
	}

	return nil
}
