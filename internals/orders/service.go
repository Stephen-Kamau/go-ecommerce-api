package orders

import (
	"context"
	"database/sql"
	"ecomApis/internals/repo"
	"ecomApis/internals/utils"
	"strconv"

	"fmt"

	"github.com/jackc/pgx/v5"
)

type OrderService struct {
	repo *repo.Queries
	db   *pgx.Conn
}

func NewOrderService(r *repo.Queries, db *pgx.Conn) *OrderService {
	return &OrderService{
		repo: r,
		db:   db,
	}
}

// Placing an order process:
// 1. get customer_ref (this is just any string that identifies the customer) and order items (product IDs and quantities)
// 2. calculate total price by fetching product prices from the products table
// 3. create order in orders table
// 4. create order items in order_items table
// 5. update product stock in products table
// We rollback if any step fails

func (s *OrderService) CreateOrder(ctx context.Context, customerRef string, items []OrderItemRequest) (repo.Order, []repo.OrderItem, error) {

	if customerRef == "" {
		return repo.Order{}, nil, &utils.ValidationError{
			Field:   "customer_ref",
			Message: "cannot be empty",
		}
	}

	if len(items) == 0 {
		return repo.Order{}, nil, &utils.ValidationError{
			Field:   "items",
			Message: "cannot be empty",
		}
	}

	// start transaction wth current context
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return repo.Order{}, nil, fmt.Errorf("begin tx: %w", err)
	}
	qtx := s.repo.WithTx(tx)

	// create order
	order, err := qtx.CreateOrder(ctx, customerRef)
	if err != nil {
		tx.Rollback(ctx)
		return repo.Order{}, nil, &utils.DatabaseError{Query: "CreateOrder", Err: err}
	}

	var total int32 = 0
	orderItems := []repo.OrderItem{}

	// each item in the order
	for _, item := range items {

		if item.Quantity <= 0 {
			tx.Rollback(ctx)
			return repo.Order{}, nil, &utils.ValidationError{
				Field:   "quantity",
				Message: fmt.Sprintf("invalid quantity for product %d", item.ProductID),
			}
		}

		// Fetch product
		product, err := qtx.FindProductByID(ctx, item.ProductID)
		if err != nil {
			if err == pgx.ErrNoRows || err == sql.ErrNoRows {
				tx.Rollback(ctx)
				return repo.Order{}, nil, &utils.NotFoundError{
					Resource: "Product",
					ID:       strconv.FormatInt(item.ProductID, 10),
				}
			}
			tx.Rollback(ctx)
			return repo.Order{}, nil, &utils.DatabaseError{Query: "FindProductByID", Err: err}
		}
		if product.Price <= 0 {
			tx.Rollback(ctx)
			return repo.Order{}, nil, &utils.ValidationError{
				Field:   "price",
				Message: fmt.Sprintf("invalid price for product %d", item.ProductID),
			}
		}
		// Check stock
		if product.Stock < item.Quantity {
			tx.Rollback(ctx)
			return repo.Order{}, nil, &utils.ValidationError{
				Field:   "stock",
				Message: fmt.Sprintf("not enough stock for product %d", item.ProductID),
			}
		}

		// Decrement stock

		_, err = qtx.UpdateProductStock(ctx, repo.UpdateProductStockParams{
			Stock: item.Quantity,
			ID:    product.ID,
		})
		if err != nil {
			tx.Rollback(ctx)
			return repo.Order{}, nil, &utils.DatabaseError{Query: "UpdateProductStock", Err: err}
		}

		// add items to order_items table
		oi, err := qtx.AddOrderItem(ctx, repo.AddOrderItemParams{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  item.Quantity,
			UnitPrice: product.Price,
		})
		if err != nil {
			tx.Rollback(ctx)
			return repo.Order{}, nil, &utils.DatabaseError{Query: "AddOrderItem", Err: err}
		}

		orderItems = append(orderItems, oi)

		// Accumulate total
		total += product.Price * item.Quantity
	}
	// Update order total
	_, err = tx.Exec(ctx, "UPDATE orders SET total_price = $1 WHERE id = $2", total, order.ID)
	if err != nil {
		tx.Rollback(ctx)
		return repo.Order{}, nil, &utils.DatabaseError{Query: "UpdateOrderTotal", Err: err}
	}

	// commit transaction
	if err := tx.Commit(ctx); err != nil {
		return repo.Order{}, nil, fmt.Errorf("commit tx: %w", err)
	}

	order.TotalPrice = total
	return order, orderItems, nil
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]repo.Order, error) {
	orders, err := s.repo.GetAllOrders(ctx)
	if err != nil {
		return nil, &utils.DatabaseError{
			Query: "GetAllOrders",
			Err:   err,
		}
	}
	return orders, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id int64) (repo.Order, []repo.OrderItem, error) {
	order, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows || err == sql.ErrNoRows {
			return repo.Order{}, nil, &utils.NotFoundError{
				Resource: "Order",
				ID:       strconv.FormatInt(id, 10),
			}
		}
		return repo.Order{}, nil, &utils.DatabaseError{
			Query: "GetOrder",
			Err:   err,
		}
	}

	// get order items
	items, err := s.repo.ListOrderItems(ctx, order.ID)
	if err != nil {
		return repo.Order{}, nil, &utils.DatabaseError{
			Query: "ListOrderItems",
			Err:   err,
		}
	}

	return order, items, nil
}

func (s *OrderService) GetOrdersByCustomerRef(ctx context.Context, customerRef string) ([]repo.Order, error) {
	orders, err := s.repo.GetOrdersByCustomerRef(ctx, customerRef)
	if err != nil {
		if err == pgx.ErrNoRows || err == sql.ErrNoRows {
			return nil, &utils.NotFoundError{
				Resource: "Orders for CustomerRef",
				ID:       customerRef,
			}
		}
		return nil, &utils.DatabaseError{
			Query: "GetOrdersByCustomerRef",
			Err:   err,
		}
	}
	return orders, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, id int64) error {

	// check if the order exists
	_, err := s.repo.GetOrder(ctx, id)
	if err != nil {

		if err == sql.ErrNoRows || err == pgx.ErrNoRows {
			return &utils.NotFoundError{
				Resource: "Order",
				ID:       strconv.FormatInt(id, 10),
			}
		}
		return &utils.DatabaseError{
			Query: "GetOrder",
			Err:   err,
		}
	}

	// delete the order items
	err = s.repo.DeleteOrderItemsByOrderID(ctx, id)
	if err != nil {
		return &utils.DatabaseError{
			Query: "DeleteOrderItemsByOrderID",
			Err:   err,
		}
	}

	// delete the order
	err = s.repo.DeleteOrder(ctx, id)
	if err != nil {
		return &utils.DatabaseError{
			Query: "DeleteOrder",
			Err:   err,
		}
	}
	return nil
}
