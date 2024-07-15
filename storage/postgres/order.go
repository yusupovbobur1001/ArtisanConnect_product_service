package postgres

import (
	"database/sql"
	pb "product_service/genproto/order"
)

type OrderRepo struct {
	Db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{Db: db}
}

// bu func sion o`xshamadi keyin koramiz
// func (o *OrderRepo) CreateOrder(req *pb.CreateOrderRequest) (*pb.OrderInfo, error) {
// 	q1 := `	INSERT INTO orders (
// 								product_id,
// 								user_id,
// 								total_amount,
// 								status,
// 								shipping_address,
// 								created_at,
// 							   )
// 			VALUES ($1, $2, $3, $4, $5, $6)
// 			RETURNING id  `
// 	for _, k := range req.Items {
// 		err := o.Db.QueryRow(q1, )
// 	}
// }

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

	err := o.Db.QueryRow(q, req.OrderId).Scan(
		&resp.Id,
		&resp.Status,
		&resp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &resp, nil

}

func (o *OrderRepo) UpdateOrderStatus(req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderRespons, error) {
	q := `
		update
			orders
		set
			status = $1, updated_at = $3
		where 
			id = $2 and deleted_at is null 
		returning id, status, updated_at `


}
