curl -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{
  "name": "Example Product",
  "description": "This is a test product",
  "price": 20,
  "stock": 23
}'




{
  "customer_ref": "cust_abc_123",
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    },
    {
      "product_id": 4,
      "quantity": 1
    },
    {
      "product_id": 7,
      "quantity": 3
    }
  ]
}


curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_ref": "cust_abc_123",
    "items": [
      { "product_id": 2, "quantity": 2 },
      { "product_id": 4, "quantity": 1 },
      { "product_id": 3, "quantity": 3 }
    ]
  }'
