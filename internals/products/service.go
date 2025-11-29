package products

import (
	"context"
	"database/sql"
	"ecomApis/internals/repo"
	"ecomApis/internals/utils"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type ProductService struct {
	repo *repo.Queries
}

func NewProductService(r *repo.Queries) *ProductService {
	return &ProductService{
		repo: r,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, arg repo.CreateProductParams) (repo.Product, error) {
	// --- Validation ---
	if arg.Name == "" {
		return repo.Product{}, &utils.ValidationError{
			Field:   "Name",
			Message: "cannot be empty",
		}
	}

	if arg.Price < 0 {
		return repo.Product{}, &utils.ValidationError{
			Field:   "Price",
			Message: "cannot be negative",
		}
	}

	if arg.Stock < 0 {
		return repo.Product{}, &utils.ValidationError{
			Field:   "Stock",
			Message: "cannot be negative",
		}
	}

	// check if product with same name exists
	exists, err := s.repo.ProductExists(ctx, arg.Name)
	if err != nil {
		return repo.Product{}, &utils.DatabaseError{
			Query: "ProductExists",
			Err:   err,
		}
	}
	if exists {
		return repo.Product{}, &utils.AlreadyExistsError{
			Resource: "Product",
			ID:       arg.Name,
		}
	}

	product, err := s.repo.CreateProduct(ctx, repo.CreateProductParams{
		Name:        arg.Name,
		Description: arg.Description,
		Price:       arg.Price,
		Stock:       arg.Stock,
	})

	if err != nil {
		return repo.Product{}, &utils.DatabaseError{
			Query: "CreateProduct",
			Err:   err,
		}
	}

	return product, nil
}

func (s *ProductService) FindProductByID(ctx context.Context, id int64) (repo.Product, error) {
	product, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows || err == pgx.ErrNoRows {
			return repo.Product{}, &utils.NotFoundError{
				Resource: "Product",
				ID:       strconv.FormatInt(id, 10),
			}
		}
		return repo.Product{}, &utils.DatabaseError{
			Query: "FindProductByID",
			Err:   err,
		}
	}

	return product, nil
}

func (s *ProductService) UpdateProductDetails(ctx context.Context, arg repo.UpdateProductDetailsParams) (repo.Product, error) {
	// --- Validation ---
	if arg.Name == "" {
		return repo.Product{}, &utils.ValidationError{
			Field:   "Name",
			Message: "cannot be empty",
		}
	}

	if arg.Price < 0 {
		return repo.Product{}, &utils.ValidationError{
			Field:   "Price",
			Message: "cannot be negative",
		}
	}

	product, err := s.repo.UpdateProductDetails(ctx, repo.UpdateProductDetailsParams{
		Name:        arg.Name,
		Description: arg.Description,
		Price:       arg.Price,
		ID:          arg.ID,
	})

	if err != nil {
		return repo.Product{}, &utils.DatabaseError{
			Query: "UpdateProductDetails",
			Err:   err,
		}
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {

	// check if the product exists
	_, err := s.repo.FindProductByID(ctx, id)
	if err != nil {

		if err == sql.ErrNoRows || err == pgx.ErrNoRows {
			return &utils.NotFoundError{
				Resource: "Product",
				ID:       strconv.FormatInt(id, 10),
			}
		}
		return &utils.DatabaseError{
			Query: "FindProductByID",
			Err:   err,
		}
	}

	// delete the product
	err = s.repo.DeleteProduct(ctx, id)
	if err != nil {
		return &utils.DatabaseError{
			Query: "DeleteProduct",
			Err:   err,
		}
	}
	return nil
}

func (s *ProductService) ListAllProducts(ctx context.Context) ([]repo.Product, error) {
	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		return nil, &utils.DatabaseError{
			Query: "ListAllProducts",
			Err:   err,
		}
	}
	return products, nil
}
