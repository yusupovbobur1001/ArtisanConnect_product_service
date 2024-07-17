package service

import (
	"context"
	"database/sql"
	"fmt"
	"product_service/config"
	pbu "product_service/genproto/auth"
	pb "product_service/genproto/orders"
	"product_service/pkg"
	"product_service/storage/postgres"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	Auth  pbu.AuthClient
	Order *postgres.OrderRepo
}

func NewOrderService(db *sql.DB, cfg *config.Config) *OrderService {
	return &OrderService{
		Order: postgres.NewOrderRepo(db),
		Auth:  pkg.CreateAuthClient(cfg),
	}
}




func (o *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderRequest, error) {
	exist, err := o.Auth.ValidateUserId(ctx, &pbu.Id{Id: req.UserId})
	fmt.Println(exist, "+++++++++++++++++++++++++++++++++++++++++++++++++++")
	if !exist.Exist || err != nil {
		return nil, err
	}

	resp, err := o.Order.CreateOrder(req)
	fmt.Println(resp, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *OrderService) CancelOrder(ctx context.Context, req *pb.Id) (*pb.CancelOrder1, error) {
	resp, err := o.Order.CancelOrder(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *OrderService) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateRespons, error) {
	resp, err := o.Order.UpdateOrderStatus(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}



func (o *OrderService) GetOrder(ctx context.Context, req *pb.Id) (*pb.OrderInfo, error) {
	resp, err := o.Order.GetOrder(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *OrderService) PayOrder(ctx context.Context, req *pb.PayOrderRequest) (*pb.PaymentInfo, error) {
	resp, err := o.Order.PayOrder(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *OrderService) GetPaymentStatus(ctx context.Context, req *pb.GetPaymentStatusRequest) (*pb.PaymentInfo, error) {
	resp, err := o.Order.GetPaymentStatus(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}


