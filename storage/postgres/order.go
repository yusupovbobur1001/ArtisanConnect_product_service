package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	pb "product_service/genproto/orders"
	"time"

	"github.com/google/uuid"
)

type OrderRepo struct {
	Db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{Db: db}
}

func (o *OrderRepo) CreateOrder(req *pb.CreateOrderRequest) (*pb.CreateOrderRequest, error) {
	id := uuid.New().String()

	var totalAmount float32
	items := req.Items
	for _, item := range items {
		totalAmount += float32(item.Quantity) * item.Price
	}

	status := "Pending"

	shippingAddressJSON, err := json.Marshal(req.ShippingAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal shipping address: %v", err)
	}

	query := `
	  INSERT INTO orders(
		id, user_id, total_amount, status, shipping_address, created_at, updated_at
	  )VALUES (
		$1, $2, $3, $4, $5, $6, $7)
	`

	_, err = o.Db.Exec(query, id, req.UserId, totalAmount, status, shippingAddressJSON, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to insert order into database: %v", err)
	}

	for _, item := range items {
		orderId := uuid.NewString()
		orderQueryItems := `insert into order_items(
				id, order_id, product_id, quantity, price
				)values(
				  $1, $2, $3, $4, $5
				)`

		_, err := o.Db.Exec(orderQueryItems, orderId, id, item.ProductId, item.Quantity, item.Price)
		if err != nil {
			return nil, err
		}
	}

	resp := &pb.CreateOrderRequest{
		Items:           items,
		ShippingAddress: &pb.Address{},
		UserId:          id,
	}

	return resp, nil
}

func (o *OrderRepo) CancelOrder(req *pb.Id) (*pb.CancelOrder1, error) {
	resp := pb.CancelOrder1{}
	q := `update
				orders
		  set  
		  		deleted_at = $1
		  where 
		  		id = $2 and status != 'shipped'
		  RETURNING 
		  		id, status, updated_at 
		`

	err := o.Db.QueryRow(q, time.Now(), req.OrderId).Scan(
		&resp.Id,
		&resp.Status,
		&resp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (o *OrderRepo) UpdateOrderStatus(req *pb.UpdateOrderStatusRequest) (*pb.UpdateRespons, error) {
	resp := pb.UpdateRespons{}
	q := `
		update
			orders
		set
			status = $1, updated_at = $3
		where 
			id = $2 and deleted_at is null 
		returning id, status, updated_at `

	err := o.Db.QueryRow(q, req.Status, req.OrderId, time.Now()).Scan(
		&resp.Id,
		&resp.Status,
		&resp.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (o *OrderRepo) GetOrder(req *pb.Id) (*pb.OrderInfo, error) {
	q := `
	  SELECT 
		id, user_id, total_amount, status, shipping_address
	  FROM 
		orders
	  WHERE 
		id = $1
	`

	var resp pb.OrderInfo
	var shippingAddress string

	err := o.Db.QueryRow(q, req.OrderId).Scan(
		&resp.Id,
		&resp.UserId,
		&resp.TotalAmount,
		&resp.Status,
		&shippingAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get order from database: %v", err)
	}

	err = json.Unmarshal([]byte(shippingAddress), &resp.ShippingAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal shipping address: %v", err)
	}

	orderItemsQuery := `
	  SELECT 
		product_id, price
	  FROM 
		order_items
	  WHERE 
		order_id = $1
	`

	rows, err := o.Db.Query(orderItemsQuery, req.OrderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items from database: %v", err)
	}
	defer rows.Close()

	var items []*pb.OrderItem

	for rows.Next() {
		var item pb.OrderItem
		err := rows.Scan(&item.ProductId, &item.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %v", err)
		}
		items = append(items, &item)
	}

	resp.Items = items

	return &resp, nil
}

func (o *OrderRepo) PayOrder(req *pb.PayOrderRequest) (*pb.PaymentInfo, error) {
	resp := pb.PaymentInfo{}
	id := uuid.New()

	q2 := `
	SELECT 
			oi.total_amount
	FROM 
		order_items AS oi
	LEFT JOIN 
		orders AS o
	ON
		oi.order_id = o.id
	WHERE 
		o.deleted_at IS NULL

	`
	var amount float32
	err := o.Db.QueryRow(q2).Scan(&amount)
	if err != nil {
		return nil, err
	}

	q1 := `
		insert into payments
							(id,
							order_id,
							payment_method,
							created_at,
							status,	
							amount)
		values($1, $2, $3, $4, $5, $6, $7) 
		returning id, order_id, created_at, status
	`

	err = o.Db.QueryRow(q1, id.String(), req.OrderId, req.PaymentMethod, req.ExpiryDate, req.Cvv, time.Now(), req.Status, amount).Scan(
		&resp.PaymentId,
		&resp.OrderId,
		&resp.CreatedAt,
		&resp.Status,
	)

	if err != nil {
		return nil, err
	}

	return &pb.PaymentInfo{
		OrderId:       resp.OrderId,
		PaymentId:     resp.PaymentId,
		Amount:        amount,
		Status:        resp.Status,
		TransactionId: "1234",
		CreatedAt:     resp.CreatedAt,
	}, nil
}

func (o *OrderRepo) GetPaymentStatus(req *pb.GetPaymentStatusRequest) (*pb.PaymentInfo, error) {
    resp := pb.PaymentInfo{}
    q := `
        SELECT 
            id,
            order_id,
            status,
            transaction_id,
            created_at
        FROM 
            payments
        WHERE 
            order_id = $1
    `
    err := o.Db.QueryRow(q, req.OrderId).Scan(
        &resp.PaymentId,
        &resp.OrderId,
        &resp.Status,
        &resp.TransactionId,
        &resp.CreatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("no payment found for order_id %s", req.OrderId)
        }
        return nil, err
    }
    return &resp, nil
}


