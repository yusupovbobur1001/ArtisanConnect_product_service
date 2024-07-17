package postgres

import (
	"fmt"
	pb "product_service/genproto/orders"
	"testing"
)

// func TestCreateOrder(t *testing.T) {
// 	db, err := ConnDB()
// 	if err != nil {
// 		panic(err)
// 	}

// 	order1 := NewOrderRepo(db)

// 	orders := []*pb.OrderItem{}
// 	order := pb.OrderItem{}

// 	order.ProductId = "b96b6aa9-57a2-4ebd-aa91-2abf1726f0d7"
// 	order.Price = 45.2
// 	orders = append(orders, &order)

// 	req := pb.CreateOrderRequest{
// 		Items:           orders,
// 		ShippingAddress: &pb.Address{},
// 		UserId:          "0e0b7e6d-3c4e-490a-8571-5097a2a7ffbf",
// 	}

// 	resp, err := order1.CreateOrder(&req)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(resp)
// }

func TestCancelOrder(t *testing.T) {
	db, err := ConnDB()
	if err != nil {
		panic(err)
	}

	order := NewOrderRepo(db)
	req := pb.Id{
		OrderId: "4497a3c5-fb1b-426f-9121-9e435ea7cb15",
	}

	resp, err := order.CancelOrder(&req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

func TestUpdateOrderStatus(t *testing.T) {
	db, err := ConnDB()
	if err != nil {
		panic(err)
	}
	order := NewOrderRepo(db)

	req := pb.UpdateOrderStatusRequest{
		OrderId: "b6c95cc4-8a22-40d9-8c12-52ee82ed68f8",
		Status:  "afdsf",
	}

	resp, err := order.UpdateOrderStatus(&req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

func TestGetOrders(t *testing.T) {
	db, err := ConnDB()
	if err != nil {
		panic(err)
	}
	order := NewOrderRepo(db)
	req := pb.Id{
		OrderId: "f48a02ab-ea9f-4ca1-abcf-311cb8a82a9b",
	}

	resp, err := order.GetOrder(&req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

func TestPayOrder(t *testing.T) {
	db, err := ConnDB()
	if err != nil {
		panic(err)
	}
	order := NewOrderRepo(db)
	req := pb.PayOrderRequest{
		OrderId:       "4497a3c5-fb1b-426f-9121-9e435ea7cb15",
		PaymentMethod: "bank",
		CardNumber:    "1111111111111111",
		ExpiryDate:    "01/28",
		Cvv:           "123",
		Status:        "dido",
	}

	resp, err := order.PayOrder(&req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

func TestGetPaymentStatus(t *testing.T) {
    db, err := ConnDB()
    if err != nil {
        panic(err)
    }
    order := NewOrderRepo(db)
    req := pb.GetPaymentStatusRequest{
        OrderId: "9bc8a809-920d-475c-8271-fa73b5c3e35d",
    }

    resp, err := order.GetPaymentStatus(&req)
    if err != nil {
        if err.Error() == "no payment found for order_id 9bc8a809-920d-475c-8271-fa73b5c3e35d" {
            fmt.Println("No payment found, as expected.")
        } else {
            panic(err)
        }
    } else {
        fmt.Println(resp)
    }
}

