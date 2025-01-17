package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"frapuccino/internal/service"
	"frapuccino/models"
)

type MenuHandler interface {
	GetMenu(w http.ResponseWriter, r *http.Request)
	PostMenu(w http.ResponseWriter, r *http.Request)
	GetMenuID(w http.ResponseWriter, r *http.Request)
	PutMenuID(w http.ResponseWriter, r *http.Request)
	DeleteMenuID(w http.ResponseWriter, r *http.Request)
}

type menuHandler struct {
	menuService service.MenuService
}

// Initializes and returns a new instance of menuHandler with the provided service
func NewMenuHandler(menuService service.MenuService) MenuHandler {
	return &menuHandler{menuService: menuService}
}

// Handles the HTTP request to add a new menu item, validating input and returning success or error
func (h *menuHandler) PostMenu(w http.ResponseWriter, r *http.Request) {
	if err := CheckContentType(r); err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	newMenu := models.MenuItem{}
	err := json.NewDecoder(r.Body).Decode(&newMenu)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	err = h.menuService.ServicePostMenu(newMenu)
	if err != nil {
		SendError(w, http.StatusNotFound, err)
		return
	}
	SendSucces(w, http.StatusCreated, "Menu item added")
}

// Handles the HTTP request to retrieve all menu items and returns them as JSON
func (h *menuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	content, err := h.menuService.ServiceGetMenuItem()
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(content)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
}

// Handles the HTTP request to retrieve a specific menu item by ID and returns it as JSON
func (h *menuHandler) GetMenuID(w http.ResponseWriter, r *http.Request) {
	r.URL.Query().Get("")
	id, err := strconv.Atoi(r.PathValue("id"))
	w.Header().Set("Content-Type", "application/json")

	newGetMenuID, err := h.menuService.ServiceGetMenuID(id)
	if err != nil {
		SendError(w, http.StatusConflict, err)
		return
	}
	err = json.NewEncoder(w).Encode(newGetMenuID)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
}

// Handles the HTTP request to update a specific menu item by ID, validating input and updating the item
func (h *menuHandler) PutMenuID(w http.ResponseWriter, r *http.Request) {
	if err := CheckContentType(r); err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	var newEdit models.MenuItem
	err := json.NewDecoder(r.Body).Decode(&newEdit)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		SendError(w, http.StatusNotFound, err)
		return
	}
	err = h.menuService.ServicePutMenuID(id, newEdit)
	if err != nil {
		SendError(w, http.StatusNotFound, err)
		return
	}
	log.Println("PUT menu ID method created")
	SendSucces(w, http.StatusOK, "Menu updated")
}

// Handles the HTTP request to delete a specific menu item by ID, removing it from the menu
func (h *menuHandler) DeleteMenuID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	err = h.menuService.ServiceDelete(id)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	SendSucces(w, http.StatusNoContent, "Menu item deleted")
}
