package products

import (
	"ecomApis/internals/repo"
	"ecomApis/internals/utils"
	"fmt"
	"net/http"

	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	service *ProductService
}

func NewProductHandler(s *ProductService) *ProductHandler {
	return &ProductHandler{
		service: s,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req repo.CreateProductParams

	err := utils.ParseJSON(r.Body, &req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	product, err := h.service.CreateProduct(ctx, req)
	if err != nil {
		if ve, ok := err.(*utils.ValidationError); ok {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": ve.Error()})
			return
		}
		if de, ok := err.(*utils.DatabaseError); ok {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": de.Error()})
			return
		}
		if ae, ok := err.(*utils.AlreadyExistsError); ok {
			utils.WriteJSON(w, http.StatusConflict, map[string]string{"error": ae.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) ListAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := h.service.ListAllProducts(ctx)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get product id from url params
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid product id"})
		return
	}

	product, err := h.service.FindProductByID(ctx, id)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid product id"})
		return
	}

	err = h.service.DeleteProduct(ctx, id)
	if err != nil {
		switch e := err.(type) {
		case *utils.NotFoundError:
			utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": e.Error()})
		case *utils.DatabaseError:
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": e.Error()})
		default:
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("product with id %d deleted", id),
	})
}
