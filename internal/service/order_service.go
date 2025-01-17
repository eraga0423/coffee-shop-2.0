package service

import (
	"errors"
	"strings"

	"frapuccino/internal/dal/orderRepo"
	"frapuccino/models"
)

type OrderService interface {
	ServicePostOrders(body models.Order) error
	ServicePutOrderID(id int, newEdit models.Order) error
	CloseOrder(id int) error
	ServiceDeleteOrdersID(id int) error
	GetOrdersService() ([]models.Order, error)
	GetIDOrdersService(id int) (models.Order, error)
	CheckBodyOrder(body models.Order) error
}

type orderService struct {
	orderRepo orderRepo.OrderRepository
}

var Id int

// Initializes and returns a new instance of orderService with the provided repository
func NewOrderService(orderRepo orderRepo.OrderRepository) OrderService {
	return &orderService{orderRepo: orderRepo}
}

// Creates a new order, validates the order details, and ensures no open orders exist
func (s orderService) ServicePostOrders(body models.Order) error {
	if err := s.CheckBodyOrder(body); err != nil {
		return err
	}
	body.Status = "open"
	if err := s.orderRepo.WriteDBNewOrder(body); err != nil {
		return err
	}
	return nil
}

// Validates the fields of an order to ensure all required information is present
func (s orderService) CheckBodyOrder(body models.Order) error {
	newbodyCustomer := strings.Trim(body.CustomerName, " ")
	if newbodyCustomer == "" {
		return errors.New("Missing customer name")
	}
	if body.Items == nil {
		return errors.New("Missing items in menu")
	}
	for _, item := range body.Items {
		if item.ProductID == 0 {
			return errors.New("Missing product id")
		}
		if item.Quantity < 1 {
			return errors.New("Quantity cannot be negative")
		}

	}
	return nil
}

// Updates an existing order by ID, ensuring it is still open and validating the new data
func (s *orderService) ServicePutOrderID(id int, body models.Order) error {
	if err := s.CheckBodyOrder(body); err != nil {
		return err
	}
	if err := s.orderRepo.UpdateOrder(id, body); err != nil {
		return err
	}
	return nil
}

// Closes an open order by ID, updates inventory quantities, and writes changes
func (s *orderService) CloseOrder(id int) error {
	err := s.orderRepo.OrderClose(id)
	if err != nil {
		return err
	}
	return nil
}

// Deletes an order by ID, returning an error if the ID is not found
func (s *orderService) ServiceDeleteOrdersID(id int) error {
	if err := s.orderRepo.DeleteOrder(id); err != nil {
		return err
	}
	return nil
}

// Retrieves all orders from the repository
func (s *orderService) GetOrdersService() ([]models.Order, error) {
	return s.orderRepo.ParseOrders()
}

// Retrieves a specific order by ID, returning an error if the ID is not found
func (s *orderService) GetIDOrdersService(id int) (models.Order, error) {
	return s.orderRepo.GetRepoId(id)
}
