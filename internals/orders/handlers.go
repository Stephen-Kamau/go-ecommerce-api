package orders

import (
	"ecomApis/internals/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	service *OrderService
}

func NewOrderHandler(s *OrderService) *OrderHandler {
	return &OrderHandler{
		service: s,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateOrderRequest
	err := utils.ParseJSON(r.Body, &req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	order, items, err := h.service.CreateOrder(ctx, req.CustomerRef, req.Items)
	if err != nil {
		if ve, ok := err.(*utils.ValidationError); ok {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": ve.Error()})
			return
		}
		if de, ok := err.(*utils.DatabaseError); ok {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": de.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"order":       order,
		"order_items": items,
	}
	utils.WriteJSON(w, http.StatusCreated, response)
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orders, err := h.service.GetAllOrders(ctx)
	if err != nil {
		if de, ok := err.(*utils.DatabaseError); ok {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": de.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid order ID"})
		return
	}

	order, order_items, err := h.service.GetOrder(ctx, id)
	if err != nil {
		if ne, ok := err.(*utils.NotFoundError); ok {
			utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": ne.Error()})
			return
		}
		if de, ok := err.(*utils.DatabaseError); ok {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": de.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	respose := map[string]interface{}{
		"order":       order,
		"order_items": order_items,
	}

	utils.WriteJSON(w, http.StatusOK, respose)
}

func (h *OrderHandler) GetOrdersByCustomerRef(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	customerRef := chi.URLParam(r, "customerRef")
	if customerRef == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "customerRef cannot be empty"})
		return
	}

	orders, err := h.service.GetOrdersByCustomerRef(ctx, customerRef)
	if err != nil {
		if ne, ok := err.(*utils.NotFoundError); ok {
			utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": ne.Error()})
			return
		}
		if de, ok := err.(*utils.DatabaseError); ok {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": de.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid order ID"})
		return
	}

	err = h.service.DeleteOrder(ctx, id)
	if err != nil {
		if ne, ok := err.(*utils.NotFoundError); ok {
			utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": ne.Error()})
			return
		}
		if de, ok := err.(*utils.DatabaseError); ok {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": de.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
