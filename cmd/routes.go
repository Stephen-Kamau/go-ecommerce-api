package main

import (
	"log"
	"net/http"

	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/rs/cors"

	"ecomApis/internals/orders"
	"ecomApis/internals/products"
	"ecomApis/internals/repo"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	// cors
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler)

	// timeout context
	r.Use(middleware.Timeout(60 * time.Second))

	// create a healthcheck endpoint
	r.Get("/health", healthCheck)

	// product routes
	productService := products.NewProductService(repo.New(app.db))
	productHandler := products.NewProductHandler(productService)

	r.Route("/products", func(r chi.Router) {
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.ListAllProducts)
		r.Get("/{id}", productHandler.GetProductById)
		r.Delete("/{id}", productHandler.DeleteProduct)

	})

	// order routes
	orderService := orders.NewOrderService(repo.New(app.db), app.db)
	orderHandler := orders.NewOrderHandler(orderService)

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", orderHandler.CreateOrder)
		r.Get("/customer/{customerRef}", orderHandler.GetOrdersByCustomerRef)
		r.Get("/", orderHandler.GetAllOrders)
		r.Get("/{id}", orderHandler.GetOrderByID)
		r.Delete("/{id}", orderHandler.DeleteOrder)
	})

	// other routes...
	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.Address,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      h,
	}
	log.Printf("E-Commerce API Server Starting....")
	log.Printf("Starting server on %s", srv.Addr)

	// start the server
	return srv.ListenAndServe()
}

type application struct {
	config appconfig
	db     *pgx.Conn
}

type appconfig struct {
	Address string
	DB      dbConfig
}
type dbConfig struct {
	DatabaseURL string
}
